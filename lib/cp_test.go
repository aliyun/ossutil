package lib 

import (
    "fmt"
    "strconv"
    "os"
    "time"
    "strings"
    
    oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestCPObject(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)
    // dest bucket not exist
    destBucket := bucketNamePrefix + randLowStr(10) 

    // put object
    s.createFile(uploadFileName, content, c)
    object := "TestCPObject_cp" 
    s.putObject(bucketName, object, uploadFileName, c)

    // get object
    s.getObject(bucketName, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, content)

    // modify uploadFile content
    data := "欢迎使用ossutil"
    s.createFile(uploadFileName, data, c)

    time.Sleep(sleepTime)
    // put to exist object
    s.putObject(bucketName, object, uploadFileName, c)

    // get to exist file
    s.getObject(bucketName, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    // get without specify dest file 
    s.getObject(bucketName, object, ".", c)
    str = s.readFile(object, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(object)

    // put without specify dest object 
    data1 := "put without specify dest object"
    s.createFile(uploadFileName, data1, c)
    s.putObject(bucketName, "", uploadFileName, c)
    s.getObject(bucketName, uploadFileName, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data1)

    // get to file in not exist directory
    notexistdir := "NOTEXISTDIR"
    s.getObject(bucketName, object, notexistdir + string(os.PathSeparator) + downloadFileName, c)
    str = s.readFile(notexistdir + string(os.PathSeparator) + downloadFileName, c) 
    c.Assert(str, Equals, data)
    _ = os.RemoveAll(notexistdir)

    // copy file
    destObject := "TestCPObject_destObject"
    s.copyObject(bucketName, object, bucketName, destObject, c)

    objectStat := s.getStat(bucketName, destObject, c)
    c.Assert(objectStat[StatACL], Equals, "default")
    
    // get dest file
    filePath := downloadFileName + "1" 
    s.getObject(bucketName, destObject, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(filePath)

    // put to not exist bucket
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString(destBucket, object), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // get not exist bucket
    showElapse, err = s.rawCP(CloudURLToString(destBucket, object), downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // get not exist object
    showElapse, err = s.rawCP(CloudURLToString(bucketName, "notexistobject"), downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy to not exist bucket
    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(destBucket, destObject), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // corse bucket copy
    s.putBucket(destBucket, c)

    s.copyObject(bucketName, object, destBucket, destObject, c)

    s.getObject(destBucket, destObject, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(filePath)

    // copy single object in directory, test the name of dest object 
    srcObject := "a/b/c/d/e"
    s.putObject(bucketName, srcObject, uploadFileName, c)
    time.Sleep(time.Second)

    s.copyObject(bucketName, srcObject, destBucket, "", c)

    s.getObject(destBucket, "e", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.copyObject(bucketName, srcObject, destBucket, "a/", c)

    s.getObject(destBucket, "a/e", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.copyObject(bucketName, srcObject, destBucket, "a", c)

    s.getObject(destBucket, "a", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    // copy without specify dest object
    s.copyObject(bucketName, uploadFileName, destBucket, "", c)
    s.getObject(destBucket, uploadFileName, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.removeBucket(bucketName, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestErrorCP(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    // error src_url
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString("", ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, CloudURLToString("", bucketName), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString("", bucketName), downloadFileName, true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString("", ""), downloadFileName, true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, "a", true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss argc
    showElapse, err = s.rawCP(CloudURLToString("", bucketName), "", true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy self
    object := "testobject"
    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(bucketName, object), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(bucketName, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), CloudURLToString(bucketName, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(bucketName, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), CloudURLToString(bucketName, object), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // err checkpoint dir, conflict with config file
    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucketName, object), false, true, true, DefaultBigFileThreshold, configFile)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestUploadErrSrc(c *C) {
    srcBucket := bucketNamePrefix + randLowStr(10) 
    destBucket := bucketNamePrefix + randLowStr(10) 
    command := "cp"
    args := []string{uploadFileName, CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, "")}
    str := ""
    ok := true
    cpDir := CheckpointDir
    thre := strconv.FormatInt(DefaultBigFileThreshold, 10)
    routines := strconv.Itoa(Routines)
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "force": &ok,
        "bigfileThreshold": &thre,
        "checkpointDir": &cpDir,
        "routines": &routines,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestBatchCPObject(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    // create local dir
    dir := randStr(10) 
    err := os.MkdirAll(dir, 0755)
    c.Assert(err, IsNil)

    // upload empty dir miss recursive
    showElapse, err := s.rawCP(dir, CloudURLToString(bucketName, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // upload empty dir
    showElapse, err = s.rawCP(dir, CloudURLToString(bucketName, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)

    // head object 
    showElapse, err = s.rawGetStat(bucketName, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawGetStat(bucketName, dir + "/")
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.RemoveAll(dir)

    // create dir in dir 
    dir = "TestBatchCPObject_dir"
    subdir := randStr(10) 
    err = os.MkdirAll(dir + string(os.PathSeparator) + subdir, 0755)
    c.Assert(err, IsNil)

    // upload dir    
    showElapse, err = s.rawCP(dir, CloudURLToString(bucketName, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true) 

    // remove object
    s.removeObjects(bucketName, subdir + "/", false, true, c)

    // create file in dir
    num := 3 
    filePaths := []string{subdir + "/"}
    for i := 0; i < num; i++ {
        filePath := fmt.Sprintf("TestBatchCPObject_%d", i) 
        s.createFile(dir + "/" + filePath, fmt.Sprintf("测试文件：%d内容", i), c)
        filePaths = append(filePaths, filePath)
    }

    // upload files
    showElapse, err = s.rawCP(dir, CloudURLToString(bucketName, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    
    time.Sleep(7*time.Second)

    // get files
    downDir := "下载目录"
    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), downDir, true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    for _, filePath := range filePaths {
        _, err := os.Stat(downDir + "/" + filePath)
        c.Assert(err, IsNil)
    }

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)

    // get to exist files
    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), downDir, true, false, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), downDir, true, false, true, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)
    //c.Assert(f.ModTime(), Equals, f1.ModTime())

    // copy files
    destBucket := bucketNamePrefix + randLowStr(10) 
    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), CloudURLToString(destBucket, "123"), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.putBucket(destBucket, c)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), CloudURLToString(destBucket, "123"), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    time.Sleep(7*time.Second)

    for _, filePath := range filePaths {
        s.getStat(destBucket, "123" + filePath, c)
    }

    // remove dir
    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(downDir)

    s.removeBucket(bucketName, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestCPObjectUpdate(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    // create older file and newer file
    oldData := "old data"
    oldFile := "oldFile"
    newData := "new data"
    newFile := "newFile"
    s.createFile(oldFile, oldData, c)
    time.Sleep(7*time.Second)
    s.createFile(newFile, newData, c)

    // put newer object
    object := "testobject"
    s.putObject(bucketName, object, newFile, c)

    // get object
    s.getObject(bucketName, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    // put old object with update
    showElapse, err := s.rawCP(oldFile, CloudURLToString(bucketName, object), false, false, true, DefaultBigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    time.Sleep(7*time.Second)

    s.getObject(bucketName, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    showElapse, err = s.rawCP(oldFile, CloudURLToString(bucketName, object), false, true, true, DefaultBigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(bucketName, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    showElapse, err = s.rawCP(oldFile, CloudURLToString(bucketName, object), false, false, false, DefaultBigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(bucketName, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    // get object with update 
    // modify downloadFile
    time.Sleep(1)
    downData := "download file has been modified locally"
    s.createFile(downloadFileName, downData, c) 

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), downloadFileName, false, false, true, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, downData)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), downloadFileName, false, true, true, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, downData)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), downloadFileName, false, false, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c)
    c.Assert(str, Equals, downData)

    // copy object with update
    destBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(destBucket, c)

    destData := "data for dest bucket"
    destFile := "destFile"
    s.createFile(destFile, destData, c)
    s.putObject(destBucket, object, destFile, c) 

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(destBucket, object), false, false, true, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, destData)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(destBucket, object), false, true, true, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, destData)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(destBucket, object), false, false, false, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), CloudURLToString(destBucket, ""), true, false, false, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(oldFile)
    _ = os.Remove(newFile)
    _ = os.Remove(destFile)

    s.removeBucket(bucketName, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestResumeCPObject(c *C) { 
    var threshold int64
    threshold = 1
    cpDir := "checkpoint目录" 

    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    data := "resume cp"
    s.createFile(uploadFileName, data, c)

    // put object
    object := "object" 
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString(bucketName, object), false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // get object
    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), downloadFileName, false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    s.createFile(downloadFileName, "-------long file which must be truncated by cp file------", c)
    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), downloadFileName, false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    // copy object
    destBucket := bucketNamePrefix + randLowStr(10) 
    s.putBucket(destBucket, c)

    destObject := "destObject" 

    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), CloudURLToString(destBucket, destObject), false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, destObject, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    s.removeBucket(bucketName, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestCPMulitSrc(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    // upload multi file 
    file1 := uploadFileName + "1"
    s.createFile(file1, file1, c)
    file2 := uploadFileName + "2"
    s.createFile(file2, file2, c)
    showElapse, err := s.rawCPWithArgs([]string{file1, file2, CloudURLToString(bucketName, "")}, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _ = os.Remove(file1)
    _ = os.Remove(file2)

    // download multi objects
    object1 := "object1"
    s.putObject(bucketName, object1, uploadFileName, c)
    object2 := "object2"
    s.putObject(bucketName, object2, uploadFileName, c)
    showElapse, err = s.rawCPWithArgs([]string{CloudURLToString(bucketName, object1), CloudURLToString(bucketName, object2), "../"}, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy multi objects
    destBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(destBucket, c)
    showElapse, err = s.rawCPWithArgs([]string{CloudURLToString(bucketName, object1), CloudURLToString(bucketName, object2), CloudURLToString(destBucket, "")}, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucketName, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestErrUpload(c *C) {
    // src file not exist
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)
    
    showElapse, err := s.rawCP("notexistfile", CloudURLToString(bucketName, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // create local dir
    dir := randStr(3) + "上传目录"
    err = os.MkdirAll(dir, 0755)
    c.Assert(err, IsNil)
    cpDir := dir + string(os.PathSeparator) + CheckpointDir 
    showElapse, err = s.rawCP(dir, CloudURLToString(bucketName, ""), true, true, true, DefaultBigFileThreshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // err object name
    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucketName, "/object"), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucketName, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    subdir := dir + string(os.PathSeparator) + "subdir"
    err = os.MkdirAll(subdir, 0755)
    c.Assert(err, IsNil)

    showElapse, err = s.rawCP(subdir, CloudURLToString(bucketName, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(subdir)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestErrDownload(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)
 
    object := "object"
    s.putObject(bucketName, object, uploadFileName, c)

    // download to dir, but dir exist as a file
    showElapse, err := s.rawCP(CloudURLToString(bucketName, object), configFile + string(os.PathSeparator), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // batch download without -r
    showElapse, err = s.rawCP(CloudURLToString(bucketName, ""), downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // download to file in not exist dir
    showElapse, err = s.rawCP(CloudURLToString(bucketName, object), configFile + string(os.PathSeparator) + downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestErrCopy(c *C) {
    srcBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(srcBucket, c)

    destBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(destBucket, c)

    // batch copy without -r
    showElapse, err := s.rawCP(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // error src object name
    showElapse, err = s.rawCP(CloudURLToString(srcBucket, "/object"), CloudURLToString(destBucket, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // err dest object
    object := "object"
    s.putObject(srcBucket, object, uploadFileName, c)
    showElapse, err = s.rawCP(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, "/object"), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, "/object"), true, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(srcBucket, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestPreparePartOption(c *C) {
    partSize, routines := copyCommand.preparePartOption(100000000000)
    c.Assert(partSize, Equals, int64(250000000))
    c.Assert(routines, Equals, 15)

    partSize, routines = copyCommand.preparePartOption(100000000000000)
    c.Assert(partSize, Equals, int64(10000000000))
    c.Assert(routines, Equals, 32)

    partSize, routines = copyCommand.preparePartOption(80485760)
    c.Assert(partSize, Equals, int64(2560000))
    c.Assert(routines, Equals, 10)

    partSize, routines = copyCommand.preparePartOption(20121443)
    c.Assert(partSize, Equals, int64(2560000))
    c.Assert(routines, Equals, 3)

    partSize, routines = copyCommand.preparePartOption(MaxInt64)
    c.Assert(partSize, Equals, int64(922337203685478))
    c.Assert(routines, Equals, 32)

    partSize, routines = copyCommand.preparePartOption(1)
    c.Assert(partSize, Equals, int64(oss.MinPartSize))
    c.Assert(routines, Equals, 1)

    partSize, routines = copyCommand.preparePartOption(0)
    c.Assert(partSize, Equals, int64(oss.MinPartSize))
    c.Assert(routines, Equals, 1)

    p := 7 
    parallel := strconv.Itoa(p) 
    copyCommand.command.options[OptionParallel] = &parallel
    partSize, routines = copyCommand.preparePartOption(1)
    c.Assert(routines, Equals, p)
    str := ""
    copyCommand.command.options[OptionParallel] = &str
}

func (s *OssutilCommandSuite) TestResumeDownloadRetry(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10) 
    bucket, err := copyCommand.command.ossBucket(bucketName)
    c.Assert(err, IsNil)

    err = copyCommand.ossResumeDownloadRetry(bucket, "", "", 0, 0)
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestCPIDKey(c *C) {
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    object := "testobject" 

    ufile := "ossutil_test.cpidkey"
    data := "欢迎使用ossutil"
    s.createFile(ufile, data, c)

    cfile := "ossutil_test.config_boto"
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucketName, "abc", bucketName, "abc") 
    s.createFile(cfile, data, c)

    command := "cp"
    str := ""
    args := []string{ufile, CloudURLToString(bucketName, object)}
    ok := true
    routines := strconv.Itoa(Routines)
    thre := strconv.FormatInt(DefaultBigFileThreshold, 10)
    cpDir := CheckpointDir
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
        "force": &ok,
        "bigfileThreshold": &thre,
        "checkpointDir": &cpDir,
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
        "force": &ok,
        "bigfileThreshold": &thre,
        "checkpointDir": &cpDir,
        "routines": &routines,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(ufile)
    _ = os.Remove(cfile)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestUploadOutputDir(c *C) {
    dir := randStr(10) 
    _ = os.RemoveAll(dir)

    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)
    object := randStr(10) 
    ufile := "ossutil_test.testoutputdir"
    data := "content" 
    s.createFile(ufile, data, c)

    // normal copy -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // NoSuchBucket err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketNameNotExist, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // SignatureDoesNotMatch err copy -> no outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s", endpoint, accessKeyID, "abc", bucketName, endpoint) 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, bucketName, "abc") 
    s.createFile(configFile, data, c)

    // err copy without -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy with -r -> outputdir
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), true, true, false, 1, dir) 
    os.Stdout = out
    str := s.readFile(resultPath, c)
    c.Assert(strings.Contains(str, "Error occurs"), Equals, true)
    c.Assert(strings.Contains(str, "See more information in file"), Equals, true)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, IsNil) 

    _ = os.Remove(configFile)
    configFile = cfile

    // get file list of outputdir
    fileList, err := s.getFileList(dir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    // get report file content
    result := s.getReportResult(fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileList[0]), c)
    c.Assert(len(result), Equals, 1)
    
    _ = os.Remove(ufile)
    _ = os.RemoveAll(dir)

    // err list with -C -> no outputdir
    udir := randStr(10) 
    err = os.MkdirAll(udir, 0755)
    c.Assert(err, IsNil)
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucketName, object), false, true, false, 1, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.RemoveAll(udir)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestBatchUploadOutputDir(c *C) {
    udir := randStr(10) 
    _ = os.RemoveAll(udir)
    err := os.MkdirAll(udir, 0755)
    c.Assert(err, IsNil)

    num := 2 
    filePaths := []string{}
    for i := 0; i < num; i++ {
        filePath := randStr(10) 
        s.createFile(udir + "/" + filePath, fmt.Sprintf("测试文件：%d内容", i), c)
        filePaths = append(filePaths, filePath)
    }

    dir := randStr(10) 
    _ = os.RemoveAll(dir)
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)

    // normal copy -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(udir, CloudURLToString(bucketName, udir + "/"), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy without -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucketName, udir + "/"), false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

// err copy -> outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n", "abc", accessKeyID, accessKeySecret) 
    s.createFile(configFile, data, c)
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucketName, udir + "/"), true, true, false, 1, dir) 
    os.Stdout = out
    str := s.readFile(resultPath, c)
    c.Assert(strings.Contains(str, "Error occurs"), Equals, true)
    c.Assert(strings.Contains(str, "See more information in file"), Equals, true)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, IsNil) 

    // get file list of outputdir
    fileList, err := s.getFileList(dir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    // get report file content
    result := s.getReportResult(fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileList[0]), c)
    c.Assert(len(result), Equals, num)
 
    _ = os.Remove(configFile)
    configFile = cfile
    _ = os.RemoveAll(dir)

    // NoSuchBucket err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucketNameNotExist, udir + "/"), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.RemoveAll(udir)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestDownloadOutputDir(c *C) {
    dir := randStr(10) 
    _ = os.RemoveAll(dir)

    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)
    object := randStr(10)
    s.putObject(bucketName, object, uploadFileName, c)

    // normal copy without -r -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(CloudURLToString(bucketName, object), downloadFileName, false, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // normal copy with -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucketName, object), downloadDir, true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucketNameNotExist, object), downloadFileName, true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy without -r -> no outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, bucketName, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucketName, object), downloadFileName, false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // list err copy with -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucketName, object), downloadDir, true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.RemoveAll(dir)
    _ = os.Remove(configFile)
    configFile = cfile

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestCopyOutputDir(c *C) { 
    dir := randStr(10) 
    _ = os.RemoveAll(dir)

    srcBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(srcBucket, c)
    destBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(destBucket, c)

    object := randStr(10)
    s.putObject(srcBucket, object, uploadFileName, c)

    // normal copy -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, object), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(bucketNameNotExist, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucketNameNotExist, object), CloudURLToString(destBucket, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // list err copy without -r -> no outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, srcBucket, "abc") 
    s.createFile(configFile, data, c)
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, object), false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.Remove(configFile)
    configFile = cfile
    _ = os.RemoveAll(dir)

    s.removeBucket(srcBucket, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestBatchCopyOutputDir(c *C) {
    dir := randStr(10) 
    _ = os.RemoveAll(dir)

    srcBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(srcBucket, c)
    destBucket := bucketNamePrefix + randLowStr(10)
    s.putBucket(destBucket, c)

    objectList := []string{}
    num := 3
    for i := 0; i < num; i++ {
        object := randStr(10)
        s.putObject(srcBucket, object, uploadFileName, c)
        objectList = append(objectList, object)
    }

    showElapse, err := s.rawCPWithOutputDir(CloudURLToString(srcBucket, objectList[0]), CloudURLToString(destBucket, ""), true, true, false, 1, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.RemoveAll(dir)

    // normal copy -> no outputdir 
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, ""), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // bucketNameNotExist err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, ""), CloudURLToString(bucketNameNotExist, ""), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // test objectStatistic err
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", srcBucket, "abc", srcBucket, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, ""), true, true, false, 1, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.Remove(configFile)
    configFile = cfile
    _ = os.RemoveAll(dir)

    s.removeBucket(srcBucket, true, c)
    s.removeBucket(destBucket, true, c)
}

func (s *OssutilCommandSuite) TestConfigOutputDir(c *C) {
    // test default outputdir
    edir := "" 
    dir := randStr(10) 
    dir1 := randStr(10) 
    _ = os.RemoveAll(DefaultOutputDir)
    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(dir1)

    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)
    object := randStr(10) 
    ufile := "ossutil_test.testoutputdir"
    data := "content" 
    s.createFile(ufile, data, c)

    // err copy -> outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, bucketName, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err := s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), true, true, false, 1, edir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(DefaultOutputDir)
    c.Assert(err, IsNil) 

    // get file list of outputdir
    fileList, err := s.getFileList(DefaultOutputDir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    _ = os.RemoveAll(DefaultOutputDir) 

    // config outputdir
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\noutputDir=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, dir, bucketName, endpoint, bucketName, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), true, true, false, 1, "") 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, IsNil) 
    _, err = os.Stat(DefaultOutputDir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // get file list of outputdir
    fileList, err = s.getFileList(dir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(DefaultOutputDir)

    // option and config
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), true, true, false, 1, dir1) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir1)
    c.Assert(err, IsNil) 
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)
    _, err = os.Stat(DefaultOutputDir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // get file list of outputdir
    fileList, err = s.getFileList(dir1)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    _ = os.Remove(configFile)
    configFile = cfile
    _ = os.RemoveAll(dir1)
    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(DefaultOutputDir)

    s.createFile(uploadFileName, content, c)
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketName, object), true, true, false, 1, uploadFileName) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestInitReportError(c *C) {
    s.createFile(uploadFileName, content, c)
    report, err := GetReporter(false, DefaultOutputDir, "")
    c.Assert(err, IsNil)
    c.Assert(report, IsNil)

    report, err = GetReporter(true, uploadFileName, "")
    c.Assert(err, NotNil)
    c.Assert(report, IsNil)
}

func (s *OssutilCommandSuite) TestCopyFunction(c *C) {
    // test fileStatistic
    copyCommand.monitor.init(operationTypePut)
    storageURL, err := StorageURLFromString("&~")
    c.Assert(err, IsNil)
    copyCommand.fileStatistic([]StorageURLer{storageURL})
    c.Assert(copyCommand.monitor.seekAheadEnd, Equals, true)
    c.Assert(copyCommand.monitor.seekAheadError, NotNil)

    // test fileProducer
    chFiles := make(chan fileInfoType, ChannelBuf)
    chListError := make(chan error, 1)
    storageURL, err = StorageURLFromString("&~")
    c.Assert(err, IsNil)
    copyCommand.fileProducer([]StorageURLer{storageURL}, chFiles, chListError)
    err = <- chListError
    c.Assert(err, NotNil)

    // test put object error
    bucketName := bucketNameNotExist
    bucket, err := copyCommand.command.ossBucket(bucketName)
    c.Assert(err, IsNil)
    err = copyCommand.ossPutObjectRetry(bucket, "object", "")
    c.Assert(err, NotNil)

    // test formatResultPrompt
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = fmt.Errorf("test error")
    copyCommand.cpOption.ctnu = true
    err = copyCommand.formatResultPrompt(err)
    c.Assert(err, IsNil)
    os.Stdout = out
    str := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(str, "succeed"), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)

    // test download file error
    err = copyCommand.ossDownloadFileRetry(bucket, "object", downloadFileName)
    c.Assert(err, NotNil)

    // test truncateFile
    err = copyCommand.truncateFile("ossutil_notexistfile", 1)
    c.Assert(err, NotNil)
    s.createFile(uploadFileName, "abc", c)
    err = copyCommand.truncateFile(uploadFileName, 1)
    c.Assert(err, IsNil)
    str = s.readFile(uploadFileName, c)
    c.Assert(str, Equals, "a")
}

func (s *OssutilCommandSuite) TestSnapshot(c *C) {
    // upload with snapshot
    bucketName := bucketNamePrefix + randLowStr(10)
    s.putBucket(bucketName, c)
    data := randStr(20)
    s.createFile(uploadFileName, data, c)
    object := randStr(10)
    spath := "ossutil.snapshot-dir"
    _ = os.RemoveAll(spath)

    err := s.initCopyWithSnapshot(uploadFileName, CloudURLToString(bucketName, object), false, false, false, DefaultBigFileThreshold, spath)
    c.Assert(err, IsNil)
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    c.Assert(copyCommand.monitor.fileNum, Equals, int64(1))
    c.Assert(copyCommand.monitor.dirNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.skipNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.errNum, Equals, int64(0))

    s.getObject(bucketName, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c)
    c.Assert(str, Equals, data)

    _, err = os.Stat(spath)
    c.Assert(err, IsNil)

    // upload again
    err = s.initCopyWithSnapshot(uploadFileName, CloudURLToString(bucketName, object), false, false, false, DefaultBigFileThreshold, spath)
    c.Assert(err, IsNil)
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    c.Assert(copyCommand.monitor.fileNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.dirNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.skipNum, Equals, int64(1))
    c.Assert(copyCommand.monitor.errNum, Equals, int64(0))

    s.getObject(bucketName, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c)
    c.Assert(str, Equals, data)

    _, err = os.Stat(spath)
    c.Assert(err, IsNil)

    // modify local and upload again
    data = randStr(21)
    s.createFile(uploadFileName, data, c)

    err = s.initCopyWithSnapshot(uploadFileName, CloudURLToString(bucketName, object), false, false, false, DefaultBigFileThreshold, spath)
    c.Assert(err, IsNil)
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    c.Assert(copyCommand.monitor.fileNum, Equals, int64(1))
    c.Assert(copyCommand.monitor.dirNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.skipNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.errNum, Equals, int64(0))

    s.getObject(bucketName, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c)
    c.Assert(str, Equals, data)

    _, err = os.Stat(spath)
    c.Assert(err, IsNil)

    // -u --snapshot-path
    time.Sleep(7*time.Second)
    s.createFile(uploadFileName, data, c)
    err = s.initCopyWithSnapshot(uploadFileName, CloudURLToString(bucketName, object), false, true, true, DefaultBigFileThreshold, spath)
    c.Assert(err, IsNil)
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    c.Assert(copyCommand.monitor.fileNum, Equals, int64(1))
    c.Assert(copyCommand.monitor.dirNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.skipNum, Equals, int64(0))
    c.Assert(copyCommand.monitor.errNum, Equals, int64(0))

    // download with snapshot
    err = s.initCopyWithSnapshot(CloudURLToString(bucketName, object), downloadFileName, false, false, false, DefaultBigFileThreshold, spath)
    c.Assert(err, IsNil)
    err = copyCommand.RunCommand()
    c.Assert(err, NotNil)

    // copy with snapshot
    err = s.initCopyWithSnapshot(CloudURLToString(bucketName, object), CloudURLToString(bucketNameDest, object), false, false, false, DefaultBigFileThreshold, spath)
    c.Assert(err, IsNil)
    err = copyCommand.RunCommand()
    c.Assert(err, NotNil)

    _ = os.RemoveAll(spath)

    // snapshot path exist and invalid 
    err = s.initCopyWithSnapshot(uploadFileName, CloudURLToString(bucketName, object), false, false, false, DefaultBigFileThreshold, uploadFileName)
    c.Assert(err, IsNil)
    err = copyCommand.RunCommand()
    c.Assert(err, NotNil)

    s.removeBucket(bucketName, true, c)
}
