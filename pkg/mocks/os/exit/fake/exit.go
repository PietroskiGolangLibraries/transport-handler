package fake

import (
	"log"
	"os"
)

func Exit(code int) {
	log.Printf("exiting with code: %v\n", code)
	os.Exit(0)
}
