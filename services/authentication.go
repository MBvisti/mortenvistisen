package services

import (
	"context"
	"encoding/gob"
	"errors"
	"log/slog"
	"net/http"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthSvc struct {
	userModel models.UserSvc
}

func hashAndPepperPassword(password, passwordPepper string) (string, error) {
	passwordBytes := []byte(password + passwordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		slog.Error("could not hash password", "error", err)
		return "", err
	}

	return string(hashedBytes), nil
}

type validatePasswordPayload struct {
	hashedpassword string
	password       string
}

func validatePassword(data validatePasswordPayload, passwordPepper string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(data.hashedpassword), []byte(data.password+passwordPepper)); err != nil {
		slog.Error("could not validate password", "error", err)
		return err
	}

	return nil
}

type AuthenticateUserPayload struct {
	Email    string
	Password string
}

func (a AuthSvc) AuthenticateUser(
	ctx context.Context,
	data AuthenticateUserPayload,
	passwordPepper string,
) (domain.User, error) {
	user, err := db.QueryUserByMail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotExist
		}

		slog.Error("could not query user by mail", "error", err)
		return domain.User{}, err
	}

	if verifiedAt := user.MailVerifiedAt; !verifiedAt.Valid {
		return domain.User{}, ErrEmailNotValidated
	}

	if err := validatePassword(validatePasswordPayload{
		hashedpassword: user.Password,
		password:       data.Password,
	}, passwordPepper); err != nil {
		return domain.User{}, ErrPasswordNotMatch
	}

	return domain.User{
		ID:        user.ID,
		CreatedAt: database.ConvertFromPGTimestamptzToTime(user.CreatedAt),
		UpdatedAt: database.ConvertFromPGTimestamptzToTime(user.UpdatedAt),
		Name:      user.Name,
		Mail:      user.Mail,
	}, nil
}

func CreateAuthenticatedSession(
	session sessions.Session,
	userID uuid.UUID,
	cfg config.Cfg,
) *sessions.Session {
	gob.Register(uuid.UUID{})

	session.Options.HttpOnly = true
	session.Options.Domain = cfg.App.AppHost
	session.Options.Secure = true
	session.Options.MaxAge = 86400

	session.Values["user_id"] = userID
	session.Values["authenticated"] = true
	if userID.String() == "d7b5e3eb-a799-4f7e-8139-905c20e8c8e9" {
		session.Values["is_admin"] = true
	} else {
		session.Values["is_admin"] = false
	}

	return &session
}

func IsAuthenticated(r *http.Request, authStore *sessions.CookieStore) (bool, uuid.UUID, error) {
	gob.Register(uuid.UUID{})
	session, err := authStore.Get(r, "ua")
	if err != nil {
		slog.Error("could not get session", "error", err)
		return false, uuid.UUID{}, err
	}

	if session.Values["authenticated"] == nil {
		return false, uuid.UUID{}, err
	}

	return session.Values["authenticated"].(bool), session.Values["user_id"].(uuid.UUID), nil
}

func IsAdmin(r *http.Request, authStore *sessions.CookieStore) (bool, error) {
	gob.Register(uuid.UUID{})
	session, err := authStore.Get(r, "ua")
	if err != nil {
		slog.Error("could not get session", "error", err)
		return false, err
	}

	return session.Values["is_admin"].(bool), nil
}
