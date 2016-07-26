package main

import (
	"flag"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port number to run the service on")
	flag.Parse()
	CreateZilla() // Tmp
	StartWebServer(port)
}
