package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("ICS - Message Server")
	err := http.ListenAndServe(":7000", enableCORS(msgServerHTTPRouter))
	if err != nil {
		log.Fatalln(err.Error())
	}
}
