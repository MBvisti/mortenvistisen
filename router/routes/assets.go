// Package routes contains all the routes used throughout the project
package routes

import (
	"fmt"
	"time"

	"mortenvistisen/internal/routing"
)

const (
	AssetsPrefix     = "/assets"
	IndexNowKeyValue = "4zd8j69sf3ju2hnfxmebr3czub8uu63m"
)

var startTime = time.Now().Unix()

var Robots = routing.NewSimpleRoute(
	"/robots.txt",
	"assets.robots",
	"",
)

var Sitemap = routing.NewSimpleRoute(
	"/sitemap.xml",
	"assets.sitemap",
	"",
)

var IndexNowKey = routing.NewSimpleRoute(
	"/"+IndexNowKeyValue+".txt",
	"assets.index_now_key",
	"",
)

var Stylesheet = routing.NewSimpleRoute(
	fmt.Sprintf("/css/%v/style.css", startTime),
	"css.stylesheet",
	AssetsPrefix,
)

var Font = routing.NewRouteWithFile(
	fmt.Sprintf("/css/%v/files/:file", startTime),
	"css.font",
	AssetsPrefix,
)

var Scripts = routing.NewSimpleRoute(
	fmt.Sprintf("/js/%v/scripts.js", startTime),
	"js.scripts",
	AssetsPrefix,
)

var Script = routing.NewRouteWithFile(
	fmt.Sprintf("/js/%v/:file", startTime),
	"js.script",
	AssetsPrefix,
)
var ViteBuild = routing.NewSimpleRoute(
	fmt.Sprintf("/dist/%v/*", startTime),
	"vite.build",
	AssetsPrefix,
)
