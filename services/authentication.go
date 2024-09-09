package services

import (
	"context"
	"encoding/gob"
	"errors"
	"log/slog"
	"net/http"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCouldNotHashPepperPW     = errors.New("could not hash password")
	ErrCouldNotValidatePassword = errors.New(
		"could not validate password",
	)
	ErrCouldNotGetAuthenticatedSession = errors.New("could not get session")
)

type authStorage interface {
	QueryUserByEmail(
		ctx context.Context,
		email string,
	) (models.User, error)
}

type Auth struct {
	cfg            config.Cfg
	passwordPepper string
	cookieStore    *sessions.CookieStore
	storage        authStorage
}

func NewAuth(cfg config.Cfg, storage authStorage) Auth {
	pwPepper := cfg.Auth.PasswordPepper

	authSessionStore := sessions.NewCookieStore(
		[]byte(cfg.Auth.SessionKey),
		[]byte(cfg.Auth.SessionEncryptionKey),
	)
	gob.Register(uuid.UUID{})

	return Auth{
		cfg,
		pwPepper,
		authSessionStore,
		storage,
	}
}

func (a Auth) HashAndPepperPassword(password string) (string, error) {
	passwordBytes := []byte(password + a.passwordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		passwordBytes,
		bcrypt.DefaultCost,
	)
	if err != nil {
		slog.Error("could not hash password", "error", err)
		return "", errors.Join(ErrCouldNotHashPepperPW, err)
	}

	return string(hashedBytes), nil
}

func (a Auth) ValidatePassword(password, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+a.passwordPepper)); err != nil {
		slog.Error(
			"could not validate password",
			"error",
			err,
			"pass",
			password,
		)
		return errors.Join(ErrCouldNotValidatePassword, err)
	}

	return nil
}

func (a Auth) AuthenticateUser(
	ctx context.Context,
	req *http.Request,
	res http.ResponseWriter,
	remember bool,
	mail string,
	password string,
) error {
	user, err := a.storage.QueryUserByEmail(ctx, mail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(ErrUserNotExist, err)
		}

		return err
	}

	if user.MailVerifiedAt.IsZero() {
		return ErrEmailNotValidated
	}

	if err := a.ValidatePassword(password, user.Password); err != nil {
		return err
	}

	return a.CreateAuthenticatedSession(req, res, user.ID, remember)
}

func (a Auth) CreateAuthenticatedSession(
	req *http.Request,
	res http.ResponseWriter,
	userID uuid.UUID,
	remember bool,
) error {
	session, err := a.cookieStore.Get(req, "mortenvistisen-ua")
	if err != nil {
		return errors.Join(ErrCouldNotGetAuthenticatedSession, err)
	}

	session.Options.HttpOnly = true
	session.Options.Domain = a.cfg.App.AppHost
	session.Options.Secure = true
	if remember {
		session.Options.MaxAge = 2 * 604800 // auth sess valid 2 week
	} else {
		session.Options.MaxAge = 604800 // auth sess valid 1 week
	}

	session.Values["user_id"] = userID
	session.Values["authenticated"] = true
	if userID.String() == "d7b5e3eb-a799-4f7e-8139-905c20e8c8e9" {
		session.Values["is_admin"] = true
	} else {
		session.Values["is_admin"] = false
	}

	return session.Save(req, res)
}

func (a Auth) IsAuthenticated(r *http.Request) (bool, uuid.UUID, error) {
	session, err := a.cookieStore.Get(r, "mortenvistisen-ua")
	if err != nil {
		slog.Error("could not get session", "error", err)
		return false, uuid.UUID{}, errors.Join(
			ErrCouldNotGetAuthenticatedSession,
			err,
		)
	}

	if session.Values["authenticated"] == nil {
		return false, uuid.UUID{}, err
	}

	return session.Values["authenticated"].(bool), session.Values["user_id"].(uuid.UUID), nil
}

func (a Auth) IsAdmin(r *http.Request) (bool, error) {
	session, err := a.cookieStore.Get(r, "mortenvistisen-ua")
	if err != nil {
		slog.Error("could not get session", "error", err)
		return false, errors.Join(ErrCouldNotGetAuthenticatedSession, err)
	}

	return session.Values["is_admin"].(bool), nil
}
