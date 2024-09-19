package main

import (
	"fmt"

	"github.com/Ananth1082/mv0/server"
)

func main() {
	fmt.Println("Hello world")
	s := server.CreateServer(":8080")
	s.StartServer()
}
