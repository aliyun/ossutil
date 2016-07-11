package lib 

import (
    "fmt"
    "strconv"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) rawSetMeta(bucket, object, meta string, update, delete, recursive, force bool) (bool, error) {
    args := []string{CloudURLToString(bucket, object), meta}
    showElapse, err := s.rawSetMetaWithArgs(args, update, delete, recursive, force) 
    return showElapse, err
}

func (s *OssutilCommandSuite) rawSetMetaWithArgs(args []string, update, delete, recursive, force bool) (bool, error) {
    command := "setmeta"
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
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) setObjectMeta(bucket, object, meta string, update, delete, recursive, force bool, c *C) {
    showElapse, err := s.rawSetMeta(bucket, object, meta, update, delete, recursive, force) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestSetBucketMeta(c *C) {
    bucket := bucketNamePrefix + "setmeta"
    s.putBucket(bucket, c)

    showElapse, err := s.rawSetMeta(bucket, "", "X-Oss-Meta-A:A", false, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestSetObjectMeta(c *C) {
    bucket := bucketNamePrefix + "setmeta"
    s.putBucket(bucket, c)

    object := "testobject" 
    s.putObject(bucket, object, uploadFileName, c)

    objectStat := s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "default") 
    _, ok := objectStat["X-Oss-Meta-A"]
    c.Assert(ok, Equals, false)

    // update
    s.setObjectMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, false, true, c)

    objectStat = s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")

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
    showElapse, err := s.rawSetMeta(bucket, object, "x-oss-object-acl:default#X-Oss-Meta-A:A", false, false, false, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss meta
    s.setObjectMeta(bucket, object, "", true, false, false, true, c)

    // delete error meta
    showElapse, err = s.rawSetMeta(bucket, object, "X-Oss-Meta-A:A", false, true, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // update error meta
    showElapse, err = s.rawSetMeta(bucket, object, "a:b", true, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private", true, false, false, true)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestSetNotExistObjectMeta(c *C) {
    bucket := bucketNamePrefix + "setmeta"
    s.putBucket(bucket, c)

    object := "testobject" 
    // set meta of not exist object
    showElapse, err := s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // batch set meta of not exist objects
    s.setObjectMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, true, true, c)

    showElapse, err = s.rawGetStat(bucket, object)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.putObject(bucket, object, uploadFileName, c)

    s.setObjectMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, false, false, true, c)

    objectStat := s.getStat(bucket, object, c) 
    c.Assert(objectStat[StatACL], Equals, "private") 
    c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A")
}

func (s *OssutilCommandSuite) TestBatchSetObjectMeta(c *C) {
    bucket := bucketNamePrefix + "setmeta"
    s.putBucket(bucket, c)

    // put objects
    num := 10
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("设置元信息：%d", i)
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
    s.setObjectMeta(bucket, "设置元信息：", "X-Oss-Meta-update", false, true, true, true, c)
   
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
    s.setObjectMeta(bucket, "非设置元信息", "X-Oss-Meta-c:c", false, false, true, true, c)
    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["X-Oss-Meta-A"], Equals, "A") 
        c.Assert(objectStat["X-Oss-Meta-B"], Equals, "b")
        _, ok := objectStat["X-Oss-Meta-C"]
        c.Assert(ok, Equals, false)
    }

    s.setObjectMeta(bucket, "设置元信息", "X-Oss-Meta-c:c", false, false, true, true, c)
    for _, object := range objectNames {
        objectStat := s.getStat(bucket, object, c) 
        c.Assert(objectStat["X-Oss-Meta-C"], Equals, "c") 
        _, ok := objectStat["X-Oss-Meta-A"]
        c.Assert(ok, Equals, false)
        _, ok = objectStat["X-Oss-Meta-B"]
        c.Assert(ok, Equals, false)
    }

    // error meta
    showElapse, err := s.rawSetMeta(bucket, "设置元信息：", "X-Oss-Meta-c:c", false, true, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, "", "a:b", true, false, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrSetMeta(c *C) {
    args := []string{"os:/", ""}
    showElapse, err := s.rawSetMetaWithArgs(args, false, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta("", "", "", false, false, false, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    bucket := bucketNamePrefix + "setmeta"
    s.putBucket(bucket, c)

    object := "notexistobject"

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", true, true, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucket, object)}, false, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMetaWithArgs([]string{CloudURLToString(bucket, object)}, true, false, false, false)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", false, false, false, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    object = "/object"
    showElapse, err = s.rawSetMeta(bucket, object, "x-oss-object-acl:private#X-Oss-Meta-A:A", false, false, true, true)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}
