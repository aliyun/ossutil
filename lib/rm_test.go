package lib 

import (
    "fmt"
    "os"
    "time"

    . "gopkg.in/check.v1"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (s *OssutilCommandSuite) rawRemove(args []string, recursive, force, bucket bool) (bool, error) {
    command := "rm"
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "recursive": &recursive,
        "force": &force,
        "bucket": &bucket,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) TestAllTypeObject(c *C) {
    bucketName := bucketNameMB
    // put object
    normal_object := "TestAllTypeObject"
    s.putObject(bucketName, normal_object, uploadFileName, c)
    time.Sleep(2*sleepTime)

    // put object
    object := "TestMultipartObject"
    s.putObject(bucketName, object, uploadFileName, c)
    time.Sleep(2*sleepTime)

    // list object
    objects := s.listObjects(bucketName, object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, object)
		
	bucket, err := copyCommand.command.ossBucket(bucketName)
    for i := 0; i < 20; i++ {
        _, err = bucket.InitiateMultipartUpload(object)
        c.Assert(err, IsNil)
    }

	lmr, e := bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 20)

    command := "rm"
    args := []string{CloudURLToString(bucketName, object)}
    str := ""
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "bucket": &ok,
        "force": &ok,
        "allType": &ok,
    }
    _, e = cm.RunCommand(command, args, options)
    c.Assert(e, IsNil)
    
	lmr, e = bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 0)

    // list normal_object
    objects = s.listObjects(bucketName, normal_object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, normal_object)
}

func (s *OssutilCommandSuite) TestMultipartObject(c *C) {
    bucketName := bucketNameMB
    // put object
    object := "TestMultipartObject"
    s.putObject(bucketName, object, uploadFileName, c)
    time.Sleep(2*sleepTime)

    // list object
    objects := s.listObjects(bucketName, object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, object)
		
	bucket, err := copyCommand.command.ossBucket(bucketName)
    for i := 0; i < 20; i++ {
        _, err = bucket.InitiateMultipartUpload(object)
        c.Assert(err, IsNil)
    }

	lmr, e := bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 20)

    command := "rm"
    args := []string{CloudURLToString(bucketName, object)}
    str := ""
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "bucket": &ok,
        "force": &ok,
        "multipart": &ok,
    }
    _, e = cm.RunCommand(command, args, options)
    c.Assert(e, IsNil)
    
	lmr, e = bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 0)
}

func (s *OssutilCommandSuite) TestMultipartObject_Prefix(c *C) {
    bucketName := bucketNameMB
    // put object
    object := "TestMultipartObject"
    s.putObject(bucketName, object, uploadFileName, c)
    time.Sleep(2*sleepTime)

    object1 := "TestMultipartObject" + "prefix" 
    s.putObject(bucketName, object1, uploadFileName, c)
    time.Sleep(2*sleepTime)

    object2 := "TestMultipartObject" + "/dir/test" 
    s.putObject(bucketName, object2, uploadFileName, c)
    time.Sleep(2*sleepTime)

    // list object
    objects := s.listObjects(bucketName, object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 3)
		
	bucket, err := copyCommand.command.ossBucket(bucketName)
    for i := 0; i < 20; i++ {
        _, err = bucket.InitiateMultipartUpload(object)
        c.Assert(err, IsNil)
    }

    for i := 0; i < 20; i++ {
        _, err = bucket.InitiateMultipartUpload(object1)
        c.Assert(err, IsNil)
    }

    for i := 0; i < 20; i++ {
        _, err = bucket.InitiateMultipartUpload(object2)
        c.Assert(err, IsNil)
    }

	lmr, e := bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 20*3)

    command := "rm"
    args := []string{CloudURLToString(bucketName, object)}
    str := ""
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "bucket": &ok,
        "recursive": &ok,
        "force": &ok,
        "multipart": &ok,
    }
    _, e = cm.RunCommand(command, args, options)
    c.Assert(e, IsNil)
    
	lmr, e = bucket.ListMultipartUploads(oss.Prefix(object))
	c.Assert(e, IsNil)
    c.Assert(len(lmr.Uploads), Equals, 0)
}

func (s *OssutilCommandSuite) TestRemoveObject(c *C) {
    bucket := bucketNameMB

    // put object
    object := "TestRemoveObject"
    s.putObject(bucket, object, uploadFileName, c)
    time.Sleep(2*sleepTime)

    // list object
    objects := s.listObjects(bucket, object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, object)

    // remove object
    s.removeObjects(bucket, object, false, true, c)

    // list object
    objects = s.listObjects(bucket, object, false, false, false, false, c)
    c.Assert(len(objects), Equals, 0)
}

func (s *OssutilCommandSuite) TestRemoveObjects(c *C) {
    bucket := bucketNamePrefix + "rmb1" 
    s.putBucket(bucket, c)
    time.Sleep(14*time.Second) 

    // put object
    num := 2
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("remove%d", i) 
        s.putObject(bucket, object, uploadFileName, c) 
        objectNames = append(objectNames, object)
    }
    time.Sleep(2*sleepTime)

    command := "rm"
    args := []string{CloudURLToString(bucket, "")}
    str := ""
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "bucket": &ok,
        "force": &ok,
    }
    _, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    // list object
    objects := s.listObjects(bucket, "", false, false, false, false, c)
    c.Assert(len(objects), Equals, num)

    // "rm oss://bucket/ -r"
    // remove object
    s.removeObjects(bucket, "", true, false, c)

    objects = s.listObjects(bucket, "", false, false, false, false, c)
    c.Assert(len(objects), Equals, num)

    // "rm oss://bucket/prefix -r -f"
    // remove object
    s.removeObjects(bucket, "re", true, true, c)
    time.Sleep(3*sleepTime)

    // list object
    objects = s.listObjects(bucket, "", false, false, false, false, c)
    c.Assert(len(objects), Equals, 0)

    //reput objects and delete bucket
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("remove%d", i) 
        s.putObject(bucket, object, uploadFileName, c) 
    }
    
    // list buckets
    buckets := s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) != -1, Equals, true)

    // error remove bucket with config 
    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    options = OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
        "recursive": &ok,
        "bucket": &ok,
        "force": &ok,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    options = OptionMapType{
        "endpoint": &endpoint,
        "accessKeyID": &accessKeyID,
        "accessKeySecret": &accessKeySecret,
        "stsToken": &str,
        "configFile": &cfile,
        "recursive": &ok,
        "bucket": &ok,
        "force": &ok,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
    time.Sleep(2*7*time.Second) 

    // list buckets
    buckets = s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) == -1, Equals, true)
}

func (s *OssutilCommandSuite) TestRemoveObjectBucketOption(c *C) {
    bucket := bucketNameExist 

    object := "test_object"
    command := "rm"
    args := []string{CloudURLToString(bucket, object)}
    str := ""
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "bucket": &ok,
        "force": &ok,
    }
    _, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    // list buckets
    buckets := s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) != -1, Equals, true)
}

func (s *OssutilCommandSuite) TestErrRemove(c *C) {
    bucket := bucketNameExist 

    showElapse, err := s.rawRemove([]string{"oss://"}, false, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawRemove([]string{"./"}, false, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawRemove([]string{CloudURLToString(bucket, "")}, false, true, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    object := "/object"
    showElapse, err = s.rawRemove([]string{CloudURLToString(bucket, object)}, false, true, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // remove bucket without force
    showElapse, err = s.rawRemove([]string{CloudURLToString(bucket, "")}, false, false, true)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    bucketStat := s.getStat(bucket, "", c)
    c.Assert(bucketStat[StatName], Equals, bucket)

    // batch delete not exist objects
    showElapse, err = s.rawRemove([]string{CloudURLToString(bucket, object)}, true, true, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // clear not exist bucket
    bucketName := bucketNamePrefix + "rmnotexist"
    showElapse, err = s.rawRemove([]string{CloudURLToString(bucketName, "")}, true, true, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrDeleteObject(c *C) {
    bucketName := bucketNameNotExist 

    bucket, err := removeCommand.command.ossBucket(bucketName)
    c.Assert(err, IsNil)

    object := "object"
    err = removeCommand.ossDeleteObjectRetry(bucket, object)
    c.Assert(err, NotNil)
    
    _, err = removeCommand.ossBatchDeleteObjectsRetry(bucket, []string{object})
    c.Assert(err, NotNil)
}

