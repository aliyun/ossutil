package lib

import (
	"encoding/xml"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
	"os"
	"strings"
	"time"
)

func (s *OssutilCommandSuite) TestAccessPointHelpInfo(c *C) {
	// mkdir command test
	options := OptionMapType{}

	mkArgs := []string{"access-point"}
	_, err := cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

	mkArgs = []string{}
	_, err = cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) TestPutBucketAccessPointError(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	accessFileName := "access-point" + randLowStr(12)
	// access point command test
	strMethod := ""
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	// method is empty
	accessArgs := []string{CloudURLToString(bucketName, ""), accessFileName}
	_, err := cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "--method value is empty")

	//method is error
	strMethod = "puttt"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// cloud url is error
	strMethod = "put"
	options[OptionMethod] = &strMethod
	accessArgs = []string{"http://mybucket", accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// local file is empty
	accessArgs = []string{CloudURLToString(bucketName, ""), accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println("1" + err.Error())

	//local file is not exist
	os.Remove(accessFileName)
	accessArgs = []string{CloudURLToString(bucketName, ""), accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println("2" + err.Error())

	// local file is dir
	err = os.MkdirAll(accessFileName, 0755)
	c.Assert(err, IsNil)
	accessArgs = []string{CloudURLToString(bucketName, ""), accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	os.Remove(accessFileName)

	//local file is empty
	s.createFile(accessFileName, "", c)
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	os.Remove(accessFileName)

	//local file is not xml file
	s.createFile(accessFileName, "aaa", c)
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	os.Remove(accessFileName)

	// StorageURLFromString error
	accessArgs = []string{"oss:///1.jpg"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// bucket name is error
	accessArgs = []string{"oss:///"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	//missing parameter
	accessArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println("2" + err.Error())

	// bucket name not exist
	accessArgs = []string{CloudURLToString("my-bucket", ""), accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	os.Remove(accessFileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestGetBucketAccessPointError(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	apName := "access-point" + randLowStr(12)

	// access point command test
	strMethod := ""
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	// method is empty
	accessArgs := []string{CloudURLToString(bucketName, ""), apName}
	_, err := cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	//method is error
	strMethod = "gett"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// cloud url is error
	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{"http://mybucket", apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// apName not exist
	accessArgs = []string{CloudURLToString(bucketName, ""), apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// StorageURLFromString error
	accessArgs = []string{"oss:///1.jpg"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// bucket name is error
	accessArgs = []string{"oss:///"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	//missing parameter
	accessArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// bucket name not exist
	accessArgs = []string{CloudURLToString("my-bucket", ""), apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestListBucketAccessPointError(c *C) {
	accessFileName := "access-point-" + randLowStr(12)
	// access point command test
	strMethod := ""
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	// method is empty
	accessArgs := []string{CloudURLToString(bucketNameExist, ""), accessFileName}
	_, err := cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	//method is error
	strMethod = "listss"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// cloud url is error
	strMethod = "list"
	options[OptionMethod] = &strMethod
	accessArgs = []string{"http://mybucket", accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// StorageURLFromString error
	accessArgs = []string{"oss:///1.jpg"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// bucket name is error
	accessArgs = []string{"oss:///"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// bucket name not exist
	accessArgs = []string{CloudURLToString("my-bucket", "")}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	os.Remove(accessFileName)
}

func (s *OssutilCommandSuite) TestDeleteBucketAccessPointError(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	apName := "access-point" + randLowStr(12)

	// access point command test
	strMethod := ""
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}

	// method is empty
	accessArgs := []string{CloudURLToString(bucketName, ""), apName}
	_, err := cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	//method is error
	strMethod = "deleteteddd"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)

	// cloud url is error
	strMethod = "delete"
	options[OptionMethod] = &strMethod
	accessArgs = []string{"http://mybucket", apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// apName not exist
	accessArgs = []string{CloudURLToString(bucketName, ""), apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// StorageURLFromString error
	accessArgs = []string{"oss:///1.jpg"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// bucket name is error
	accessArgs = []string{"oss:///"}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	//missing ap name
	accessArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// bucket name not exist
	accessArgs = []string{CloudURLToString("my-bucket", ""), apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestAccessPointPolicyError(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	if accountID == "" {
		c.ExpectFailure("account ID is empty!")
	}
	apName := "access-point-" + randLowStr(8)
	tmp := strings.SplitN(endpoint, ".", 3)
	region := tmp[0]
	policy := `{
   "Version":"1",
   "Statement":[
   {
     "Action":[
       "oss:PutObject",
       "oss:GetObject"
    ],
    "Effect":"Deny",
    "Principal":["1234567890"],
    "Resource":[
       "acs:oss:` + region + `:` + accountID + `:accesspoint/` + apName + `",
       "acs:oss:` + region + `:` + accountID + `:accesspoint/` + apName + `/object/*"
     ]
   }
  ]
 }`
	accessFileName := "ap-policy-" + randLowStr(8)
	s.createFile(accessFileName, policy, c)

	// access point policy command test
	strMethod := "put"
	strItem := "policy"
	get := "get"
	del := "delete"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
		"item":            &strItem,
	}

	// missing ap name
	accessArgs := []string{CloudURLToString(bucketName, "")}
	_, err := cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())
	c.Assert(strings.Contains(err.Error(), "put access point policy need at least 3 parameters"), Equals, true)

	// missing json file
	accessArgs = []string{CloudURLToString(bucketName, ""), apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "put access point policy need at least 3 parameters"), Equals, true)

	// apName not exist
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// StorageURLFromString error
	accessArgs = []string{"oss:///1.jpg", apName, accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// bucket name is error
	accessArgs = []string{"oss:///", apName, accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	// bucket name not exist
	accessArgs = []string{CloudURLToString("my-bucket", ""), apName, accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	errItem := "policysss"
	options[OptionItem] = &errItem
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessFileName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "--item value is not in the optional value:policy"), Equals, true)

	accessResult := "get-policy-" + randLowStr(8)
	options[OptionMethod] = &get
	accessArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "--item value is not in the optional value:policy"), Equals, true)

	options[OptionItem] = &strItem
	accessArgs = []string{"oss:///1.jpg", apName, accessResult}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	accessArgs = []string{"oss:///", apName, accessResult}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	accessArgs = []string{CloudURLToString("my-bucket", ""), apName, accessResult}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	options[OptionItem] = &errItem
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessResult}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())
	c.Assert(strings.Contains(err.Error(), "--item value is not in the optional value:policy"), Equals, true)

	options[OptionMethod] = &del
	accessArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "--item value is not in the optional value:policy"), Equals, true)

	options[OptionItem] = &strItem
	accessArgs = []string{"oss:///1.jpg", apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	accessArgs = []string{"oss:///", apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	accessArgs = []string{CloudURLToString("my-bucket", ""), apName}
	_, err = cm.RunCommand("access-point", accessArgs, options)
	c.Assert(err, NotNil)
	testLogger.Println(err.Error())

	os.Remove(accessFileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestBucketAccessPoint(c *C) {
	if accountID == "" {
		c.ExpectFailure("account ID is empty!")
	}
	apName := "ap1-" + randLowStr(8)
	accessXml := `<?xml version="1.0" encoding="UTF-8"?>
<CreateAccessPointConfiguration>
    <AccessPointName>` + apName + `</AccessPointName>
    <NetworkOrigin>internet</NetworkOrigin>
</CreateAccessPointConfiguration>`

	accessFileName := "ap-config-" + randLowStr(8)
	s.createFile(accessFileName, accessXml, c)

	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	// access point command test
	str := ""
	strMethod := "put"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}
	accessResult := "ap-result-" + randLowStr(8)
	command := "access-point"
	accessArgs := []string{CloudURLToString(bucketName, ""), accessFileName, accessResult}
	_, err := cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)

	_, err = os.Stat(accessResult)
	c.Assert(err, IsNil)
	accessBody := s.readFile(accessResult, c)

	var putOut oss.CreateBucketAccessPointResult
	err = xml.Unmarshal([]byte(accessBody), &putOut)
	c.Assert(err, IsNil)

	os.Remove(accessFileName)
	os.Remove(accessResult)

	time.Sleep(3 * time.Second)

	apName1 := "ap2-" + randLowStr(8)
	accessXml = `<?xml version="1.0" encoding="UTF-8"?>
<CreateAccessPointConfiguration>
    <AccessPointName>` + apName1 + `</AccessPointName>
    <NetworkOrigin>vpc</NetworkOrigin>
<VpcConfiguration>
      <VpcId>vpc-1234567890</VpcId>
    </VpcConfiguration>
</CreateAccessPointConfiguration>`

	s.createFile(accessFileName, accessXml, c)

	accessArgs = []string{CloudURLToString(bucketName, ""), accessFileName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	os.Remove(accessFileName)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessGetResult := "ap-get-" + randLowStr(8)
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessGetResult}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)

	// check access point DownName
	_, err = os.Stat(accessGetResult)
	c.Assert(err, IsNil)

	accessBody = s.readFile(accessGetResult, c)

	var out oss.GetBucketAccessPointResult
	err = xml.Unmarshal([]byte(accessBody), &out)
	c.Assert(err, IsNil)
	c.Assert(out.AccessPointName, Equals, apName)
	accessPointArn := out.AccessPointArn
	aliasName := out.Alias
	os.Remove(accessFileName)
	os.Remove(accessGetResult)
	time.Sleep(3 * time.Second)

	accessListResult := "ap-list-" + randLowStr(8)
	strMethod = "list"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(bucketName, ""), accessListResult}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	time.Sleep(3 * time.Second)
	_, err = os.Stat(accessListResult)
	c.Assert(err, IsNil)
	accessBody = s.readFile(accessListResult, c)
	var list oss.ListBucketAccessPointsResult
	err = xml.Unmarshal([]byte(accessBody), &list)
	c.Assert(err, IsNil)
	c.Assert(len(list.AccessPoints), Equals, 2)

	c.Assert(list.AccessPoints[1].AccessPointName, Equals, apName1)
	c.Assert(list.AccessPoints[1].Bucket, Equals, bucketName)
	c.Assert(list.AccessPoints[1].Alias != "", Equals, true)
	c.Assert(list.AccessPoints[1].NetworkOrigin, Equals, "vpc")
	c.Assert(list.AccessPoints[1].Status != "", Equals, true)
	os.Remove(accessListResult)

	policy := `{
   "Version":"1",
   "Statement":[
   {
     "Action":[
       	"oss:*"
    ],
    "Effect": "Allow",
    "Principal":["` + accountID + `"],
    "Resource":[
		"` + accessPointArn + `",
		"` + accessPointArn + `/object/*"
     ]
   }
  ]
 }`
	accessJsonName := "ap-policy-" + randLowStr(6)
	s.createFile(accessJsonName, policy, c)
	strMethod = "put"
	strItem := "policy"
	options[OptionMethod] = &strMethod
	options[OptionItem] = &strItem
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	os.Remove(accessJsonName)
	time.Sleep(3 * time.Second)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(bucketName, ""), apName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	time.Sleep(3 * time.Second)

	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	time.Sleep(3 * time.Second)

	_, err = os.Stat(accessJsonName)
	c.Assert(err, IsNil)
	accessBody = s.readFile(accessJsonName, c)
	//c.Assert(accessBody,Equals,policy)
	os.Remove(accessJsonName)

	policy1 := `{
   "Version":"1",
   "Statement":[
   {
     "Action":[
       	"oss:*"
    ],
    "Effect": "Allow",
    "Principal":["` + accountID + `"],
    "Resource":[
		"` + accessPointArn + `",
		"` + accessPointArn + `/object/*"
     ]
   },
	{
     "Action":[
       "oss:PutObject",
       "oss:GetObject"
    ],
    "Effect":"Deny",
    "Principal":["123456"],
    "Resource":[
       "` + accessPointArn + `",
       "` + accessPointArn + `/object/*"
     ]
   }
  ]
 }`
	for {
		strMethod = "get"
		options[OptionMethod] = &strMethod
		options[OptionItem] = &str
		accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessGetResult}
		_, err = cm.RunCommand(command, accessArgs, options)
		c.Assert(err, IsNil)
		_, err = os.Stat(accessGetResult)
		c.Assert(err, IsNil)
		accessBody = s.readFile(accessGetResult, c)
		var out oss.GetBucketAccessPointResult
		err = xml.Unmarshal([]byte(accessBody), &out)
		c.Assert(err, IsNil)
		os.Remove(accessGetResult)
		if out.Status == "enable" {
			break
		}
		time.Sleep(3 * time.Second)
	}
	s.createFile(accessJsonName, policy1, c)
	strMethod = "put"
	options[OptionItem] = &strItem
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(aliasName, ""), apName, accessJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(aliasName, ""), apName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)

	downJsonName := "ap-down-" + randLowStr(5)
	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(aliasName, ""), apName, downJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	os.Remove(downJsonName)

	strMethod = "delete"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(aliasName, ""), apName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)

	strMethod = "put"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	os.Remove(accessJsonName)

	strMethod = "delete"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(bucketName, ""), apName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)

	for {
		strMethod = "list"
		options[OptionMethod] = &strMethod
		options[OptionItem] = &str
		accessArgs = []string{CloudURLToString(bucketName, ""), accessListResult}
		_, err = cm.RunCommand(command, accessArgs, options)
		c.Assert(err, IsNil)
		accessBody = s.readFile(accessListResult, c)
		var list2 oss.ListBucketAccessPointsResult
		err = xml.Unmarshal([]byte(accessBody), &list2)
		c.Assert(err, IsNil)
		os.Remove(accessListResult)
		if len(list2.AccessPoints) > 0 {
			for _, point := range list2.AccessPoints {
				if point.Status == "enable" {
					strMethod = "delete"
					options[OptionMethod] = &strMethod
					accessArgs = []string{CloudURLToString(bucketName, ""), point.AccessPointName}
					_, err = cm.RunCommand(command, accessArgs, options)
					c.Assert(err, IsNil)
				}
			}
		} else {
			break
		}
		time.Sleep(60 * time.Second)
	}
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestBucketAccessPointOtherError(c *C) {
	if accountID == "" {
		c.ExpectFailure("account ID is empty!")
	}
	apName := "ap1-" + randLowStr(8)
	accessXml := `<?xml version="1.0" encoding="UTF-8"?>
<CreateAccessPointConfiguration>
    <AccessPointName>` + apName + `</AccessPointName>
    <NetworkOrigin>internet</NetworkOrigin>
</CreateAccessPointConfiguration>`

	accessFileName := "ap-config-" + randLowStr(8)
	s.createFile(accessFileName, accessXml, c)

	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	// access point command test
	str := ""
	strMethod := "put"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}
	accessResult := "ap-result-" + randLowStr(8)
	command := "access-point"
	accessArgs := []string{CloudURLToString(bucketName, ""), accessFileName, accessResult}
	_, err := cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)

	_, err = os.Stat(accessResult)
	c.Assert(err, IsNil)
	accessBody := s.readFile(accessResult, c)

	var putOut oss.CreateBucketAccessPointResult
	err = xml.Unmarshal([]byte(accessBody), &putOut)
	c.Assert(err, IsNil)
	os.Remove(accessFileName)
	os.Remove(accessResult)

	time.Sleep(3 * time.Second)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessGetResult := "ap-get-" + randLowStr(8)
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessGetResult}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	accessBody = s.readFile(accessGetResult, c)
	var out oss.GetBucketAccessPointResult
	err = xml.Unmarshal([]byte(accessBody), &out)
	c.Assert(err, IsNil)
	c.Assert(out.AccessPointName, Equals, apName)
	accessPointArn := out.AccessPointArn
	os.Remove(accessGetResult)
	time.Sleep(3 * time.Second)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(bucketName, ""), "not-exist-ap", accessGetResult}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, NotNil)
	time.Sleep(3 * time.Second)

	policy := `{
   "Version":"1",
   "Statement":[
   {
     "Action":[
       	"oss:*"
    ],
    "Effect": "Allow",
    "Principal":["` + accountID + `"],
    "Resource":[
		"` + accessPointArn + `",
		"` + accessPointArn + `/object/*"
     ]
   }
  ]
 }`
	accessJsonName := "ap-policy-" + randLowStr(6)
	s.createFile(accessJsonName, policy, c)
	strMethod = "put"
	strItem := "policy"
	options[OptionMethod] = &strMethod
	options[OptionItem] = &strItem
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	os.Remove(accessJsonName)
	time.Sleep(3 * time.Second)

	s.createFile(accessJsonName, "error-policy", c)
	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, NotNil)
	os.Remove(accessJsonName)
	time.Sleep(3 * time.Second)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(bucketName, ""), apName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	time.Sleep(3 * time.Second)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString(bucketName, ""), "not-exist-ap-name"}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, NotNil)
	time.Sleep(3 * time.Second)

	accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, IsNil)
	time.Sleep(3 * time.Second)

	_, err = os.Stat(accessJsonName)
	c.Assert(err, IsNil)
	accessBody = s.readFile(accessJsonName, c)
	os.Remove(accessJsonName)

	for {
		strMethod = "get"
		options[OptionMethod] = &strMethod
		options[OptionItem] = &str
		accessArgs = []string{CloudURLToString(bucketName, ""), apName, accessGetResult}
		_, err = cm.RunCommand(command, accessArgs, options)
		c.Assert(err, IsNil)
		_, err = os.Stat(accessGetResult)
		c.Assert(err, IsNil)
		accessBody = s.readFile(accessGetResult, c)
		var out oss.GetBucketAccessPointResult
		err = xml.Unmarshal([]byte(accessBody), &out)
		c.Assert(err, IsNil)
		os.Remove(accessGetResult)
		if out.Status == "enable" {
			break
		}
		time.Sleep(3 * time.Second)
	}

	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString("error-alias", ""), apName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, NotNil)

	downJsonName := "ap-down-" + randLowStr(5)
	strMethod = "get"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString("error-alias", ""), apName, downJsonName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, NotNil)

	strMethod = "delete"
	options[OptionMethod] = &strMethod
	accessArgs = []string{CloudURLToString("error-alias", ""), apName}
	_, err = cm.RunCommand(command, accessArgs, options)
	c.Assert(err, NotNil)

	accessListResult := "ap-list-" + randLowStr(8)
	for {
		strMethod = "list"
		options[OptionMethod] = &strMethod
		options[OptionItem] = &str
		accessArgs = []string{CloudURLToString(bucketName, ""), accessListResult}
		_, err = cm.RunCommand(command, accessArgs, options)
		c.Assert(err, IsNil)
		accessBody = s.readFile(accessListResult, c)
		var list2 oss.ListBucketAccessPointsResult
		err = xml.Unmarshal([]byte(accessBody), &list2)
		c.Assert(err, IsNil)
		os.Remove(accessListResult)
		if len(list2.AccessPoints) > 0 {
			for _, point := range list2.AccessPoints {
				if point.Status == "enable" {
					strMethod = "delete"
					options[OptionMethod] = &strMethod
					accessArgs = []string{CloudURLToString(bucketName, ""), point.AccessPointName}
					_, err = cm.RunCommand(command, accessArgs, options)
					c.Assert(err, IsNil)
				}
			}
		} else {
			break
		}
		time.Sleep(60 * time.Second)
	}
	s.removeBucket(bucketName, true, c)
}
