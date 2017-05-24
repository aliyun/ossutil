package lib

import (
	"fmt"
	"net/url"
	"strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestSignurlGet(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	data := "签名url"
	s.createFile(uploadFileName, data, c)

	object := randStr(10)
	s.putObject(bucketName, object, uploadFileName, c)

	cmdline := CloudURLToString(bucketName, object)
	str := s.signURL(cmdline, DefaultMethod, DefaultTimeout, c)
	c.Assert(strings.Contains(str, "Expires"), Equals, true)
	c.Assert(strings.Contains(str, "OSSAccessKeyId"), Equals, true)
	c.Assert(strings.Contains(str, "Signature"), Equals, true)
	c.Assert(strings.Contains(str, bucketName), Equals, true)
	c.Assert(strings.Contains(str, object), Equals, true)

	bucket, err := signURLCommand.command.ossBucket(bucketName)
	c.Assert(err, IsNil)

	// get object with url
	err = bucket.GetObjectToFileWithURL(str, downloadFileName)
	c.Assert(err, IsNil)
	str = s.readFile(downloadFileName, c)
	c.Assert(str, Equals, data)

	// not exist object
	object = "中文"
	urlObject := url.QueryEscape(object)
	c.Assert(object != urlObject, Equals, true)

	cmdline = fmt.Sprintf("%s --encoding-type url", CloudURLToString(bucketName, urlObject))
	str = s.signURL(cmdline, DefaultMethod, DefaultTimeout, c)
	c.Assert(strings.Contains(str, "Expires"), Equals, true)
	c.Assert(strings.Contains(str, "OSSAccessKeyId"), Equals, true)
	c.Assert(strings.Contains(str, "Signature"), Equals, true)
	c.Assert(strings.Contains(str, bucketName), Equals, true)

	// get object with url
	err = bucket.GetObjectToFileWithURL(str, downloadFileName)
	c.Assert(err, NotNil)

	// reput object
	data = randStr(100)
	s.createFile(uploadFileName, data, c)
	s.putObject(bucketName, object, uploadFileName, c)

	// sign url
	cmdline = fmt.Sprintf("%s --encoding-type url", CloudURLToString(bucketName, urlObject))
	str = s.signURL(cmdline, DefaultMethod, DefaultTimeout, c)
	c.Assert(strings.Contains(str, "Expires"), Equals, true)
	c.Assert(strings.Contains(str, "OSSAccessKeyId"), Equals, true)
	c.Assert(strings.Contains(str, "Signature"), Equals, true)
	c.Assert(strings.Contains(str, bucketName), Equals, true)

	err = bucket.GetObjectToFileWithURL(str, downloadFileName)
	c.Assert(err, IsNil)
	str = s.readFile(downloadFileName, c)
	c.Assert(str, Equals, data)

	// timeout = -1
	cmdline = fmt.Sprintf("%s --encoding-type url", CloudURLToString(bucketName, urlObject))
	err = s.initSignURL(cmdline, DefaultMethod, -1)
	c.Assert(err, IsNil)
	err = signURLCommand.RunCommand()
	c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestSignurlPut(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	object := randStr(10)

	bucket, err := signURLCommand.command.ossBucket(bucketName)
	c.Assert(err, IsNil)

	// sign url for put
	cmdline := fmt.Sprintf("%s %s --encoding-type url", CloudURLToString(bucketName, object), "x-oss-object-acl:private#X-Oss-Meta-A:A#content-type:text")
	str := s.signURL(cmdline, MethodPut, DefaultTimeout, c)
	c.Assert(strings.Contains(str, "Expires"), Equals, true)
	c.Assert(strings.Contains(str, "OSSAccessKeyId"), Equals, true)
	c.Assert(strings.Contains(str, "Signature"), Equals, true)
	c.Assert(strings.Contains(str, bucketName), Equals, true)

	data := randStr(30)
	// put object with error url
	options := []oss.Option{oss.ContentType("text")}
	err = bucket.PutObjectWithURL(str, strings.NewReader(data), options...)
	c.Assert(err, NotNil)
	c.Assert(err.(oss.ServiceError).Code, Equals, "SignatureDoesNotMatch")

	// put object with url
	options = []oss.Option{oss.ContentType("text"), oss.Meta("A", "A"), oss.ObjectACL(oss.ACLPrivate)}
	err = bucket.PutObjectWithURL(str, strings.NewReader(data), options...)
	c.Assert(err, IsNil)

	// get object meta
	meta, err := bucket.GetObjectDetailedMeta(object)
	c.Assert(err, IsNil)
	c.Assert(meta.Get(oss.HTTPHeaderContentType), Equals, "text")
	c.Assert(meta.Get("X-Oss-Meta-A"), Equals, "A")

	acl, err := bucket.GetObjectACL(object)
	c.Assert(err, IsNil)
	c.Assert(acl.ACL, Equals, string(oss.ACLPrivate))
}

func (s *OssutilCommandSuite) TestSignurlErr(c *C) {
	bucketName := bucketNamePrefix + randLowStr(10)
	s.putBucket(bucketName, c)

	data := "签名url"
	s.createFile(uploadFileName, data, c)

	object := randStr(10)
	s.putObject(bucketName, object, uploadFileName, c)

	cmdline := CloudURLToString("", object)
	err := s.initSignURL(cmdline, DefaultMethod, DefaultTimeout)
	c.Assert(err, IsNil)
	err = signURLCommand.RunCommand()
	c.Assert(err, NotNil)

	cmdline = CloudURLToString(bucketName, "")
	err = s.initSignURL(cmdline, DefaultMethod, DefaultTimeout)
	c.Assert(err, IsNil)
	err = signURLCommand.RunCommand()
	c.Assert(err, NotNil)

	// error header
	cmdline = fmt.Sprintf("%s %s", CloudURLToString(bucketName, object), "Content-Disposition:a")
	err = s.initSignURL(cmdline, DefaultMethod, DefaultTimeout)
	c.Assert(err, IsNil)
	err = signURLCommand.RunCommand()
	c.Assert(err, NotNil)
}
