package lib

import (
	"fmt"
	"os"
	"strconv"
    "net/url"

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
