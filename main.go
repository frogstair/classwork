package main

import (
	"classwork/backend"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go backend.Run(&wg)

	wg.Wait()
}
