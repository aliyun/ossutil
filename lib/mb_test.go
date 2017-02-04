package lib 

import (
    "fmt"
    "os"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestMakeBucket(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    // put bucket already exists
    s.putBucket(bucketName, c)

    // get bucket stat 
    bucketStat := s.getStat(bucketName, "", c) 
    c.Assert(bucketStat[StatName], Equals, bucketName)
    c.Assert(bucketStat[StatACL], Equals, "private")

    // put bucket with ACL
    for _, acl := range []string{"public-read-write"} {
        showElapse, err := s.putBucketWithACL(bucketName, acl)
        c.Assert(err, IsNil)
        c.Assert(showElapse, Equals, true)

        bucketStat := s.getStat(bucketName, "", c) 
        c.Assert(bucketStat[StatName], Equals, bucketName)
        c.Assert(bucketStat[StatACL], Equals, acl)
    }

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestMakeBucketErrorName(c *C) {
    for _, bucketName := range []string{"中文测试", "a"} {
        command := "mb"
        args := []string{CloudURLToString(bucketName, "")}
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

        showElapse, err = s.rawGetStat(bucketName, "")
        c.Assert(err, NotNil)
        c.Assert(showElapse, Equals, false)
    }
}

func (s *OssutilCommandSuite) TestMakeBucketErrorACL(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10) 
    for _, language := range []string{DefaultLanguage, EnglishLanguage, LEnglishLanguage, "unknown"} {
        for _, acl := range []string{"default", "def", "erracl"} {
            showElapse, err := s.rawPutBucketWithACLLanguage([]string{CloudURLToString(bucketName, "")}, acl, language)
            c.Assert(err, NotNil)
            c.Assert(showElapse, Equals, false)

            showElapse, err = s.rawGetStat(bucketName, "")
            c.Assert(err, NotNil)
            c.Assert(showElapse, Equals, false)
        }
    }
    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestMakeBucketErrorOption(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    command := "mb"
    args := []string{CloudURLToString(bucketName, "")}
    str := ""
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "shortFormat": &ok, 
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrMakeBucket(c *C) {
    acl := "private"
    showElapse, err := s.rawPutBucketWithACL([]string{"os://"}, acl)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawPutBucketWithACL([]string{CloudURLToString("", "")}, acl)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestMakeBucketIDKey(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s", "abc", "def", "ghi", bucketName, "abc") 
    s.createFile(cfile, data, c)

    command := "mb"
    str := ""
    args := []string{CloudURLToString(bucketName, "")}
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

    s.removeBucket(bucketName, false, c)
}
