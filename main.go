package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gummiboll/mongokaos/handlers"
	"github.com/gummiboll/mongokaos/middleware"
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
	mux.Handle("POST /action/{action}",
		middleware.LoggingMiddleware(
			middleware.AuthMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handlers.ApiHandler(w, r)
				}),
			),
		))
	http.ListenAndServe(fmt.Sprintf(":%s", state.Config.ListenPort), mux)
}
