package lib

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
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
		if strings.EqualFold(v, elem) {
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

type sysInfo struct {
	name    string // 操作系统名称windows/Linux
	release string // 操作系统版本 2.6.32-220.23.2.ali1089.el5.x86_64等
	machine string // 机器类型amd64/x86_64
}

// Get　system info
// 获取操作系统信息、机器类型
func getSysInfo() sysInfo {
	name := runtime.GOOS
	release := "-"
	machine := runtime.GOARCH
	if out, err := exec.Command("uname", "-s").CombinedOutput(); err == nil {
		name = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command("uname", "-r").CombinedOutput(); err == nil {
		release = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command("uname", "-m").CombinedOutput(); err == nil {
		machine = string(bytes.TrimSpace(out))
	}
	return sysInfo{name: name, release: release, machine: machine}
}

func getUserAgent() string {
	sys := getSysInfo()
	return fmt.Sprintf("aliyun-sdk-go/%s (%s/%s/%s;%s)/%s-%s", oss.Version, sys.name, sys.release, sys.machine, runtime.Version(), Package, Version)
}

func utcToLocalTime(utc time.Time) time.Time {
	return utc.In(time.Local)
}

func max(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func getSizeString(size int64) string {
	prefix := ""
	str := fmt.Sprintf("%d", size)
	if size < 0 {
		prefix = "-"
		str = str[1:]
	}
	len := len(str)
	strList := []string{}
	i := len % 3
	if i != 0 {
		strList = append(strList, str[0:i])
	}
	for ; i < len; i += 3 {
		strList = append(strList, str[i:i+3])
	}
	return fmt.Sprintf("%s%s", prefix, strings.Join(strList, ","))
}
