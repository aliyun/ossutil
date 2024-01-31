package lib

import (
	"encoding/xml"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
	"os"
)

func (s *OssutilCommandSuite) TestArchiveDirectReadHelpInfo(c *C) {
	// mkdir command test
	options := OptionMapType{}

	mkArgs := []string{"archive-direct-read"}
	_, err := cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

	mkArgs = []string{}
	_, err = cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) TestPutBucketArchiveDirectReadError(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	archiveFileName := "archive-direct-read" + randLowStr(12)

	strMethod := ""
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	// method is empty
	archiveArgs := []string{CloudURLToString(bucketName, ""), archiveFileName}
	_, err := cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	//method is error
	strMethod = "puttt"
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	// cloudurl is error
	strMethod = "put"
	archiveArgs = []string{"http://mybucket", archiveFileName}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	// local file is emtpy
	archiveArgs = []string{CloudURLToString(bucketName, ""), archiveFileName}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	//local file is not exist
	os.Remove(archiveFileName)
	archiveArgs = []string{CloudURLToString(bucketName, ""), archiveFileName}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	// local file is dir
	err = os.MkdirAll(archiveFileName, 0755)
	c.Assert(err, IsNil)
	archiveArgs = []string{CloudURLToString(bucketName, ""), archiveFileName}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)
	os.Remove(archiveFileName)

	//local file is empty
	s.createFile(archiveFileName, "", c)
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)
	os.Remove(archiveFileName)

	//local file is not xml file
	s.createFile(archiveFileName, "aaa", c)
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)
	os.Remove(archiveFileName)

	// StorageURLFromString error
	archiveArgs = []string{"oss:///1.jpg"}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	// bucketname is error
	archiveArgs = []string{"oss:///"}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	//missing parameter
	archiveArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	// bucketname not exist
	archiveArgs = []string{CloudURLToString("my-bucket", "")}
	_, err = cm.RunCommand("archive-direct-read", archiveArgs, options)
	c.Assert(err, NotNil)

	os.Remove(archiveFileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestPutBucketArchiveDirectRead(c *C) {
	archiveXml := `<?xml version="1.0" encoding="UTF-8"?>
<ArchiveDirectReadConfiguration>
  <Enabled>true</Enabled>
</ArchiveDirectReadConfiguration>`

	archiveConfigSrc := oss.PutBucketArchiveDirectRead{}
	err := xml.Unmarshal([]byte(archiveXml), &archiveConfigSrc)
	c.Assert(err, IsNil)

	archiveFileName := randLowStr(12)
	s.createFile(archiveFileName, archiveXml, c)

	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	strMethod := "put"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	testLogger.Println(endpoint)

	command := "archive-direct-read"
	archiveArgs := []string{CloudURLToString(bucketName, ""), archiveFileName}
	_, err = cm.RunCommand(command, archiveArgs, options)
	c.Assert(err, IsNil)

	archiveDownName := archiveFileName + "-down"
	strMethod = "get"
	options[OptionMethod] = &strMethod

	archiveArgs = []string{CloudURLToString(bucketName, ""), archiveDownName}
	_, err = cm.RunCommand(command, archiveArgs, options)
	c.Assert(err, IsNil)

	_, err = os.Stat(archiveDownName)
	c.Assert(err, IsNil)

	accessBody := s.readFile(archiveDownName, c)

	var out oss.GetBucketArchiveDirectReadResult
	err = xml.Unmarshal([]byte(accessBody), &out)
	c.Assert(err, IsNil)

	c.Assert(archiveConfigSrc.Enabled, Equals, out.Enabled)

	os.Remove(archiveFileName)
	os.Remove(archiveDownName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestGetBucketArchiveDirectReadConfirm(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	archiveXml := `<?xml version="1.0" encoding="UTF-8"?>
<ArchiveDirectReadConfiguration>
  <Enabled>true</Enabled>
</ArchiveDirectReadConfiguration>`

	archiveFileName := inputFileName + randLowStr(5)
	s.createFile(archiveFileName, archiveXml, c)

	strMethod := "put"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	command := "archive-direct-read"
	archiveArgs := []string{CloudURLToString(bucketName, ""), archiveFileName}
	_, err := cm.RunCommand(command, archiveArgs, options)
	c.Assert(err, IsNil)

	archiveDownName := archiveFileName + "-down"
	strMethod = "get"
	options[OptionMethod] = &strMethod

	archiveArgs = []string{CloudURLToString(bucketName, ""), archiveDownName}
	_, err = cm.RunCommand(command, archiveArgs, options)
	c.Assert(err, IsNil)

	archiveArgs = []string{CloudURLToString(bucketName, ""), archiveDownName}
	_, err = cm.RunCommand(command, archiveArgs, options)
	c.Assert(err, IsNil)

	archiveArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand(command, archiveArgs, options)
	c.Assert(err, IsNil)

	os.Remove(archiveFileName)
	os.Remove(archiveDownName)
	s.removeBucket(bucketName, true, c)
}
