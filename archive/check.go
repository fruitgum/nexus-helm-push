package archive

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func CheckFileExtension(filename string) (string, error) {
	var err error
	err = errors.New("badExtension")
	extension := "tgz"
	fileExtension := filepath.Ext(filename)
	if strings.HasSuffix(fileExtension, extension) == false {
		return "Chart archive extension is wrong or missing", err
	}
	return "", nil
}

func CheckFileExists(filename string) (string, error) {

	checkExtension, err := CheckFileExtension(filename)
	if err != nil {
		return checkExtension, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return "Cannot access file: " + err.Error(), err
	}

	defer file.Close()

	return "", nil
}
