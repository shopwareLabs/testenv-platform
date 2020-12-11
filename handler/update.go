package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/hashicorp/go-version"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

var newestShopwareVersion = ""

func getNewestShopwareImage() string {
	if len(newestShopwareVersion) == 0 {
		return "shopware/testenv:6.3.4"
	}

	return newestShopwareVersion
}

func PullImageUpdatesTask() {
	for {
		PullImageUpdates()
		time.Sleep(time.Hour * 24)
	}
}

func PullImageUpdates() {
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

	versions := make([]*version.Version, len(dockerHubResponse.Results))

	for i, tag := range dockerHubResponse.Results {
		v, _ := version.NewVersion(tag.Name)
		versions[i] = v

		imageName := fmt.Sprintf("docker.io/shopware/testenv:%s", tag.Name)
		log.Printf("Pullling image %s", imageName)
		outputReader, err := dClient.ImagePull(context.Background(), imageName, types.ImagePullOptions{})

		if err != nil {
			log.Println(err)
		}

		text, _ := ioutil.ReadAll(outputReader)
		fmt.Println(string(text))
	}

	sort.Sort(version.Collection(versions))
	newestShopwareVersion = versions[len(versions)-1].String()
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
