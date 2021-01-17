package main

import (
	"classwork/backend"
	"classwork/fileserver"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not find .env file!")
	}

	rand.Seed(time.Now().UnixNano())

	wg := sync.WaitGroup{}
	wg.Add(2)

	go backend.Run(&wg)
	go fileserver.Run(&wg)

	wg.Wait()
}
