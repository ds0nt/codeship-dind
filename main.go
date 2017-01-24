package main

import (
	"log"
	"time"

	"github.com/pkg/errors"
)

func main() {
	for i := 0; i < 5; i++ {
		DoStuff(i)
		time.Sleep(1 * time.Second)
	}
}

func DoStuff(num int) error {
	if num > 10 {
		return errors.Errorf("Our number %d was greater than 10", num)
	}
	log.Printf("Doing Stuff!!! %d", num)
	return nil
}
