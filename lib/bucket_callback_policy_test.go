package lib

import (
	. "gopkg.in/check.v1"
	"os"
	"strings"
)

func (s *OssutilCommandSuite) TestBucketCallbackPolicySuccess(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	// put success
	strMethod := "put"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}
	fileName := "test-ossutil-callback-" + randLowStr(5)
	callbackContent := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<BucketCallbackPolicy>\n  <PolicyItem>\n    <PolicyName>first</PolicyName>\n    <Callback>eyJjYWxsYmFja1VybCI6Imh0dHA6Ly93d3cuYWxpeXVuY3MuY29tIiwgImNhbGxiYWNrQm9keSI6ImJ1Y2tldD0ke2J1Y2tldH0mb2JqZWN0PSR7b2JqZWN0fSJ9</Callback>\n    <CallbackVar></CallbackVar>\n  </PolicyItem>\n  <PolicyItem>\n    <PolicyName>second</PolicyName>\n    <Callback>eyJjYWxsYmFja1VybCI6Imh0dHA6Ly93d3cuYWxpeXVuLmNvbSIsICJjYWxsYmFja0JvZHkiOiJidWNrZXQ9JHtidWNrZXR9Jm9iamVjdD0ke29iamVjdH0mYT0ke3g6YX0mYj0ke3g6Yn0ifQ==</Callback>\n    <CallbackVar>eyJ4OmEiOiJhIiwgIng6YiI6ImIifQ==</CallbackVar>\n  </PolicyItem>\n</BucketCallbackPolicy>\n"
	s.createFile(fileName, callbackContent, c)

	callbackArgs := []string{CloudURLToString(bucketName, ""), fileName}
	_, err := cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, IsNil)

	// get success
	// output to file
	outputFile := "test-file-" + randLowStr(5)
	testResultFile, err = os.OpenFile(outputFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	c.Assert(err, IsNil)

	oldStdout := os.Stdout
	os.Stdout = testResultFile

	strMethod = "get"
	options[OptionMethod] = &strMethod
	callbackArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, IsNil)
	testResultFile.Close()
	os.Stdout = oldStdout
	outBody := s.readFile(outputFile, c)
	c.Assert(strings.ReplaceAll(outBody, "\n", ""), Equals, strings.ReplaceAll(callbackContent, "\n", ""))

	strMethod = "delete"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, IsNil)

	os.Remove(outputFile)
	os.Remove(fileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestBucketCallbackPolicyError(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	// method is empty
	strMethod := "put"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	callbackArgs := []string{CloudURLToString(bucketName, "")}
	_, err := cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, NotNil)

	// cloud url is error
	strMethod = "put"
	callbackArgs = []string{"http://test-bucket"}
	_, err = cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, NotNil)

	strMethod = "putt"
	callbackArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, NotNil)

	strMethod = "deletee"
	callbackArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, NotNil)

	// no callback policy
	strMethod = "get"
	callbackArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, NotNil)

	// bucket name is empty
	callbackArgs = []string{"oss://"}
	_, err = cm.RunCommand("callback-policy", callbackArgs, options)
	c.Assert(err, NotNil)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestBucketCallbackPolicyHelpInfo(c *C) {
	options := OptionMapType{}

	mkArgs := []string{"callback-policy"}
	_, err := cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

	mkArgs = []string{}
	_, err = cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)
}
