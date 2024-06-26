package authentication

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/controllers/misc"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func CreateAuthenticatedSession(ctx echo.Context) error {
	return authentication.LoginPage(authentication.LoginPageProps{
		CsrfToken: csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

type UserLoginPayload struct {
	Mail       string `form:"email"`
	Password   string `form:"password"`
	RememberMe string `form:"remember_me"` // TODO impl
}

func StoreAuthenticatedSession(
	ctx echo.Context,
	authService services.Auth,
) error {
	var payload UserLoginPayload
	if err := ctx.Bind(&payload); err != nil {
		slog.ErrorContext(ctx.Request().Context(), "could not parse UserLoginPayload", "error", err)

		return authentication.LoginResponse(true).Render(views.ExtractRenderDeps(ctx))
	}

	_, err := authService.AuthenticateUser(
		ctx.Request().Context(),
		ctx.Request(),
		false,
		payload.Mail,
		payload.Password,
	)
	if err != nil {
		var errMsg string

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
	userModel models.UserService,
	emailService services.Email,
	tknService services.TokenSvc,
	cfg config.Cfg,
) error {
	var payload StorePasswordResetPayload
	if err := ctx.Bind(&payload); err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}

	user, err := userModel.ByEmail(ctx.Request().Context(), payload.Mail)
	if err != nil {
		failureOccurred := true
		if errors.Is(err, pgx.ErrNoRows) {
			failureOccurred = false
		}

		return authentication.ForgottenPasswordSuccess(failureOccurred).
			Render(views.ExtractRenderDeps(ctx))
	}

	resetPWToken, err := tknService.CreateResetPasswordToken(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	if err := emailService.SendPasswordReset(ctx.Request().Context(), user.Mail, fmt.Sprintf(
		"%s://%s/reset-password?token=%s",
		cfg.App.AppScheme,
		cfg.App.AppHost,
		resetPWToken)); err != nil {
		return err
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
	userModel models.UserService,
	authService services.Auth,
	tknService services.TokenSvc,
) error {
	var payload ResetPasswordPayload
	if err := ctx.Bind(&payload); err != nil {
		return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
			HasError: true,
			Msg:      "An error occurred while trying to reset your password. Please try again.",
		}).Render(views.ExtractRenderDeps(ctx))
	}

	if err := tknService.Validate(ctx.Request().Context(), payload.Token); err != nil {
		return err
	}

	userID, err := tknService.GetAssociatedUserID(ctx.Request().Context(), payload.Token)

	user, err := userModel.UpdatePassword(
		ctx.Request().Context(),
		userID,
		payload.Password,
		payload.ConfirmPassword,
	)
	if err != nil {
		return err
	}

	_, err = authService.CreateAuthenticatedSession(ctx.Request(), user.ID, false)
	if err != nil {
		return err
	}

	// 	return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
	// 		HasError: true,
	// 		Msg:      "An error occurred while trying to reset your password. Please try again.",
	// 	}).Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// if database.ConvertFromPGTimestamptzToTime(token.ExpiresAt).Before(time.Now()) &&
	// 	token.Scope != tokens.ScopeResetPassword {
	// 	return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
	// 		HasError: true,
	// 		Msg:      "The token has expired. Please request a new one.",
	// 	}).Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// user, err := db.QueryUser(ctx.Request().Context(), token.UserID)
	// if err != nil {
	// 	return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
	// 		HasError: true,
	// 		Msg:      "An error occurred while trying to reset your password. Please try again.",
	// 	}).Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// _, err = services.UpdateUser(ctx.Request().Context(), domain.UpdateUser{
	// 	Name:            user.Name,
	// 	Mail:            user.Mail,
	// 	Password:        payload.Password,
	// 	ConfirmPassword: payload.ConfirmPassword,
	// 	ID:              user.ID,
	// }, &db, v, cfg.Auth.PasswordPepper)
	// if err != nil {
	// 	e, ok := err.(validator.ValidationErrors)
	// 	if !ok {
	// 		slog.ErrorContext(ctx.Request().Context(), "could not infer type ValidationErrors")
	// 	}
	//
	// 	if len(e) == 0 {
	// 		return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
	// 			HasError: true,
	// 			Msg:      "An error occurred while trying to reset your password. Please try again.",
	// 		}).Render(views.ExtractRenderDeps(ctx))
	// 	}
	//
	// 	props := authentication.ResetPasswordFormProps{
	// 		CsrfToken:  csrf.Token(ctx.Request()),
	// 		ResetToken: token.Hash,
	// 	}
	//
	// 	for _, validationError := range e {
	// 		switch validationError.StructField() {
	// 		case "Password", "ConfirmPassword":
	// 			props.Password = validation.InputField{
	// 				Invalid:    true,
	// 				InvalidMsg: validationError.Param(),
	// 			}
	// 			props.ConfirmPassword = validation.InputField{
	// 				Invalid:    true,
	// 				InvalidMsg: validationError.Param(),
	// 			}
	// 		}
	// 	}
	//
	// 	return authentication.ResetPasswordForm(props).Render(views.ExtractRenderDeps(ctx))
	// }

	// if err := db.DeleteToken(ctx.Request().Context(), token.ID); err != nil {
	// 	ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
	// 	ctx.Response().Writer.Header().Add("PreviousLocation", "/login")
	//
	// 	slog.ErrorContext(ctx.Request().Context(), "could not delete token", "error", err)
	// 	return misc.InternalError(ctx)
	// }

	if err := tknService.Delete(ctx.Request().Context(), payload.Token); err != nil {
		return err
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
	userModel models.UserService,
	tknService services.TokenSvc,
	authService services.Auth,
) error {
	var tkn VerifyEmail
	if err := ctx.Bind(&tkn); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return misc.InternalError(ctx)
	}

	if err := tknService.Validate(ctx.Request().Context(), tkn.Token); err != nil {
		return err
	}

	userID, err := tknService.GetAssociatedUserID(ctx.Request().Context(), tkn.Token)
	if err != nil {
		return err
	}

	user, err := userModel.ConfirmEmail(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}

	if _, err := authService.CreateAuthenticatedSession(ctx.Request(), user.ID, false); err != nil {
		return err
	}

	if err := tknService.Delete(ctx.Request().Context(), tkn.Token); err != nil {
		return err
	}

	return authentication.VerifyEmailPage(false, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}
