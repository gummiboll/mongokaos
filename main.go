package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gummiboll/mongokaos/handlers"
	"github.com/gummiboll/mongokaos/state"
)

func main() {
	state := state.GetAppState()
	mux := http.NewServeMux()

	log.Printf("Listening on :%s", state.Config.ListenPort)

	mux.HandleFunc("POST /action/{action}", func(w http.ResponseWriter, r *http.Request) { handlers.ApiHandler(w, r, state) })

	http.ListenAndServe(fmt.Sprintf(":%s", state.Config.ListenPort), mux)
}
