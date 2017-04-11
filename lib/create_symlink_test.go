package lib

import (
	"fmt"
	"net/url"

	. "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestCreateSymlink(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	notExistBucketName := bucketNamePrefix + randLowStr(10)

	symObject := bucketNamePrefix + randStr(5) + "符号链接"
	targetObject := bucketNamePrefix + randStr(5) + "符号链接目标"
	targetObject1 := bucketNamePrefix + randStr(5) + "target"

	// put symlink to not exist bucket
	cmdline := fmt.Sprintf("%s %s", CloudURLToString(bucketName, symObject), targetObject)
	err := s.initCreateSymlink(cmdline)
	c.Assert(err, IsNil)
	err = createSymlinkCommand.RunCommand()
	c.Assert(err, NotNil)

	s.putBucket(bucketName, c)

	data := "中文内容"
	s.createFile(uploadFileName, data, c)
	s.putObject(bucketName, targetObject, uploadFileName, c)

	data1 := "english"
	s.createFile(uploadFileName, data1, c)
	s.putObject(bucketName, targetObject1, uploadFileName, c)

	// put symlink to different bucket
	cmdline = fmt.Sprintf("%s %s", CloudURLToString(bucketName, symObject), CloudURLToString(notExistBucketName, targetObject))
	err = s.initCreateSymlink(cmdline)
	c.Assert(err, IsNil)
	err = createSymlinkCommand.RunCommand()
	c.Assert(err, NotNil)

	cmdline = fmt.Sprintf("%s %s", CloudURLToString(bucketName, symObject), targetObject)
	err = s.initCreateSymlink(cmdline)
	c.Assert(err, IsNil)
	err = createSymlinkCommand.RunCommand()
	c.Assert(err, IsNil)

	s.getObject(bucketName, symObject, downloadFileName, c)
	str := s.readFile(downloadFileName, c)
	c.Assert(str, Equals, data)

	// put symlink again
	cmdline = fmt.Sprintf("%s %s", CloudURLToString(bucketName, symObject), targetObject1)
	err = s.initCreateSymlink(cmdline)
	c.Assert(err, IsNil)
	err = createSymlinkCommand.RunCommand()
	c.Assert(err, IsNil)

	s.getObject(bucketName, symObject, downloadFileName, c)
	str = s.readFile(downloadFileName, c)
	c.Assert(str, Equals, data1)

	// error put symlink
	cmdline = fmt.Sprintf("%s %s", CloudURLToString(bucketName, symObject), targetObject1, "abc")
	err = s.initCreateSymlink(cmdline)
	c.Assert(err, NotNil)

	// put symlink with urlencoding
	urlTarget := url.QueryEscape(targetObject)
	c.Assert(urlTarget != targetObject, Equals, true)

	cmdline = fmt.Sprintf("%s %s --encoding-type url", CloudURLToString(bucketName, url.QueryEscape(symObject)), urlTarget)
	err = s.initCreateSymlink(cmdline)
	c.Assert(err, IsNil)
	err = createSymlinkCommand.RunCommand()
	c.Assert(err, IsNil)

	s.getObject(bucketName, symObject, downloadFileName, c)
	str = s.readFile(downloadFileName, c)
	c.Assert(str, Equals, data)
}
