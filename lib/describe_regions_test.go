package lib

import (
	"encoding/xml"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
	"os"
)

func (s *OssutilCommandSuite) TestRegionsGetSuccess(c *C) {
	strMethod := "get"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	options[OptionMethod] = &strMethod

	region := "oss-cn-hangzhou"
	regionsArgs := []string{region}
	_, err := cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, IsNil)

	regionsFileName := randLowStr(12)
	regionsDownName := regionsFileName + "-down.xml"

	regionsArgs = []string{region, regionsDownName}
	_, err = cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, IsNil)

	_, err = os.Stat(regionsDownName)
	c.Assert(err, IsNil)

	websiteBody := s.readFile(regionsDownName, c)

	list := oss.DescribeRegionsResult{}
	err = xml.Unmarshal([]byte(websiteBody), &list)
	c.Assert(err, IsNil)
	c.Assert(len(list.Regions), Equals, 1)
	c.Assert(list.Regions[0].Region, Equals, region)
	c.Assert(list.Regions[0].InternetEndpoint, Equals, "oss-cn-hangzhou.aliyuncs.com")
	c.Assert(list.Regions[0].InternalEndpoint, Equals, "oss-cn-hangzhou-internal.aliyuncs.com")
	c.Assert(list.Regions[0].AccelerateEndpoint, Equals, "oss-accelerate.aliyuncs.com")

	os.Remove(regionsDownName)
}

func (s *OssutilCommandSuite) TestRegionsGetError(c *C) {
	strMethod := ""
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	options[OptionMethod] = &strMethod

	region := "oss-cn-hangzhou"
	regionsArgs := []string{region}
	_, err := cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, NotNil)

	strMethod = "gets"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, NotNil)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	regionsArgs = []string{}
	_, err = cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestRegionsListError(c *C) {
	strMethod := ""
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	options[OptionMethod] = &strMethod

	region := ""
	regionsArgs := []string{region}
	_, err := cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, NotNil)

	strMethod = "lists"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestRegionsListSuccess(c *C) {
	strMethod := "list"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	options[OptionMethod] = &strMethod

	regionsArgs := []string{}
	_, err := cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, IsNil)

	regionsFileName := randLowStr(12)
	regionsDownName := regionsFileName + "-down.xml"

	regionsArgs = []string{regionsDownName}
	_, err = cm.RunCommand("regions", regionsArgs, options)
	c.Assert(err, IsNil)

	_, err = os.Stat(regionsDownName)
	c.Assert(err, IsNil)

	websiteBody := s.readFile(regionsDownName, c)
	list := oss.DescribeRegionsResult{}
	err = xml.Unmarshal([]byte(websiteBody), &list)
	c.Assert(err, IsNil)
	c.Assert(len(list.Regions) > 0, Equals, true)

	os.Remove(regionsDownName)
}

func (s *OssutilCommandSuite) TestRegionsHelpInfo(c *C) {
	// mkdir command test
	options := OptionMapType{}

	mkArgs := []string{"regions"}
	_, err := cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

	mkArgs = []string{}
	_, err = cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)
}
