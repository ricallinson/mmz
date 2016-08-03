package main

import (
	"fmt"
	"github.com/goforgery/forgery2"
	"github.com/goforgery/mustache"
	"github.com/goforgery/static"
)

// Starts the web server.
func StartWebServer(port int, zilla *Zilla) {

	app := f.CreateApp()
	app.Use(static.Create())
	app.Engine(".html", mustache.Create())

	app.Get("/", func(req *f.Request, res *f.Response, next func()) {
		res.Render("index.html", ReadLatestFromDataStream())
	})

	app.Get("/datastream", func(req *f.Request, res *f.Response, next func()) {
		res.Send(ReadLatestFromDataStream())
	})

	app.Get("/settings", func(req *f.Request, res *f.Response, next func()) {
		res.Render("settings.html", zilla)
	})

	fmt.Printf("The Manzanita Micro Zilla interface is now running on port '%d'.\n", port)

	app.Listen(port)
}
