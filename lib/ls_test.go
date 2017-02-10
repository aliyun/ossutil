package lib

import (
	"fmt"
	"os"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
)

// test list buckets
func (s *OssutilCommandSuite) TestListLoadConfig(c *C) {
	command := "ls"
	var args []string
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	args = []string{"oss://"}
	options = OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}
	showElapse, err = cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestListNotExistConfigFile(c *C) {
	command := "ls"
	var args []string
	str := ""
	cfile := "notexistfile"
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestListErrConfigFile(c *C) {
	cfile := randStr(10)
	s.createFile(cfile, content, c)

	command := "ls"
	var args []string
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListConfigFile(c *C) {
	cfile := randStr(10)
	data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\nretryTimes=%s", endpoint, accessKeyID, accessKeySecret, "errretry")
	s.createFile(cfile, data, c)

	command := "ls"
	var args []string
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListWithBucketEndpoint(c *C) {
	bucketName := bucketNameExist

	cfile := randStr(10)
	data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s", "abc", accessKeyID, accessKeySecret, bucketName, endpoint)
	s.createFile(cfile, data, c)

	command := "ls"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListWithBucketCname(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	cfile := randStr(10)
	data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s\n[Bucket-Cname]\n%s=%s", "abc", accessKeyID, accessKeySecret, bucketName, "abc", bucketName, bucketName+"."+endpoint)
	s.createFile(cfile, data, c)

	command := "ls"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	os.Remove(cfile)
	s.removeBucket(bucketName, true, c)
}

// list objects with not exist bucket
func (s *OssutilCommandSuite) TestListObjectsBucketNotExist(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	command := "ls"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)
}

// list objects
func (s *OssutilCommandSuite) TestListObjects(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	// put objects
	num := 3
	for i := 0; i < num; i++ {
		object := fmt.Sprintf("lstest:#%d", i)
		s.putObject(bucketName, object, uploadFileName, c)
	}

	object := "another_object"
	s.putObject(bucketName, object, uploadFileName, c)

	objectStat := s.getStat(bucketName, object, c)
	c.Assert(objectStat[StatACL], Equals, "default")
	c.Assert(len(objectStat["Etag"]), Equals, 32)
	c.Assert(objectStat["Last-Modified"] != "", Equals, true)
	c.Assert(objectStat[StatOwner] != "", Equals, true)

	//put directories
	num1 := 2
	for i := 0; i < num1; i++ {
		object := fmt.Sprintf("lstest:#%d/", i)
		s.putObject(bucketName, object, uploadFileName, c)

		object = fmt.Sprintf("lstest:#%d/%d/", i, i)
		s.putObject(bucketName, object, uploadFileName, c)
	}

	// "ls oss://bucket -s"
	//objects := s.listObjects(bucketName, "", true, false, false, false, c)
	//c.Assert(len(objects), Equals, num + 2*num1 + 1)

	// "ls oss://bucket/prefix -s"
	objects := s.listObjects(bucketName, "lstest:", "ls -s", c)
	c.Assert(len(objects), Equals, num+2*num1)

	// "ls oss://bucket/prefix"
	objects = s.listObjects(bucketName, "lstest:#", "ls - ", c)
	c.Assert(len(objects), Equals, num+2*num1)

	// "ls oss://bucket/prefix -d"
	objects = s.listObjects(bucketName, "lstest:#", "ls -d", c)
	c.Assert(len(objects), Equals, num+num1)

	objects = s.listObjects(bucketName, "lstest:#1/", "ls -d", c)
	c.Assert(len(objects), Equals, 2)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestErrList(c *C) {
	showElapse, err := s.rawList([]string{"../"}, "ls -s")
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	// not exist bucket
	bucketName := bucketNamePrefix + randLowStr(10)
	showElapse, err = s.rawList([]string{CloudURLToString(bucketName, "")}, "ls -d")
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)

	// list buckets with -d
	showElapse, err = s.rawList([]string{"oss://"}, "ls -d")
	c.Assert(err, NotNil)
	c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestListIDKey(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	cfile := randStr(10)
	data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucketName, "abc", bucketName, "abc")
	s.createFile(cfile, data, c)

	command := "ls"
	str := ""
	args := []string{"oss://"}
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)

	options = OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err = cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	os.Remove(cfile)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestListBucketIDKey(c *C) {
	bucketName := bucketNameExist

	cfile := randStr(10)
	data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucketName, "abc", bucketName, "abc")
	s.createFile(cfile, data, c)

	command := "ls"
	str := ""
	args := []string{CloudURLToString(bucketName, "")}
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err := cm.RunCommand(command, args, options)
	c.Assert(err, NotNil)

	options = OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"stsToken":        &str,
		"configFile":      &cfile,
	}
	showElapse, err = cm.RunCommand(command, args, options)
	c.Assert(err, IsNil)
	c.Assert(showElapse, Equals, true)

	os.Remove(cfile)
}

// list multipart
func (s *OssutilCommandSuite) TestListMultipartUploads(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)
	// "rm -arf oss://bucket/"
	command := "rm"
	args := []string{CloudURLToString(bucketName, "")}
	str := ""
	ok := true
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
		"recursive":       &ok,
		"force":           &ok,
		"allType":         &ok,
	}
	cm.RunCommand(command, args, options)

	object := "TestMultipartObjectLs"
	s.putObject(bucketName, object, uploadFileName, c)

	// list object
	objects := s.listObjects(bucketName, object, "ls - ", c)
	c.Assert(len(objects), Equals, 1)
	c.Assert(objects[0], Equals, object)

	bucket, err := copyCommand.command.ossBucket(bucketName)

	lmr_origin, e := bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)

	for i := 0; i < 20; i++ {
		_, err = bucket.InitiateMultipartUpload(object)
		c.Assert(err, IsNil)
	}

	lmr, e := bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
	c.Assert(len(lmr.Uploads), Equals, 20+len(lmr_origin.Uploads))

	// list multipart: ls oss://bucket/object
	objects = s.listObjects(bucketName, object, "ls - ", c)
	c.Assert(len(objects), Equals, 1)
	c.Assert(objects[0], Equals, object)

	// list multipart: ls -m oss://bucket/object
	objects = s.listObjects(bucketName, object, "ls -m", c)
	c.Assert(len(objects), Equals, 20)

	// list all type object: ls -a oss://bucket/object
	objects = s.listObjects(bucketName, object, "ls -a", c)
	c.Assert(len(objects), Equals, 21)

	// list multipart: ls -am oss://bucket/object
	objects = s.listObjects(bucketName, object, "ls -am", c)
	c.Assert(len(objects), Equals, 21)

	// list multipart: ls -ms oss://bucket/object
	objects = s.listObjects(bucketName, object, "ls -ms", c)
	c.Assert(len(objects), Equals, 20)

	// list multipart: ls -as oss://bucket/object
	objects = s.listObjects(bucketName, object, "ls -as", c)
	c.Assert(len(objects), Equals, 21)

	lmr, e = bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
	c.Assert(len(lmr.Uploads), Equals, 20+len(lmr_origin.Uploads))

	s.removeBucket(bucketName, true, c)
}
