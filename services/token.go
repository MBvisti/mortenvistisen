package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hash"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	ScopeEmailVerification           = "email_verification"
	ScopeSubscriberEmailVerification = "subscriber_email_verification"
	ScopeUnsubscribe                 = "unsubscribe"
	ScopeResetPassword               = "password_reset"
)

type tokenServiceStorage interface {
	InsertSubscriberToken(
		ctx context.Context,
		hash, scope string, expiresAt time.Time,
		subscriberID uuid.UUID,
	) error
	InsertToken(
		ctx context.Context,
		hash, scope string, expiresAt time.Time,
		userID uuid.UUID,
	) error
	QueryTokenByHash(ctx context.Context, hash string) (database.Token, error)
	QuerySubscriberTokenByHash(ctx context.Context, hash string) (database.SubscriberToken, error)
	DeleteTokenByHash(ctx context.Context, hash string) error
	DeleteTokenBySubID(ctx context.Context, id uuid.UUID) error
}

type TokenSvc struct {
	storage tokenServiceStorage
	hasher  hash.Hash
}

func NewTokenSvc(
	storage tokenServiceStorage,
	tokenSigningKey string,
) *TokenSvc {
	h := hmac.New(sha256.New, []byte(tokenSigningKey))

	return &TokenSvc{
		storage,
		h,
	}
}

type tokenPair struct {
	plain  string
	hashed string
}

func (svc *TokenSvc) create() (tokenPair, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return tokenPair{}, err
	}

	plainText := base64.URLEncoding.EncodeToString(b)
	hashedToken := svc.hash(plainText)

	return tokenPair{
		plainText,
		hashedToken,
	}, nil
}

func (svc *TokenSvc) hash(token string) string {
	svc.hasher.Reset()
	svc.hasher.Write([]byte(token))
	b := svc.hasher.Sum(nil)

	return base64.URLEncoding.EncodeToString(b)
}

func (svc *TokenSvc) CreateSubscriptionToken(
	ctx context.Context,
	subscriberID uuid.UUID,
) (string, error) {
	tokenPair, err := svc.create()
	if err != nil {
		return "", err
	}

	tkn := NewToken(
		ScopeSubscriberEmailVerification,
		time.Now().Add(72*time.Hour),
		tokenPair.hashed,
		tokenPair.plain,
	)

	if err := svc.storage.InsertSubscriberToken(ctx, tkn.Hash, tkn.GetScope(), tkn.GetExpirationTime(), subscriberID); err != nil {
		slog.ErrorContext(
			ctx,
			"could not insert a subscriber token",
			"error",
			err,
			"subscriber_id",
			subscriberID,
		)
		return "", err
	}

	return tkn.GetPlainText(), nil
}

func (svc *TokenSvc) CreateEmailVerificationToken(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	tokenPair, err := svc.create()
	if err != nil {
		return "", err
	}

	tkn := NewToken(
		ScopeEmailVerification,
		time.Now().Add(48*time.Hour),
		tokenPair.hashed,
		tokenPair.plain,
	)

	if err := svc.storage.InsertToken(ctx, tkn.Hash, tkn.GetScope(), tkn.GetExpirationTime(), userID); err != nil {
		slog.ErrorContext(
			ctx,
			"could not insert a token",
			"error",
			err,
			"user_id",
			userID,
		)
		return "", err
	}

	return tkn.GetPlainText(), nil
}

func (svc *TokenSvc) CreateResetPasswordToken(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	tokenPair, err := svc.create()
	if err != nil {
		return "", err
	}

	tkn := NewToken(
		ScopeResetPassword,
		time.Now().Add(24*time.Hour),
		tokenPair.hashed,
		tokenPair.plain,
	)

	if err := svc.storage.InsertToken(ctx, tkn.Hash, tkn.GetScope(), tkn.GetExpirationTime(), userID); err != nil {
		return "", err
	}

	return tkn.GetPlainText(), nil
}

func (svc *TokenSvc) CreateUnsubscribeToken(
	ctx context.Context,
	subscriberID uuid.UUID,
) (string, error) {
	tokenPair, err := svc.create()
	if err != nil {
		return "", err
	}

	tkn := NewToken(
		ScopeUnsubscribe,
		time.Now().Add(168*time.Hour), // allow 7 days for an unsubscribe link to be valid
		tokenPair.hashed,
		tokenPair.plain,
	)

	if err := svc.storage.InsertSubscriberToken(ctx, tkn.Hash, tkn.GetScope(), tkn.GetExpirationTime(), subscriberID); err != nil {
		return "", err
	}

	return tkn.GetPlainText(), nil
}

func (svc *TokenSvc) Validate(ctx context.Context, token string) error {
	tkn, err := svc.storage.QueryTokenByHash(ctx, svc.hash(token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.InfoContext(ctx, "a token was requested that could not be found", "token", tkn)

			slog.ErrorContext(ctx, "could not query token by hash", "error", err)
			return errors.Join(ErrTokenNotExist, err)
		}

		return err
	}

	if time.Now().After(tkn.ExpiresAt.Time) {
		return ErrTokenExpired
	}

	return nil
}

func (svc *TokenSvc) ValidateSubscriber(ctx context.Context, token string) error {
	tkn, err := svc.storage.QuerySubscriberTokenByHash(ctx, svc.hash(token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.InfoContext(ctx, "a token was requested that could not be found", "token", tkn)

			slog.ErrorContext(ctx, "could not query token by hash", "error", err)
			return errors.Join(ErrTokenNotExist, err)
		}

		return err
	}

	if time.Now().After(tkn.ExpiresAt.Time) {
		return ErrTokenExpired
	}

	return nil
}

func (svc *TokenSvc) GetAssociatedUserID(ctx context.Context, token string) (uuid.UUID, error) {
	tkn, err := svc.storage.QueryTokenByHash(ctx, svc.hash(token))
	if err != nil {
		return uuid.UUID{}, err
	}

	return tkn.UserID, nil
}

func (svc *TokenSvc) GetAssociatedSubscriberID(
	ctx context.Context,
	token string,
) (uuid.UUID, error) {
	tkn, err := svc.storage.QuerySubscriberTokenByHash(ctx, svc.hash(token))
	if err != nil {
		return uuid.UUID{}, err
	}

	return tkn.SubscriberID, nil
}

func (svc *TokenSvc) Delete(ctx context.Context, token string) error {
	err := svc.storage.DeleteTokenByHash(ctx, svc.hash(token))
	if err != nil {
		return err
	}

	return nil
}

func (svc *TokenSvc) DeleteSubscriberToken(ctx context.Context, subscriberID uuid.UUID) error {
	err := svc.storage.DeleteTokenBySubID(ctx, subscriberID)
	if err != nil {
		return err
	}

	return nil
}

type Token struct {
	scope     string
	expiresAt time.Time
	Hash      string
	plain     string
}

func NewToken(
	scope string,
	expiresAt time.Time,
	Hash string,
	plain string,
) Token {
	return Token{
		scope,
		expiresAt,
		Hash,
		plain,
	}
}

func (t *Token) GetPlainText() string {
	return t.plain
}

func (t *Token) GetExpirationTime() time.Time {
	return t.expiresAt
}

func (t *Token) GetScope() string {
	return t.scope
}

func (t *Token) IsValid() bool {
	return time.Now().Before(t.expiresAt)
}
