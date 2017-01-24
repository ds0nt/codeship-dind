package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

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
		fmt.Println("\t", image.ID, image.Labels, image.RepoTags)
		fmt.Printf("\t\t size: %d\n", image.Size)
		// r, err := cli.(context.Background(), []string{image.ID})
		// if err != nil {
		// 	return errors.Wrapf(err, "Could not save image %s to tar", image.ID)
		// }
		// outFile, err := os.Create(image.ID + ".tar")
		// // handle err
		// defer outFile.Close()
		// _, err = io.Copy(outFile, r)
		// if err != nil {
		// 	return errors.Wrapf(err, "io.Copy error %s", image.ID)
		// }
	}

	fmt.Println(exec.Command("ls -l").Output())
	return nil
}
