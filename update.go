package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/shopwareLabs/testenv-platform/handler"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func PullImageUpdatesTask() {
	for {
		PullImageUpdates()
		time.Sleep(time.Hour * 24)
	}
}

func PullImageUpdates() {
	client := handler.GetDocker()

	resp, err := http.Get("https://hub.docker.com/v2/repositories/shopware/testenv/tags/?page_size=25&page=1")
	if err != nil {
		log.Println(err)
		return
	}

	respContent, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var dockerHubResponse DockerHubTagsResponse

	err = json.Unmarshal(respContent, &dockerHubResponse)

	if err != nil {
		log.Println(err)
		return
	}

	for _, tag := range dockerHubResponse.Results {
		imageName := fmt.Sprintf("docker.io/shopware/testenv:%s", tag.Name)
		log.Printf("Pullling image %s", imageName)
		outputReader, err := client.ImagePull(context.Background(), imageName, types.ImagePullOptions{})

		if err != nil {
			log.Println(err)
		}

		text, _ := ioutil.ReadAll(outputReader)
		fmt.Println(string(text))
	}
}

type DockerHubTagsResponse struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Creator int         `json:"creator"`
		ID      int         `json:"id"`
		ImageID interface{} `json:"image_id"`
		Images  []struct {
			Architecture string      `json:"architecture"`
			Features     string      `json:"features"`
			Variant      interface{} `json:"variant"`
			Digest       string      `json:"digest"`
			Os           string      `json:"os"`
			OsFeatures   string      `json:"os_features"`
			OsVersion    interface{} `json:"os_version"`
			Size         int         `json:"size"`
			Status       string      `json:"status"`
			LastPulled   interface{} `json:"last_pulled"`
			LastPushed   interface{} `json:"last_pushed"`
		} `json:"images"`
		LastUpdated         time.Time   `json:"last_updated"`
		LastUpdater         int         `json:"last_updater"`
		LastUpdaterUsername string      `json:"last_updater_username"`
		Name                string      `json:"name"`
		Repository          int         `json:"repository"`
		FullSize            int         `json:"full_size"`
		V2                  bool        `json:"v2"`
		TagStatus           string      `json:"tag_status"`
		TagLastPulled       interface{} `json:"tag_last_pulled"`
		TagLastPushed       time.Time   `json:"tag_last_pushed"`
	} `json:"results"`
}
