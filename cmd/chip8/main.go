package main

import (
	"fmt"
	"log"

	"github.com/ejholmes/chip8"
)

func main() {
	c, err := chip8.NewCPU(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)
}
