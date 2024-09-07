package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func ReturnResult(w http.ResponseWriter, result interface{}, debug bool) {
	ejson, err := bson.MarshalExtJSON(result, false, false)
	if err != nil {
		log.Println("Error marshaling result: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if debug {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, ejson, "", "  ")
		log.Printf("Response:\n%s", prettyJSON.String())
	}
	w.Header().Set("Content-Type", "application/ejson")
	w.Write(ejson)
}

func Authenticate(r *http.Request, apiKey string) bool {
	if r.Header.Get("api-key") != apiKey {
		log.Printf("Unauthorized request from %s", r.RemoteAddr)
		return false
	}
	return true
}
