package handler

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/client"
	"math/rand"
	"net/http"
	"os"
)

var ctx context.Context
var dClient *client.Client

func init() {
	ctx = context.Background()

	var err error

	dClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func apiResponse(w http.ResponseWriter, aResp interface{}, statusCode int) {
	res, _ := json.Marshal(aResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(res)
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
