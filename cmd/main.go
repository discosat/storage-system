package main

import (
	"fmt"
	"github.com/discosat/storage-system/cmd/dam"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	fmt.Println("starting DIM-DAM system backend")

	//Initialize environment variables
	err := godotenv.Load("cmd/dam/.env")
	if err != nil {
		log.Fatalf("NewMinioStore: Cant find env - %v", err)
	}
	// Initialize the database connection
	dam.InitDB()

	// Initialize DIM and DAM services
	go dam.Start()
	//go dim.Start()

	fmt.Println("DIM-DAM up and running")

	select {}
}
