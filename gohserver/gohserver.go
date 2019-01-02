package gohserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abilioesteves/goh/gohtypes"
	"github.com/sirupsen/logrus"
)

// HandleError handles unexpected errors, keeping the response message clean
// Use it by deferring on first line of any http handler
func HandleError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		if err, ok := r.(gohtypes.Error); ok {
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
	gohtypes.PanicIfError(fmt.Sprintf("Not possible to write %v response", status), 500, err)

	logrus.Infof("200 Response sent. Payload: %s", payload)
}
