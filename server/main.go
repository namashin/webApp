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
	gateway := flag.String("gateway", "http:/127.0.0.1:5000", "web server gateway")
	flag.Parse()

	app := NewWebServer(uint16(*port), *gateway)
	app.Run()

}
