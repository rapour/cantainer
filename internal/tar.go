package internal

import (
	"archive/tar"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rapour/cantainer/images"
)

func Extract(destination string) {

	r, err := images.Images.Open(images.Alpine)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	tr := tar.NewReader(r)

	for {

		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}

		if header == nil {
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(destination, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
					log.Fatal(err)
				}
			}

		// if it's a file create it
		case tar.TypeReg:

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				log.Fatal(err)
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				log.Fatal(err)
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()

		case tar.TypeSymlink:
			os.Symlink(header.Linkname, target)

		default:
			fmt.Printf("Unsupported type: %c in file %s\n", header.Typeflag, target)

		}
	}
}
