package paths

import (
	"context"

	"github.com/MBvisti/mortenvistisen/views/contexts"
	"github.com/a-h/templ"
)

type Route string

const (
	HomePage               Route = "homePage"
	AboutPage              Route = "aboutPage"
	LoginPage              Route = "loginPage"
	Login                  Route = "login"
	ForgotPasswordPage     Route = "forgotPasswordPage"
	ForgotPassword         Route = "forgotPassword"
	ResetPasswordPage      Route = "resetPasswordPage"
	ResetPassword          Route = "resetPassword"
	RegisterPage           Route = "registerPage"
	RegisterUser           Route = "registerUser"
	VerifyEmailPage        Route = "verifyEmailPage"
	ArticlePage            Route = "articlePage"
	ArticlesPage           Route = "articlesPage"
	ProjectsPage           Route = "projectsPage"
	NewslettersPage        Route = "newslettersPage"
	SubscribeEvent         Route = "subscribeEvent"
	VerifySubEvent         Route = "verifySubEvent"
	UnsubscribeEvent       Route = "unsubscribeEvent"
	DashboardHomePage      Route = "dashbordHomePage"
	DashboardNewsletter    Route = "dashboardNewsletter"
	DashboardNewsletterNew Route = "dashboardNewsletterNew"
)

func Get(ctx context.Context, route Route) string {
	return contexts.ExtractApp(ctx).Routes[string(route)]
}

func GetSafeURL(ctx context.Context, route Route) templ.SafeURL {
	return templ.SafeURL(contexts.ExtractApp(ctx).Routes[string(route)])
}
