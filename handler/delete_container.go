package handler

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func DeleteContainer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'id' is missing")
		return
	}

	id := string(keys[0])

	err := dClient.ContainerKill(ctx, id, "SIGKILL")

	if err != nil {
		log.Println(err)
	}

	err = dClient.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: true})

	if err != nil {
		log.Println(err)
	}

	res, _ := json.Marshal(map[string]bool{"success": true})

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}
