package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/shopwareLabs/testenv-platform/handler"
	"log"
	"net/http"
)

func main() {
	go PullImageUpdatesTask()

	router := httprouter.New()

	// New Routes
	router.GET("/environments", handler.ListContainer)
	router.POST("/environments", handler.CreateEnvironment)
	router.DELETE("/environments", handler.DeleteContainer)

	log.Println("Go!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
