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
	if Authenticate(r, appState.Config.APIKey) == false {
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

	switch result.(type) {
	case *mongo.SingleResult:
		var resultData bson.M
		var res types.SingleResult
		_ = result.(*mongo.SingleResult).Decode(&resultData)
		res.Document = resultData
		ReturnResult(w, res, state.GetAppState().Config.Debug)
		return

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
		ReturnResult(w, res, state.GetAppState().Config.Debug)
		return

	default:
		ReturnResult(w, result, state.GetAppState().Config.Debug)
		return
	}
}
