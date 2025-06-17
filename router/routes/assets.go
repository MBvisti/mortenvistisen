package routes

import (
	"net/http"
)

const (
	AssetsRoutePrefix = "/assets"
	assetsNamePrefix  = "assets"
)

var Assets = []Route{
	Robots,
	Sitemap,
	LLM,
	CssEntrypoint,
	AllCss,
	JsEntrypoint,
	JsEasyMDE,
	JsNav,
	AllJs,
	Favicon,
	Favicon16,
	Favicon32,
	FaviconAppletouch,
	FaviconSiteManifest,
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

var CssEntrypoint = Route{
	Name:        assetsNamePrefix + "css.entry",
	Path:        assetsNamePrefix + "/css/styles.css",
	Method:      http.MethodGet,
	HandlerName: "Styles",
}

var AllCss = Route{
	Name:        assetsNamePrefix + "css.all",
	Path:        assetsNamePrefix + "/css/:file",
	Method:      http.MethodGet,
	HandlerName: "AllCss",
}

var JsEntrypoint = Route{
	Name:        assetsNamePrefix + "js.entry",
	Path:        assetsNamePrefix + "/js/script.js",
	Method:      http.MethodGet,
	HandlerName: "Scripts",
}

var JsEasyMDE = Route{
	Name:        assetsNamePrefix + "js.easymde",
	Path:        assetsNamePrefix + "/js/easymde.js",
	Method:      http.MethodGet,
	HandlerName: "IndividualScript",
}

var JsNav = Route{
	Name:        assetsNamePrefix + "js.nav",
	Path:        assetsNamePrefix + "/js/nav.js",
	Method:      http.MethodGet,
	HandlerName: "IndividualScript",
}

var AllJs = Route{
	Name:        assetsNamePrefix + "js.all",
	Path:        assetsNamePrefix + "/js/:file",
	Method:      http.MethodGet,
	HandlerName: "AllJs",
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
