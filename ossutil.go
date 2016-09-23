package main

import (
    "os"
    "fmt"

	"github.com/aliyun/ossutil/lib"
)

func main() {
    if err := lib.ParseAndRunCommand(); err != nil {
        fmt.Printf("Error: %s!\n", err)
        os.Exit(1)
    }
    os.Exit(0)
}


