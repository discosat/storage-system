package main

import (
	"github.com/discosat/storage-system/cmd/dim"
	"log"
)

func main() {
	log.Println("starting DIM-DAM system backend")

	//go dam.Start()
	go dim.Start()

	select {}
}
