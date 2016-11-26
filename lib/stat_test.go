package lib 

import (
    "fmt"
    "os"
    "time"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) rawGetStat(bucket, object string) (bool, error) {
    args := []string{CloudURLToString(bucket, object)}
    showElapse, err := s.rawGetStatWithArgs(args)
    return showElapse, err 
}

func (s *OssutilCommandSuite) rawGetStatWithArgs(args []string) (bool, error) {
    command := "stat"
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err 
}

func (s *OssutilCommandSuite) TestStatErrArgc(c *C) {
    bucket := bucketNameExist 

    command := "stat"
    args := []string{CloudURLToString(bucket, ""), CloudURLToString(bucket, "")}
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

func (s *OssutilCommandSuite) TestGetBucketStat(c *C) {
    bucket := bucketNameDest 

    // get bucket stat 
    bucketStat := s.getStat(bucket, "", c) 
    c.Assert(bucketStat[StatName], Equals, bucket)
    c.Assert(bucketStat[StatLocation] != "", Equals, true)
    c.Assert(bucketStat[StatCreationDate] != "", Equals, true)
    c.Assert(bucketStat[StatExtranetEndpoint] != "", Equals, true)
    c.Assert(bucketStat[StatIntranetEndpoint] != "", Equals, true)
    c.Assert(bucketStat[StatACL], Equals, "private")
    c.Assert(bucketStat[StatOwner] != "", Equals, true)
}

func (s *OssutilCommandSuite) TestGetObjectStat(c *C) {
    bucket := bucketNameExist 

    object := "TestGetObjectStat"
    s.putObject(bucket, object, uploadFileName, c)

    objectStat := s.getStat(bucket, object, c)
    c.Assert(objectStat[StatACL], Equals, "default")
    c.Assert(len(objectStat["Etag"]), Equals, 32)
    c.Assert(objectStat["Last-Modified"] != "", Equals, true)
    c.Assert(objectStat[StatOwner] != "", Equals, true)
}

func (s *OssutilCommandSuite) TestGetStatNotExist(c *C) {
    bucket := bucketNameNotExist 
    showElapse, err := s.rawGetStat(bucket, "")
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    bucket = bucketNameExist
    s.removeObjects(bucket, "", true, true, c)
    time.Sleep(sleepTime)

    showElapse, err = s.rawGetStat(bucket, "")
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    object := "testobject_for_getstat_not_exist"
    showElapse, err = s.rawGetStat(bucket, object)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    object = "testobject_exist"
    s.putObject(bucket, object, uploadFileName, c)
    showElapse, err = s.rawGetStat(bucket, object)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestGetStatRetryTimes(c *C) {
    bucket := bucketNameExist 

    command := "stat"
    args := []string{CloudURLToString(bucket, "")}
    str := ""
    retryTimes := "1"
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "retryTimes": &retryTimes,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestGetStatErrSrc(c *C) {
    showElapse, err := s.rawGetStat("", "")
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawGetStatWithArgs([]string{"../"})
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestStatIDKey(c *C) {
    bucket := bucketNameExist 

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "stat"
    str := ""
    args := []string{CloudURLToString(bucket, "")}
    retryTimes := "1"
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
        "retryTimes": &retryTimes,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    options = OptionMapType{
        "endpoint": &endpoint,
        "accessKeyID": &accessKeyID,
        "accessKeySecret": &accessKeySecret,
        "stsToken": &str,
        "configFile": &cfile,
        "retryTimes": &retryTimes,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
}
