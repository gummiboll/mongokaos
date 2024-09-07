package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gummiboll/mongokaos/handlers"
	"github.com/gummiboll/mongokaos/state"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	state := state.GetAppState()
	log.Printf("Listening on :%s", state.Config.ListenPort)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /action/{action}", func(w http.ResponseWriter, r *http.Request) { handlers.ApiHandler(w, r, state) })

	http.ListenAndServe(fmt.Sprintf(":%s", state.Config.ListenPort), mux)
}
