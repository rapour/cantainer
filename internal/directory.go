package internal

import (
	"log"
	"os"
)

const baseDIR = "/tmp/cantainer"

func CreateTempDir() string {

	if _, err := os.Stat(baseDIR); err != nil {
		if err := os.MkdirAll(baseDIR, 0755); err != nil {
			log.Fatal(err)
		}
	}

	dir, err := os.MkdirTemp("/tmp/cantainer", "*")
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
