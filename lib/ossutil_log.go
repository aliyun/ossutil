package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var utilLog = &UtilLog{}
var maxLogSize int64 = 1024 * 1024 * 50
var logName = "ossutil.log"
var logLevel = oss.LogOff
var lock = &sync.Mutex{}

type UtilLog struct {
}

func (utilLog UtilLog) Write(buf []byte) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = ""
	}
	absLogName := dir + string(os.PathSeparator) + logName
	absLogNameBak := absLogName + ".bak"

	statInfo, err := os.Stat(absLogName)
	if err == nil && statInfo.Size() >= (maxLogSize) {
		os.Rename(absLogName, absLogNameBak)
	}

	f, err := os.OpenFile(absLogName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Printf("error,open log file %s error.\n", absLogName)
		return 0, nil
	}
	defer f.Close()
	n, err = io.WriteString(f, fmt.Sprintf("%s", buf))
	return
}

func writeLog(level int, format string, a ...interface{}) (n int, err error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = ""
	}
	absLogName := dir + string(os.PathSeparator) + logName
	absLogNameBak := absLogName + ".bak"

	statInfo, err := os.Stat(absLogName)
	if err == nil && statInfo.Size() >= (maxLogSize) {
		os.Rename(absLogName, absLogNameBak)
	}

	f, err := os.OpenFile(absLogName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Printf("error,open log file %s error.\n", absLogName)
		return 0, nil
	}
	defer f.Close()

	nowT := time.Now()
	nowS := time.Unix(nowT.Unix(), 0).Format("2006-01-02 15:04:05")

	if level == oss.ErrorLevel {
		n, err = io.WriteString(f, fmt.Sprintf("[%s][error]", nowS))
	} else if level == oss.NormalLevel {
		n, err = io.WriteString(f, fmt.Sprintf("[%s][normal]", nowS))
	} else if level == oss.InfoLevel {
		n, err = io.WriteString(f, fmt.Sprintf("[%s][info]", nowS))
	} else if level == oss.DebugLevel {
		n, err = io.WriteString(f, fmt.Sprintf("[%s][debug]", nowS))
	}

	if err != nil {
		return n, err
	}

	n1, err1 := io.WriteString(f, fmt.Sprintf(format, a...))
	if err1 != nil {
		err = err1
	} else {
		n += n1
	}
	return n, err
}

func InitLogger(level int, name string, len int64) {
	logLevel = level
	logName = name
	maxLogSize = len
}

func LogError(format string, a ...interface{}) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if logLevel < oss.ErrorLevel {
		return
	}
	return writeLog(oss.ErrorLevel, format, a...)
}

func LogNormal(format string, a ...interface{}) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if logLevel < oss.NormalLevel {
		return
	}
	return writeLog(oss.NormalLevel, format, a...)
}

func LogInfo(format string, a ...interface{}) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if logLevel < oss.InfoLevel {
		return
	}
	return writeLog(oss.InfoLevel, format, a...)
}

func LogDebug(format string, a ...interface{}) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if logLevel < oss.DebugLevel {
		return
	}
	return writeLog(oss.DebugLevel, format, a...)
}
