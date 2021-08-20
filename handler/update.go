package handler

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"github.com/hashicorp/go-version"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
)

var newestShopwareVersion = ""

func getNewestShopwareImage() string {
	if len(newestShopwareVersion) == 0 {
		return "ghcr.io/shopwarelabs/testenv:6.3.4"
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
	log.Println("Pulling images")
	outputReader, err := dClient.ImagePull(context.Background(), "ghcr.io/shopwarelabs/testenv", types.ImagePullOptions{All: true})

	if err != nil {
		log.Println(err)
		return
	}

	text, _ := ioutil.ReadAll(outputReader)
	fmt.Println(string(text))

	log.Println("Pulled new images")
	log.Println("Detecting newest version")

	opts := types.ImageListOptions{}

	opts.Filters = filters.NewArgs()
	opts.Filters.Add("reference", "ghcr.io/shopwarelabs/testenv")

	images, err := dClient.ImageList(context.Background(), opts)

	if err != nil {
		log.Println(err)
		return
	}

	versions := make([]*version.Version, 0)

	for _, image := range images {
		if len(image.RepoTags) == 0 {
			continue
		}

		v, err := version.NewVersion(strings.Replace(image.RepoTags[0], "ghcr.io/shopwarelabs/testenv:", "", 1))

		if err != nil {
			log.Println(err)
			continue
		}

		versions = append(versions, v)
	}

	sort.Sort(version.Collection(versions))
	newestShopwareVersion = fmt.Sprintf("ghcr.io/shopwarelabs/testenv:%s", versions[len(versions)-1].String())

	log.Println("Completed update task")
	log.Printf("New newest Shopware version is: %s\n", newestShopwareVersion)
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
