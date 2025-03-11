package main

import (
	"github.com/discosat/storage-system/cmd/dam"
	"log"
)

func main() {
	log.Println("starting DIM-DAM system backend")

	go dam.Start()

	select {}
}
