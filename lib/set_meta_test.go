package lib 

import (
    "fmt"
    "os"
    "strconv"
    "time"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) rawSetMeta(bucket, object, meta string, update, delete, recursive, force bool, language string) (bool, error) {
    args := []string{CloudURLToString(bucket, object), meta}
    showElapse, err := s.rawSetMetaWithArgs(args, update, delete, recursive, force, language) 
    return showElapse, err
}

func (s *OssutilCommandSuite) rawSetMetaWithArgs(args []string, update, delete, recursive, force bool, language string) (bool, error) {
    command := "set-meta"
    str := ""
    routines := strconv.Itoa(Routines)
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "update": &update,
        "delete": &delete,
        "recursive": &recursive,
        "force": &force,
        "routines": &routines,
        "language": &language,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    time.Sleep(2*time.Second)
    return showElapse, err
}

func (s *OssutilCommandSuite) setObjectMeta(bucket, object, meta string, update, delete, recursive, force bool, c *C) {
    showElapse, err := s.rawSetMeta(bucket, object, meta, update, delete, recursive, force, DefaultLanguage) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestSetBucketMeta(c *C) {
    bucket := bucketNameExist 

    showElapse, err := s.rawSetMeta(bucket, "", "X-Oss-Meta-A:A", false, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestSetObjectMeta(c *C) {
    bucket := bucketNameExist 
    s.removeObjects(bucket, "", true, true, c)
    time.Sleep(2*sleepTime)

    object := "TestSetObjectMeta_testobject" 
    s.putObject(bucket, object, uploadFileName, c)

    objectStat := s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "default") 
    _, ok := objectStat["X-Oss-Meta-A"]
    c.Assert(ok, Equals, false)

    // update
    s.setObjectMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A#Expires:2006-01-02T15:04:05Z", true, false, false, true, c)

    objectStat = s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")
    c.Assert(objectStat["Expires"], Equals, "Mon, 02 Jan 2006 15:04:05 GMT")

    // error expires
    showElapse, err := s.rawSetMeta(bucket, object, "Expires:2006-01", true, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    objectStat = s.getStat(bucket, object, c) 
    c.Assert(objectStat["Expires"], Equals, "Mon, 02 Jan 2006 15:04:05 GMT")

    // delete
    s.setObjectMeta(bucket, object, "x-oss-object-acl#X-Oss-Meta-A", false, true, false, true, c)
    objectStat = s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    _, ok = objectStat["X-Oss-Meta-A"]
    c.Assert(ok, Equals, false)

    s.setObjectMeta(bucket, object, "X-Oss-Meta-A:A#x-oss-meta-B:b", true, false, false, true, c)
    objectStat = s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")
    c.Assert(objectStat["X-Oss-Meta-B"], Equals, "b")

    s.setObjectMeta(bucket, object, "X-Oss-Meta-c:c", false, false, false, true, c)
    objectStat = s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    _, ok = objectStat["X-Oss-Meta-A"]
    c.Assert(ok, Equals, false)
    _, ok = objectStat["X-Oss-Meta-B"]
    c.Assert(ok, Equals, false)
    c.Assert(objectStat["X-Oss-Meta-C"], Equals, "c")

    // without force
    s.setObjectMeta(bucket, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", true, false, false, false, c)

    objectStat = s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "default") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")

    // without update, delete and force
    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", false, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", false, false, false, false, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", false, false, false, false, LEnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss meta
    s.setObjectMeta(bucket, object, "", true, false, false, true, c)

    showElapse, err = s.rawSetMeta(bucket, object, "", true, false, false, true, EnglishLanguage)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    showElapse, err = s.rawSetMeta(bucket, object, "", true, false, false, true, LEnglishLanguage)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // delete error meta
    showElapse, err = s.rawSetMeta(bucket, object, "X-Oss-Meta-A:A", false, true, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // update error meta
    showElapse, err = s.rawSetMeta(bucket, object, "a:b", true, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private", true, false, false, true, DefaultLanguage)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestSetNotExistObjectMeta(c *C) {
    bucket := bucketNameExist 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

    object := "testobject-notexistone" 
    // set meta of not exist object
    showElapse, err := s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // batch set meta of not exist objects
    s.setObjectMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, true, true, c)

    showElapse, err = s.rawGetStat(bucket, object)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    object = "testobject"
    s.putObject(bucket, object, uploadFileName, c)

    s.setObjectMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, false, true, c)

    objectStat := s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")
}

func (s *OssutilCommandSuite) TestBatchSetObjectMeta(c *C) {
    bucket := bucketNameExist 
    s.removeObjects(bucket, "", true, true, c)
    time.Sleep(2*sleepTime)

    // put objects
    num := 10
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestBatchSetObjectMeta_设置元信息：%d", i)
        s.putObject(bucket, object, uploadFileName, c)
        objectNames = append(objectNames, object)
    }

    // update without force
    s.setObjectMeta(bucket, "", "content-type:abc#X-Oss-Meta-Update:update", true, false, true, false, c)

    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["Content-Type"] != "abc", Equals, true) 
        _, ok := objectStat["X-Oss-Meta-Update"]
        c.Assert(ok, Equals, false)
    }

    // update
    s.setObjectMeta(bucket, "", "content-type:abc#X-Oss-Meta-update:update", true, false, true, true, c)

    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["Content-Type"], Equals, "abc") 
        c.Assert(objectStat["X-Oss-Meta-Update"], Equals, "update")
    }

     // delete
    s.setObjectMeta(bucket, "TestBatchSetObjectMeta_设置元信息：", "X-Oss-Meta-update", false, true, true, true, c)
   
    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        _, ok := objectStat["X-Oss-Meta-Update"]
        c.Assert(ok, Equals, false)
    }

    s.setObjectMeta(bucket, "", "X-Oss-Meta-A:A#x-oss-meta-B:b", true, false, true, true, c)
    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A") 
        c.Assert(objectStat["X-Oss-Meta-B"], Equals, "b")
    }

    // set all
    s.setObjectMeta(bucket, "no设置元信息", "X-Oss-Meta-M:c", false, false, true, true, c)

    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A") 
        c.Assert(objectStat["X-Oss-Meta-B"], Equals, "b")
        _, ok := objectStat["X-Oss-Meta-M"]
        c.Assert(ok, Equals, false)
    }

    s.setObjectMeta(bucket, "TestBatchSetObjectMeta_设置元信息", "X-Oss-Meta-c:c", false, false, true, true, c)
    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["X-Oss-Meta-C"], Equals, "c") 
        _, ok := objectStat["X-Oss-Meta-A"]
        c.Assert(ok, Equals, false)
        _, ok = objectStat["X-Oss-Meta-B"]
        c.Assert(ok, Equals, false)
    }

    // error meta
    showElapse, err := s.rawSetMeta(bucket, "TestBatchSetObjectMeta_设置元信息：", "X-Oss-Meta-c:c", false, true, true, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, "", "a:b", true, false, true, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrBatchSetMeta(c *C) {
    bucket := bucketNameExist 

    // put objects
    num := 10
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestErrBatchSetMeta_设置元信息：%d", i)
        s.putObject(bucket, object, uploadFileName, c)
        objectNames = append(objectNames, object)
    }

    // update without force
    meta := "content-type:abc#X-Oss-Meta-Update:update" 
    args := []string{CloudURLToString(bucket, ""), meta}
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
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["Content-Type"] != "abc", Equals, true) 
        _, ok := objectStat["X-Oss-Meta-Update"]
        c.Assert(ok, Equals, false)
    }

}

func (s *OssutilCommandSuite) TestErrSetMeta(c *C) {
    args := []string{"os:/", ""}
    showElapse, err := s.rawSetMetaWithArgs(args, false, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta("", "", "", false, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    bucket := bucketNameExist 

    object := "notexistobject"

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, true, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucket, object)}, false, false, false, true, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucket, object)}, true, false, false, false, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucket, object)}, true, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", false, false, false, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    object = "/object"
    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", false, false, true, true, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString("", "")}, true, false, false, false, DefaultLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "unknown:a", true, false, false, true, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "Expires:a", true, false, false, true, EnglishLanguage)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestGetOSSOption(c *C) {
    _, err := getOSSOption("unknown", "a")
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestSetMetaIDKey(c *C) {
    bucket := bucketNameExist 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

    object := "testobject" 
    s.putObject(bucket, object, uploadFileName, c)

    cfile := "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "set-meta"
    str := ""
    args := []string{CloudURLToString(bucket, object), "x-oss-object-acl:private#X-Oss-Meta-A:A#Expires:2006-01-02T15:04:05Z"}
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

    _ = os.Remove(cfile)
}
