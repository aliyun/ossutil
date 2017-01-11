package lib 

import (
    "fmt"
    "os"
    "time"

    . "gopkg.in/check.v1"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (s *OssutilCommandSuite) rawList(args []string, shortFormat, directory bool, multipart, allType bool) (bool, error) {
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
        "multipart": &multipart,
        "allType": &allType,
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
    bucket := bucketNameExist 

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
    bucket := bucketNamePrefix + "ls1"
    s.putBucket(bucket, c)

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s\n[Bucket-Cname]\n%s=%s", "abc", accessKeyID, accessKeySecret, bucket, "abc", bucket, bucket + "." + endpoint) 
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
    s.removeBucket(bucket, true, c)
    time.Sleep(7*time.Second)
}

func (s *OssutilCommandSuite) TestListBuckets(c *C) {
    // "ls" 
    bucket := bucketNamePrefix + "ls2" 
    // put bucket
    s.putBucket(bucket, c)
    time.Sleep(10*time.Second)

    // get result
    buckets := s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) != -1, Equals, true)
    bucketNum := len(buckets)

    // remove empty bucket
    s.removeBucket(bucket, false, c)
    time.Sleep(10*time.Second)

    // get result
    buckets = s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) == -1, Equals, true)
    c.Assert(len(buckets) <= bucketNum, Equals, true)
}

// list objects with not exist bucket 
func (s *OssutilCommandSuite) TestListObjectsBucketNotExist(c *C) {
    bucket := bucketNameNotExist 
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
    bucket := bucketNameList 

    // put objects
    num := 3 
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("lstest:#%d", i) 
        s.putObject(bucket, object, uploadFileName, c) 
    }

    object := "another_object"
    s.putObject(bucket, object, uploadFileName, c)
    time.Sleep(sleepTime)

    objectStat := s.getStat(bucket, object, c)
    c.Assert(objectStat[StatACL], Equals, "default")
    c.Assert(len(objectStat["Etag"]), Equals, 32)
    c.Assert(objectStat["Last-Modified"] != "", Equals, true)
    c.Assert(objectStat[StatOwner] != "", Equals, true)

    //put directories
    num1 := 2 
    for i := 0; i < num1; i++ {
        object := fmt.Sprintf("lstest:#%d/", i) 
        s.putObject(bucket, object, uploadFileName, c) 

        object = fmt.Sprintf("lstest:#%d/%d/", i, i) 
        s.putObject(bucket, object, uploadFileName, c) 
    }

    // "ls oss://bucket -s"
    //objects := s.listObjects(bucket, "", true, false, false, false, c)
    //c.Assert(len(objects), Equals, num + 2*num1 + 1)

    // "ls oss://bucket/prefix -s"
    objects := s.listObjects(bucket, "lstest:", true, false, false, false, c)
    c.Assert(len(objects), Equals, num + 2*num1)


    // "ls oss://bucket/prefix"
    objects = s.listObjects(bucket, "lstest:#", false, false, false, false, c)
    c.Assert(len(objects), Equals, num + 2*num1)

    // "ls oss://bucket/prefix -d"
    objects = s.listObjects(bucket, "lstest:#", false, true, false, false, c)
    c.Assert(len(objects), Equals, num + num1)

    objects = s.listObjects(bucket, "lstest:#1/", false, true, false, false, c)
    c.Assert(len(objects), Equals, 2)
}

func (s *OssutilCommandSuite) TestErrList(c *C) {
    showElapse, err := s.rawList([]string{"../"}, true, false, false, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    bucket := bucketNameNotExist 
    showElapse, err = s.rawList([]string{CloudURLToString(bucket, "")}, false, true, false, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // list buckets with -d
    showElapse, err = s.rawList([]string{"oss://"}, false, true, false, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestListIDKey(c *C) {
    bucket := bucketNamePrefix + "lsidkey"

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "ls"
    str := ""
    args := []string{"oss://"}
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    options = OptionMapType{
        "endpoint": &endpoint,
        "accessKeyID": &accessKeyID,
        "accessKeySecret": &accessKeySecret,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestListBucketIDKey(c *C) {
    bucket := bucketNameExist 

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "ls"
    str := ""
    args := []string{CloudURLToString(bucket, "")}
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    options = OptionMapType{
        "endpoint": &endpoint,
        "accessKeyID": &accessKeyID,
        "accessKeySecret": &accessKeySecret,
        "stsToken": &str,
        "configFile": &cfile,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
}

// list multipart 
func (s *OssutilCommandSuite) TestListMultipartObjects(c *C) {
    
    bucketName := bucketNameDest
    // "rm -arf oss://bucket/"
    command := "rm"
    args := []string{CloudURLToString(bucketName, "")}
    str := ""
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "recursive": &ok,
        "force": &ok,
        "allType": &ok,
    }
    _, e := cm.RunCommand(command, args, options)
    c.Assert(e, IsNil)

    object := "TestMultipartObjectLs"
    s.putObject(bucketName, object, uploadFileName, c)
    time.Sleep(5*sleepTime)

    // list object
    objects := s.listObjects(bucketName, object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, object)
		
	bucket, err := copyCommand.command.ossBucket(bucketName)
	
    lmr_origin, e := bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    
    for i := 0; i < 20; i++ {
        _, err = bucket.InitiateMultipartUpload(object)
        c.Assert(err, IsNil)
    }

    time.Sleep(2*sleepTime)
	lmr, e := bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 20 + len(lmr_origin.Uploads))

    // list multipart: ls oss://bucket/object
    objects = s.listObjects(bucketName, object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, object)

    // list multipart: ls -m oss://bucket/object
    objects = s.listObjects(bucketName, object, false, false, true, false, c)
    c.Assert(len(objects), Equals, 20)

    // list all type object: ls -a oss://bucket/object
    objects = s.listObjects(bucketName, object, false, false, false, true, c)
    c.Assert(len(objects), Equals, 21)

    // list multipart: ls -am oss://bucket/object
    objects = s.listObjects(bucketName, object, false, false, true, true, c)
    c.Assert(len(objects), Equals, 21)

    // list multipart: ls -ms oss://bucket/object
    objects = s.listObjects(bucketName, object, false, false, true, false, c)
    c.Assert(len(objects), Equals, 20)

    // list multipart: ls -as oss://bucket/object
    objects = s.listObjects(bucketName, object, false, false, true, true, c)
    c.Assert(len(objects), Equals, 21)

	lmr, e = bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 20 + len(lmr_origin.Uploads))
}

