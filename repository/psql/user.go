package psql

import (
	"context"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/repository/psql/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p Postgres) QueryUserByID(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	user, err := p.db.QueryUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
		Name:           user.Name,
		Mail:           user.Mail,
		MailVerifiedAt: user.MailVerifiedAt.Time,
	}, nil
}

func (p Postgres) QueryUserByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	user, err := p.db.QueryUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
		Name:           user.Name,
		Mail:           user.Mail,
		MailVerifiedAt: user.MailVerifiedAt.Time,
	}, nil
}

func (p Postgres) InsertUser(
	ctx context.Context,
	data domain.User,
) (domain.User, error) {
	createdAt := pgtype.Timestamptz{
		Time:  data.CreatedAt,
		Valid: true,
	}
	updatedAt := pgtype.Timestamptz{
		Time:  data.UpdatedAt,
		Valid: true,
	}

	_, err := p.db.InsertUser(ctx, database.InsertUserParams{
		ID:        data.ID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      data.Name,
		Mail:      data.Mail,
		Password:  data.GetPassword(),
	})
	if err != nil {
		return domain.User{}, err
	}

	return data, nil
}

func (p Postgres) UpdateUser(
	ctx context.Context,
	data domain.User,
) (domain.User, error) {
	updatedAt := pgtype.Timestamptz{
		Time:  data.UpdatedAt,
		Valid: true,
	}

	_, err := p.db.UpdateUser(ctx, database.UpdateUserParams{
		ID:        data.ID,
		UpdatedAt: updatedAt,
		Name:      data.Name,
		Mail:      data.Mail,
		Password:  data.GetPassword(),
	})
	if err != nil {
		return domain.User{}, err
	}

	return data, nil
}
