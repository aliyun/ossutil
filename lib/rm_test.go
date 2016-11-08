package lib 

import (
    "fmt"
    "os"
    "time"

    . "gopkg.in/check.v1"
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

func (s *OssutilCommandSuite) TestRemoveObject(c *C) {
    s.SetUpBucketEnv(c)
    bucket := bucketNamePrefix + "rmb" 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime) 

    // put object
    object := "test_object"
    s.putObject(bucket, object, uploadFileName, c)

    // list object
    objects := s.listObjects(bucket, "", false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, object)

    // remove object
    s.removeObjects(bucket, object, false, true, c)

    // list object
    objects = s.listObjects(bucket, "", false, false, c)
    c.Assert(len(objects), Equals, 0)

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 
}

func (s *OssutilCommandSuite) TestRemoveObjects(c *C) {
    bucket := bucketNamePrefix + "rmb" 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime) 

    // put object
    num := 15
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("^$%d$^工具", i) 
        s.putObject(bucket, object, uploadFileName, c) 
        objectNames = append(objectNames, object)
    }

    // list object
    objects := s.listObjects(bucket, "", false, false, c)
    c.Assert(len(objects), Equals, num)

    // "rm oss://bucket/ -r"
    // remove object
    s.removeObjects(bucket, "", true, false, c)

    objects = s.listObjects(bucket, "", false, false, c)
    c.Assert(len(objects), Equals, num)

    // "rm oss://bucket/prefix -r -f"
    // remove object
    s.removeObjects(bucket, "^$", true, true, c)

    // list object
    objects = s.listObjects(bucket, "", false, false, c)
    c.Assert(len(objects), Equals, 0)

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 
}

func (s *OssutilCommandSuite) TestRemoveBucket(c *C) {
    bucket := bucketNamePrefix + "rmb"
    s.putBucket(bucket, c)
    time.Sleep(sleepTime) 

    // put object
    num := 15
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("^$%d$^工具", i) 
        s.putObject(bucket, object, uploadFileName, c) 
        objectNames = append(objectNames, object)
    }
    
    // list buckets
    buckets := s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) != -1, Equals, true)

    // rm bucket
    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 

    // list buckets
    buckets = s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) == -1, Equals, true)
}

func (s *OssutilCommandSuite) TestRemoveEmptyBucket(c *C) {
    bucket := bucketNamePrefix + "rmb"
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

    // list buckets
    buckets := s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) != -1, Equals, true)

    // rm bucket
    s.removeBucket(bucket, false, c)
    time.Sleep(sleepTime)

    // list buckets
    buckets = s.listBuckets(false, c)
    c.Assert(FindPos(bucket, buckets) == -1, Equals, true)
}

func (s *OssutilCommandSuite) TestRemoveNonEmptyBucket(c *C) {
    bucket := bucketNamePrefix + "rmb" 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

    // put object
    object := "test_object_for_rm"
    s.putObject(bucket, object, uploadFileName, c)

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
    objects := s.listObjects(bucket, "", false, false, c)
    c.Assert(len(objects), Equals, 1)
    c.Assert(objects[0], Equals, object)

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 
}

func (s *OssutilCommandSuite) TestRemoveObjectBucketOption(c *C) {
    bucket := bucketNamePrefix + "rmb"
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

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

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 
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

func (s *OssutilCommandSuite) TestRemoveIDKey(c *C) {
    bucket := bucketNamePrefix + "rmidkey"
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "rm"
    str := ""
    ok := true
    args := []string{CloudURLToString(bucket, "")}
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
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
        "bucket": &ok,
        "force": &ok,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime)
}
