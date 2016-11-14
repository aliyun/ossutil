package lib 

import (
    "fmt"
    "strconv"
    "os"
    "time"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) rawCP(srcURL, destURL string, recursive, force, update bool, threshold int64, cpDir string) (bool, error) {
    args := []string{srcURL, destURL}
    showElapse, err := s.rawCPWithArgs(args, recursive, force, update, threshold, cpDir)
    return showElapse, err
}

func (s *OssutilCommandSuite) rawCPWithArgs(args []string, recursive, force, update bool, threshold int64, cpDir string) (bool, error) {
    command := "cp"
    str := ""
    thre := strconv.FormatInt(threshold, 10)
    routines := strconv.Itoa(Routines)
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "recursive": &recursive,
        "force": &force,
        "update": &update,
        "bigfileThreshold": &thre,
        "checkpointDir": &cpDir,
        "routines": &routines,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) TestCPObject(c *C) {
    s.SetUpBucketEnv(c)
    bucket := bucketNamePrefix + "cp1"
    s.putBucket(bucket, c)
    time.Sleep(sleepTime)

    destBucket := bucketNameNotExist 

    // put object
    object := "中文cp" 
    s.putObject(bucket, object, uploadFileName, c)

    // get object
    s.getObject(bucket, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, content)

    // modify uploadFile content
    data := "欢迎使用ossutil"
    s.createFile(uploadFileName, data, c)

    // put to exist object
    s.putObject(bucket, object, uploadFileName, c)

    // get to exist file
    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    // put without specify dest object 
    data1 := "put without specify dest object"
    s.createFile(uploadFileName, data1, c)
    s.putObject(bucket, "", uploadFileName, c)
    s.getObject(bucket, uploadFileName, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data1)

    // get without specify dest file 
    s.getObject(bucket, object, ".", c)
    str = s.readFile(object, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(object)

    // get to file in not exist directory
    notexistdir := "不存在的目录"
    s.getObject(bucket, object, notexistdir + string(os.PathSeparator) + downloadFileName, c)
    str = s.readFile(notexistdir + string(os.PathSeparator) + downloadFileName, c) 
    c.Assert(str, Equals, data)
    _ = os.RemoveAll(notexistdir)

    // copy file
    destObject := "destObject"
    s.copyObject(bucket, object, bucket, destObject, c)

    objectStat := s.getStat(bucket, destObject, c)
    c.Assert(objectStat[StatACL], Equals, "default")
    
    // get dest file
    filePath := downloadFileName + "1" 
    s.getObject(bucket, destObject, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(filePath)

    // put to not exist bucket
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString(destBucket, object), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // get not exist bucket
    showElapse, err = s.rawCP(CloudURLToString(destBucket, object), downloadFileName, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // get not exist object
    showElapse, err = s.rawCP(CloudURLToString(bucket, "notexistobject"), downloadFileName, false, true, false, BigFileThreshold, CheckpointDir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy to not exist bucket
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, destObject), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // corse bucket copy
    destBucket = bucketNameDest

    s.copyObject(bucket, object, destBucket, destObject, c)

    s.getObject(destBucket, destObject, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(filePath)

    // copy single object in directory, test the name of dest object 
    srcObject := "a/b/c/d/e"
    s.putObject(bucket, srcObject, uploadFileName, c)

    s.copyObject(bucket, srcObject, destBucket, "", c)

    s.getObject(destBucket, "e", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.copyObject(bucket, srcObject, destBucket, "a/", c)

    s.getObject(destBucket, "a/e", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.copyObject(bucket, srcObject, destBucket, "a", c)

    s.getObject(destBucket, "a", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    // copy without specify dest object
    s.copyObject(bucket, uploadFileName, destBucket, "", c)
    s.getObject(destBucket, uploadFileName, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 
}

func (s *OssutilCommandSuite) TestErrorCP(c *C) {
    bucket := bucketNameExist 

    // error src_url
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString("", ""), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, CloudURLToString("", bucket), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString("", bucket), downloadFileName, true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString("", ""), downloadFileName, true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, "a", true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss argc
    showElapse, err = s.rawCP(CloudURLToString("", bucket), "", true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy self
    object := "testobject"
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(bucket, object), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(bucket, ""), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(bucket, ""), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(bucket, ""), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(bucket, object), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // err checkpoint dir, conflict with config file
    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucket, object), false, true, true, BigFileThreshold, configFile)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestUploadErrSrc(c *C) {
    srcBucket := bucketNamePrefix + "uploadsrc"
    destBucket := bucketNameNotExist 
    command := "cp"
    args := []string{uploadFileName, CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, "")}
    str := ""
    ok := true
    cpDir := CheckpointDir
    thre := strconv.FormatInt(BigFileThreshold, 10)
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
    bucket := bucketNamePrefix + "bcp" 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime) 

    // create local dir
    dir := "上传目录"
    err := os.MkdirAll(dir, 0777)
    c.Assert(err, IsNil)

    // upload empty dir miss recursive
    showElapse, err := s.rawCP(dir, CloudURLToString(bucket, ""), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // upload empty dir
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, false, BigFileThreshold, CheckpointDir)

    // head object 
    showElapse, err = s.rawGetStat(bucket, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawGetStat(bucket, dir + "/")
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // create dir in dir 
    subdir := "子目录"
    err = os.MkdirAll(dir + "/" + subdir, 0777)
    c.Assert(err, IsNil)

    // upload dir    
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true) 

    time.Sleep(2*sleepTime)

    s.getStat(bucket, subdir + "/", c)

    // remove object
    s.removeObjects(bucket, subdir + "/", false, true, c)

    // create file in dir
    num := 10
    filePaths := []string{subdir + "/"}
    for i := 0; i < num; i++ {
        filePath := fmt.Sprintf("测试文件：%d", i) 
        s.createFile(dir + "/" + filePath, fmt.Sprintf("测试文件：%d内容", i), c)
        filePaths = append(filePaths, filePath)
    }

    os.Stdout = out 
    os.Stderr = errout 

    // upload files
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    
    os.Stdout = testLogFile 
    os.Stderr = testLogFile 

    time.Sleep(10*time.Second)
/*
    for _, filePath := range filePaths {
        s.getStat(bucket, filePath, c)
    }
*/
    // get files
    downDir := "下载目录"
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downDir, true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    for _, filePath := range filePaths {
        _, err := os.Stat(downDir + "/" + filePath)
        c.Assert(err, IsNil)
    }

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)

    // get to exist files
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downDir, true, false, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downDir, true, false, true, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)
    //c.Assert(f.ModTime(), Equals, f1.ModTime())

    // copy files
    destBucket := bucketNameNotExist 
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(destBucket, "123"), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    destBucket = bucketNameDest

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(destBucket, "123"), true, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    time.Sleep(sleepTime)

    for _, filePath := range filePaths {
        s.getStat(destBucket, "123" + filePath, c)
    }

    // remove dir
    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(downDir)

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 
}

func (s *OssutilCommandSuite) TestCPObjectUpdate(c *C) {
    bucket := bucketNamePrefix + "cpupdate" 
    s.putBucket(bucket, c)
    time.Sleep(sleepTime) 

    // create older file and newer file
    oldData := "old data"
    oldFile := "oldFile"
    newData := "new data"
    newFile := "newFile"
    s.createFile(oldFile, oldData, c)
    time.Sleep(1)
    s.createFile(newFile, newData, c)

    // put newer object
    object := "testobject"
    s.putObject(bucket, object, newFile, c)

    // get object
    s.getObject(bucket, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    // put old object with update
    showElapse, err := s.rawCP(oldFile, CloudURLToString(bucket, object), false, false, true, BigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    time.Sleep(2*sleepTime)

    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    showElapse, err = s.rawCP(oldFile, CloudURLToString(bucket, object), false, true, true, BigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    showElapse, err = s.rawCP(oldFile, CloudURLToString(bucket, object), false, false, false, BigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    // get object with update 
    // modify downloadFile
    time.Sleep(1)
    downData := "download file has been modified locally"
    s.createFile(downloadFileName, downData, c) 

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, false, true, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, downData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, true, true, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, downData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, false, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c)
    c.Assert(str, Equals, downData)

    // copy object with update
    //destBucket := bucketNamePrefix + "updatedest"  
    destBucket := bucketNameDest 

    destData := "data for dest bucket"
    destFile := "destFile"
    s.createFile(destFile, destData, c)
    s.putObject(destBucket, object, destFile, c) 

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, object), false, false, true, BigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, destData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, object), false, true, true, BigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, destData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, object), false, false, false, BigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(destBucket, ""), true, false, false, BigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(oldFile)
    _ = os.Remove(newFile)
    _ = os.Remove(destFile)

    s.removeBucket(bucket, true, c)
    time.Sleep(sleepTime) 
}

func (s *OssutilCommandSuite) TestResumeCPObject(c *C) { 
    var threshold int64
    threshold = 1
    cpDir := "checkpoint目录" 

    bucket := bucketNameExist 

    data := "resume cp"
    s.createFile(uploadFileName, data, c)

    // put object
    object := "object" 
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString(bucket, object), false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // get object
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    s.createFile(downloadFileName, "-------long file which must be truncated by cp file------", c)
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    // copy object
    destBucket := bucketNameDest 
    s.putBucket(destBucket, c)

    destObject := "destObject" 

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, destObject), false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, destObject, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)
}

func (s *OssutilCommandSuite) TestCPMulitSrc(c *C) {
    bucket := bucketNameExist 

    // upload multi file 
    file1 := uploadFileName + "1"
    s.createFile(file1, file1, c)
    file2 := uploadFileName + "2"
    s.createFile(file2, file2, c)
    showElapse, err := s.rawCPWithArgs([]string{file1, file2, CloudURLToString(bucket, "")}, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _ = os.Remove(file1)
    _ = os.Remove(file2)

    // download multi objects
    object1 := "object1"
    s.putObject(bucket, object1, uploadFileName, c)
    object2 := "object2"
    s.putObject(bucket, object2, uploadFileName, c)
    showElapse, err = s.rawCPWithArgs([]string{CloudURLToString(bucket, object1), CloudURLToString(bucket, object2), "../"}, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy multi objects
    destBucket := bucketNameDest 
    showElapse, err = s.rawCPWithArgs([]string{CloudURLToString(bucket, object1), CloudURLToString(bucket, object2), CloudURLToString(destBucket, "")}, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrUpload(c *C) {
    // src file not exist
    bucket := bucketNameExist 
    
    showElapse, err := s.rawCP("notexistfile", CloudURLToString(bucket, ""), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // create local dir
    dir := "上传目录"
    err = os.MkdirAll(dir, 0777)
    c.Assert(err, IsNil)
    cpDir := dir + string(os.PathSeparator) + CheckpointDir 
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, true, BigFileThreshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // err object name
    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucket, "/object"), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucket, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    subdir := dir + string(os.PathSeparator) + "subdir"
    err = os.MkdirAll(subdir, 0777)
    c.Assert(err, IsNil)

    showElapse, err = s.rawCP(subdir, CloudURLToString(bucket, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(subdir)
}

func (s *OssutilCommandSuite) TestErrDownload(c *C) {
    bucket := bucketNameExist 
 
    object := "object"
    s.putObject(bucket, object, uploadFileName, c)

    // download to dir, but dir exist as a file
    showElapse, err := s.rawCP(CloudURLToString(bucket, object), configFile + string(os.PathSeparator), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // batch download without -r
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downloadFileName, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // download to file in not exist dir
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), configFile + string(os.PathSeparator) + downloadFileName, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrCopy(c *C) {
    srcBucket := bucketNameExist 

    destBucket := bucketNameDest 

    // batch copy without -r
    showElapse, err := s.rawCP(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, ""), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // error src object name
    showElapse, err = s.rawCP(CloudURLToString(srcBucket, "/object"), CloudURLToString(destBucket, ""), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // err dest object
    object := "object"
    s.putObject(srcBucket, object, uploadFileName, c)
    showElapse, err = s.rawCP(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, "/object"), false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, "/object"), true, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestPreparePartOption(c *C) {
    partSize, routines := copyCommand.preparePartOption(100000000000)
    c.Assert(partSize, Equals, int64(250000000))
    c.Assert(routines, Equals, 5)

    partSize, routines = copyCommand.preparePartOption(80485760)
    c.Assert(partSize, Equals, int64(12816225))
    c.Assert(routines, Equals, 2)

    partSize, routines = copyCommand.preparePartOption(MaxInt64)
    c.Assert(partSize, Equals, int64(922337203685478))
    c.Assert(routines, Equals, 10)

}

func (s *OssutilCommandSuite) TestResumeDownloadRetry(c *C) {
    bucketName := bucketNamePrefix + "cpnotexist"
    bucket, err := copyCommand.command.ossBucket(bucketName)
    c.Assert(err, IsNil)

    err = copyCommand.ossResumeDownloadRetry(bucket, "", "", 0, 0)
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestCPIDKey(c *C) {
    bucket := bucketNameExist 

    object := "testobject" 

    ufile := "ossutil_test.cpidkey"
    data := "欢迎使用ossutil"
    s.createFile(ufile, data, c)

    cfile := "ossutil_test.config_boto"
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "cp"
    str := ""
    args := []string{ufile, CloudURLToString(bucket, object)}
    ok := true
    routines := strconv.Itoa(Routines)
    thre := strconv.FormatInt(BigFileThreshold, 10)
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
}
