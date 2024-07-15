package lib

import (
	"encoding/xml"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"

	. "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestHttpsConfigPutSuccess(c *C) {
	httpsXml := `<?xml version="1.0" encoding="UTF-8"?>
<HttpsConfiguration>  
  <TLS>
    <Enable>true</Enable>   
    <TLSVersion>TLSv1.2</TLSVersion>
    <TLSVersion>TLSv1.3</TLSVersion>
  </TLS>
</HttpsConfiguration>`

	httpsFileName := randLowStr(12)
	s.createFile(httpsFileName, httpsXml, c)

	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	// https config command test
	var str string
	strMethod := "put"
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"configFile":      &configFile,
		"method":          &strMethod,
	}

	httpsArgs := []string{CloudURLToString(bucketName, ""), httpsFileName}
	_, err := cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, IsNil)

	// check,get cors
	httpsDownName := httpsFileName + "-down"
	strMethod = "get"

	httpsArgs = []string{CloudURLToString(bucketName, ""), httpsDownName}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, IsNil)

	// check httpsDownName
	_, err = os.Stat(httpsDownName)
	c.Assert(err, IsNil)

	httpsBody := s.readFile(httpsDownName, c)

	httpsConfigDest := oss.GetBucketHttpsConfigResult{}
	err = xml.Unmarshal([]byte(httpsBody), &httpsConfigDest)
	c.Assert(err, IsNil)
	c.Assert(httpsConfigDest.TLS.Enable, Equals, true)
	c.Assert(httpsConfigDest.TLS.TLSVersion[0], Equals, "TLSv1.2")
	c.Assert(httpsConfigDest.TLS.TLSVersion[1], Equals, "TLSv1.3")
	os.Remove(httpsFileName)
	os.Remove(httpsDownName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestHttpsConfigPutError(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	httpsFileName := "httpsFile-" + randLowStr(12)

	// cors command test
	var str string
	strMethod := ""
	options := OptionMapType{
		"endpoint":        &str,
		"accessKeyID":     &str,
		"accessKeySecret": &str,
		"configFile":      &configFile,
		"method":          &strMethod,
	}

	// method is empty
	httpsArgs := []string{CloudURLToString(bucketName, ""), httpsFileName}
	_, err := cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	//method is error
	strMethod = "puttt"
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	// cloudurl is error
	strMethod = "put"
	httpsArgs = []string{"http://mybucket", httpsFileName}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	// local file is emtpy
	httpsArgs = []string{CloudURLToString(bucketName, ""), httpsFileName}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	//local file is not exist
	os.Remove(httpsFileName)
	httpsArgs = []string{CloudURLToString(bucketName, ""), httpsFileName}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	// localfile is dir
	err = os.MkdirAll(httpsFileName, 0755)
	c.Assert(err, IsNil)
	httpsArgs = []string{CloudURLToString(bucketName, ""), httpsFileName}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)
	os.Remove(httpsFileName)

	//local file is emtpy
	s.createFile(httpsFileName, "", c)
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)
	os.Remove(httpsFileName)

	//local file is not xml file
	s.createFile(httpsFileName, "aaa", c)
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)
	os.Remove(httpsFileName)

	// StorageURLFromString error
	httpsArgs = []string{"oss:///1.jpg"}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	// bucketname is error
	httpsArgs = []string{"oss:///"}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	//missing parameter
	httpsArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	// bucketname not exist
	httpsArgs = []string{CloudURLToString("my-bucket", "")}
	_, err = cm.RunCommand("https-config", httpsArgs, options)
	c.Assert(err, NotNil)

	os.Remove(httpsFileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestHttpsConfigHelpInfo(c *C) {
	// mkdir command test
	options := OptionMapType{}

	mkArgs := []string{"https-config"}
	_, err := cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

	mkArgs = []string{}
	_, err = cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

}
