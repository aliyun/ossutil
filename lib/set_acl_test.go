package lib 

import (
    "fmt"
    "os"
    "strconv"
    "time"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) rawSetBucketACL(bucket, acl string, force bool) (bool, error) {
    args := []string{CloudURLToString(bucket, ""), acl}
    showElapse, err := s.rawSetACLWithArgs(args, false, true, force)
    return showElapse, err
}

func (s *OssutilCommandSuite) rawSetObjectACL(bucket, object, acl string, recursive, force bool) (bool, error) {
    args := []string{CloudURLToString(bucket, object), acl}
    showElapse, err := s.rawSetACLWithArgs(args, recursive, false, force)
    return showElapse, err
}

func (s *OssutilCommandSuite) rawSetACLWithArgs(args []string, recursive, bucket, force bool) (bool, error) {
    command := "set-acl"
    str := ""
    routines := strconv.Itoa(Routines)
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "routines": &routines,
        "recursive": &recursive,
        "bucket": &bucket,
        "force": &force,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    time.Sleep(2*sleepTime)
    return showElapse, err
}

func (s *OssutilCommandSuite) TestSetBucketACL(c *C) {
    bucket := bucketNameExist 

    // set acl
    for _, acl := range []string{"private", "public-read", "public-read-write"} {
        s.setBucketACL(bucket, acl, c)
        bucketStat := s.getStat(bucket, "", c)
        c.Assert(bucketStat[StatACL], Equals, acl)
    }
}

func (s *OssutilCommandSuite) TestSetBucketErrorACL(c *C) {
    bucket := bucketNameDest 

    for _, acl := range []string{"default", "def", "erracl", "私有"} {
        showElapse, err := s.rawSetBucketACL(bucket, acl, false)
        c.Assert(err, NotNil)
        c.Assert(showElapse, Equals, false)

        showElapse, err = s.rawSetBucketACL(bucket, acl, true)
        c.Assert(err, NotNil)
        c.Assert(showElapse, Equals, false)

        bucketStat := s.getStat(bucket, "", c)
        c.Assert(bucketStat[StatACL], Equals, "private")
    }
}

func (s *OssutilCommandSuite) TestSetNotExistBucketACL(c *C) {
    bucket := bucketNamePrefix + "noexistsetacl" 

    showElapse, err := s.rawGetStat(bucket, "")
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // set acl and create bucket
    showElapse, err = s.rawSetBucketACL(bucket, "public-read", true)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    time.Sleep(sleepTime)
    bucketStat := s.getStat(bucket, "", c)
    c.Assert(bucketStat[StatACL], Equals, "public-read")

    s.removeBucket(bucket, true, c)
    time.Sleep(3*sleepTime)

    // invalid bucket name
    bucket = "a"
    showElapse, err = s.rawSetBucketACL(bucket, "public-read", true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawGetStat(bucket, "")
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestSetBucketEmptyACL(c *C) {
    bucket := bucketNamePrefix + "acl3" 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

    object := "test"
    s.putObject(bucket, object, uploadFileName, c)

    showElapse, err := s.rawSetObjectACL(bucket, object, "", false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime)
}

func (s *OssutilCommandSuite) TestSetObjectACL(c *C) {
    bucket := bucketNameSetACL 
    s.removeObjects(bucket, "", true, true, c)
    time.Sleep(3*sleepTime)

    object := "TestSetObjectACL"

    // set acl to not exist object
    showElapse, err := s.rawSetObjectACL(bucket, object, "default", false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    object = "setacl-oldobject"
    s.putObject(bucket, object, uploadFileName, c)

    //get object acl
    objectStat := s.getStat(bucket, object, c)
    c.Assert(objectStat[StatACL], Equals, "default")

    object = "setacl-newobject"
    s.putObject(bucket, object, uploadFileName, c)

    // set acl
    for _, acl := range []string{"default", "private", "public-read", "public-read-write"} {
        s.setObjectACL(bucket, object, acl, false, true, c)
        objectStat = s.getStat(bucket, object, c)
        c.Assert(objectStat[StatACL], Equals, acl)
    }

    s.setObjectACL(bucket, object, "private", false, true, c)

    // set error acl
    for _, acl := range []string{"public_read", "erracl", "私有"} {
        showElapse, err = s.rawSetObjectACL(bucket, object, acl, false, false)
        c.Assert(showElapse, Equals, false)
        c.Assert(err, NotNil)

        objectStat = s.getStat(bucket, object, c)
        c.Assert(objectStat[StatACL], Equals, "private")
    }
}

func (s *OssutilCommandSuite) TestBatchSetObjectACL(c *C) {
    bucket := bucketNameSetACL1 
    s.removeObjects(bucket, "", true, true, c)
    time.Sleep(3*sleepTime)

    // put objects
    num := 2 
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestBatchSetObjectACL_setacl%d", i)
        s.putObject(bucket, object, uploadFileName, c)
        objectNames = append(objectNames, object)
    }
    time.Sleep(time.Second)

    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c)
        c.Assert(objectStat[StatACL], Equals, "default")
    }

    // without --force option
    s.setObjectACL(bucket, "", "public-read-write", true, false, c)

    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c)
        c.Assert(objectStat[StatACL], Equals, "default")
    }

    for _, acl := range []string{"public-read", "private", "public-read-write", "default"} {
        s.setObjectACL(bucket, "TestBatchSetObjectACL_setacl", acl, true, true, c)
        time.Sleep(sleepTime)

        for _, object := range objectNames {
            objectStat := s.getStat(bucket, object, c)
            c.Assert(objectStat[StatACL], Equals, acl)
        }
    }

    showElapse, err := s.rawSetObjectACL(bucket, "TestBatchSetObjectACL_setacl", "erracl", true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrSetACL(c *C) {
    acl := "private"
    args := []string{"os://", acl}
    showElapse, err := s.rawSetACLWithArgs(args, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    args = []string{"oss://", acl} 
    showElapse, err = s.rawSetACLWithArgs(args, false, false, true)

    // error set bucket acl
    bucket := bucketNamePrefix + "set-acl"
    object := "testobject"
    args = []string{CloudURLToString(bucket, object), acl}
    showElapse, err = s.rawSetACLWithArgs(args, false, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    args = []string{CloudURLToString(bucket, ""), acl}
    showElapse, err = s.rawSetACLWithArgs(args, true, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss acl
    args = []string{CloudURLToString(bucket, "")}
    showElapse, err = s.rawSetACLWithArgs(args, false, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    args = []string{CloudURLToString(bucket, object)}
    showElapse, err = s.rawSetACLWithArgs(args, false, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    args = []string{CloudURLToString(bucket, object)}
    showElapse, err = s.rawSetACLWithArgs(args, true, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss object
    args = []string{CloudURLToString(bucket, ""), acl}
    showElapse, err = s.rawSetACLWithArgs(args, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // bad prefix
    showElapse, err = s.rawSetObjectACL(bucket, "/object", acl, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrBatchSetACL(c *C) {
    bucket := bucketNameExist  

    // put objects
    num := 10
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestErrBatchSetACL_setacl:%d", i)
        s.putObject(bucket, object, uploadFileName, c)
        objectNames = append(objectNames, object)
    }

    command := "set-acl"
    str := ""
    str1 := "abc"
    args := []string{CloudURLToString(bucket, ""), "public-read-write"}
    routines := strconv.Itoa(Routines)
    ok := true
    options := OptionMapType{
        "endpoint": &str1,
        "accessKeyID": &str1,
        "accessKeySecret": &str1,
        "stsToken": &str,
        "routines": &routines,
        "recursive": &ok,
        "force": &ok,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c)
        c.Assert(objectStat[StatACL], Equals, "default")
    }
}

func (s *OssutilCommandSuite) TestSetACLIDKey(c *C) {
    bucket := bucketNameExist 

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "set-acl"
    str := ""
    args := []string{CloudURLToString(bucket, ""), "public-read"}
    routines := strconv.Itoa(Routines)
    ok := true
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
        "routines": &routines,
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
        "routines": &routines,
        "bucket": &ok,
        "force": &ok,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(cfile)
}
