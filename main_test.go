package main

import (
	"log"
	"testing"
)

func TestDockerImages(t *testing.T) {
	err := DockerImages(1)
	if err != nil {
		panic("NOooooes!")
	}
	log.Println("Hooray!")
}
