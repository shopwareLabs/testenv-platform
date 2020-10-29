package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/shopwareLabs/testenv-platform/handler"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.GET("/", handler.ListContainer)
	router.DELETE("/", handler.DeleteContainer)
	log.Println("Go!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
