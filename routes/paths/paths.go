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

func GP(
	ctx context.Context,
	name Name,
	params Params,
	query QueryParams,
) string {
	// appCtx, ok := ctx.Value(contexts.AppKey{}).(*contexts.App)
	// if !ok {
	// 	return ""
	// }
	//
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

	// p := appCtx.Routes[string(name)]
	//
	// slog.Info("##########################", "p", p)

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

// type Path struct {
// 	name    string
// 	pattern string
// }
//
// func (p Path) Name() string {
// 	return p.name
// }
//
// // Raw should only be used to setup the path in router
// func (p Path) Raw() string {
// 	return p.pattern
// }
//
// type Param map[string]string
//
// // WithParams figure out a better name for this
// func (p Path) WithParams(params Param) string {
// 	if len(params) == 0 || !strings.Contains(p.pattern, ":") {
// 		return p.pattern
// 	}
//
// 	path := p.pattern
// 	for key, value := range params {
// 		path = strings.Replace(path, ":"+key, value, 1)
// 	}
//
// 	return path
// }
