package lib 

import (
    "fmt"
    "os"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) rawList(args []string, shortFormat, directory bool) (bool, error) {
    command := "ls"
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "shortFormat": &shortFormat,
        "directory": &directory,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

// test list buckets
func (s *OssutilCommandSuite) TestListLoadConfig(c *C) {
    command := "ls"
    var args []string
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    args = []string{"oss://"}
    options = OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
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
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestListErrConfigFile(c *C) {
    cfile := "ossutil_test.config_boto"
    s.createFile(cfile, content, c)

    command := "ls"
    var args []string
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListConfigFile(c *C) {
    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\nretryTimes=%s", endpoint, accessKeyID, accessKeySecret, "errretry") 
    s.createFile(cfile, data, c)

    command := "ls"
    var args []string
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListWithBucketEndpoint(c *C) {
    bucket := bucketNamePrefix + "ls"
    s.putBucket(bucket, c)

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s", "abc", accessKeyID, accessKeySecret, bucket, endpoint) 
    s.createFile(cfile, data, c)

    command := "ls"
    args := []string{CloudURLToString(bucket, "")}
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListWithBucketCname(c *C) {
    bucket := bucketNamePrefix + "ls"
    s.putBucket(bucket, c)

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s\n[Bucket-Cname]\n%s=%s", "abc", accessKeyID, accessKeySecret, bucket, "abc", bucket, bucket + "." +endpoint) 
    s.createFile(cfile, data, c)

    command := "ls"
    args := []string{CloudURLToString(bucket, "")}
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    //_ = os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListBuckets(c *C) {
    // "ls" 
    buckets := s.listBuckets(false, c)
    bucketNum := len(buckets)

    // "ls -s"
    buckets = s.listBuckets(true, c)
    c.Assert(len(buckets), Equals, bucketNum)

    // put bucket
    bucket := bucketNamePrefix + "ls" 
    s.putBucket(bucket, c)

    // get result
    buckets = s.listBuckets(false, c)
    c.Assert(len(buckets), Equals, bucketNum + 1)
    c.Assert(FindPos(bucket, buckets) != -1, Equals, true)

    // remove bucket
    s.removeBucket(bucket, true, c)

    // get result
    buckets = s.listBuckets(false, c)
    c.Assert(len(buckets), Equals, bucketNum)
    c.Assert(FindPos(bucket, buckets) == -1, Equals, true)
}

// list objects with not exist bucket 
func (s *OssutilCommandSuite) TestListObjectsBucketNotExist(c *C) {
    bucket := bucketNamePrefix + "notexist"
    command := "ls"
    args := []string{CloudURLToString(bucket, "")}
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

// list objects
func (s *OssutilCommandSuite) TestListObjects(c *C) {
    bucket := bucketNamePrefix + "ls"
    s.putBucket(bucket, c)

    // "ls oss://bucket"
    objects := s.listObjects(bucket, "", false, false, c)
    c.Assert(len(objects), Equals, 0)

    // put objects
    num := 10
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("中文测试:#%d", i) 
        s.putObject(bucket, object, uploadFileName, c) 
        objectNames = append(objectNames, object)
    }

    object := "another_object"
    s.putObject(bucket, object, uploadFileName, c)
    objectNames = append(objectNames, object)

    // "ls oss://bucket -s"
    objects = s.listObjects(bucket, "", true, false, c)
    c.Assert(len(objects), Equals, len(objectNames))

    // "ls oss://bucket/prefix -s"
    objects = s.listObjects(bucket, "中文测试:", true, false, c)
    c.Assert(len(objects), Equals, len(objectNames) - 1)

    //put directories
    num = 5 
    objectNames = []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("中文测试:#%d/", i) 
        s.putObject(bucket, object, uploadFileName, c) 

        object = fmt.Sprintf("中文测试:#%d/%d/", i, i) 
        s.putObject(bucket, object, uploadFileName, c) 
        objectNames = append(objectNames, object)
    }

    // "ls oss://bucket/prefix"
    objects = s.listObjects(bucket, "中文测试:#", false, false, c)
    c.Assert(len(objects), Equals, 20)

    // "ls oss://bucket/prefix -d"
    objects = s.listObjects(bucket, "中文测试:#", false, true, c)
    c.Assert(len(objects), Equals, 15)

    objects = s.listObjects(bucket, "中文测试:#1/", false, true, c)
    c.Assert(len(objects), Equals, 2)
}

func (s *OssutilCommandSuite) TestErrList(c *C) {
    showElapse, err := s.rawList([]string{"../"}, true, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    bucket := bucketNamePrefix + "ls"
    showElapse, err = s.rawList([]string{CloudURLToString(bucket, "")}, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}
