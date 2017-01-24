package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

var (
	dockerClient  *client.Client
	postgresImage = "clkao/postgres-plv8"
)

type killFunc func() error

func main() {
	var err error
	dockerClient, err = client.NewEnvClient()
	if err != nil {
		log.Println("Error:", err)
	}
	id, kill := spinPostgres(5432)
	defer kill()

	fmt.Println("Postgres Container Available: ", id)
}

func printAnyWarnings(warnings []string) {
	if len(warnings) > 0 {
		for _, warn := range warnings {
			fmt.Println("ContainerCreate Warning:", warn)
		}
	}
}

func spinPostgres(pgPort int) (string, killFunc) {
	ctx := context.Background()
	out, err := exec.Command("/bin/sh", "-c", "/sbin/ip route|awk '/default/ { print $3 }'").Output()
	if err != nil {
		panic(err)
	}
	log.Println("Running from", string(out))

	r, err := dockerClient.ImagePull(ctx, postgresImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, r)
	if err != nil {
		panic(err)
	}

	err = r.Close()
	if err != nil {
		panic(err)
	}

	res, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Cmd:   []string{"postgres"},
			Image: postgresImage,
			ExposedPorts: nat.PortSet{
				"5432": struct{}{},
			},
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", os.Getenv("POSTGRES_USER")),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", os.Getenv("POSTGRES_PASSWORD")),
				fmt.Sprintf("POSTGRES_DB=%s", os.Getenv("POSTGRES_DB")),
				fmt.Sprintf("POSTGRES_SSLMODE=%s", os.Getenv("POSTGRES_SSLMODE")),
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"5432": []nat.PortBinding{
					nat.PortBinding{
						HostIP:   "",
						HostPort: fmt.Sprintf("%d", pgPort),
					},
				},
			},
		},
		&network.NetworkingConfig{},
		"",
	)
	if err != nil {
		panic(err)
	}
	inspect, err := dockerClient.ContainerInspect(ctx, res.ID)
	if err != nil {
		panic(err)
	}
	ipAddr := inspect.NetworkSettings.IPAddress
	fmt.Printf("Created container %s\n", res.ID)
	fmt.Printf("Container IP %s", ipAddr)

	printAnyWarnings(res.Warnings)
	err = dockerClient.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	// wait till connect
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(35 * time.Second)
		timeout <- true
	}()

	kFn := func() error {
		fmt.Println("Killing container", res.ID)
		killErr := dockerClient.ContainerRemove(ctx, res.ID, types.ContainerRemoveOptions{Force: true})
		if killErr != nil {
			panic(killErr)
		}
		return killErr
	}

	logsResp, err := dockerClient.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			_, err := io.Copy(os.Stdout, logsResp)
			if err != nil {
				return
			}

		}
	}()
	log.Println("Testing Postgres Connection")
	var pg *pgx.Conn
	for {
		pg, err = pgx.Connect(pgx.ConnConfig{
			Host:     ipAddr,
			Port:     uint16(pgPort),
			Database: os.Getenv("POSTGRES_DB"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		})
		if err == nil {
			break
		}
		select {
		case <-timeout:
			kFn()
			panic(errors.Wrap(err, "postgres connection timeout to container"))
		default:
			log.Println("Failed to connect to ephemeral PG instance", err)
			time.Sleep(time.Second)
			continue
		}
	}
	defer pg.Close()
	log.Println("Postgres Connection Available")
	return ipAddr, kFn
}

func DockerImages() error {

	images, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{
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
