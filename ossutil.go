package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/ossutil/lib"
)

func main() {
	if err := lib.ParseAndRunCommand(); err != nil {
		fmt.Printf("Error: %s!\n", err)
		if strings.Contains(err.Error(), "ErrorCode=NoSuchUpload") {
			fmt.Printf("Will remove checkpoint dir '%s' automatically. Please try again.", lib.CheckpointDir)
			os.RemoveAll(lib.CheckpointDir)
		}
		os.Exit(1)
	}
	os.Exit(0)
}
