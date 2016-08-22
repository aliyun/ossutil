package lib 

import (
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
    bucket := bucketNamePrefix + "mb" 
    s.putBucket(bucket, c)

    // put bucket already exists
    s.putBucket(bucket, c)

    // get bucket stat 
    bucketStat := s.getStat(bucket, "", c) 
    c.Assert(bucketStat[StatName], Equals, bucket)
    c.Assert(bucketStat[StatACL], Equals, "private")

    // put bucket with ACL
    for _, acl := range []string{"private", "public-read", "public-read-write"} {
        showElapse, err := s.putBucketWithACL(bucket, acl)
        c.Assert(err, IsNil)
        c.Assert(showElapse, Equals, true)

        bucketStat := s.getStat(bucket, "", c) 
        c.Assert(bucketStat[StatName], Equals, bucket)
        c.Assert(bucketStat[StatACL], Equals, acl)
    }

    result := []string{"private", "public-read", "public-read-write"}
    for i, str := range []string{"pri", "pr", "prw"} {
        showElapse, err := s.putBucketWithACL(bucket, str)
        c.Assert(err, IsNil)
        c.Assert(showElapse, Equals, true)

        bucketStat := s.getStat(bucket, "", c) 
        c.Assert(bucketStat[StatName], Equals, bucket)
        c.Assert(bucketStat[StatACL], Equals, result[i])
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
    bucket := bucketNamePrefix + "mb" 
    for _, language := range []string{DefaultLanguage, EnglishLanguage, LEnglishLanguage, "unknown"} {
        for _, acl := range []string{"default", "def", "erracl"} {
            showElapse, err := s.rawPutBucketWithACLLanguage([]string{CloudURLToString(bucket, "")}, acl, language)
            c.Assert(err, NotNil)
            c.Assert(showElapse, Equals, false)

            showElapse, err = s.rawGetStat(bucket, "")
            c.Assert(err, NotNil)
            c.Assert(showElapse, Equals, false)
        }
    }
}

func (s *OssutilCommandSuite) TestMakeBucketErrorOption(c *C) {
    bucket := bucketNamePrefix + "mb"
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
