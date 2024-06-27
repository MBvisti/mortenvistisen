package models

import (
	"context"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/google/uuid"
)

type userStorage interface {
	QueryUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)
	QueryUserByEmail(
		ctx context.Context,
		email string,
	) (domain.User, error)
	InsertUser(
		ctx context.Context,
		data domain.User,
	) (domain.User, error)
	UpdateUser(
		ctx context.Context,
		data domain.User,
	) (domain.User, error)
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

func (us UserService) ByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := us.userStorage.QueryUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (us UserService) UpdatePassword(
	ctx context.Context,
	userID uuid.UUID,
	password, confirmPassword string,
) (domain.User, error) {
	user, err := us.userStorage.QueryUserByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (us UserService) ConfirmEmail(
	ctx context.Context,
	userID uuid.UUID,
) (domain.User, error) {
	user, err := us.userStorage.QueryUserByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	user.ConfirmEmail()

	updatedUser, err := us.userStorage.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return updatedUser, nil
}

func (us UserService) New(
	ctx context.Context,
	name, mail, password, confirmPassword string,
) (domain.User, error) {
	user, err := domain.NewUser(name, mail, password, confirmPassword)
	if err != nil {
		return domain.User{}, err
	}

	hashedPassword, err := us.auth.HashAndPepperPassword(password)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = hashedPassword

	if _, err := us.userStorage.InsertUser(ctx, user); err != nil {
		slog.ErrorContext(ctx, "could not insert user to database", "error", err)
		return domain.User{}, err
	}

	return user, nil
}

// func (u UserSvc) New(
// 	ctx context.Context,
// 	data domain.NewUser,
// 	db userDatabase,
// 	v *validator.Validate,
// 	passwordPepper string,
// ) (domain.User, error) {
// 	mailAlreadyRegistered, err := db.DoesMailExists(ctx, data.Mail)
// 	if err != nil {
// 		telemetry.Logger.Error("could not check if email exists", "error", err)
// 		return domain.User{}, err
// 	}
//
// 	newUserData := NewUserValidation{
// 		ConfirmPassword: data.ConfirmPassword,
// 		Name:            data.Name,
// 		Mail:            data.Mail,
// 		MailRegistered:  mailAlreadyRegistered,
// 		Password:        data.Password,
// 	}
//
// 	if err := v.Struct(newUserData); err != nil {
// 		return domain.User{}, err
// 	}
//
// 	hashedPassword, err := u.auth.HashAndPepperPassword(newUserData.Password)
// 	if err != nil {
// 		telemetry.Logger.Error("error hashing and peppering password", "error", err)
// 		return domain.User{}, err
// 	}
//
// 	user, err := db.InsertUser(ctx, database.InsertUserParams{
// 		ID:        uuid.New(),
// 		CreatedAt: database.ConvertToPGTimestamptz(time.Now()),
// 		UpdatedAt: database.ConvertToPGTimestamptz(time.Now()),
// 		Name:      newUserData.Name,
// 		Mail:      newUserData.Mail,
// 		Password:  hashedPassword,
// 	})
// 	if err != nil {
// 		telemetry.Logger.Error("could not insert user", "error", err)
// 		return domain.User{}, err
// 	}
//
// 	return domain.User{
// 		ID:        user.ID,
// 		CreatedAt: database.ConvertFromPGTimestamptzToTime(user.CreatedAt),
// 		UpdatedAt: database.ConvertFromPGTimestamptzToTime(user.UpdatedAt),
// 		Name:      user.Name,
// 		Mail:      user.Mail,
// 	}, nil
// }
//
// type UpdateUserValidation struct {
// 	ConfirmPassword string `validate:"required,gte=8"`
// 	Password        string `validate:"required,gte=8"`
// 	Name            string `validate:"required,gte=2"`
// 	Mail            string `validate:"required,email"`
// }
//
// func ResetPasswordMatchValidation(sl validator.StructLevel) {
// 	data := sl.Current().Interface().(UpdateUserValidation)
//
// 	if data.ConfirmPassword != data.Password {
// 		sl.ReportError(
// 			data.ConfirmPassword,
// 			"",
// 			"ConfirmPassword",
// 			"",
// 			"confirm password must match password",
// 		)
// 	}
// }
//
// func (u UserSvc) UpdateUser(
// 	ctx context.Context,
// 	data domain.UpdateUser,
// 	db userDatabase,
// 	v *validator.Validate,
// 	passwordPepper string,
// ) (domain.User, error) {
// 	validatedData := UpdateUserValidation{
// 		ConfirmPassword: data.ConfirmPassword,
// 		Password:        data.Password,
// 		Name:            data.Name,
// 		Mail:            data.Mail,
// 	}
//
// 	if err := v.Struct(validatedData); err != nil {
// 		return domain.User{}, err
// 	}
//
// 	hashedPassword, err := u.auth.HashAndPepperPassword(validatedData.Password)
// 	if err != nil {
// 		telemetry.Logger.Error("error hashing and peppering password", "error", err)
// 		return domain.User{}, err
// 	}
//
// 	updatedUser, err := db.UpdateUser(ctx, database.UpdateUserParams{
// 		UpdatedAt: database.ConvertToPGTimestamptz(time.Now()),
// 		Name:      data.Name,
// 		Mail:      data.Mail,
// 		Password:  hashedPassword,
// 		ID:        data.ID,
// 	})
// 	if err != nil {
// 		telemetry.Logger.Error("could not insert user", "error", err)
// 		return domain.User{}, err
// 	}
//
// 	return domain.User{
// 		ID:        updatedUser.ID,
// 		CreatedAt: database.ConvertFromPGTimestamptzToTime(updatedUser.CreatedAt),
// 		UpdatedAt: database.ConvertFromPGTimestamptzToTime(updatedUser.UpdatedAt),
// 		Name:      updatedUser.Name,
// 		Mail:      updatedUser.Mail,
// 	}, nil
// }

// func (a AuthSvc) AuthenticateUser(
// 	ctx context.Context,
// 	data AuthenticateUserPayload,
// 	passwordPepper string,
// ) (domain.User, error) {
// 	user, err := db.QueryUserByMail(ctx, data.Email)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return domain.User{}, ErrUserNotExist
// 		}
//
// 		slog.Error("could not query user by mail", "error", err)
// 		return domain.User{}, err
// 	}
//
// 	if verifiedAt := user.MailVerifiedAt; !verifiedAt.Valid {
// 		return domain.User{}, ErrEmailNotValidated
// 	}
//
// 	if err := validatePassword(validatePasswordPayload{
// 		hashedpassword: user.Password,
// 		password:       data.Password,
// 	}, passwordPepper); err != nil {
// 		return domain.User{}, ErrPasswordNotMatch
// 	}
//
// 	return domain.User{
// 		ID:        user.ID,
// 		CreatedAt: database.ConvertFromPGTimestamptzToTime(user.CreatedAt),
// 		UpdatedAt: database.ConvertFromPGTimestamptzToTime(user.UpdatedAt),
// 		Name:      user.Name,
// 		Mail:      user.Mail,
// 	}, nil
// }
