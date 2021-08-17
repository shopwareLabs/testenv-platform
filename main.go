package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/shopwareLabs/testenv-platform/handler"
	"log"
	"net/http"
	"os"
)

var (
	Github_PAT = ""
)

func main() {
	flag.StringVar(&Github_PAT, "github-token", LookupEnvOrString("GITHUB_TOKEN", ""), "Github token for authentication for automatic updates")

	if len(Github_PAT) > 0 {
		go handler.PullImageUpdatesTask(Github_PAT)
	}

	router := httprouter.New()

	router.GET("/", handler.Info)

	// New Routes
	router.GET("/environments", handler.ListContainer)
	router.POST("/environments", handler.CreateEnvironment)
	router.DELETE("/environments", handler.DeleteContainer)

	log.Println("Go!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
