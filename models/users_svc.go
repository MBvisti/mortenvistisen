package models

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/validation"
	"github.com/google/uuid"
)

type userStorage interface {
	QueryUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (User, error)
	QueryUserByEmail(
		ctx context.Context,
		email string,
	) (User, error)
	InsertUser(
		ctx context.Context,
		data User,
	) (User, error)
	UpdateUser(
		ctx context.Context,
		data User,
	) (User, error)
}

type authService interface {
	HashAndPepperPassword(password string) (string, error)
	ValidatePassword(password, hashedPassword string) error
}

type UserService struct {
	auth        authService
	userStorage userStorage
}

func NewUserSvc(auth authService, usrStorage userStorage) UserService {
	return UserService{auth, usrStorage}
}

func (us UserService) ByEmail(ctx context.Context, email string) (User, error) {
	user, err := us.userStorage.QueryUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (us UserService) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	password, confirmPassword string,
) (User, error) {
	user, err := us.userStorage.QueryUserByID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (us UserService) ConfirmEmail(
	ctx context.Context,
	userID uuid.UUID,
) (User, error) {
	user, err := us.userStorage.QueryUserByID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	now := time.Now()

	user.MailVerifiedAt = now
	user.UpdatedAt = now

	updatedUser, err := us.userStorage.UpdateUser(ctx, user)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (us UserService) New(
	ctx context.Context,
	name, mail, password, confirmPassword string,
) (User, error) {
	t := time.Now()
	user := User{uuid.New(), t, t, name, mail, time.Time{}, password}
	if err := validation.Validate(user, CreateUserValidations(confirmPassword)); err != nil {
		return User{}, errors.Join(ErrFailValidation, err)
	}

	hashedPassword, err := us.auth.HashAndPepperPassword(password)
	if err != nil {
		return User{}, err
	}

	user.Password = hashedPassword

	if _, err := us.userStorage.InsertUser(ctx, user); err != nil {
		slog.ErrorContext(ctx, "could not insert user to database", "error", err)
		return User{}, err
	}

	return user, nil
}
