package legacy

import (
	"fmt"
	"log"
)

func legacyMain() {
	fmt.Println("ICS - Message Server")
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}
