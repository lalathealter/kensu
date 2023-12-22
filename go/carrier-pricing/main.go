package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lalathealter/kensu/controllers"
)

const (
	HOSTkey = "host"
	PORTkey = "port"
)

func main() {
	h, p := os.Getenv(HOSTkey), os.Getenv(PORTkey)
	c := h + ":" + p
	fmt.Println("Hello, World: Serving at", c)

	mux := http.NewServeMux()
	mux.HandleFunc("/quotes", controllers.HandleQuotes)
	muxM := controllers.LogErrors(mux)
	log.Fatal(http.ListenAndServe(c, muxM))
}
