package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func main() {
	for i := 0; i < 5; i++ {
		DockerImages(i)
		time.Sleep(1 * time.Second)
	}
}

func DockerImages(num int) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "Error creating NewEnvClient")
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return errors.Wrap(err, "Error creating NewEnvClient")
	}

	fmt.Println("List of Docker Images")
	for _, image := range images {
		fmt.Println("\t", image)
	}

	return nil
}
