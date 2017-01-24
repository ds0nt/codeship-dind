package main

import (
	"log"
	"testing"
)

func TestDoStuff(t *testing.T) {
	err := DoStuff(1)
	if err != nil {
		panic("NOooooes!")
	}
	err = DoStuff(12)
	if err != nil {
		panic("NOooooes!")
	}
	log.Println("Hooray!")
}
