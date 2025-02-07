package seeds

import (
	"math/rand"
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/models/internal/db"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/net/context"
)

type userSeedData struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Mail           string
	MailVerifiedAt time.Time
	IsAdmin        bool
}

type userSeedOption func(*userSeedData)

func WithUserID(id uuid.UUID) userSeedOption {
	return func(usd *userSeedData) {
		usd.ID = id
	}
}

func WithUserCreatedAt(createdAt time.Time) userSeedOption {
	return func(usd *userSeedData) {
		usd.CreatedAt = createdAt
	}
}

func WithUserUpdatedAt(updatedAt time.Time) userSeedOption {
	return func(usd *userSeedData) {
		usd.UpdatedAt = updatedAt
	}
}

func WithUserEmail(email string) userSeedOption {
	return func(usd *userSeedData) {
		usd.Mail = email
	}
}

func WithUserEmailVerifiedAt(emailVerifiedAt time.Time) userSeedOption {
	return func(usd *userSeedData) {
		usd.MailVerifiedAt = emailVerifiedAt
	}
}

func WithUserIsAdmin(isAdmin bool) userSeedOption {
	return func(usd *userSeedData) {
		usd.IsAdmin = isAdmin
	}
}

func (s Seeder) PlantUser(
	ctx context.Context,
	opts ...userSeedOption,
) (models.User, error) {
	trueOrFalse := rand.Float32() < 0.5
	var emailVerifiedAt time.Time
	if trueOrFalse {
		emailVerifiedAt = time.Now()
	}

	data := &userSeedData{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Mail:           faker.Email(),
		MailVerifiedAt: emailVerifiedAt,
		IsAdmin:        trueOrFalse,
	}

	for _, opt := range opts {
		opt(data)
	}

	user, err := models.NewUser(ctx, models.NewUserPayload{
		Email:           data.Mail,
		Password:        "password",
		ConfirmPassword: "password",
	}, s.dbtx, services.HashAndPepperPassword)
	if err != nil {
		return models.User{}, err
	}

	if !data.MailVerifiedAt.IsZero() {
		if err := db.Stmts.VerifyUserMail(ctx, s.dbtx, db.VerifyUserMailParams{
			Mail: data.Mail,
			UpdatedAt: pgtype.Timestamptz{
				Time:  data.UpdatedAt,
				Valid: true,
			},
			MailVerifiedAt: pgtype.Timestamptz{
				Time:  data.MailVerifiedAt,
				Valid: true,
			},
		}); err != nil {
			return models.User{}, err
		}
	}

	if data.IsAdmin {
		if _, err := db.Stmts.UpdateUserIsAdmin(ctx, s.dbtx, db.UpdateUserIsAdminParams{
			ID: user.ID,
			// IsAdmin: data.IsAdmin,
			UpdatedAt: pgtype.Timestamptz{
				Time:  data.UpdatedAt,
				Valid: true,
			},
		}); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (s Seeder) PlantUsers(
	ctx context.Context,
	amount int,
) ([]models.User, error) {
	users := make([]models.User, amount)

	for i := range amount {
		usr, err := s.PlantUser(ctx)
		if err != nil {
			return nil, err
		}

		users[i] = usr
	}

	return users, nil
}
