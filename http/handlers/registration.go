package handlers

import (
	"errors"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

type Registration struct {
	authSvc services.Auth
	db      psql.Postgres
	email   services.Mail
}

func newRegistration(
	authSvc services.Auth,
	db psql.Postgres,
	email services.Mail,
) Registration {
	return Registration{authSvc, db, email}
}

func (r *Registration) NewUser(c echo.Context) error {
	return authentication.RegisterPage(authentication.RegisterFormProps{
		CsrfToken: csrf.Token(c.Request()),
	}).Render(renderArgs(c))
}

type StoreUserPayload struct {
	UserName        string `form:"username"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
}

func (r *Registration) CreateUser(c echo.Context) error {
	var payload StoreUserPayload
	if err := c.Bind(&payload); err != nil {
		return views.ErrorPage().Render(renderArgs(c))
	}

	err := r.authSvc.RegisterUser(
		c.Request().
			Context(),
		payload.UserName,
		payload.Email,
		payload.Password,
		payload.ConfirmPassword,
	)
	if err != nil {
		if errors.Is(err, services.ErrUnrecoverable) {
			return views.ErrorPage().Render(renderArgs(c))
		}

		if errors.Is(err, models.ErrDomainValidation) {
			var validationErrors validator.ValidationErrors
			if ok := errors.As(err, &validationErrors); !ok {
				return views.ErrorPage().Render(renderArgs(c))
			}

			fields := make(
				map[string]components.InputFieldProps,
				len(validationErrors),
			)
			for _, validationError := range validationErrors {
				fields[validationError.StructField()] = components.InputFieldProps{
					Value:     validationError.Value().(string),
					ErrorMsgs: []string{validationError.Error()},
				}
			}

			props := authentication.RegisterFormProps{
				SuccessRegister: false,
				Fields:          fields,
				CsrfToken:       csrf.Token(c.Request()),
			}
			return authentication.RegisterForm(props).
				Render(renderArgs(c))
		}
	}

	props := authentication.RegisterFormProps{
		SuccessRegister: true,
		CsrfToken:       csrf.Token(c.Request()),
	}
	return authentication.RegisterForm(props).
		Render(renderArgs(c))
}

type verificationTokenPayload struct {
	Token string `query:"token"`
}

func (r *Registration) VerifyUserEmail(c echo.Context) error {
	var payload verificationTokenPayload
	if err := c.Bind(&payload); err != nil {
		return views.ErrorPage().Render(renderArgs(c))
	}

	// if err := r.tknService.Validate(c.Request().Context(), payload.Token, services.ScopeEmailVerification); err != nil {
	// 	if err := c.Bind(&payload); err != nil {
	// 		c.Response().Writer.Header().Add("HX-Redirect", "/500")
	// 		c.Response().
	// 			Writer.Header().
	// 			Add("PreviousLocation", "/user/create")
	//
	// 		return views.ErrorPage().Render(renderArgs(c))
	// 	}
	// }

	// userID, err := r.tknService.GetAssociatedUserID(
	// 	c.Request().Context(),
	// 	payload.Token,
	// )
	// if err != nil {
	// 	if err := c.Bind(&payload); err != nil {
	// 		c.Response().Writer.Header().Add("HX-Redirect", "/500")
	// 		c.Response().
	// 			Writer.Header().
	// 			Add("PreviousLocation", "/user/create")
	//
	// 		return r.InternalError(c)
	// 	}
	// }
	//
	// user, err := r.db.QueryUserByID(c.Request().Context(), userID)
	// if err != nil {
	// 	return r.InternalError(c)
	// }

	// if err := r.userModel.VerifyEmail(c.Request().Context(), user.Email); err != nil {
	// 	return r.InternalError(c)
	// }
	//
	// _, err = r.authService.NewUserSession(
	// 	c.Request(),
	// 	c.Response(),
	// 	user.ID,
	// )
	// if err != nil {
	// 	return r.InternalError(c)
	// }

	return authentication.VerifyEmailPage(false).
		Render(renderArgs(c))
}
