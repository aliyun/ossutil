package lib

import (
	"io/ioutil"
	"os"
	"strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestCatObjectSuccess(c *C) {
	// create client and bucket
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	c.Assert(err, IsNil)

	bucketName := bucketNamePrefix + randLowStr(5)
	err = client.CreateBucket(bucketName)
	c.Assert(err, IsNil)

	bucket, err := client.Bucket(bucketName)
	c.Assert(err, IsNil)

	// put object
	//first:upload a object
	textBuffer := randStr(1024)
	objectName := randStr(10)
	err = bucket.PutObject(objectName, strings.NewReader(textBuffer))
	c.Assert(err, IsNil)

	// begin cat
	var str string
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}

	// output to file
	fileName := "test-file-" + randLowStr(5)
	testResultFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)

	oldStdout := os.Stdout
	os.Stdout = testResultFile

	catArgs := []string{CloudURLToString(bucketName, objectName)}
	_, err = cm.RunCommand("cat", catArgs, options)
	c.Assert(err, IsNil)
	testResultFile.Close()
	os.Stdout = oldStdout

	// check file content
	catBody := s.readFile(fileName, c)
	c.Assert(strings.Contains(catBody, textBuffer), Equals, true)

	//remove file
	os.Remove(fileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestCatObjectError(c *C) {
	// create client and bucket
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	c.Assert(err, IsNil)

	bucketName := bucketNamePrefix + randLowStr(5)
	err = client.CreateBucket(bucketName)
	c.Assert(err, IsNil)

	// begin cat
	var str string
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}

	// object is empty
	catArgs := []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("cat", catArgs, options)
	c.Assert(err, NotNil)

	// object is not exist
	catArgs = []string{CloudURLToString(bucketName, randLowStr(5))}
	_, err = cm.RunCommand("cat", catArgs, options)
	c.Assert(err, NotNil)

	// cloud url is error
	catArgs = []string{"oss:///1.jpg"}
	_, err = cm.RunCommand("cat", catArgs, options)
	c.Assert(err, NotNil)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestCatObjecEndpointEmptyError(c *C) {
	// create client and bucket
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	c.Assert(err, IsNil)

	bucketName := bucketNamePrefix + randLowStr(5)
	err = client.CreateBucket(bucketName)
	c.Assert(err, IsNil)

	// begin cat
	var str string
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"stsToken":        &str,
		"configFile":      &configFile,
	}

	// ossclient error
	//set endpoint emtpy
	oldConfigStr, err := ioutil.ReadFile(configFile)
	c.Assert(err, IsNil)

	fd, _ := os.OpenFile(configFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	configStr := "[Credentials]" + "\n" + "language=CH" + "\n" + "accessKeyID=123" + "\n" + "accessKeySecret=456" + "\n" + "endpoint="
	fd.WriteString(configStr)
	fd.Close()

	catArgs := []string{CloudURLToString(bucketName, randLowStr(5))}
	_, err = cm.RunCommand("cat", catArgs, options)
	c.Assert(err, NotNil)

	err = ioutil.WriteFile(configFile, []byte(oldConfigStr), 0664)
	c.Assert(err, IsNil)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestCatObjectHelpInfo(c *C) {
	// mkdir command test
	options := OptionMapType{}

	mkArgs := []string{"cat"}
	_, err := cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

	mkArgs = []string{}
	_, err = cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)
}
