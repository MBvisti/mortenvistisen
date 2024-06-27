package authentication

import (
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

// CreateUser method    shows the form to create the user
func CreateUser(ctx echo.Context) error {
	return authentication.RegisterPage(authentication.RegisterFormProps{
		CsrfToken: csrf.Token(ctx.Request()),
	}, views.Head{}).Render(views.ExtractRenderDeps(ctx))
}

type StoreUserPayload struct {
	UserName        string `form:"user_name"`
	Mail            string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
}

// StoreUser method    stores the new user
func StoreUser(
	ctx echo.Context,
	userModel models.UserService,
	tokenService services.TokenSvc,
	emailService services.Email,
) error {
	var payload StoreUserPayload
	if err := ctx.Bind(&payload); err != nil {
		return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
			Render(views.ExtractRenderDeps(ctx))
	}

	user, err := userModel.New(
		ctx.Request().Context(),
		payload.UserName,
		payload.Mail,
		payload.Password,
		payload.ConfirmPassword,
	)
	if err != nil {
		return err
	}

	activationToken, err := tokenService.CreateEmailVerificationToken(
		ctx.Request().Context(),
		user.ID,
	)
	if err != nil {
		return err
	}

	if err := emailService.SendUserSignup(ctx.Request().Context(), user.Mail, activationToken); err != nil {
		return err
	}

	// user, err := services.NewUser(ctx.Request().Context(), domain.NewUser{
	// 	Name:            payload.UserName,
	// 	Mail:            payload.Mail,
	// 	Password:        payload.Password,
	// 	ConfirmPassword: payload.ConfirmPassword,
	// }, &db, v, cfg.Auth.PasswordPepper)
	// if err != nil {
	// 	telemetry.Logger.Info("error", "err", err)
	// 	e, ok := err.(validator.ValidationErrors)
	// 	if !ok {
	// 		telemetry.Logger.WarnContext(
	// 			ctx.Request().Context(),
	// 			"an unrecoverable error occurred",
	// 			"error",
	// 			err,
	// 		)
	//
	// 		return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
	// 			Render(views.ExtractRenderDeps(ctx))
	// 	}
	//
	// 	if len(e) == 0 {
	// 		telemetry.Logger.WarnContext(
	// 			ctx.Request().Context(),
	// 			"an unrecoverable error occurred",
	// 			"error",
	// 			err,
	// 		)
	//
	// 		return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
	// 			Render(views.ExtractRenderDeps(ctx))
	// 	}
	//
	// 	props := authentication.RegisterFormProps{
	// 		NameInput: validation.InputField{
	// 			OldValue: payload.UserName,
	// 		},
	// 		EmailInput: validation.InputField{
	// 			OldValue: payload.Mail,
	// 		},
	// 		CsrfToken: csrf.Token(ctx.Request()),
	// 	}
	//
	// 	for _, validationError := range e {
	// 		switch validationError.StructField() {
	// 		case "Name":
	// 			props.NameInput.Invalid = true
	// 			props.NameInput.InvalidMsg = validationError.Param()
	// 		case "MailRegistered":
	// 			props.EmailInput.Invalid = true
	// 			props.EmailInput.InvalidMsg = validationError.Param()
	// 		case "Password", "ConfirmPassword":
	// 			props.PasswordInput = validation.InputField{
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
	// 	return authentication.RegisterForm(props).Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// generatedTkn, err := tknManager.GenerateToken()
	// if err != nil {
	// 	telemetry.Logger.ErrorContext(ctx.Request().Context(), "could not query user", "error", err)
	//
	// 	return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// activationToken := tokens.CreateActivationToken(
	// 	generatedTkn.PlainTextToken,
	// 	generatedTkn.HashedToken,
	// )
	//
	// if err := db.StoreToken(ctx.Request().Context(), database.StoreTokenParams{
	// 	ID:        uuid.New(),
	// 	CreatedAt: database.ConvertToPGTimestamptz(time.Now()),
	// 	Hash:      activationToken.Hash,
	// 	ExpiresAt: database.ConvertToPGTimestamptz(activationToken.GetExpirationTime()),
	// 	Scope:     activationToken.GetScope(),
	// 	UserID:    user.ID,
	// }); err != nil {
	// 	telemetry.Logger.ErrorContext(ctx.Request().Context(), "could not query user", "error", err)
	//
	// 	return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// userSignupMail := templates.UserSignupWelcomeMail{
	// 	ConfirmationLink: fmt.Sprintf(
	// 		"%s://%s/verify-email?token=%s",
	// 		cfg.App.AppScheme,
	// 		cfg.App.AppHost,
	// 		activationToken.GetPlainText(),
	// 	),
	// }
	// textVersion, err := userSignupMail.GenerateTextVersion()
	// if err != nil {
	// 	telemetry.Logger.ErrorContext(ctx.Request().Context(), "could not query user", "error", err)
	//
	// 	return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }
	// htmlVersion, err := userSignupMail.GenerateHtmlVersion()
	// if err != nil {
	// 	telemetry.Logger.ErrorContext(ctx.Request().Context(), "could not query user", "error", err)
	//
	// 	return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// _, err = queueClient.Insert(ctx.Request().Context(), queue.EmailJobArgs{
	// 	To:          user.Mail,
	// 	From:        cfg.App.DefaultSenderSignature,
	// 	Subject:     "Thanks for signing up!",
	// 	TextVersion: textVersion,
	// 	HtmlVersion: htmlVersion,
	// }, nil)
	// if err != nil {
	// 	telemetry.Logger.ErrorContext(ctx.Request().Context(), "could not query user", "error", err)
	//
	// 	return authentication.RegisterResponse("An error occurred", "Please refresh the page an try again.", true).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }

	return authentication.RegisterResponse("You're now registered", "You should receive an email soon to validate your account.", false).
		Render(views.ExtractRenderDeps(ctx))
}
