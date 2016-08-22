package lib

import (
	"fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// Output print input string to stdout and add '\n'
func Output(str string) {
	fmt.Println(str)
}

// FindPos find the elem position in a string array
func FindPos(elem string, elemArray []string) int {
	for p, v := range elemArray {
		if v == elem {
			return p
		}
	}
	return -1
}

// FindPosCaseInsen find the elem position in a string array, ignore case
func FindPosCaseInsen(elem string, elemArray []string) int {
	for p, v := range elemArray {
		if strings.ToLower(v) == strings.ToLower(elem) {
			return p
		}
	}
	return -1
}

func getBinaryPath() (string, string) {
    filePath, _ := exec.LookPath(os.Args[0])
    if path, err := os.Readlink(filePath); err == nil {
        filePath = path
    }

    fileName := filepath.Base(filePath)
    renameFilePath := ".temp_" + fileName
    return filePath, renameFilePath
}
