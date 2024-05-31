package authentication

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/controllers/misc"
	"github.com/MBvisti/mortenvistisen/domain"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/MBvisti/mortenvistisen/views/validation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
)

func CreateAuthenticatedSession(ctx echo.Context) error {
	return authentication.LoginPage(authentication.LoginPageProps{
		CsrfToken: csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

type UserLoginPayload struct {
	Mail       string `form:"email"`
	Password   string `form:"password"`
	RememberMe string `form:"remember_me"`
}

func StoreAuthenticatedSession(
	ctx echo.Context,
	db database.Queries,
	cfg config.Cfg,
	authStore *sessions.CookieStore,
) error {
	var payload UserLoginPayload
	if err := ctx.Bind(&payload); err != nil {
		slog.ErrorContext(ctx.Request().Context(), "could not parse UserLoginPayload", "error", err)

		return authentication.LoginResponse(true).Render(views.ExtractRenderDeps(ctx))
	}

	authenticatedUser, err := services.AuthenticateUser(
		ctx.Request().Context(), services.AuthenticateUserPayload{
			Email:    payload.Mail,
			Password: payload.Password,
		}, &db, cfg.Auth.PasswordPepper)
	if err != nil {
		errMsg := "An error occurred while trying to authenticate you. Please try again."

		switch err {
		case services.ErrPasswordNotMatch, services.ErrUserNotExist:
			errMsg = "The password you entered is incorrect."
		case services.ErrEmailNotValidated:
			errMsg = "You need to verify your email before you can log in. Please check your inbox for a verification email."
		}
		return authentication.LoginForm(csrf.Token(ctx.Request()), authentication.LoginFormProps{
			HasError: true,
			ErrMsg:   errMsg,
		}).Render(views.ExtractRenderDeps(ctx))
	}

	session, err := authStore.Get(ctx.Request(), "ua")
	if err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not get auth session", "error", err)
		return misc.InternalError(ctx)
	}

	authSession := services.CreateAuthenticatedSession(*session, authenticatedUser.ID, cfg)
	if err := authSession.Save(ctx.Request(), ctx.Response()); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not save auth session", "error", err)
		return misc.InternalError(ctx)
	}

	return authentication.LoginResponse(false).Render(views.ExtractRenderDeps(ctx))
}

func CreatePasswordReset(ctx echo.Context) error {
	return authentication.ForgottenPasswordPage(authentication.ForgottenPasswordPageProps{
		CsrfToken: csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

type StorePasswordResetPayload struct {
	Mail string `form:"email"`
}

func StorePasswordReset(
	ctx echo.Context,
	db database.Queries,
	tknManager tokens.Manager,
	cfg config.Cfg,
	queueClient *river.Client[pgx.Tx],
) error {
	var payload StorePasswordResetPayload
	if err := ctx.Bind(&payload); err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}

	user, err := db.QueryUserByMail(ctx.Request().Context(), payload.Mail)
	if err != nil {
		failureOccurred := true
		if errors.Is(err, pgx.ErrNoRows) {
			failureOccurred = false
		}

		return authentication.ForgottenPasswordSuccess(failureOccurred).
			Render(views.ExtractRenderDeps(ctx))
	}

	generatedTkn, err := tknManager.GenerateToken()
	if err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}

	resetPWToken := tokens.CreateResetPasswordToken(
		generatedTkn.PlainTextToken,
		generatedTkn.HashedToken,
	)

	if err := db.StoreToken(ctx.Request().Context(), database.StoreTokenParams{
		ID:        uuid.New(),
		CreatedAt: database.ConvertToPGTimestamptz(time.Now()),
		Hash:      resetPWToken.Hash,
		ExpiresAt: database.ConvertToPGTimestamptz(resetPWToken.GetExpirationTime()),
		Scope:     resetPWToken.GetScope(),
		UserID:    user.ID,
	}); err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}

	// TODO fix this error flow
	pwResetMail := &templates.PasswordResetMail{
		ResetPasswordLink: fmt.Sprintf(
			"%s://%s/reset-password?token=%s",
			cfg.App.AppScheme,
			cfg.App.AppHost,
			resetPWToken.GetPlainText(),
		),
	}

	textVersion, err := pwResetMail.GenerateTextVersion()
	if err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}
	htmlVersion, err := pwResetMail.GenerateHtmlVersion()
	if err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}

	_, err = queueClient.Insert(ctx.Request().Context(), queue.EmailJobArgs{
		To:          user.Mail,
		From:        cfg.App.DefaultSenderSignature,
		Subject:     "Password Reset Request",
		TextVersion: textVersion,
		HtmlVersion: htmlVersion,
	}, nil)
	if err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}

	return authentication.ForgottenPasswordSuccess(false).Render(views.ExtractRenderDeps(ctx))
}

type PasswordResetToken struct {
	Token string `query:"token"`
}

func CreateResetPassword(ctx echo.Context) error {
	var passwordResetToken PasswordResetToken
	if err := ctx.Bind(&passwordResetToken); err != nil {
		return misc.InternalError(ctx)
	}

	return authentication.ResetPasswordPage(authentication.ResetPasswordPageProps{
		ResetToken: passwordResetToken.Token,
		CsrfToken:  csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

type ResetPasswordPayload struct {
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
	Token           string `form:"token"`
}

func StoreResetPassword(
	ctx echo.Context,
	db database.Queries,
	tknManager tokens.Manager,
	cfg config.Cfg,
	v *validator.Validate,
) error {
	var payload ResetPasswordPayload
	if err := ctx.Bind(&payload); err != nil {
		return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
			HasError: true,
			Msg:      "An error occurred while trying to reset your password. Please try again.",
		}).Render(views.ExtractRenderDeps(ctx))
	}

	hashedToken := tknManager.Hash(payload.Token)

	token, err := db.QueryTokenByHash(ctx.Request().Context(), hashedToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
				HasError: true,
				Msg:      "The token is invalid. Please request a new one.",
			}).Render(views.ExtractRenderDeps(ctx))
		}

		return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
			HasError: true,
			Msg:      "An error occurred while trying to reset your password. Please try again.",
		}).Render(views.ExtractRenderDeps(ctx))
	}

	if database.ConvertFromPGTimestamptzToTime(token.ExpiresAt).Before(time.Now()) &&
		token.Scope != tokens.ScopeResetPassword {
		return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
			HasError: true,
			Msg:      "The token has expired. Please request a new one.",
		}).Render(views.ExtractRenderDeps(ctx))
	}

	user, err := db.QueryUser(ctx.Request().Context(), token.UserID)
	if err != nil {
		return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
			HasError: true,
			Msg:      "An error occurred while trying to reset your password. Please try again.",
		}).Render(views.ExtractRenderDeps(ctx))
	}

	_, err = services.UpdateUser(ctx.Request().Context(), domain.UpdateUser{
		Name:            user.Name,
		Mail:            user.Mail,
		Password:        payload.Password,
		ConfirmPassword: payload.ConfirmPassword,
		ID:              user.ID,
	}, &db, v, cfg.Auth.PasswordPepper)
	if err != nil {
		e, ok := err.(validator.ValidationErrors)
		if !ok {
			slog.ErrorContext(ctx.Request().Context(), "could not infer type ValidationErrors")
		}

		if len(e) == 0 {
			return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
				HasError: true,
				Msg:      "An error occurred while trying to reset your password. Please try again.",
			}).Render(views.ExtractRenderDeps(ctx))
		}

		props := authentication.ResetPasswordFormProps{
			CsrfToken:  csrf.Token(ctx.Request()),
			ResetToken: token.Hash,
		}

		for _, validationError := range e {
			switch validationError.StructField() {
			case "Password", "ConfirmPassword":
				props.Password = validation.InputField{
					Invalid:    true,
					InvalidMsg: validationError.Param(),
				}
				props.ConfirmPassword = validation.InputField{
					Invalid:    true,
					InvalidMsg: validationError.Param(),
				}
			}
		}

		return authentication.ResetPasswordForm(props).Render(views.ExtractRenderDeps(ctx))
	}

	if err := db.DeleteToken(ctx.Request().Context(), token.ID); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not delete token", "error", err)
		return misc.InternalError(ctx)
	}

	return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
		HasError: false,
	}).Render(views.ExtractRenderDeps(ctx))
}

type VerifyEmail struct {
	Token string `query:"token"`
}

func UserEmailVerification(
	ctx echo.Context,
	db database.Queries,
	tknManager tokens.Manager,
	cfg config.Cfg,
	authStore *sessions.CookieStore,
) error {
	var tkn VerifyEmail
	if err := ctx.Bind(&tkn); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return misc.InternalError(ctx)
	}

	hashedToken := tknManager.Hash(tkn.Token)

	token, err := db.QueryTokenByHash(ctx.Request().Context(), hashedToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return authentication.VerifyEmailPage(
				true,
				views.Head{},
			).Render(views.ExtractRenderDeps(ctx))
		}

		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not query token by hash", "error", err)
		return misc.InternalError(ctx)
	}

	if database.ConvertFromPGTimestamptzToTime(
		token.ExpiresAt).Before(time.Now()) &&
		token.Scope != tokens.ScopeEmailVerification {
		return authentication.VerifyEmailPage(true, views.Head{}).
			Render(views.ExtractRenderDeps(ctx))
	}

	confirmTime := time.Now()
	user, err := db.ConfirmUserEmail(ctx.Request().Context(), database.ConfirmUserEmailParams{
		ID:             token.UserID,
		UpdatedAt:      database.ConvertToPGTimestamptz(confirmTime),
		MailVerifiedAt: database.ConvertToPGTimestamptz(confirmTime),
	})
	if err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not confirm user email", "error", err)
		return misc.InternalError(ctx)
	}

	if err := db.DeleteToken(ctx.Request().Context(), token.ID); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not deleted token", "error", err)
		return misc.InternalError(ctx)
	}

	session, err := authStore.Get(ctx.Request(), "ua")
	if err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(
			ctx.Request().Context(),
			"could not get authenticated session",
			"error",
			err,
		)
		return misc.InternalError(ctx)
	}

	authSession := services.CreateAuthenticatedSession(*session, user.ID, cfg)
	if err := authSession.Save(ctx.Request(), ctx.Response()); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not save auth session", "error", err)
		return misc.InternalError(ctx)
	}

	return authentication.VerifyEmailPage(false, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

// VerifySubscriberEmail method    verifies the email the subscriber provided during signup
func SubscriberEmailVerification(
	ctx echo.Context,
	db database.Queries,
	tknManager tokens.Manager,
) error {
	var tkn VerifyEmail
	if err := ctx.Bind(&tkn); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return misc.InternalError(ctx)
	}

	hashedToken := tknManager.Hash(tkn.Token)

	token, err := db.QuerySubscriberTokenByHash(ctx.Request().Context(), hashedToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return authentication.VerifyEmailPage(true, views.Head{}).
				Render(views.ExtractRenderDeps(ctx))
		}

		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not query subscriber token", "error", err)
		return misc.InternalError(ctx)
	}

	if database.ConvertFromPGTimestamptzToTime(token.ExpiresAt).Before(time.Now()) &&
		token.Scope != tokens.ScopeEmailVerification {
		return authentication.VerifyEmailPage(true, views.Head{}).
			Render(views.ExtractRenderDeps(ctx))
	}

	if err := db.ConfirmSubscriberEmail(ctx.Request().Context(), database.ConfirmSubscriberEmailParams{
		ID:        token.SubscriberID,
		UpdatedAt: database.ConvertToPGTimestamptz(time.Now()),
	}); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/login")

		slog.ErrorContext(ctx.Request().Context(), "could not confirm email", "error", err)
		return misc.InternalError(ctx)
	}

	return authentication.VerifySubscriberEmailPage(false, views.Head{}).
		Render(views.ExtractRenderDeps(ctx))
}
