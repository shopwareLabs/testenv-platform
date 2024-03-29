package handler

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CreateEnvironment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := randSeq(5)
	request, err := getPluginInformationFromRequest(id, r)
	scheme := "http"

	if os.Getenv("USE_HTTPS") == "1" {
		scheme = "https"
	}

	if err != nil {
		log.Printf("%s", err)
		apiResponse(w, ApiResponse{Success: false, Message: fmt.Sprintf("%s", err)}, http.StatusInternalServerError)
		return
	}

	log.Printf("Requested environment for version: %s and plugin: %s", request.InstallVersion, request.Name)

	instanceName := strings.ToLower(fmt.Sprintf("%s-%s", request.Name, id))
	host := strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s.%s", instanceName, os.Getenv("BASE_HOST"))), "_", "-")

	imageName, err := getImage(request)
	if err != nil {
		log.Printf("%s", err)
		apiResponse(w, ApiResponse{Success: false, Message: fmt.Sprintf("%s", err)}, http.StatusInternalServerError)
		return
	}

	appUrl := fmt.Sprintf("%s://%s/shop/public", scheme, host)
	labels := map[string]string{
		"testenv":        "1",
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routers.http-%s.rule", instanceName):        fmt.Sprintf("Host(`%s`)", host),
		fmt.Sprintf("traefik.http.routers.http-%s.entrypoints", instanceName): "web",
	}

	if len(os.Getenv("SSL_PROVIDER")) > 0 {
		labels[fmt.Sprintf("traefik.http.routers.http-%s.middlewares", instanceName)] = "web-redirect@file"
		labels[fmt.Sprintf("traefik.http.routers.https-%s.middlewares", instanceName)] = "compress@file"
		labels[fmt.Sprintf("traefik.http.routers.https-%s.rule", instanceName)] = fmt.Sprintf("Host(`%s`)", host)
		labels[fmt.Sprintf("traefik.http.routers.https-%s.entrypoints", instanceName)] = "websecure"
		labels[fmt.Sprintf("traefik.http.routers.https-%s.tls", instanceName)] = "true"
		labels[fmt.Sprintf("traefik.http.routers.https-%s.tls.certresolver", instanceName)] = os.Getenv("SSL_PROVIDER")
		labels[fmt.Sprintf("traefik.http.routers.https-%s.tls.domains[0].main", instanceName)] = os.Getenv("BASE_HOST")
		labels[fmt.Sprintf("traefik.http.routers.https-%s.tls.domains[0].sans", instanceName)] = fmt.Sprintf("*.%s", os.Getenv("BASE_HOST"))
	}

	cConfig := &container.Config{
		Image: imageName,
		Env: []string{
			fmt.Sprintf("PLUGIN_NAME=%s", request.Name),
			fmt.Sprintf("VIRTUAL_HOST=%s", host),
			fmt.Sprintf("APP_URL=%s", appUrl),
			fmt.Sprintf("SHOPWARE_DEMO_USER_PASSWORD=%s", request.ShopwarePassword),
		},
		Labels: labels,
	}

	cHost := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:%s", request.VolumeFolder, fmt.Sprintf("/var/www/shop/%s", request.MountFolder)),
		},
	}
	cNetwork := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"docker_default": {},
		},
	}

	cBody, err := dClient.ContainerCreate(ctx, cConfig, cHost, cNetwork, nil, instanceName)

	if err != nil {
		log.Printf("%s", err)
		apiResponse(w, ApiResponse{Success: false, Message: fmt.Sprintf("%s", err)}, http.StatusInternalServerError)
		return
	}

	err = dClient.ContainerStart(ctx, cBody.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Printf("%s", err)
		apiResponse(w, ApiResponse{Success: false, Message: fmt.Sprintf("%s", err)}, http.StatusInternalServerError)
		return
	}

	apiResponse(w, EnvironmentCreated{
		ID:               cBody.ID,
		URL:              appUrl,
		InstallVersion:   request.InstallVersion,
		ShopwarePassword: request.ShopwarePassword,
	}, http.StatusOK)
}

func getPluginInformationFromRequest(id string, r *http.Request) (*PluginInformation, error) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	var env EnvironmentRequest
	err = json.Unmarshal(body, &env)

	if err != nil {
		return nil, err
	}

	pluginZipContentByte, err := base64.StdEncoding.DecodeString(env.PluginZipEncoded)

	if err != nil {
		return nil, err
	}

	pluginZipContent := string(pluginZipContentByte)

	zipReader, err := zip.NewReader(strings.NewReader(pluginZipContent), int64(len(pluginZipContent)))

	if err != nil {
		return nil, err
	}

	result := PluginInformation{EnvironmentRequest: env}

	if len(zipReader.File) == 0 {
		return nil, fmt.Errorf("Zip is empty")
	}

	names := strings.Split(zipReader.File[0].Name, "/")

	result.Name = names[0]

	if env.InstallVersion == "app" {
		result.MountFolder = "custom/apps/"
	} else if result.Name == "Backend" || result.Name == "Core" || result.Name == "Frontend" {
		result.MountFolder = "engine/Shopware/Plugins/Local/"
	} else {
		result.MountFolder = "custom/plugins/"
	}

	result.VolumeFolder, err = ioutil.TempDir("", "plugin")
	if err != nil {
		return nil, err
	}

	err = Unzip(zipReader, result.VolumeFolder)

	err = os.Chown(result.VolumeFolder, 1000, 1000)
	if err != nil {
		return nil, err
	}

	if LookupEnvOrString("GENERATE_INSTANCE_PASSWORDS", "0") == "1" {
		result.ShopwarePassword = randSeq(10)
	} else {
		result.ShopwarePassword = "demo"
	}

	return &result, nil
}

func Unzip(r *zip.Reader, dest string) error {
	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}

		err = os.Chown(fpath, 1000, 1000)
		if err != nil {
			return err
		}
	}

	return nil
}

func getImage(info *PluginInformation) (string, error) {
	if info.InstallVersion == "app" {
		return getNewestShopwareImage(), nil
	}

	var v1, v2, v3 int
	var imageTag string

	_, err := fmt.Sscanf(info.InstallVersion, "%d.%d.%d", &v1, &v2, &v3)
	if err != nil {
		return "", err
	}

	if v2 >= 3 {
		imageTag = fmt.Sprintf("%d.%d.%d", v1, v2, v3)
	} else {
		imageTag = fmt.Sprintf("%d.%d", v1, v2)
	}

	imageName := fmt.Sprintf("ghcr.io/shopwarelabs/testenv:%s", imageTag)
	opts := types.ImageListOptions{}
	opts.Filters = filters.NewArgs()
	opts.Filters.Add("before", imageName)

	_, err = dClient.ImageList(ctx, opts)
	if err != nil {
		return "", err
	}

	return imageName, nil
}

type EnvironmentRequest struct {
	InstallVersion   string `json:"installVersion"`
	PluginZipEncoded string `json:"plugin"`
	ShopwarePassword string
}

type PluginInformation struct {
	Name         string
	VolumeFolder string
	MountFolder  string
	EnvironmentRequest
}

type EnvironmentCreated struct {
	ID               string `json:"id"`
	URL              string `json:"domain"`
	InstallVersion   string `json:"installVersion"`
	ShopwarePassword string `json:"shopwarePassword"`
}
