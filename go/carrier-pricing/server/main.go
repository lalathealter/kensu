package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lalathealter/kensu/server/controllers"
	"github.com/lalathealter/kensu/server/db"
)

const (
	HOSTkey = "host"
	PORTkey = "port"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	h, p := os.Getenv(HOSTkey), os.Getenv(PORTkey)
	c := h + ":" + p
	fmt.Println("Hello, World: Serving at", c)

	mux := http.NewServeMux()
	mux.HandleFunc("/quotes", controllers.HandleQuotes)
	muxM := controllers.LogErrors(mux)
	db.InitDB()

	log.Fatal(http.ListenAndServe(c, muxM))
}
