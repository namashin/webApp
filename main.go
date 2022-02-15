package main

import (
	"flag"
	"log"
	"sync"

	"webApp/server"
)

var once sync.Once

func init() {
	log.SetPrefix("Server: ")
	go once.Do(server.Start)
}

func main() {

	port := flag.Uint("port", 8000, "Tcp Port Number for web server")
	gateway := flag.String("gateway", "http:/127.0.0.1:5000", "web server gateway")
	flag.Parse()

	app := server.NewWebServer(uint16(*port), *gateway)
	app.Run()

}
