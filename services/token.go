package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"time"

	"github.com/google/uuid"
)

const (
	ScopeEmailVerification           = "email_verification"
	ScopeSubscriberEmailVerification = "subscriber_email_verification"
	ScopeUnsubscribe                 = "unsubscribe"
	ScopeResetPassword               = "password_reset"
)

type TokenServiceStorage interface {
	InsertSubscriberToken(
		ctx context.Context,
		hash, scope string, expiresAt time.Time,
		subscriberID uuid.UUID,
	) error
}

type TokenSvc struct {
	storage TokenServiceStorage
	hasher  hash.Hash
}

func NewTokenSvc(
	storage TokenServiceStorage,
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
		return "", err
	}

	return tkn.GetPlainText(), nil
}

func (svc *TokenSvc) CreateResetPasswordToken() (string, error) {
	tokenPair, err := svc.create()
	if err != nil {
		return "", err
	}

	return tokenPair.plain, nil
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
