package lib 

import (
    "fmt"
    "os"
    "time"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) putBucketWithACL(bucket string, acl string) (bool, error) {
    args := []string{CloudURLToString(bucket, "")}
    showElapse, err := s.rawPutBucketWithACL(args, acl)
    return showElapse, err
}

func (s *OssutilCommandSuite) rawPutBucketWithACL(args []string, acl string) (bool, error) {
    command := "mb"
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "acl": &acl,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) rawPutBucketWithACLLanguage(args []string, acl, language string) (bool, error) {
    command := "mb"
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "acl": &acl,
        "language": &language,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) TestMakeBucket(c *C) {
    bucket := bucketNameMB 

    // put bucket already exists
    s.putBucket(bucket, c)

    // get bucket stat 
    bucketStat := s.getStat(bucketNameDest, "", c) 
    c.Assert(bucketStat[StatName], Equals, bucketNameDest)
    c.Assert(bucketStat[StatACL], Equals, "private")

    // put bucket with ACL
    for _, acl := range []string{"public-read-write"} {
        showElapse, err := s.putBucketWithACL(bucket, acl)
        c.Assert(err, IsNil)
        c.Assert(showElapse, Equals, true)
        time.Sleep(3*7*time.Second)

        bucketStat := s.getStat(bucket, "", c) 
        c.Assert(bucketStat[StatName], Equals, bucket)
        c.Assert(bucketStat[StatACL], Equals, acl)
    }
}

func (s *OssutilCommandSuite) TestMakeBucketErrorName(c *C) {
    for _, bucket := range []string{"中文测试", "a"} {
        command := "mb"
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

        showElapse, err = s.rawGetStat(bucket, "")
        c.Assert(err, NotNil)
        c.Assert(showElapse, Equals, false)
    }
}

func (s *OssutilCommandSuite) TestMakeBucketErrorACL(c *C) {
    bucket := bucketNamePrefix + "mb1" 
    for _, language := range []string{DefaultLanguage, EnglishLanguage, LEnglishLanguage, "unknown"} {
        for _, acl := range []string{"default", "def", "erracl"} {
            showElapse, err := s.rawPutBucketWithACLLanguage([]string{CloudURLToString(bucket, "")}, acl, language)
            c.Assert(err, NotNil)
            c.Assert(showElapse, Equals, false)
            time.Sleep(7*time.Second)

            showElapse, err = s.rawGetStat(bucket, "")
            c.Assert(err, NotNil)
            c.Assert(showElapse, Equals, false)
        }
    }
}

func (s *OssutilCommandSuite) TestMakeBucketErrorOption(c *C) {
    bucket := bucketNamePrefix + "mb2"
    command := "mb"
    args := []string{CloudURLToString(bucket, "")}
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
    bucket := bucketNamePrefix + "assembleoptions"

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s", "abc", "def", "ghi", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "mb"
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

    s.removeBucket(bucket, true, c)
}
