package main

import (
	"fmt"
	"github.com/goforgery/forgery2"
	"github.com/goforgery/mustache"
	"github.com/goforgery/static"
	"strconv"
)

// Starts the web server.
func StartWebServer(port int, zilla *Zilla) {

	app := f.CreateApp()
	app.Use(static.Create())
	app.Engine(".html", mustache.Create())

	app.Get("/", func(req *f.Request, res *f.Response, next func()) {
		res.Render("index.html", zilla.GetLiveData())
	})

	app.Get("/livedata", func(req *f.Request, res *f.Response, next func()) {
		res.Send(zilla.GetLiveData())
	})

	app.Get("/set/:attribute", func(req *f.Request, res *f.Response, next func()) {
		var status bool
		attribute := req.Param("attribute")
		value, _ := strconv.Atoi(req.Query("value"))
		switch attribute {
		case "BatteryAmpLimit":
			status = zilla.SetBatteryAmpLimit(value)
		}
		res.Send(map[string]interface{}{"status": status, "attribute": attribute, "value": value})
	})

	app.Get("/settings", func(req *f.Request, res *f.Response, next func()) {
		zilla.Refresh()
		res.Render("settings.html", zilla)
	})

	fmt.Printf("The Manzanita Micro Zilla interface is now running on port '%d'.\n", port)

	app.Listen(port)
}
