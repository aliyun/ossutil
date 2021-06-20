package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/ossutil/lib"
)

func main() {
	if err := lib.ParseAndRunCommand(); err != nil {
		fmt.Printf("Error: %s\n", err)
		if strings.Contains(err.Error(), "ErrorCode=NoSuchUpload") {
			fmt.Printf("Will remove checkpoint dir '%s' automatically. Please try again.\n", lib.CheckpointDir)
			os.RemoveAll(lib.CheckpointDir)
		}
		if strings.Contains(err.Error(), ": EOF,") {
			fmt.Printf("Connection has been closed by remote peer. Please check the network. If you download/upload large file, You can reduce concurrency with the --parallel option and reduce part-size with --part-size (it must greater than the file size divided by 10000. By default, it will retry 10 times when failed, you can increse the retry times with --retry-times option.).\n")
		}
		os.Exit(1)
	}
	os.Exit(0)
}
