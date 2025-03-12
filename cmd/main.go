package main

import (
	"github.com/discosat/storage-system/cmd/dam"
	"github.com/discosat/storage-system/cmd/dim"
	"log"
)

func main() {
	log.Println("starting DIM-DAM system backend")

	go dam.Start()
	go dim.Start()

	log.Println("DIM-DAM up and running")

	select {}
}
