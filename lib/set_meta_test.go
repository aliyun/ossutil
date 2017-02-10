package lib 

import (
    "fmt"
    "os"
    "strconv"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestSetBucketMeta(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    showElapse, err := s.rawSetMeta(bucketName, "", "X-Oss-Meta-A:A", false, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestSetObjectMeta(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    object := "TestSetObjectMeta_testobject" 
    s.putObject(bucketName, object, uploadFileName, c)

    objectStat := s.getStat(bucketName, object, c) 
    c.Assert(objectStat[StatACL], Equals, "default") 
    _, ok := objectStat["X-Oss-Meta-A"]
    c.Assert(ok, Equals, false)

    // update
    s.setObjectMeta(bucketName, object, "x-oss-object-acl:private#X-Oss-Meta-A:A#Expires:2006-01-02T15:04:05Z", true, false, false, true, c)

    objectStat = s.getStat(bucketName, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")
    c.Assert(objectStat["Expires"], Equals, "Mon, 02 Jan 2006 15:04:05 GMT")

    // error expires
    showElapse, err := s.rawSetMeta(bucketName, object, "Expires:2006-01", true, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    objectStat = s.getStat(bucketName, object, c) 
    c.Assert(objectStat["Expires"], Equals, "Mon, 02 Jan 2006 15:04:05 GMT")

    // delete
    s.setObjectMeta(bucketName, object, "x-oss-object-acl#X-Oss-Meta-A", false, true, false, true, c)
    objectStat = s.getStat(bucketName, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    _, ok = objectStat["X-Oss-Meta-A"]
    c.Assert(ok, Equals, false)

    s.setObjectMeta(bucketName, object, "X-Oss-Meta-A:A#x-oss-meta-B:b", true, false, false, true, c)

    s.setObjectMeta(bucketName, object, "X-Oss-Meta-c:c", false, false, false, true, c)
    objectStat = s.getStat(bucketName, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 

    // without force
    s.setObjectMeta(bucketName, object, "x-oss-object-acl:public-read#X-Oss-Meta-A:A", true, false, false, false, c)

    // without update, delete and force
    showElapse, err = s.rawSetMeta(bucketName, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", false, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucketName, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", false, false, false, false, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucketName, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", false, false, false, false, LEnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss meta
    s.setObjectMeta(bucketName, object, "", true, false, false, true, c)

    showElapse, err = s.rawSetMeta(bucketName, object, "", true, false, false, true, EnglishLanguage)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    showElapse, err = s.rawSetMeta(bucketName, object, "", true, false, false, true, LEnglishLanguage)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // delete error meta
    showElapse, err = s.rawSetMeta(bucketName, object, "X-Oss-Meta-A:A", false, true, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // update error meta
    showElapse, err = s.rawSetMeta(bucketName, object, "a:b", true, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucketName, object, "x-oss-object-acl:private", true, false, false, true, DefaultLanguage)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    //batch
    s.setObjectMeta(bucketName, "", "content-type:abc#X-Oss-Meta-Update:update", true, false, true, false, c)

    s.setObjectMeta(bucketName, "", "content-type:abc#X-Oss-Meta-Update:update", true, false, true, true, c)

    objectStat = s.getStat(bucketName, object, c)
    c.Assert(objectStat["Content-Type"], Equals, "abc")
    c.Assert(objectStat["X-Oss-Meta-Update"], Equals, "update")

    s.setObjectMeta(bucketName, "", "X-Oss-Meta-update", false, true, true, true, c)

    s.setObjectMeta(bucketName, "", "X-Oss-Meta-A:A#x-oss-meta-B:b", true, false, true, true, c)

    s.setObjectMeta(bucketName, "nosetmeta", "X-Oss-Meta-M:c", false, false, true, true, c)

    s.setObjectMeta(bucketName, "", "X-Oss-Meta-C:c", false, false, true, true, c)

    objectStat = s.getStat(bucketName, object, c)
    c.Assert(objectStat["X-Oss-Meta-C"], Equals, "c")

    showElapse, err = s.rawSetMeta(bucketName, "", "X-Oss-Meta-c:c", false, true, true, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucketName, "", "a:b", true, false, true, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestSetNotExistObjectMeta(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    object := "testobject-notexistone" 
    // set meta of not exist object
    showElapse, err := s.rawSetMeta(bucketName, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // batch set meta of not exist objects
    s.setObjectMeta(bucketName, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, true, true, c)

    showElapse, err = s.rawGetStat(bucketName, object)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    object = "testobject"
    s.putObject(bucketName, object, uploadFileName, c)

    s.setObjectMeta(bucketName, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, false, true, c)

    objectStat := s.getStat(bucketName, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")

    s.removeBucket(bucketName, true, c)
}
    
func (s *OssutilCommandSuite) TestErrBatchSetMeta(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    // put objects
    num := 10
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestErrBatchSetMeta_setmeta:%d", i)
        s.putObject(bucketName, object, uploadFileName, c)
        objectNames = append(objectNames, object)
    }

    // update without force
    meta := "content-type:abc#X-Oss-Meta-Update:update" 
    args := []string{CloudURLToString(bucketName, ""), meta}
    command := "set-meta"
    str := ""
    str1 := "abc"
    ok := true
    routines := strconv.Itoa(Routines)
    options := OptionMapType{
        "endpoint": &str1,
        "accessKeyID": &str1,
        "accessKeySecret": &str1,
        "stsToken": &str,
        "update": &ok,
        "recursive": &ok,
        "force": &ok,
        "routines": &routines,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    for _, object := range objectNames {
        objectStat := s.getStat(bucketName, object, c) 
        c.Assert(objectStat["Content-Type"] != "abc", Equals, true) 
        _, ok := objectStat["X-Oss-Meta-Update"]
        c.Assert(ok, Equals, false)
    }

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestErrSetMeta(c *C) {
    args := []string{"os:/", ""}
    showElapse, err := s.rawSetMetaWithArgs(args, false, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta("", "", "", false, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    object := "notexistobject"

    showElapse, err = s.rawSetMeta(bucketName, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, true, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucketName, object)}, false, false, false, true, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucketName, object)}, true, false, false, false, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucketName, object)}, true, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucketName, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", false, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{"oss:///object", "x-oss-object-acl:private#X-Oss-Meta-A:A"}, false, false, true, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString("", "")}, true, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucketName, object, "unknown:a", true, false, false, true, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucketName, object, "Expires:a", true, false, false, true, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestGetOSSOption(c *C) {
    _, err := getOSSOption("unknown", "a")
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestSetMetaIDKey(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    object := "testobject" 
    s.putObject(bucketName, object, uploadFileName, c)

    cfile := randStr(10) 
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucketName, "abc", bucketName, "abc") 
    s.createFile(cfile, data, c)

    command := "set-meta"
    str := ""
    args := []string{CloudURLToString(bucketName, object), "x-oss-object-acl:private#X-Oss-Meta-A:A#Expires:2006-01-02T15:04:05Z"}
    ok := true
    routines := strconv.Itoa(Routines)
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
        "update": &ok,
        "force": &ok,
        "routines": &routines,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    options = OptionMapType{
        "endpoint": &endpoint,
        "accessKeyID": &accessKeyID,
        "accessKeySecret": &accessKeySecret,
        "stsToken": &str,
        "configFile": &cfile,
        "update": &ok,
        "force": &ok,
        "routines": &routines,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    os.Remove(cfile)

    s.removeBucket(bucketName, true, c)
}
