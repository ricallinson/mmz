package main

import (
    "fmt"
    "github.com/goforgery/forgery2"
)

// Starts the web server.
func StartWebServer(port int) {

    app := f.CreateApp()

    app.Get("/", func(req *f.Request, res *f.Response, next func()) {
        res.Send("Empty page.")
    })

    fmt.Printf("The Manzanita Micro Zilla interface is now running on port '%d'.\n", port)

    app.Listen(port)
}
