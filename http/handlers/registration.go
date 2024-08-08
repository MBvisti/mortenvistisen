package handlers

import (
	"log/slog"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

type Registration struct {
	base         Base
	authService  services.Auth
	userModel    models.UserService
	tokenService services.Token
	emailService services.Email
}

func NewRegistration(
	base Base,
	authService services.Auth,
	userModel models.UserService,
	tokenService services.Token,
	emailService services.Email,
) Registration {
	return Registration{base, authService, userModel, tokenService, emailService}
}

func (r Registration) CreateUser(ctx echo.Context) error {
	return authentication.RegisterPage(authentication.RegisterFormProps{
		CsrfToken: csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

func (r Registration) StoreUser(ctx echo.Context) error {
	type storeUserPayload struct {
		UserName        string `form:"user_name"`
		Mail            string `form:"email"`
		Password        string `form:"password"`
		ConfirmPassword string `form:"confirm_password"`
	}

	var payload storeUserPayload
	if err := ctx.Bind(&payload); err != nil {
		return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
			Render(views.ExtractRenderDeps(ctx))
	}

	user, err := r.userModel.New(
		ctx.Request().Context(),
		payload.UserName,
		payload.Mail,
		payload.Password,
		payload.ConfirmPassword,
	)
	if err != nil {
		return err
	}

	activationToken, err := r.tokenService.CreateUserEmailVerification(
		ctx.Request().Context(),
		user.ID,
	)
	if err != nil {
		return err
	}

	if err := r.emailService.SendUserSignup(ctx.Request().Context(), user.Mail, activationToken); err != nil {
		return err
	}

	return authentication.RegisterResponse("You're now registered", "You should receive an email soon to validate your account.", false).
		Render(views.ExtractRenderDeps(ctx))
}

func (r Registration) UserEmailVerification(c echo.Context) error {
	ctx, span := r.base.Tracer.CreateSpan(
		c.Request().Context(),
		"RegistrationHandler/UserEmailVerification",
	)
	span.AddEvent("UserEmailVerification/start")
	type verifyEmail struct {
		Token string `query:"token"`
	}

	var payload verifyEmail
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(ctx, "could not bind verify email", "error", err)
		c.Response().Writer.Header().Add("HX-Redirect", "/500")
		c.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return r.base.InternalError(c)
	}

	if err := r.tokenService.Validate(ctx, payload.Token, services.ScopeEmailVerification); err != nil {
		return err
	}

	userID, err := r.tokenService.GetAssociatedUserID(ctx, payload.Token)
	if err != nil {
		return err
	}

	// TODO: add span
	user, err := r.userModel.ConfirmEmail(ctx, userID)
	if err != nil {
		return err
	}

	if err := r.authService.CreateAuthenticatedSession(c.Request(), c.Response(), user.ID, false); err != nil {
		return err
	}

	tokenDeleteCtx, tokenDeleteSpan := r.base.Tracer.CreateChildSpan(
		ctx,
		span,
		"TokenService/Delete",
	)
	tokenDeleteSpan.AddEvent("Delete/Start")
	if err := r.tokenService.Delete(tokenDeleteCtx, payload.Token); err != nil {
		return err
	}
	tokenDeleteSpan.End()
	if err := r.tokenService.Delete(tokenDeleteCtx, payload.Token); err != nil {
		return err
	}

	span.End()

	return authentication.VerifyEmailPage(false, views.Head{}).Render(views.ExtractRenderDeps(c))
}
