package handlers

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Authentication struct {
	base         Base
	authService  services.Auth
	userModel    models.UserService
	emailService services.Email
	tknService   services.Token
	cfg          config.Cfg
}

func NewAuthentication(
	base Base,
	authService services.Auth,
	userModel models.UserService,
	emailService services.Email,
	tknService services.Token,
	cfg config.Cfg,
) Authentication {
	return Authentication{base, authService, userModel, emailService, tknService, cfg}
}

func (a Authentication) CreateSession(ctx echo.Context) error {
	return authentication.LoginPage(authentication.LoginPageProps{
		CsrfToken: csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

func (a Authentication) StoreAuthenticatedSession(ctx echo.Context) error {
	type userLoginPayload struct {
		Mail       string `form:"email"`
		Password   string `form:"password"`
		RememberMe string `form:"remember_me"` // TODO impl
	}

	var payload userLoginPayload
	if err := ctx.Bind(&payload); err != nil {
		slog.ErrorContext(ctx.Request().Context(), "could not parse UserLoginPayload", "error", err)

		return authentication.LoginResponse(true).Render(views.ExtractRenderDeps(ctx))
	}

	if err := a.authService.AuthenticateUser(
		ctx.Request().Context(),
		ctx.Request(),
		ctx.Response(),
		false,
		payload.Mail,
		payload.Password,
	); err != nil {
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

func (a Authentication) CreatePasswordReset(ctx echo.Context) error {
	return authentication.ForgottenPasswordPage(authentication.ForgottenPasswordPageProps{
		CsrfToken: csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

func (a Authentication) StorePasswordReset(ctx echo.Context) error {
	type storePasswordResetPayload struct {
		Mail string `form:"email"`
	}

	var payload storePasswordResetPayload
	if err := ctx.Bind(&payload); err != nil {
		return authentication.ForgottenPasswordSuccess(true).Render(views.ExtractRenderDeps(ctx))
	}

	user, err := a.userModel.ByEmail(ctx.Request().Context(), payload.Mail)
	if err != nil {
		failureOccurred := true
		if errors.Is(err, pgx.ErrNoRows) {
			failureOccurred = false
		}

		return authentication.ForgottenPasswordSuccess(failureOccurred).
			Render(views.ExtractRenderDeps(ctx))
	}

	resetPWToken, err := a.tknService.CreateResetPasswordToken(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	if err := a.emailService.SendPasswordReset(ctx.Request().Context(), user.Mail, fmt.Sprintf(
		"%s://%s/reset-password?token=%s",
		a.cfg.App.AppScheme,
		a.cfg.App.AppHost,
		resetPWToken)); err != nil {
		return err
	}

	return authentication.ForgottenPasswordSuccess(false).Render(views.ExtractRenderDeps(ctx))
}

func (a Authentication) CreateResetPassword(ctx echo.Context) error {
	type passwordResetToken struct {
		Token string `query:"token"`
	}

	var payload passwordResetToken
	if err := ctx.Bind(&payload); err != nil {
		return a.base.InternalError(ctx)
	}

	return authentication.ResetPasswordPage(authentication.ResetPasswordPageProps{
		ResetToken: payload.Token,
		CsrfToken:  csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

func (a Authentication) StoreResetPassword(c echo.Context) error {
	ctx, span := a.base.Tracer.Start(
		c.Request().Context(),
		"AuthenticationHandler/StoreResetPassword",
	)
	span.AddEvent("StoreResetPassword/start")
	type resetPasswordPayload struct {
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirm_password"`
		Token           string `form:"token"`
	}

	var payload resetPasswordPayload
	if err := c.Bind(&payload); err != nil {
		return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
			HasError: true,
			Msg:      "An error occurred while trying to reset your password. Please try again.",
		}).Render(views.ExtractRenderDeps(c))
	}

	if err := a.tknService.Validate(ctx, payload.Token, services.ScopeResetPassword); err != nil {
		return err
	}

	userID, err := a.tknService.GetAssociatedUserID(ctx, payload.Token)
	if err != nil {
		return err
	}

	user, err := a.userModel.UpdatePassword(
		ctx,
		userID,
		payload.Password,
		payload.ConfirmPassword,
	)
	if err != nil {
		return err
	}

	if err = a.authService.CreateAuthenticatedSession(c.Request(), c.Response(), user.ID, false); err != nil {
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

	if err := a.tknService.Delete(ctx, span, payload.Token); err != nil {
		return err
	}

	return authentication.ResetPasswordResponse(authentication.ResetPasswordResponseProps{
		HasError: false,
	}).Render(views.ExtractRenderDeps(c))
}
