package paths

import (
	"context"
	"log/slog"
	"strings"

	"github.com/a-h/templ"
)

type (
	Name        string
	Params      map[string]string
	QueryParams map[string]string
)

func (n Name) String() string {
	return string(n)
}

type paths []Name

var Paths = []Name{
	APIHealth,
	APICollect,

	Home,
	Articles,
	Article,
	About,
	Projects,
	Newsletters,
	Newsletter,
	CreateSubscription,
	UnSubscribe,
	VerifySubscriber,

	Login,
	StoreAuthenticatedSession,
	ForgotPassword,
	StoreForgotPassword,
	ResetPassword,
	StoreResetPassword,

	NewUser,
	CreateUser,
	VerifyEmail,

	Dashboard,
	DashboardSubscribers,
	DashboardShowSubscriber,
	DashboardUpdateSubscriber,
	DashboardDeleteSubscriber,
	DashboardNewsletters,
	DashboardNewNewsletter,
	DashboardShowNewsletter,
	DashboardDeleteNewsletter,
	DashboardCreateNewsletter,

	Robots,
	CssTrix,
	CssBootstrap,
	CssBootstrapOverrides,
	JsThemeSwitcher,
	JsAlpine,
	JsAnalytics,
	JsHtmx,
	JsTrix,
	JsPopper,
	JsBootstrap,
	JsScript,
	Sitemap,
}

func GP(
	ctx context.Context,
	name Name,
	params Params,
	query QueryParams,
) string {
	p, ok := ctx.Value(name).(string)
	if !ok {
		return ""
	}

	if len(params) == 0 || !strings.Contains(p, ":") {
		return p
	}

	for key, value := range params {
		p = strings.Replace(p, ":"+key, value, 1)
	}

	if len(query) != 0 {
		var queryParams string

		for key, value := range query {
			q := key + "=" + value
			if queryParams == "" {
				queryParams = q
			}
			if queryParams != "" {
				queryParams = queryParams + "&" + q
			}
		}

		p = p + "?" + queryParams
	}

	return p
}

func GSP(
	ctx context.Context,
	name Name,
	params Params,
	query QueryParams,
) templ.SafeURL {
	slog.Info("CTX", "ctx", ctx)
	p, ok := ctx.Value(name).(string)
	if !ok {
		return ""
	}

	if len(params) == 0 || !strings.Contains(p, ":") {
		return templ.SafeURL(p)
	}

	for key, value := range params {
		p = strings.Replace(p, ":"+key, value, 1)
	}

	if len(query) != 0 {
		var queryParams string
		for key, value := range query {
			q := key + "=" + value
			if queryParams == "" {
				queryParams = q
			}
			if queryParams != "" {
				queryParams = queryParams + "&" + q
			}
		}

		p = p + "?" + queryParams
	}

	return templ.SafeURL(p)
}
