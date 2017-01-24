package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/jackc/pgx"
)

var (
	dockerClient  *client.Client
	postgresImage = "clkao/postgres-plv8"
)

type killFunc func() error

func main() {
	// var err error
	// dockerClient, err = client.NewEnvClient()
	// if err != nil {
	// 	log.Println("Error:", err)
	// }
	// id, kill := spinPostgres(5432)
	// defer kill()

	ok := waitPGConn()
	if !ok {
		log.Println("Did not connect")
		return
	}

	log.Println("Postgres Container Available")
}

// func printAnyWarnings(warnings []string) {
// 	if len(warnings) > 0 {
// 		for _, warn := range warnings {
// 			fmt.Println("ContainerCreate Warning:", warn)
// 		}
// 	}
// }

// func spinPostgres(pgPort int) (string, killFunc) {
// 	ctx := context.Background()
// 	out, err := exec.Command("/bin/sh", "-c", "/sbin/ip route|awk '/default/ { print $3 }'").Output()
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println("Running from", string(out))

// 	r, err := dockerClient.ImagePull(ctx, postgresImage, types.ImagePullOptions{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	io.Copy(os.Stdout, r)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = r.Close()
// 	if err != nil {
// 		panic(err)
// 	}

// 	res, err := dockerClient.ContainerCreate(
// 		ctx,
// 		&container.Config{
// 			Cmd:   []string{"postgres"},
// 			Image: postgresImage,
// 			ExposedPorts: nat.PortSet{
// 				"5432": struct{}{},
// 			},
// 			Env: []string{
// 				fmt.Sprintf("POSTGRES_USER=%s", os.Getenv("POSTGRES_USER")),
// 				fmt.Sprintf("POSTGRES_PASSWORD=%s", os.Getenv("POSTGRES_PASSWORD")),
// 				fmt.Sprintf("POSTGRES_DB=%s", os.Getenv("POSTGRES_DB")),
// 				fmt.Sprintf("POSTGRES_SSLMODE=%s", os.Getenv("POSTGRES_SSLMODE")),
// 			},
// 		},
// 		&container.HostConfig{
// 			PortBindings: nat.PortMap{
// 				"5432": []nat.PortBinding{
// 					nat.PortBinding{
// 						HostIP:   "",
// 						HostPort: fmt.Sprintf("%d", pgPort),
// 					},
// 				},
// 			},
// 		},
// 		&network.NetworkingConfig{},
// 		"",
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	inspect, err := dockerClient.ContainerInspect(ctx, res.ID)
// 	if err != nil {
// 		panic(err)
// 	}
// 	ipAddr := inspect.NetworkSettings.IPAddress
// 	fmt.Printf("Created container %s\n", res.ID)
// 	fmt.Printf("Container IP %s", ipAddr)

// 	printAnyWarnings(res.Warnings)
// 	err = dockerClient.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	kFn := func() error {
// 		fmt.Println("Killing container", res.ID)
// 		killErr := dockerClient.ContainerRemove(ctx, res.ID, types.ContainerRemoveOptions{Force: true})
// 		if killErr != nil {
// 			panic(killErr)
// 		}
// 		return killErr
// 	}

// 	logsResp, err := dockerClient.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{
// 		ShowStdout: true,
// 		ShowStderr: true,
// 		Follow:     true,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	go func() {
// 		for {
// 			_, err := io.Copy(os.Stdout, logsResp)
// 			if err != nil {
// 				return
// 			}

// 		}
// 	}()
// 	log.Println("Testing Postgres Connection")

// 	ok, err := waitPGConn()
// 	if err != nil {

// 	}
// 	return ipAddr, kFn
// }

func waitPGConn() (ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	for {
		pg, err := pgx.Connect(pgx.ConnConfig{
			Host:     "postgres",
			Port:     5432,
			Database: os.Getenv("POSTGRES_DB"),
			User:     "postgres",
			Password: os.Getenv("POSTGRES_PASSWORD"),
		})
		if err == nil {
			pg.Close()
			return true
		}

		select {
		case <-ctx.Done():
			return false
		default:
			log.Println("Still waiting for connection:", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

// func DockerImages() error {

// 	images, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{
// 		All: true,
// 	})
// 	if err != nil {
// 		return errors.Wrap(err, "Error creating NewEnvClient")
// 	}

// 	fmt.Println("List of Docker Images")
// 	for _, image := range images {
// 		fmt.Println("\t", image.ID, image.Labels, image.RepoTags)
// 		fmt.Printf("\t\t size: %d\n", image.Size)
// 		// r, err := cli.(context.Background(), []string{image.ID})
// 		// if err != nil {
// 		// 	return errors.Wrapf(err, "Could not save image %s to tar", image.ID)
// 		// }
// 		// outFile, err := os.Create(image.ID + ".tar")
// 		// // handle err
// 		// defer outFile.Close()
// 		// _, err = io.Copy(outFile, r)
// 		// if err != nil {
// 		// 	return errors.Wrapf(err, "io.Copy error %s", image.ID)
// 		// }
// 	}

// 	fmt.Println(exec.Command("ls -l").Output())
// 	return nil
// }
