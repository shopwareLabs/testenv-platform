package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/hashicorp/go-version"
)

var newestShopwareVersion = ""

func getNewestShopwareImage() string {
	if len(newestShopwareVersion) == 0 {
		return "ghcr.io/shopwarelabs/testenv:6.3.4"
	}

	return newestShopwareVersion
}

func PullImageUpdatesTask(token string) {
	for {
		PullImageUpdates(token)
		time.Sleep(time.Hour * 24)

		log.Println("Completed update task")
	}
}

func PullImageUpdates(token string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/users/shopwareLabs/packages/container/testenv/versions?per_page=100", nil)
	req.Header.Set("User-Agent", "shopware/testenv Client")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	respContent, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var githubApiResponse GithubApiResponse

	err = json.Unmarshal(respContent, &githubApiResponse)

	if err != nil {
		log.Println(err)
		return
	}

	versionCount := 0

	for _, image := range githubApiResponse {
		if len(image.Metadata.Container.Tags) == 0 {
			continue
		}
		versionCount = versionCount + 1
	}

	versions := make([]*version.Version, versionCount)

	for i, image := range githubApiResponse {
		if len(image.Metadata.Container.Tags) == 0 {
			continue
		}

		tag := image.Metadata.Container.Tags[0]

		v, _ := version.NewVersion(tag)
		versions[i] = v

		imageName := fmt.Sprintf("ghcr.io/shopwarelabs/testenv:%s", tag)
		log.Printf("Pullling image %s", imageName)
		outputReader, err := dClient.ImagePull(context.Background(), imageName, types.ImagePullOptions{})

		if err != nil {
			log.Println(err)
		}

		text, _ := ioutil.ReadAll(outputReader)
		fmt.Println(string(text))
	}

	sort.Sort(version.Collection(versions))
	newestShopwareVersion = fmt.Sprintf("ghcr.io/shopwarelabs/testenv:%s", versions[len(versions)-1].String())
}

type GithubApiResponse []struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	PackageHTMLURL string    `json:"package_html_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	HTMLURL        string    `json:"html_url"`
	Metadata       struct {
		PackageType string `json:"package_type"`
		Container   struct {
			Tags []string `json:"tags"`
		} `json:"container"`
	} `json:"metadata"`
}
