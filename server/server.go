package gohserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/abilioesteves/goh/types"
)

// HandleError handles unexpected errors, keeping the response message clean
// Use it by deferring on first line of any http handler
func HandleError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		if err, ok := r.(types.Error); ok {
			logrus.Error(err)
			http.Error(w, err.Message, err.Code)
		} else {
			logrus.Error(r)
			http.Error(w, "Internal Error", 500)
		}
	}
}

// WriteJSONResponse writes the response to be sent
func WriteJSONResponse(payload interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	types.PanicIfError(types.Error{Message: fmt.Sprintf("Not possible to write %v response", status), Code: 500, Err: err})

	logrus.Infof("200 Response sent. Payload: %s", payload)
}
