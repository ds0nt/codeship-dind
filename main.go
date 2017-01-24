package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func main() {
	err := DockerImages()
	if err != nil {
		log.Println("Error:", err)
	}
}

func DockerImages() error {
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
