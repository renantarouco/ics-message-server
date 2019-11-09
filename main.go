package main

import (
	"log"
	"os"
)

func main() {
	if err := Init(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	if err := Run(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
