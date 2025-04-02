package main

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/dam"
)

func main() {
	fmt.Println("starting DIM-DAM system backend")

	go dam.Start()
	//go dim.Start()

	fmt.Println("DIM-DAM up and running")

	select {}
}
