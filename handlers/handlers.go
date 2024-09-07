package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gummiboll/mongokaos/mongodb"
	"github.com/gummiboll/mongokaos/state"
	"github.com/gummiboll/mongokaos/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ApiHandler(w http.ResponseWriter, r *http.Request, appState *types.AppState) {
	var reqData types.RequestData
	ctx := context.Background()
	action := r.PathValue("action")

	log.Printf(`%s %s`, r.Method, r.URL)
	if r.Header.Get("api-key") != appState.Config.APIKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	if state.GetAppState().Config.Debug {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, bodyBytes, "", "  ")
		log.Printf("Request:\n%s", prettyJSON.String())
	}

	err = bson.UnmarshalExtJSON(bodyBytes, false, &reqData)
	if err != nil {
		log.Println("Error unmarshaling json: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := appState.Mongo.Database(reqData.Database).Collection(reqData.Collection)

	result, err := mongodb.ExecuteAction(action, ctx, collection, reqData)
	if err != nil {
		log.Println("Error: ", err)
	}

	var reply interface{}
	switch result.(type) {
	case *mongo.SingleResult:
		var resultData bson.M
		var res types.SingleResult
		_ = result.(*mongo.SingleResult).Decode(&resultData)
		res.Document = resultData
		reply = res

	case *mongo.Cursor:
		var resultData []bson.M
		var res types.MultipleResults

		defer result.(*mongo.Cursor).Close(ctx)
		if err := result.(*mongo.Cursor).All(ctx, &resultData); err != nil {
			log.Println("Error: ", err)
		}
		if resultData == nil {
			resultData = []bson.M{}
		}

		res.Documents = resultData
		reply = res

	default:
		reply = result
	}

	ejson, err := bson.MarshalExtJSON(reply, false, false)
	if err != nil {
		log.Println("Error marshaling result: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if state.GetAppState().Config.Debug {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, ejson, "", "  ")
		log.Printf("Response:\n%s", prettyJSON.String())
	}

	w.Header().Set("Content-Type", "application/ejson")
	w.Write(ejson)
}
