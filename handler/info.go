package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Info(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{\"message\": \"API server up and running\"}"))
}
