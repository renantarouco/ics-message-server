package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("starting message server")
	if err := server.httpServer.ListenAndServe(); err != nil {
		log.Fatalln(err.Error())
	}
}
