package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Server: ")
}

func main() {

	port := flag.Uint("port", 8000, "Tcp Port Number for web server")
	flag.Parse()

	app := NewWebServer(uint16(*port))
	app.Run()
}
