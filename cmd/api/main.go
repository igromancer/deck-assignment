package main

import (
	"fmt"

	"github.com/igromancer/deck-assignment/internal/server"
)

func main() {
	fmt.Println("running API")
	s, err := server.NewServer()
	if err != nil {
		panic(fmt.Errorf("failed to create the server: %s", err.Error()))
	}
	err = s.Listen()
	if err != nil {
		panic(fmt.Errorf("failed to start the server: %s", err.Error()))
	}
}
