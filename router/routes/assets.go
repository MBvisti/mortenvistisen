package routes

import (
	"net/http"
	"strings"
)

const (
	AssetsRoutePrefix = "/assets"
	assetsNamePrefix  = "assets"
)

var Assets = []Route{
	Robots,
	Sitemap,
	LLM,
	CSSEntrypoint,
	CSSFile,
	JsEntrypoint,
	JavascriptFile.Route,
	Favicon,
	Favicon16,
	Favicon32,
	FaviconAppletouch,
	FaviconSiteManifest,
	IndexNow,
}

var Robots = Route{
	Name:        ".robots",
	Path:        "/robots.txt",
	Method:      http.MethodGet,
	HandlerName: "Robots",
}

var Sitemap = Route{
	Name:        ".sitemap",
	Path:        "/sitemap.xml",
	Method:      http.MethodGet,
	HandlerName: "Sitemap",
}

var LLM = Route{
	Name:        ".llm",
	Path:        "/llm.txt",
	Method:      http.MethodGet,
	HandlerName: "LLM",
}

var IndexNow = Route{
	Name:        ".index_now",
	Path:        "/4zd8j69sf3ju2hnfxmebr3czub8uu63m.txt",
	Method:      http.MethodGet,
	HandlerName: "IndexNow",
}

var CSSEntrypoint = Route{
	Name:        assetsNamePrefix + "css.entry",
	Path:        assetsNamePrefix + "/css/:version/styles.css",
	Method:      http.MethodGet,
	HandlerName: "Styles",
}

var CSSFile = Route{
	Name:        assetsNamePrefix + "css.file",
	Path:        assetsNamePrefix + "/css/:version/:file",
	Method:      http.MethodGet,
	HandlerName: "IndividualCSSFile",
}

var JsEntrypoint = Route{
	Name:        assetsNamePrefix + "js.entry",
	Path:        assetsNamePrefix + "/js/:version/script.js",
	Method:      http.MethodGet,
	HandlerName: "Scripts",
}

var JsDashboardEntrypoint = Route{
	Name:        assetsNamePrefix + "js.dashboard_entry",
	Path:        assetsNamePrefix + "/js/:version/dashboard_script.js",
	Method:      http.MethodGet,
	HandlerName: "ScriptsDashboard",
}

var JavascriptFile = javascriptFile{
	Route: Route{
		Name:        assetsNamePrefix + "js.file",
		Path:        assetsNamePrefix + "/js/:version/:file",
		Method:      http.MethodGet,
		HandlerName: "IndividualScript",
	},
}

type javascriptFile struct {
	Route
}

func (j javascriptFile) GetPath(file string) string {
	return strings.Replace(j.Path, ":file", file, 1)
}

var Favicon = Route{
	Name:        assetsNamePrefix + ".favicon",
	Path:        "/favicon.ico",
	Method:      http.MethodGet,
	HandlerName: "Favicon",
}

var Favicon16 = Route{
	Name:        assetsNamePrefix + ".favicon_16",
	Path:        assetsNamePrefix + "/favicon-16x16.png",
	Method:      http.MethodGet,
	HandlerName: "Favicon16",
}

var Favicon32 = Route{
	Name:        assetsNamePrefix + ".favicon_32",
	Path:        assetsNamePrefix + "/favicon-32x32.png",
	Method:      http.MethodGet,
	HandlerName: "Favicon32",
}

var FaviconAppletouch = Route{
	Name:        assetsNamePrefix + ".favicon_appletouch",
	Path:        assetsNamePrefix + "/apple-touch-icon.png",
	Method:      http.MethodGet,
	HandlerName: "FaviconAppleTouch",
}

var FaviconSiteManifest = Route{
	Name:        assetsNamePrefix + ".favicon_site_manifest",
	Path:        assetsNamePrefix + "/site.webmanifest",
	Method:      http.MethodGet,
	HandlerName: "FaviconSiteManifest",
}
