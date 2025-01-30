package models

import (
	"context"
	"errors"
	"time"

	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Email           string
	EmailVerifiedAt time.Time
	IsAdmin         bool
}

func (ue User) IsVerified() bool {
	return !ue.EmailVerifiedAt.IsZero()
}

func GetUserByEmail(
	ctx context.Context,
	email string,
	dbtx db.DBTX,
) (User, error) {
	user, err := db.Stmts.QueryUserByMail(ctx, dbtx, email)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:              user.ID,
		CreatedAt:       user.CreatedAt.Time,
		UpdatedAt:       user.UpdatedAt.Time,
		Email:           user.Mail,
		EmailVerifiedAt: user.MailVerifiedAt.Time,
		IsAdmin:         false, // TODO: FIX
	}, nil
}

type NewUserPayload struct {
	Email           string `validate:"required,email"`
	Password        string `validate:"required,gte=6"`
	ConfirmPassword string `validate:"required,gte=6"`
}

func GetUser(
	ctx context.Context,
	id uuid.UUID,
	dbtx db.DBTX,
) (User, error) {
	usr, err := db.Stmts.QueryUserByID(ctx, dbtx, id)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:              id,
		CreatedAt:       usr.CreatedAt.Time,
		UpdatedAt:       usr.UpdatedAt.Time,
		Email:           usr.Mail,
		EmailVerifiedAt: usr.MailVerifiedAt.Time,
		IsAdmin:         false,
	}, nil
}

func NewUser(
	ctx context.Context,
	data NewUserPayload,
	dbtx db.DBTX,
	hash func(password string) (string, error),
) (User, error) {
	if err := validate.Struct(data); err != nil {
		return User{}, errors.Join(ErrDomainValidation, err)
	}

	usr := User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     data.Email,
	}

	hashedPassword, err := hash(data.Password)
	if err != nil {
		return User{}, err
	}

	_, err = db.Stmts.InsertUser(ctx, dbtx, db.InsertUserParams{
		ID:        usr.ID,
		CreatedAt: pgtype.Timestamptz{Time: usr.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: usr.UpdatedAt, Valid: true},
		Mail:      usr.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return User{}, err
	}

	return usr, nil
}

type UpdateUserPayload struct {
	ID             uuid.UUID `validate:"required,uuid"`
	UpdatedAt      time.Time `validate:"required"`
	Name           string    `validate:"required,gte=2,lte=25"`
	Email          string    `validate:"required,email"`
	EmailUpdatedAt time.Time
}

func UpdateUser(
	ctx context.Context,
	data UpdateUserPayload,
	dbtx db.DBTX,
) (User, error) {
	// validate payload
	if err := validate.Struct(data); err != nil {
		return User{}, errors.Join(ErrDomainValidation, err)
	}

	updatedUsr, err := db.Stmts.UpdateUser(ctx, dbtx, db.UpdateUserParams{
		ID: data.ID,
		UpdatedAt: pgtype.Timestamptz{
			Time:  data.UpdatedAt,
			Valid: true,
		},
		Mail: data.Email,
	})
	if err != nil {
		return User{}, err
	}

	return User{
		ID:              updatedUsr.ID,
		CreatedAt:       updatedUsr.CreatedAt.Time,
		UpdatedAt:       updatedUsr.UpdatedAt.Time,
		Email:           updatedUsr.Mail,
		EmailVerifiedAt: updatedUsr.MailVerifiedAt.Time,
		IsAdmin:         false,
	}, nil
}

type UpdateUserPasswordPayload struct {
	ID        uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
	Password  string    `validate:"required"`
}

func UpdateUserPassword(
	ctx context.Context,
	data UpdateUserPasswordPayload,
	q func(
		ctx context.Context,
		userID uuid.UUID,
		newPassword string,
		updatedAt time.Time,
	) error,
) error {
	// validate payload
	if err := validate.Struct(data); err != nil {
		return errors.Join(ErrDomainValidation, err)
	}

	return q(ctx, data.ID, data.Password, data.UpdatedAt)
}

type MakeUserAdminPayload struct {
	UserID    uuid.UUID `validate:"required,uuid"`
	UpdatedAt time.Time `validate:"required"`
	ActorID   uuid.UUID `validate:"required,uuid"`
}

func MakeUserAdmin(
	ctx context.Context,
	data MakeUserAdminPayload,
	dbtx db.DBTX,
) (User, error) {
	if err := validate.Struct(data); err != nil {
		return User{}, errors.Join(ErrDomainValidation, err)
	}

	actor, err := GetUser(ctx, data.ActorID, dbtx)
	if err != nil {
		return User{}, err
	}
	if !actor.IsAdmin {
		return User{}, ErrMustBeAdmin
	}

	user, err := db.Stmts.UpdateUserIsAdmin(
		ctx,
		dbtx,
		db.UpdateUserIsAdminParams{
			ID: data.UserID,
			// IsAdmin:   true,
			UpdatedAt: pgtype.Timestamptz{Time: data.UpdatedAt, Valid: true},
		},
	)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:              user.ID,
		CreatedAt:       user.CreatedAt.Time,
		UpdatedAt:       user.UpdatedAt.Time,
		Email:           user.Mail,
		EmailVerifiedAt: user.MailVerifiedAt.Time,
		IsAdmin:         true,
	}, nil
}
