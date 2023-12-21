package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lalathealter/kensu/controllers"
)

func main() {
	h := "localhost:8080"
	fmt.Println("Hello, World: Serving at", h)

	mux := http.NewServeMux()
	mux.HandleFunc("/quotes", controllers.HandleQuotes)
	muxM := controllers.LogErrors(mux)
	log.Fatal(http.ListenAndServe(h, muxM))
}
