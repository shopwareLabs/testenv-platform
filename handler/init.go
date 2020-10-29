package handler

import (
	"context"
	"github.com/docker/docker/client"
)

var ctx context.Context
var dClient *client.Client

func init() {
	ctx = context.Background()

	var err error

	dClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}
