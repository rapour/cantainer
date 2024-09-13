package cantainer

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

const (
	baseDIR = "/tmp/cantainer"
	length  = 10
)

func CreateTempDir() (dir string, identifier string) {

	hasher := sha1.New()
	hasher.Write([]byte(time.Now().String()))
	identifier = string([]rune(hex.EncodeToString(hasher.Sum(nil)))[:length])

	if _, err := os.Stat(baseDIR); err != nil {
		if err := os.MkdirAll(baseDIR, 0755); err != nil {
			log.Fatal(err)
		}
	}

	dir = fmt.Sprintf("%s/%s", baseDIR, identifier)

	err := os.Mkdir(dir, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return
}
