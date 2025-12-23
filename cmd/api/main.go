package main

import (
	"fmt"

	"github.com/igromancer/deck-assignment/internal/server"
)

func main() {
	fmt.Println("running API")
	server.Listen()
}
