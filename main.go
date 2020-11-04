package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/shopwareLabs/testenv-platform/handler"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()

	// New Routes
	router.GET("/environments", handler.ListContainer)
	router.POST("/environments", handler.CreateEnvironment)
	router.DELETE("/environments", handler.DeleteContainer)

	// Remove when SBP has fixed the link :(
	router.GET("/index.php", handler.ListContainer)
	router.POST("/index.php", handler.CreateEnvironment)
	router.DELETE("/index.php", handler.DeleteContainer)

	log.Println("Go!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
