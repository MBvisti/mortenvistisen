package paths

import (
	"context"

	"github.com/a-h/templ"
)

type Route string

const (
	HomePage                 Route = "homePage"
	AboutPage                Route = "aboutPage"
	LoginPage                Route = "loginPage"
	Login                    Route = "login"
	ForgotPasswordPage       Route = "forgotPasswordPage"
	ForgotPassword           Route = "forgotPassword"
	ResetPasswordPage        Route = "resetPasswordPage"
	ResetPassword            Route = "resetPassword"
	RegisterPage             Route = "registerPage"
	RegisterUser             Route = "registerUser"
	VerifyEmailPage          Route = "verifyEmailPage"
	ArticlePage              Route = "articlePage"
	ArticlesPage             Route = "articlesPage"
	ProjectsPage             Route = "projectsPage"
	NewslettersPage          Route = "newslettersPage"
	NewsletterPage           Route = "newsletterPage"
	SubscribeEvent           Route = "subscribeEvent"
	VerifySubEvent           Route = "verifySubEvent"
	UnsubscribeEvent         Route = "unsubscribeEvent"
	DashboardHomePage        Route = "dashbordHomePage"
	DashboardSubscriberPage  Route = "dashbordSubscriberPage"
	DashboardNewsletter      Route = "dashboardNewsletter"
	DashboardNewsletterNew   Route = "dashboardNewsletterNew"
	DashboardNewsletterStore Route = "dashboardNewsletterStore"
)

func (r Route) ToString() string {
	return string(r)
}

func Get(ctx context.Context, route Route) string {
	return ctx.Value(route).(string)
}

func GetSafeURL(ctx context.Context, route Route) templ.SafeURL {
	return templ.SafeURL(ctx.Value(route).(string))
}
