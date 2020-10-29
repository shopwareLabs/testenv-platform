package handler

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func ListContainer(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	opts := types.ContainerListOptions{}
	opts.Filters = filters.NewArgs()
	opts.Filters.Add("label", "testenv=1")

	containers, err := dClient.ContainerList(ctx, opts)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Cannot list containers")
	}

	jData, _ := json.Marshal(containers)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jData)
}
