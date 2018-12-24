package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	nologLevel = iota
	errorLevel
	normalLevel
	infoLevel
	debugLevel
)

type UtilLog struct {
}

func (utilLog UtilLog) Write(buf []byte) (n int, err error) {
	_, err = LogDebug(string(buf))
	if err == nil {
		return len(buf), nil
	} else {
		return 0, err
	}
}

var utilLog = &UtilLog{}
var maxLogSize int64 = 1024 * 1024 * 50
var logName = "ossutil.log"
var logLevel = nologLevel
var lock = &sync.Mutex{}

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

	if level == errorLevel {
		n, err = io.WriteString(f, fmt.Sprintf("[%s][error]", nowS))
	} else if level == normalLevel {
		n, err = io.WriteString(f, fmt.Sprintf("[%s][normal]", nowS))
	} else if level == infoLevel {
		n, err = io.WriteString(f, fmt.Sprintf("[%s][info]", nowS))
	} else if level == debugLevel {
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

	if logLevel < errorLevel {
		return
	}
	return writeLog(errorLevel, format, a...)
}

func LogNormal(format string, a ...interface{}) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if logLevel < normalLevel {
		return
	}
	return writeLog(normalLevel, format, a...)
}

func LogInfo(format string, a ...interface{}) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if logLevel < infoLevel {
		return
	}
	return writeLog(infoLevel, format, a...)
}

func LogDebug(format string, a ...interface{}) (n int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if logLevel < debugLevel {
		return
	}
	return writeLog(debugLevel, format, a...)
}
