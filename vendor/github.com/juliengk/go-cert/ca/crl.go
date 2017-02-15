package ca

import (
	"os"
)

func CreateCRLFile(filename string) {
	os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0664)
}
