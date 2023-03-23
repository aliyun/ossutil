package lib

import (
	"encoding/xml"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	. "gopkg.in/check.v1"
	"os"
	"strings"
)

func (s *OssutilCommandSuite) TestReservedCapacityHelp(c *C) {
	// mkdir command test
	options := OptionMapType{}

	mkArgs := []string{"reserved-capacity"}
	_, err := cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)

	mkArgs = []string{}
	_, err = cm.RunCommand("help", mkArgs, options)
	c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) TestReservedCapacity(c *C) {
	strRedundancy := "LRS"
	name := "ossutil-" + randLowStr(8)
	createXml := `<ReservedCapacityConfiguration>
  <Name>` + name + `</Name>
  <DataRedundancyType>` + strRedundancy + `</DataRedundancyType>
  <ReservedCapacity>10240</ReservedCapacity>
</ReservedCapacityConfiguration>`
	createFileName := "create-reserved-" + randLowStr(4) + ".xml"
	s.createFile(createFileName, createXml, c)

	// test create reserved capacity
	strMethod := "create"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
		"method":          &strMethod,
	}
	createArgs := []string{createFileName}
	_, err := cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(err, IsNil)

	// test list reserved capacity
	reservedDownName := "list-reserved-" + randLowStr(4) + "-down" + ".xml"
	strMethod = "list"
	options[OptionMethod] = &strMethod
	listArgs := []string{reservedDownName}
	_, err = cm.RunCommand("reserved-capacity", listArgs, options)
	c.Assert(err, IsNil)

	_, err = os.Stat(reservedDownName)
	c.Assert(err, IsNil)
	var out oss.ListReservedCapacityResult
	listBody := s.readFile(reservedDownName, c)
	err = xml.Unmarshal([]byte(listBody), &out)
	c.Assert(err, IsNil)
	c.Assert(len(out.ReservedCapacityRecord) > 0, Equals, true)

	id := ""
	for _, record := range out.ReservedCapacityRecord {
		if record.Name == name {
			id = record.InstanceId
		}
	}

	// test get reserved capacity
	getReservedDownName := "get-reserved-" + randLowStr(4) + "-down" + ".xml"
	strMethod = "get"
	options[OptionMethod] = &strMethod
	getArgs := []string{id, getReservedDownName}
	_, err = cm.RunCommand("reserved-capacity", getArgs, options)
	c.Assert(err, IsNil)
	var info oss.GetReservedCapacityResult
	body := s.readFile(getReservedDownName, c)
	err = xml.Unmarshal([]byte(body), &info)
	c.Assert(err, IsNil)
	c.Assert(info.DataRedundancyType, Equals, "LRS")
	c.Assert(info.ReservedCapacity, Equals, int64(10240))
	c.Assert(info.Name, Equals, name)

	// test update reserved capacity
	size := "100"
	maxSize := "20480"
	updateXml := `<ReservedCapacityConfiguration>
  <Status>Enabled</Status>
  <ReservedCapacity>10240</ReservedCapacity>
  <AutoExpansionSize>` + size + `</AutoExpansionSize>
  <AutoExpansionMaxSize>` + maxSize + `</AutoExpansionMaxSize>
</ReservedCapacityConfiguration>`
	updateReservedName := "update-reserved-" + randLowStr(4) + ".xml"
	s.createFile(updateReservedName, updateXml, c)
	strMethod = "update"
	options[OptionMethod] = &strMethod
	updateArgs := []string{id, updateReservedName}
	_, err = cm.RunCommand("reserved-capacity", updateArgs, options)
	c.Assert(err, IsNil)

	// test get reserved capacity
	os.Remove(getReservedDownName)
	getReservedDownName = "get-reserved-" + randLowStr(4) + "-down" + ".xml"
	strMethod = "get"
	options[OptionMethod] = &strMethod
	getArgs = []string{id, getReservedDownName}
	_, err = cm.RunCommand("reserved-capacity", getArgs, options)
	c.Assert(err, IsNil)

	body = s.readFile(getReservedDownName, c)
	err = xml.Unmarshal([]byte(body), &info)
	c.Assert(err, IsNil)
	c.Assert(info.DataRedundancyType, Equals, "LRS")
	c.Assert(info.ReservedCapacity, Equals, int64(10240))
	c.Assert(info.Name, Equals, name)
	c.Assert(info.AutoExpansionSize, Equals, int64(100))
	c.Assert(info.AutoExpansionMaxSize, Equals, int64(20480))

	// test list bucket reserved capacity
	listBucketDownName := "list-bucket-reserved-" + randLowStr(4) + "-down" + ".xml"
	strMethod = "list-bucket"
	options[OptionMethod] = &strMethod
	getArgs = []string{id, listBucketDownName}
	_, err = cm.RunCommand("reserved-capacity", getArgs, options)
	c.Assert(err, IsNil)

	var listBucket oss.ListBucketWithReservedCapacityResult
	body = s.readFile(listBucketDownName, c)
	err = xml.Unmarshal([]byte(body), &listBucket)
	c.Assert(err, IsNil)

	c.Assert(len(listBucket.BucketList) == 0, Equals, true)

	// test create bucket under the reserved capacity

	bucketName := bucketNamePrefix + randLowStr(10)
	strStorageClass := "ReservedCapacity"
	options = OptionMapType{
		"endpoint":                   &endpoint,
		"accessKeyID":                &accessKeyID,
		"accessKeySecret":            &accessKeySecret,
		"redundancyType":             &strRedundancy,
		"storageClass":               &strStorageClass,
		"reservedCapacityInstanceId": &id,
	}
	getArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("mb", getArgs, options)
	c.Assert(err, IsNil)

	os.Remove(listBucketDownName)
	listBucketDownName = "list-bucket-reserved-" + randLowStr(4) + "-down" + ".xml"
	strMethod = "list-bucket"
	options = OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	options[OptionMethod] = &strMethod
	getArgs = []string{id, listBucketDownName}
	_, err = cm.RunCommand("reserved-capacity", getArgs, options)
	c.Assert(err, IsNil)

	body = s.readFile(listBucketDownName, c)
	err = xml.Unmarshal([]byte(body), &listBucket)
	c.Assert(err, IsNil)

	c.Assert(len(listBucket.BucketList) > 0, Equals, true)
	c.Assert(listBucket.BucketList[0], Equals, bucketName)
	c.Assert(listBucket.InstanceId, Equals, id)

	outputFile := "stat-bucket-" + randLowStr(4) + "-down" + ".txt"
	testResultFile, err = os.OpenFile(outputFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
	c.Assert(err, IsNil)

	oldStdout := os.Stdout
	os.Stdout = testResultFile

	options = OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	getArgs = []string{CloudURLToString(bucketName, "")}
	_, err = cm.RunCommand("stat", getArgs, options)
	c.Assert(err, IsNil)
	testResultFile.Close()
	os.Stdout = oldStdout

	body = s.readFile(outputFile, c)
	c.Assert(strings.Contains(body, "ReservedCapacityInstanceId: "+id), Equals, true)
	s.removeBucket(bucketName, true, c)

	os.Remove(createXml)
	os.Remove(reservedDownName)
	os.Remove(getReservedDownName)
	os.Remove(updateReservedName)
	os.Remove(listBucketDownName)
	os.Remove(outputFile)
}

func (s *OssutilCommandSuite) TestReservedCapacityError(c *C) {
	// test create reserved capacity
	strMethod := "create"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	options[OptionMethod] = &strMethod
	createArgs := []string{}
	_, err := cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "create reserved capacity need at least 1 parameters,the local xml file is empty"), Equals, true)

	strMethod = "creatert"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "--method value is not in the optional value:create|update|get|list|list-bucket"), Equals, true)

	// test get
	strMethod = "gets"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "--method value is not in the optional value:create|update|get|list|list-bucket"), Equals, true)

	strMethod = "get"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "get reserved capacity need at least 1 parameters,the id is empty"), Equals, true)

	// test update
	strMethod = "updats"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "--method value is not in the optional value:create|update|get|list|list-bucket"), Equals, true)

	strMethod = "update"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "update reserved capacity need at least 2 parameters,the local xml file is empty"), Equals, true)

	updateArgs := []string{"id"}
	_, err = cm.RunCommand("reserved-capacity", updateArgs, options)
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "update reserved capacity need at least 2 parameters,the local xml file is empty"), Equals, true)

	// test list reserved capacity
	strMethod = "listss"
	options[OptionMethod] = &strMethod
	createArgs = []string{}
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(strings.Contains(err.Error(), "--method value is not in the optional value:create|update|get|list|list-bucket"), Equals, true)

	// test list bucket
	strMethod = "list-bucket-s"
	options[OptionMethod] = &strMethod
	createArgs = []string{}
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(strings.Contains(err.Error(), "--method value is not in the optional value:create|update|get|list|list-bucket"), Equals, true)

	strMethod = "list-bucket"
	options[OptionMethod] = &strMethod
	_, err = cm.RunCommand("reserved-capacity", createArgs, options)
	c.Assert(strings.Contains(err.Error(), "list bucket under reserved capacity need at least 1 parameters,the id is empty"), Equals, true)
}

func (s *OssutilCommandSuite) TestReservedCapacityConfirm(c *C) {
	strMethod := "list"
	options := OptionMapType{
		"endpoint":        &endpoint,
		"accessKeyID":     &accessKeyID,
		"accessKeySecret": &accessKeySecret,
	}
	options[OptionMethod] = &strMethod
	listXml := `<ReservedCapacityConfiguration>
  <Status>Enabled</Status>
  <ReservedCapacity>10240</ReservedCapacity>
  <AutoExpansionSize></AutoExpansionSize>
  <AutoExpansionMaxSize></AutoExpansionMaxSize>
</ReservedCapacityConfiguration>`
	listReservedName := "list-reserved-" + randLowStr(4) + ".xml"
	s.createFile(listReservedName, listXml, c)
	cmdArgs := []string{listReservedName}
	_, err := cm.RunCommand("reserved-capacity", cmdArgs, options)
	c.Assert(err, IsNil)

	os.Remove(listXml)
}
