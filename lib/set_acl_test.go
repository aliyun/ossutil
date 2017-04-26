package lib

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestSetBucketACL(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	// set acl
	for _, acl := range []string{"private", "public-read", "public-read-write"} {
		s.setBucketACL(bucketName, acl, c)
		s.getStat(bucketName, "", c)
	}
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestSetBucketErrorACL(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	for _, acl := range []string{"default", "def", "erracl", "私有"} {
		showElapse, err := s.rawSetBucketACL(bucketName, acl, false)
		c.Assert(err, NotNil)
		c.Assert(showElapse, Equals, false)

		showElapse, err = s.rawSetBucketACL(bucketName, acl, true)
		c.Assert(err, NotNil)
		c.Assert(showElapse, Equals, false)

		bucketStat := s.getStat(bucketName, "", c)
		c.Assert(bucketStat[StatACL], Equals, "private")
	}
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestSetNotExistBucketACL(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)

	// set acl and create bucket
	showElapse, err := s.rawSetBucketACL(bucketName, "public-read", true)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	s.removeBucket(bucketName, true, c)

	// invalid bucket name
	bucketName = "a"
	showElapse, err = s.rawSetBucketACL(bucketName, "public-read", true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	showElapse, err = s.rawGetStat(bucketName, "")
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestSetObjectACL(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	object := "TestSetObjectACL"

	// set acl to not exist object
	showElapse, err := s.rawSetObjectACL(bucketName, object, "default", false, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	object = "setacl-oldobject"
	s.putObject(bucketName, object, uploadFileName, c)

	//get object acl
	objectStat := s.getStat(bucketName, object, c)
	c.Assert(objectStat[StatACL], Equals, "default")

	object = "setacl-newobject"
	s.putObject(bucketName, object, uploadFileName, c)

	// set acl
	for _, acl := range []string{"default"} {
		s.setObjectACL(bucketName, object, acl, false, true, c)
		objectStat = s.getStat(bucketName, object, c)
		c.Assert(objectStat[StatACL], Equals, acl)
	}

	s.setObjectACL(bucketName, object, "private", false, true, c)

	// set error acl
	for _, acl := range []string{"public_read", "erracl", "私有", ""} {
		showElapse, err = s.rawSetObjectACL(bucketName, object, acl, false, false)
		c.Assert(showElapse, Equals, false)
		c.Assert(err, NotNil)
	}

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestBatchSetObjectACL(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	// put objects
	num := 2
	objectNames := []string{}
	for i := 0; i < num; i++ {
		object := fmt.Sprintf("TestBatchSetObjectACL_setacl%d", i)
		s.putObject(bucketName, object, uploadFileName, c)
		objectNames = append(objectNames, object)
	}

	// without --force option
	s.setObjectACL(bucketName, "", "public-read-write", true, false, c)

	s.setObjectACL(bucketName, "TestBatchSetObjectACL_setacl", "public-read", true, true, c)

	for _, object := range objectNames {
		objectStat := s.getStat(bucketName, object, c)
		c.Assert(objectStat[StatACL], Equals, "public-read")
	}

	showElapse, err := s.rawSetObjectACL(bucketName, "TestBatchSetObjectACL_setacl", "erracl", true, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestErrSetACL(c *C) {
	acl := "private"
	args := []string{"os://", acl}
	showElapse, err := s.rawSetACLWithArgs(args, false, false, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	args = []string{"oss://", acl}
	showElapse, err = s.rawSetACLWithArgs(args, false, false, true)

	// error set bucket acl
	bucketName := bucketNamePrefix + randLowStr(10)
	object := "testobject"
	args = []string{CloudURLToString(bucketName, object), acl}
	showElapse, err = s.rawSetACLWithArgs(args, false, true, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	args = []string{CloudURLToString(bucketName, ""), acl}
	showElapse, err = s.rawSetACLWithArgs(args, true, true, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	// miss acl
	args = []string{CloudURLToString(bucketName, "")}
	showElapse, err = s.rawSetACLWithArgs(args, false, true, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	args = []string{CloudURLToString(bucketName, object)}
	showElapse, err = s.rawSetACLWithArgs(args, false, true, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	args = []string{CloudURLToString(bucketName, object)}
	showElapse, err = s.rawSetACLWithArgs(args, true, false, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	// miss object
	args = []string{CloudURLToString(bucketName, ""), acl}
	showElapse, err = s.rawSetACLWithArgs(args, false, false, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	// bad prefix
	showElapse, err = s.rawSetObjectACL(bucketName, "/object", acl, true, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrBatchSetACL(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	// put objects
	num := 10
	objectNames := []string{}
	for i := 0; i < num; i++ {
		object := fmt.Sprintf("TestErrBatchSetACL_setacl:%d", i)
		s.putObject(bucketName, object, uploadFileName, c)
		objectNames = append(objectNames, object)
	}

	command := "set-acl"
	str := ""
	str1 := "abc"
	args := []string{CloudURLToString(bucketName, ""), "public-read-write"}
	routines := strconv.Itoa(Routines)
	ok := true
	options := OptionMapType{
		"endpoint":        &str1,
		"accessKeyID":     &str1,
		"accessKeySecret": &str1,
		"stsToken":        &str,
		"routines":        &routines,
		"recursive":       &ok,
		"force":           &ok,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	for _, object := range objectNames {
		objectStat := s.getStat(bucketName, object, c)
		c.Assert(objectStat[StatACL], Equals, "default")
	}

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestSetACLIDKey(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	cfile := randStr(10)
	data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucketName, "abc", bucketName, "abc")
	s.createFile(cfile, data, c)

	command := "set-acl"
	str := ""
	args := []string{CloudURLToString(bucketName, ""), "public-read"}
	routines := strconv.Itoa(Routines)
	ok := true
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
		"routines":        &routines,
		"bucket":          &ok,
		"force":           &ok,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)

	options = OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"stsToken":        &str,
		"configFile":      &cfile,
		"routines":        &routines,
		"bucket":          &ok,
		"force":           &ok,
	}
	showElapse, err = cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	os.Remove(cfile)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestSetACLURLEncoding(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	object := "^M特殊字符 加上空格 test"
	s.putObject(bucketName, object, uploadFileName, c)

	urlObject := url.QueryEscape(object)

	showElapse, err := s.rawSetObjectACL(bucketName, urlObject, "default", false, true)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	command := "set-acl"
	str := ""
	args := []string{CloudURLToString(bucketName, urlObject), "public-read"}
	routines := strconv.Itoa(Routines)
	ok := true
	encodingType := URLEncodingType
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"stsToken":        &str,
		"configFile":      &configFile,
		"routines":        &routines,
		"force":           &ok,
		"encodingType":    &encodingType,
	}
	showElapse, err = cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestSetACLErrArgs(c *C) {
	object := randStr(20)

	err := s.initSetACLWithArgs([]string{CloudURLToString("", object), "private"}, "", DefaultOutputDir)
	c.Assert(err, IsNil)
	err = setACLCommand.RunCommand()
	c.Assert(err, NotNil)

	err = s.initSetACLWithArgs([]string{CloudURLToString("", ""), "private"}, "", DefaultOutputDir)
	c.Assert(err, IsNil)
	err = setACLCommand.RunCommand()
	c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestBatchSetACLNotExistBucket(c *C) {
	// set acl notexist bucket
	err := s.initSetACLWithArgs([]string{CloudURLToString(bucketNamePrefix+randLowStr(10), ""), "private"}, "-rf", DefaultOutputDir)
	c.Assert(err, IsNil)
	err = setACLCommand.RunCommand()
	c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestBatchSetACLErrorContinue(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	// put object to archive bucket
	num := 2
	objectNames := []string{}
	for i := 0; i < num; i++ {
		object := fmt.Sprintf("设置权限object:%d%s", i, randStr(5))
		s.putObject(bucketName, object, uploadFileName, c)
		objectNames = append(objectNames, object)
	}

	// set acl prepare
	acl := oss.ACLPrivate

	err := s.initSetACLWithArgs([]string{CloudURLToString(bucketName, ""), string(acl)}, "-rf", DefaultOutputDir)
	c.Assert(err, IsNil)

	bucket, err := setACLCommand.command.ossBucket(bucketName)
	c.Assert(err, IsNil)
	c.Assert(bucket, NotNil)

	setACLCommand.monitor.init("Setted acl on")
	setACLCommand.saOption.ctnu = true

	// init reporter
	setACLCommand.saOption.reporter, err = GetReporter(setACLCommand.saOption.ctnu, DefaultOutputDir, commandLine)
	c.Assert(err, IsNil)

	defer setACLCommand.saOption.reporter.Clear()

	var routines int64
	routines = 3
	chObjects := make(chan string, ChannelBuf)
	chError := make(chan error, routines+1)
	chListError := make(chan error, 1)

	chObjects <- objectNames[0]
	chObjects <- "notexistobject" + randStr(3)
	chObjects <- objectNames[1]
	chListError <- nil
	close(chObjects)

	for i := 0; int64(i) < routines; i++ {
		setACLCommand.setObjectACLConsumer(bucket, acl, chObjects, chError)
	}

	err = setACLCommand.waitRoutinueComplete(chError, chListError, routines)
	c.Assert(err, IsNil)

	str := setACLCommand.monitor.progressBar(false, normalExit)
	c.Assert(str, Equals, "")
	str = setACLCommand.monitor.progressBar(false, errExit)
	c.Assert(str, Equals, "")
	str = setACLCommand.monitor.progressBar(true, normalExit)
	c.Assert(str, Equals, "")
	str = setACLCommand.monitor.progressBar(true, errExit)
	c.Assert(str, Equals, "")

	snap := setACLCommand.monitor.getSnapshot()
	c.Assert(snap.okNum, Equals, int64(2))
	c.Assert(snap.errNum, Equals, int64(1))
	c.Assert(snap.dealNum, Equals, int64(3))

	setACLCommand.monitor.seekAheadEnd = true
	setACLCommand.monitor.seekAheadError = nil
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(normalExit))
	c.Assert(strings.Contains(str, "finishwitherror:"), Equals, true)
	c.Assert(strings.Contains(str, "succeed:"), Equals, false)
	c.Assert(strings.Contains(str, "error"), Equals, true)
	setACLCommand.monitor.seekAheadEnd = false
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(normalExit))
	c.Assert(strings.Contains(str, "finishwitherror:"), Equals, true)
	c.Assert(strings.Contains(str, "succeed:"), Equals, false)
	c.Assert(strings.Contains(str, "error"), Equals, true)

	setACLCommand.monitor.seekAheadEnd = true
	setACLCommand.monitor.seekAheadError = nil
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(errExit))
	c.Assert(strings.Contains(str, "when error happens."), Equals, true)
	c.Assert(strings.Contains(str, "total"), Equals, true)
	setACLCommand.monitor.seekAheadEnd = false
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(errExit))
	c.Assert(strings.Contains(str, "when error happens."), Equals, true)
	c.Assert(strings.Contains(str, "scanned"), Equals, true)

	for _, object := range objectNames {
		objectStat := s.getStat(bucketName, object, c)
		c.Assert(objectStat[StatACL], Equals, string(acl))
	}

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestBatchSetACLErrorBreak(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucketWithStorageClass(bucketName, StorageArchive, c)

	// put object to archive bucket
	num := 2
	objectNames := []string{}
	for i := 0; i < num; i++ {
		object := fmt.Sprintf("设置权限object:%d%s", i, randStr(5))
		s.putObject(bucketName, object, uploadFileName, c)
		objectNames = append(objectNames, object)
	}

	// prepare
	acl := oss.ACLPrivate

	err := s.initSetACLWithArgs([]string{CloudURLToString(bucketName, ""), string(acl)}, "-rf", DefaultOutputDir)
	c.Assert(err, IsNil)

	// make error bucket with error id
	bucket := s.getErrorOSSBucket(bucketName, c)
	c.Assert(bucket, NotNil)

	setACLCommand.monitor.init("Setted acl on")
	setACLCommand.saOption.ctnu = true

	// init reporter
	setACLCommand.saOption.reporter, err = GetReporter(setACLCommand.saOption.ctnu, DefaultOutputDir, commandLine)
	c.Assert(err, IsNil)

	defer setACLCommand.saOption.reporter.Clear()

	var routines int64
	routines = 3
	chObjects := make(chan string, ChannelBuf)
	chError := make(chan error, routines+1)
	chListError := make(chan error, 1)

	chObjects <- objectNames[0]
	chObjects <- objectNames[1]
	chListError <- nil
	close(chObjects)

	for i := 0; int64(i) < routines; i++ {
		setACLCommand.setObjectACLConsumer(bucket, acl, chObjects, chError)
	}

	err = setACLCommand.waitRoutinueComplete(chError, chListError, routines)
	c.Assert(err, NotNil)

	str := setACLCommand.monitor.progressBar(false, normalExit)
	c.Assert(str, Equals, "")
	str = setACLCommand.monitor.progressBar(false, errExit)
	c.Assert(str, Equals, "")
	str = setACLCommand.monitor.progressBar(true, normalExit)
	c.Assert(str, Equals, "")
	str = setACLCommand.monitor.progressBar(true, errExit)
	c.Assert(str, Equals, "")

	snap := setACLCommand.monitor.getSnapshot()
	c.Assert(snap.okNum, Equals, int64(0))
	c.Assert(snap.errNum, Equals, int64(2))
	c.Assert(snap.dealNum, Equals, int64(2))

	setACLCommand.monitor.seekAheadEnd = true
	setACLCommand.monitor.seekAheadError = nil
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(normalExit))
	c.Assert(strings.Contains(str, "finishwitherror:"), Equals, true)
	c.Assert(strings.Contains(str, "succeed:"), Equals, false)
	c.Assert(strings.Contains(str, "error"), Equals, true)
	setACLCommand.monitor.seekAheadEnd = false
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(normalExit))
	c.Assert(strings.Contains(str, "finishwitherror:"), Equals, true)
	c.Assert(strings.Contains(str, "succeed:"), Equals, false)
	c.Assert(strings.Contains(str, "error"), Equals, true)

	setACLCommand.monitor.seekAheadEnd = true
	setACLCommand.monitor.seekAheadError = nil
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(errExit))
	c.Assert(strings.Contains(str, "when error happens."), Equals, true)
	c.Assert(strings.Contains(str, "total"), Equals, true)
	setACLCommand.monitor.seekAheadEnd = false
	str = strings.ToLower(setACLCommand.monitor.getFinishBar(errExit))
	c.Assert(strings.Contains(str, "when error happens."), Equals, true)
	c.Assert(strings.Contains(str, "scanned"), Equals, true)

	for _, object := range objectNames {
		objectStat := s.getStat(bucketName, object, c)
		c.Assert(objectStat[StatACL], Equals, "default")
	}

	s.removeBucket(bucketName, true, c)
}
