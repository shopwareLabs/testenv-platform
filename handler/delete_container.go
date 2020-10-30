package handler

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"strings"
)

func DeleteContainer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	keys, ok := r.URL.Query()["id"]

	defer deleteResponse(w)

	if !ok || len(keys[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'id' is missing")
		return
	}

	id := string(keys[0])

	container, err := dClient.ContainerInspect(ctx, id)

	if err != nil {
		log.Println(err)
		return
	}

	err = dClient.ContainerKill(ctx, id, "SIGKILL")

	if err != nil {
		log.Println(err)
		return
	}

	err = dClient.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: true})

	if err != nil {
		log.Println(err)
		return
	}

	for _, bind := range container.HostConfig.Binds {
		_ = os.RemoveAll(strings.Split(bind, ":")[0])
	}
}

func deleteResponse(w http.ResponseWriter)  {
	res, _ := json.Marshal(map[string]bool{"success": true})

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}