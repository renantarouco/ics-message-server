package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("ICS - Message Server")
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}
