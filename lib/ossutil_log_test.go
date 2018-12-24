package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "gopkg.in/check.v1"
)

type OssUtilLogSuite struct {
	testLogSize  int64
	testLogName  string
	testLogLevel int
}

var _ = Suite(&OssUtilLogSuite{})

// Run once when the suite starts running
func (s *OssUtilLogSuite) SetUpSuite(c *C) {
}

// Run before each test or benchmark starts running
func (s *OssUtilLogSuite) TearDownSuite(c *C) {

}

// Run after each test or benchmark runs
func (s *OssUtilLogSuite) SetUpTest(c *C) {
	s.testLogSize = maxLogSize
	s.testLogName = logName
	s.testLogLevel = logLevel

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = ""
	}
	absLogName := dir + string(os.PathSeparator) + logName
	os.Remove(absLogName)
}

// Run once after all tests or benchmarks have finished running
func (s *OssUtilLogSuite) TearDownTest(c *C) {
	maxLogSize = s.testLogSize
	logName = s.testLogName
	logLevel = s.testLogLevel

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = ""
	}
	absLogName := dir + string(os.PathSeparator) + logName
	os.Remove(absLogName)
}

// test "config"
func (s *OssUtilLogSuite) TestLogLevel(c *C) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = ""
	}
	absLogName := dir + string(os.PathSeparator) + logName

	// nologLevel
	logLevel = nologLevel
	errorContext := "i am error log.\n"
	LogError(errorContext)
	LogNormal(errorContext)
	LogInfo(errorContext)
	LogDebug(errorContext)

	contents, err := ioutil.ReadFile(absLogName)
	LogContent := string(contents)
	c.Assert(strings.Contains(LogContent, "[error]"+errorContext), Equals, false)
	c.Assert(strings.Contains(LogContent, "[normal]"+errorContext), Equals, false)
	c.Assert(strings.Contains(LogContent, "[info]"+errorContext), Equals, false)
	c.Assert(strings.Contains(LogContent, "[debug]"+errorContext), Equals, false)
	os.Remove(absLogName)

	// errorLevel
	logLevel = errorLevel
	LogError(errorContext)
	LogNormal(errorContext)
	LogInfo(errorContext)
	LogDebug(errorContext)

	contents, err = ioutil.ReadFile(absLogName)
	LogContent = string(contents)
	c.Assert(strings.Contains(LogContent, "[error]"+errorContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[normal]"+errorContext), Equals, false)
	c.Assert(strings.Contains(LogContent, "[info]"+errorContext), Equals, false)
	c.Assert(strings.Contains(LogContent, "[debug]"+errorContext), Equals, false)
	os.Remove(absLogName)

	// normalLevel
	logLevel = normalLevel
	normalContext := "i am normal log.\n"
	LogError(normalContext)
	LogNormal(normalContext)
	LogInfo(normalContext)
	LogDebug(normalContext)

	contents, err = ioutil.ReadFile(absLogName)
	LogContent = string(contents)
	c.Assert(strings.Contains(LogContent, "[error]"+normalContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[normal]"+normalContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[info]"+normalContext), Equals, false)
	c.Assert(strings.Contains(LogContent, "[debug]"+normalContext), Equals, false)
	os.Remove(absLogName)

	// infolevel
	logLevel = infoLevel
	infoContext := "i am info log.\n"
	LogError(infoContext)
	LogNormal(infoContext)
	LogInfo(infoContext)
	LogDebug(infoContext)

	contents, err = ioutil.ReadFile(absLogName)
	LogContent = string(contents)
	c.Assert(strings.Contains(LogContent, "[error]"+infoContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[normal]"+infoContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[info]"+infoContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[debug]"+infoContext), Equals, false)
	os.Remove(absLogName)

	// debuglevel
	logLevel = debugLevel
	debugContext := "i am debug log.\n"
	LogError(debugContext)
	LogNormal(debugContext)
	LogInfo(debugContext)
	LogDebug(debugContext)

	contents, err = ioutil.ReadFile(absLogName)
	LogContent = string(contents)
	c.Assert(strings.Contains(LogContent, "[error]"+debugContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[normal]"+debugContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[info]"+debugContext), Equals, true)
	c.Assert(strings.Contains(LogContent, "[debug]"+debugContext), Equals, true)
	os.Remove(absLogName)
}

func (s *OssUtilLogSuite) TestLogFileSwitch(c *C) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		dir = ""
	}
	absLogName := dir + string(os.PathSeparator) + logName
	absLogNameBak := absLogName + ".bak"

	os.Remove(absLogName)
	os.Remove(absLogNameBak)

	maxLogSize = 5 * 1024

	// debuglevel
	logLevel = debugLevel
	debugContext := "i am debug log.\n"

	for i := 1; i < 1024; i++ {
		LogDebug(debugContext)
	}

	f, err := os.Stat(absLogNameBak)
	bLarge := (f.Size() > maxLogSize)
	c.Assert(bLarge, Equals, true)

	f, err = os.Stat(absLogName)
	c.Assert(err, IsNil)

	contents, err := ioutil.ReadFile(absLogName)
	LogContent := string(contents)
	c.Assert(strings.Contains(LogContent, "[debug]"+debugContext), Equals, true)

	os.Remove(absLogName)
	os.Remove(absLogNameBak)
}
