package lib 

import (
    "os"
    "fmt"
    //"time"
    "strings"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestUploadProgressBar(c *C) {
    // single file
    udir := randStr(11) 
    _ = os.RemoveAll(udir)
    err := os.MkdirAll(udir, 0755)
    c.Assert(err, IsNil)
    object := "TestUploadProgressBar" + randStr(10)

    num := 2 
    len := 0
    for i := 0; i < num; i++ {
        filePath := randStr(10) 
        s.createFile(udir + string(os.PathSeparator) + filePath, randStr((i+3)*30*num), c)
        len += (i+3)*30*num 
    }

    bucket := bucketNameExist

    // init copyCommand
    err = s.initCopyCommand(udir, CloudURLToString(bucket, object), true, true, false, DefaultBigFileThreshold, CheckpointDir, DefaultOutputDir)
    c.Assert(err, IsNil)

    // check output
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)
    
    str := copyCommand.monitor.progressBar(false, normalExit) 
    c.Assert(str, Equals, "") 
    str = copyCommand.monitor.progressBar(false, errExit) 
    c.Assert(str, Equals, "") 
    str = copyCommand.monitor.progressBar(true, normalExit)
    c.Assert(str, Equals, "")
    str = copyCommand.monitor.progressBar(true, errExit)
    c.Assert(str, Equals, "")

    snap := copyCommand.monitor.getSnapshot()
    c.Assert(snap.transferSize, Equals, int64(len)) 
    c.Assert(snap.skipSize, Equals, int64(0))
    c.Assert(snap.dealSize, Equals, int64(len))
    c.Assert(snap.fileNum, Equals, int64(num))
    c.Assert(snap.dirNum, Equals, int64(0))
    c.Assert(snap.skipNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.okNum, Equals, int64(num))
    c.Assert(snap.dealNum, Equals, int64(num))
    c.Assert(copyCommand.monitor.getPrecent(snap), Equals, 100)

    str = strings.ToLower(copyCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", num)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)

    str = strings.ToLower(copyCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)

    str = strings.ToLower(copyCommand.monitor.getFinishBar(errExit))
    c.Assert(strings.Contains(str, "when error happens."), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)

    // mkdir a subdir in dir
    subdir := udir + string(os.PathSeparator) + "subdir" 
    subdir1 := udir + string(os.PathSeparator) + "subdir1" 
    err = os.MkdirAll(subdir, 0755)
    c.Assert(err, IsNil)
    err = os.MkdirAll(subdir1, 0755)
    c.Assert(err, IsNil)

    // put file to subdir
    num1 := 2
    len1 := 0
    for i := 0; i < num1; i++ {
        filePath := randStr(10) 
        s.createFile(subdir + string(os.PathSeparator) + filePath, randStr((i+1)*20*num1), c)
        len1 += (i+1)*20*num1 
    }

    // init copyCommand
    err = s.initCopyCommand(udir, CloudURLToString(bucket, object), true, true, true, DefaultBigFileThreshold, CheckpointDir, DefaultOutputDir)
    c.Assert(err, IsNil)

    // copy with update
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)
 
    str = copyCommand.monitor.progressBar(false, normalExit) 
    c.Assert(str, Equals, "") 

    snap = copyCommand.monitor.getSnapshot()
    c.Assert(snap.transferSize, Equals, int64(len1)) 
    c.Assert(snap.skipSize, Equals, int64(len))
    c.Assert(snap.dealSize, Equals, int64(len+len1))
    c.Assert(snap.fileNum, Equals, int64(num1))
    c.Assert(snap.dirNum, Equals, int64(2))
    c.Assert(snap.skipNum, Equals, int64(num))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.okNum, Equals, int64(num+num1+2))
    c.Assert(snap.dealNum, Equals, int64(num+num1+2))
    c.Assert(copyCommand.monitor.getPrecent(snap), Equals, 100)

    str = strings.ToLower(copyCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", snap.dealNum)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)
    c.Assert(strings.Contains(str, "skip"), Equals, true)
    c.Assert(strings.Contains(str, "directories"), Equals, true)

    str = strings.ToLower(copyCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(str, "skip"), Equals, true)
    c.Assert(strings.Contains(str, "directories"), Equals, true)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)

    _ = os.RemoveAll(udir)
}

func (s *OssutilCommandSuite) TestDownloadProgressBar(c *C) {
    s.createFile(uploadFileName, "", c)
    bucket := bucketNameExist
    object := randStr(10)
    s.putObject(bucket, object, uploadFileName, c)

    // normal download single object of content length 0
    err := s.initCopyCommand(CloudURLToString(bucket, object), downloadDir, true, true, false, DefaultBigFileThreshold, CheckpointDir, DefaultOutputDir)
    c.Assert(err, IsNil)

    // check output
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    str := copyCommand.monitor.progressBar(false, normalExit)
    c.Assert(str, Equals, "")

    snap := copyCommand.monitor.getSnapshot()
    c.Assert(snap.transferSize, Equals, int64(0)) 
    c.Assert(snap.skipSize, Equals, int64(0))
    c.Assert(snap.dealSize, Equals, int64(0))
    c.Assert(snap.fileNum, Equals, int64(1))
    c.Assert(snap.dirNum, Equals, int64(0))
    c.Assert(snap.skipNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.okNum, Equals, int64(1))
    c.Assert(snap.dealNum, Equals, int64(1))
    c.Assert(copyCommand.monitor.getPrecent(snap), Equals, 100)

    str = strings.ToLower(copyCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", 1)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)

    str = strings.ToLower(copyCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)
}

func (s *OssutilCommandSuite) TestCopyProgressBar(c *C) {
    s.createFile(uploadFileName, randStr(15), c)
    srcBucket := bucketNameExist
    destBucket := bucketNameDest
    num := 2
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestCopyProgressBar%d", i)
        s.putObject(srcBucket, object, uploadFileName, c)
    }

    // normal download single object of content length 0
    err := s.initCopyCommand(CloudURLToString(srcBucket, "TestCopyProgressBar"), CloudURLToString(destBucket, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir, DefaultOutputDir)
    c.Assert(err, IsNil)

    // check output
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    str := copyCommand.monitor.progressBar(false, normalExit)
    c.Assert(str, Equals, "")

    snap := copyCommand.monitor.getSnapshot()
    c.Assert(snap.transferSize, Equals, int64(30)) 
    c.Assert(snap.skipSize, Equals, int64(0))
    c.Assert(snap.dealSize, Equals, int64(30))
    c.Assert(snap.fileNum, Equals, int64(2))
    c.Assert(snap.dirNum, Equals, int64(0))
    c.Assert(snap.skipNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.okNum, Equals, int64(2))
    c.Assert(snap.dealNum, Equals, int64(2))
    c.Assert(copyCommand.monitor.getPrecent(snap), Equals, 100)

    str = strings.ToLower(copyCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", 2)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)

    str = strings.ToLower(copyCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)
}

func (s *OssutilCommandSuite) TestProgressBarStatisticErr(c *C) {
    // test batch download err 
    s.createFile(uploadFileName, "TestProgressBarStatisticErr", c)
    bucket := bucketNameExist
    num := 2
    for i := 0; i < num; i++ {
        object := randStr(10)
        s.putObject(bucket, object, uploadFileName, c)
    }

    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(configFile, data, c)

    err := s.initCopyCommand(CloudURLToString(bucket, ""), downloadDir, true, true, false, DefaultBigFileThreshold, CheckpointDir, DefaultOutputDir)
    c.Assert(err, IsNil)

    // check output
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = copyCommand.RunCommand()
    c.Assert(err, NotNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, false)

    _ = os.Remove(configFile)
    configFile = cfile

    snap := copyCommand.monitor.getSnapshot()
    c.Assert(snap.transferSize, Equals, int64(0)) 
    c.Assert(snap.skipSize, Equals, int64(0))
    c.Assert(snap.dealSize, Equals, int64(0))
    c.Assert(snap.fileNum, Equals, int64(0))
    c.Assert(snap.dirNum, Equals, int64(0))
    c.Assert(snap.skipNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.okNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))

    str := strings.ToLower(copyCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("scanned num: %d", snap.dealNum)), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, false)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, false)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, false)

    str = strings.ToLower(copyCommand.monitor.getFinishBar(errExit))
    c.Assert(strings.Contains(str, "succeed"), Equals, false)
    c.Assert(strings.Contains(str, "scanned"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, false)
    c.Assert(strings.Contains(str, "when error happens"), Equals, true)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, false)

    str1 := strings.ToLower(copyCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str1)), Equals, false)
}

func (s *OssutilCommandSuite) TestProgressBarContinueErr(c *C) {
    bucket := bucketNameExist
    udir := randStr(11) 
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

    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n", "abc", accessKeyID, accessKeySecret) 
    s.createFile(configFile, data, c)

    err = s.initCopyCommand(udir, CloudURLToString(bucket, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir, DefaultOutputDir)
    c.Assert(err, IsNil)

    _ = os.Remove(configFile)
    configFile = cfile

    // check output
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = copyCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "error"), Equals, true)
    c.Assert(strings.Contains(pstr, "succeed"), Equals, false)

    snap := copyCommand.monitor.getSnapshot()
    c.Assert(snap.transferSize, Equals, int64(0)) 
    c.Assert(snap.skipSize, Equals, int64(0))
    c.Assert(snap.dealSize, Equals, int64(0))
    c.Assert(snap.fileNum, Equals, int64(0))
    c.Assert(snap.dirNum, Equals, int64(0))
    c.Assert(snap.skipNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(num))
    c.Assert(snap.okNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(num))

    str := strings.ToLower(copyCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", snap.dealNum)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, true)
    c.Assert(strings.Contains(str, "progress"), Equals, true)
    c.Assert(strings.Contains(str, "skip"), Equals, false)
    c.Assert(strings.Contains(str, "directories"), Equals, false)

    str = strings.ToLower(copyCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "finishwitherror"), Equals, true)
    c.Assert(strings.Contains(str, "succeed"), Equals, false)
    c.Assert(strings.Contains(str, fmt.Sprintf("error num: %d", snap.errNum)), Equals, true)

    _ = os.RemoveAll(udir)
}

func (s *OssutilCommandSuite) TestSingleFileProgress(c *C) {
    bucket := bucketNameExist
    object := randStr(10)
    destObject := randStr(10)

    // single large file
    data := strings.Repeat("a", 10240)
    s.createFile(uploadFileName, data, c)

    for threshold := range []int64{1024, DefaultBigFileThreshold} {
        // init copyCommand
        err := s.initCopyCommand(uploadFileName, CloudURLToString(bucket, object), false, true, false, int64(threshold), CheckpointDir, DefaultOutputDir)
        c.Assert(err, IsNil)
        copyCommand.monitor.init(operationTypePut)

        snap := copyCommand.monitor.getSnapshot()
        c.Assert(snap.transferSize, Equals, int64(0)) 
        c.Assert(snap.skipSize, Equals, int64(0))
        c.Assert(snap.dealSize, Equals, int64(0))
        c.Assert(snap.fileNum, Equals, int64(0))
        c.Assert(snap.dirNum, Equals, int64(0))
        c.Assert(snap.skipNum, Equals, int64(0))
        c.Assert(snap.errNum, Equals, int64(0))
        c.Assert(snap.okNum, Equals, int64(0))
        c.Assert(snap.dealNum, Equals, int64(0))

        str := strings.ToLower(copyCommand.monitor.getProgressBar())
        c.Assert(strings.Contains(str, "total num"), Equals, false)
        c.Assert(strings.Contains(str, "scanned"), Equals, true)
        c.Assert(strings.Contains(str, "error"), Equals, false)
        c.Assert(strings.Contains(str, "progress"), Equals, false)
        c.Assert(strings.Contains(str, "skip"), Equals, false)
        c.Assert(strings.Contains(str, "directories"), Equals, false)
        c.Assert(strings.Contains(str, "upload"), Equals, false)
        c.Assert(strings.Contains(str, "download"), Equals, false)
        c.Assert(strings.Contains(str, "copy"), Equals, false)

        // check output
        testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
        out := os.Stdout
        os.Stdout = testResultFile
        err = copyCommand.RunCommand()
        c.Assert(err, IsNil)
        os.Stdout = out
        pstr := strings.ToLower(s.readFile(resultPath, c))
        c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
        c.Assert(strings.Contains(pstr, "error"), Equals, false)
     
        snap = copyCommand.monitor.getSnapshot()
        c.Assert(snap.transferSize, Equals, int64(10240)) 
        c.Assert(snap.skipSize, Equals, int64(0))
        c.Assert(snap.dealSize, Equals, int64(10240))
        c.Assert(snap.fileNum, Equals, int64(1))
        c.Assert(snap.dirNum, Equals, int64(0))
        c.Assert(snap.skipNum, Equals, int64(0))
        c.Assert(snap.errNum, Equals, int64(0))
        c.Assert(snap.okNum, Equals, int64(1))
        c.Assert(snap.dealNum, Equals, int64(1))

        str = strings.ToLower(copyCommand.monitor.getProgressBar())
        c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", 1)), Equals, true)
        c.Assert(strings.Contains(str, "error"), Equals, false)
        c.Assert(strings.Contains(str, "progress"), Equals, true)
        c.Assert(strings.Contains(str, "skip"), Equals, false)
        c.Assert(strings.Contains(str, "directories"), Equals, false)
        c.Assert(strings.Contains(str, "upload"), Equals, true)
        c.Assert(strings.Contains(str, "download"), Equals, false)
        c.Assert(strings.Contains(str, "copy"), Equals, false)

        // download
        err = s.initCopyCommand(CloudURLToString(bucket, object), downloadFileName, false, true, false, 1024, CheckpointDir, DefaultOutputDir)
        c.Assert(err, IsNil)
        copyCommand.monitor.init(operationTypeGet)

        testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
        out = os.Stdout
        os.Stdout = testResultFile
        err = copyCommand.RunCommand()
        c.Assert(err, IsNil)
        os.Stdout = out
        pstr = strings.ToLower(s.readFile(resultPath, c))
        c.Assert(strings.Contains(pstr, "error"), Equals, false)
     
        snap = copyCommand.monitor.getSnapshot()
        c.Assert(snap.transferSize, Equals, int64(10240)) 
        c.Assert(snap.skipSize, Equals, int64(0))
        c.Assert(snap.dealSize, Equals, int64(10240))
        c.Assert(snap.fileNum, Equals, int64(1))
        c.Assert(snap.dirNum, Equals, int64(0))
        c.Assert(snap.skipNum, Equals, int64(0))
        c.Assert(snap.errNum, Equals, int64(0))
        c.Assert(snap.okNum, Equals, int64(1))
        c.Assert(snap.dealNum, Equals, int64(1))
        c.Assert(copyCommand.monitor.getPrecent(snap), Equals, 100)

        str = strings.ToLower(copyCommand.monitor.getProgressBar())
        c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", 1)), Equals, true)
        c.Assert(strings.Contains(str, "error"), Equals, false)
        c.Assert(strings.Contains(str, "progress"), Equals, true)
        c.Assert(strings.Contains(str, "skip"), Equals, false)
        c.Assert(strings.Contains(str, "directories"), Equals, false)
        c.Assert(strings.Contains(str, "upload"), Equals, false)
        c.Assert(strings.Contains(str, "download"), Equals, true)
        c.Assert(strings.Contains(str, "copy"), Equals, false)

        // copy
        err = s.initCopyCommand(CloudURLToString(bucket, object), CloudURLToString(bucket, destObject), false, true, false, 1024, CheckpointDir, DefaultOutputDir)
        c.Assert(err, IsNil)
        copyCommand.monitor.init(operationTypeCopy)

        testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
        out = os.Stdout
        os.Stdout = testResultFile
        err = copyCommand.RunCommand()
        c.Assert(err, IsNil)
        os.Stdout = out
        pstr = strings.ToLower(s.readFile(resultPath, c))
        c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
        c.Assert(strings.Contains(pstr, "error"), Equals, false)
     
        snap = copyCommand.monitor.getSnapshot()
        c.Assert(snap.transferSize, Equals, int64(10240)) 
        c.Assert(snap.skipSize, Equals, int64(0))
        c.Assert(snap.dealSize, Equals, int64(10240))
        c.Assert(snap.fileNum, Equals, int64(1))
        c.Assert(snap.dirNum, Equals, int64(0))
        c.Assert(snap.skipNum, Equals, int64(0))
        c.Assert(snap.errNum, Equals, int64(0))
        c.Assert(snap.okNum, Equals, int64(1))
        c.Assert(snap.dealNum, Equals, int64(1))
        c.Assert(copyCommand.monitor.getPrecent(snap), Equals, 100)

        str = strings.ToLower(copyCommand.monitor.getProgressBar())
        c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", 1)), Equals, true)
        c.Assert(strings.Contains(str, "error"), Equals, false)
        c.Assert(strings.Contains(str, "progress"), Equals, true)
        c.Assert(strings.Contains(str, "skip"), Equals, false)
        c.Assert(strings.Contains(str, "directories"), Equals, false)
        c.Assert(strings.Contains(str, "upload"), Equals, false)
        c.Assert(strings.Contains(str, "download"), Equals, false)
        c.Assert(strings.Contains(str, "copy"), Equals, true)

        // copy skip
        err = s.initCopyCommand(CloudURLToString(bucket, object), CloudURLToString(bucket, destObject), false, true, true, 1024, CheckpointDir, DefaultOutputDir)
        c.Assert(err, IsNil)
        copyCommand.monitor.init(operationTypeCopy)

        testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
        out = os.Stdout
        os.Stdout = testResultFile
        err = copyCommand.RunCommand()
        c.Assert(err, IsNil)
        os.Stdout = out
        pstr = strings.ToLower(s.readFile(resultPath, c))
        c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
        c.Assert(strings.Contains(pstr, "error"), Equals, false)
     
        snap = copyCommand.monitor.getSnapshot()
        c.Assert(snap.transferSize, Equals, int64(0)) 
        c.Assert(snap.skipSize, Equals, int64(10240))
        c.Assert(snap.dealSize, Equals, int64(10240))
        c.Assert(snap.fileNum, Equals, int64(0))
        c.Assert(snap.dirNum, Equals, int64(0))
        c.Assert(snap.skipNum, Equals, int64(1))
        c.Assert(snap.errNum, Equals, int64(0))
        c.Assert(snap.okNum, Equals, int64(1))
        c.Assert(snap.dealNum, Equals, int64(1))
        c.Assert(copyCommand.monitor.getPrecent(snap), Equals, 100)

        str = strings.ToLower(copyCommand.monitor.getProgressBar())
        c.Assert(strings.Contains(str, fmt.Sprintf("total num: %d", 1)), Equals, true)
        c.Assert(strings.Contains(str, "error"), Equals, false)
        c.Assert(strings.Contains(str, "progress"), Equals, true)
        c.Assert(strings.Contains(str, "skip"), Equals, true)
        c.Assert(strings.Contains(str, "directories"), Equals, false)
        c.Assert(strings.Contains(str, "upload"), Equals, false)
        c.Assert(strings.Contains(str, "download"), Equals, false)
    }
}

func (s *OssutilCommandSuite) TestSetACLProgress(c *C) {
    bucket := bucketNameExist

    num := 2
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestSetACLProgress%d", i)
        s.putObject(bucket, object, uploadFileName, c)
        objectNames = append(objectNames, object)
    }

    // set object acl without -r -> no progress
    err := s.initSetACL(bucket, objectNames[0], "private", false, false, true)
    c.Assert(err, IsNil)

    // check output
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = setACLCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, false)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    snap := setACLCommand.monitor.getSnapshot()
    c.Assert(snap.okNum, Equals, int64(0)) 
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))

    // batch set object acl -> progress
    err = s.initSetACL(bucket, "TestSetACLProgress", "private", true, false, true)
    c.Assert(err, IsNil)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = setACLCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    snap = setACLCommand.monitor.getSnapshot()
    c.Assert(snap.okNum, Equals, int64(num)) 
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(num))
    c.Assert(setACLCommand.monitor.getPrecent(snap), Equals, 100)

    str := strings.ToLower(setACLCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", 2)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)

    str = strings.ToLower(setACLCommand.monitor.getFinishBar())
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)

    // batch set acl list error
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n", endpoint, accessKeyID, "") 
    s.createFile(configFile, data, c)

    err = s.initSetACL(bucket, "TestSetACLProgress", "private", true, false, true)
    c.Assert(err, IsNil)

    _ = os.Remove(configFile)
    configFile = cfile

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = setACLCommand.RunCommand()
    c.Assert(err, NotNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, false)

    snap = setACLCommand.monitor.getSnapshot()
    c.Assert(snap.okNum, Equals, int64(0)) 
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))

    str = strings.ToLower(setACLCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("scanned %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, false)

    str = strings.ToLower(setACLCommand.monitor.getFinishBar())
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, false)

    setACLCommand.monitor.init("Setted acl on") 
    setACLCommand.command.updateMonitor(err, &setACLCommand.monitor)
    c.Assert(setACLCommand.monitor.errNum, Equals, int64(1))
    c.Assert(setACLCommand.monitor.okNum, Equals, int64(0))

    str = strings.ToLower(setACLCommand.monitor.getFinishBar())
    c.Assert(strings.Contains(str, "succeed:"), Equals, false)
    c.Assert(strings.Contains(str, "when error happens"), Equals, true)
    c.Assert(strings.Contains(str, "setted acl on 0 objects"), Equals, true)
}

func (s *OssutilCommandSuite) TestSetMetaProgress(c *C) {
    bucket := bucketNameExist

    num := 2
    objectNames := []string{}
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestSetMetaProgress%d", i)
        s.putObject(bucket, object, uploadFileName, c)
        objectNames = append(objectNames, object)
    }

    // set object meta without -r -> no progress
    err := s.initSetMeta(bucket, objectNames[0], "x-oss-object-acl:default#X-Oss-Meta-A:A", true, false, false, true, DefaultLanguage)
    c.Assert(err, IsNil)

    // check output
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = setMetaCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, false)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    snap := setMetaCommand.monitor.getSnapshot()
    c.Assert(snap.okNum, Equals, int64(0)) 
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))

    // batch set object acl -> progress
    err = s.initSetMeta(bucket, "TestSetMetaProgress", "x-oss-object-acl:default#X-Oss-Meta-A:A", true, false, true, true, DefaultLanguage)
    c.Assert(err, IsNil)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = setMetaCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    snap = setMetaCommand.monitor.getSnapshot()
    c.Assert(snap.okNum, Equals, int64(num)) 
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(num))
    c.Assert(setMetaCommand.monitor.getPrecent(snap), Equals, 100)

    str := strings.ToLower(setMetaCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", 2)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)

    str = strings.ToLower(setMetaCommand.monitor.getFinishBar())
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)

    // batch set acl list error
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n", endpoint, accessKeyID, "") 
    s.createFile(configFile, data, c)

    err = s.initSetMeta(bucket, "TestSetMetaProgress", "x-oss-object-acl:default#X-Oss-Meta-A:A", true, false, true, true, DefaultLanguage)
    c.Assert(err, IsNil)

    _ = os.Remove(configFile)
    configFile = cfile

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = setMetaCommand.RunCommand()
    c.Assert(err, NotNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, false)

    snap = setMetaCommand.monitor.getSnapshot()
    c.Assert(snap.okNum, Equals, int64(0)) 
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))

    str = strings.ToLower(setMetaCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("scanned %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, false)

    str = strings.ToLower(setMetaCommand.monitor.getFinishBar())
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, "total"), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, false)

    setMetaCommand.monitor.init("Setted meta on") 
    setMetaCommand.command.updateMonitor(err, &setMetaCommand.monitor)
    c.Assert(setMetaCommand.monitor.errNum, Equals, int64(1))
    c.Assert(setMetaCommand.monitor.okNum, Equals, int64(0))

    str = strings.ToLower(setMetaCommand.monitor.getFinishBar())
    c.Assert(strings.Contains(str, "succeed:"), Equals, false)
    c.Assert(strings.Contains(str, "when error happens"), Equals, true)
    c.Assert(strings.Contains(str, "setted meta on 0 objects"), Equals, true)
}

func (s *OssutilCommandSuite) TestRemoveSingleProgress(c *C) {
    bucket := bucketNameExist

    // remove single not exist object
    object := randStr(10)
    err := s.initRemove(bucket, object, false, true, false)
    c.Assert(err, IsNil)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = removeCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, fmt.Sprintf("total %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    c.Assert(int64(removeCommand.monitor.op), Equals, int64(rmObject))
    c.Assert(removeCommand.monitor.removedBucket, Equals, "")

    snap := removeCommand.monitor.getSnapshot()
    c.Assert(snap.objectNum, Equals, int64(0)) 
    c.Assert(snap.uploadIdNum, Equals, int64(0))
    c.Assert(snap.errObjectNum, Equals, int64(0))
    c.Assert(snap.errUploadIdNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(removeCommand.monitor.getPrecent(snap), Equals, 100)

    str := strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)

    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)

    // remove single exist object
    s.putObject(bucket, object, uploadFileName, c)

    err = s.initRemove(bucket, object, false, true, false)
    c.Assert(err, IsNil)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = removeCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, fmt.Sprintf("total %d objects", 1)), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    snap = removeCommand.monitor.getSnapshot()
    c.Assert(snap.objectNum, Equals, int64(1)) 
    c.Assert(snap.uploadIdNum, Equals, int64(0))
    c.Assert(snap.errObjectNum, Equals, int64(0))
    c.Assert(snap.errUploadIdNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(1))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(removeCommand.monitor.getPrecent(snap), Equals, 100)

    str = strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", 1)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)

    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", 1)), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d objects", 1)), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)
}

func (s *OssutilCommandSuite) TestBatchRemoveProgress(c *C) {
    bucket := bucketNameExist

    // batch remove not exist objects
    err := s.initRemove(bucket, "TestBatchRemoveProgresssss", true, true, false)
    c.Assert(err, IsNil)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = removeCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, fmt.Sprintf("total %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    c.Assert(int64(removeCommand.monitor.op), Equals, int64(rmObject))
    c.Assert(removeCommand.monitor.removedBucket, Equals, "")

    snap := removeCommand.monitor.getSnapshot()
    c.Assert(snap.objectNum, Equals, int64(0)) 
    c.Assert(snap.uploadIdNum, Equals, int64(0))
    c.Assert(snap.errObjectNum, Equals, int64(0))
    c.Assert(snap.errUploadIdNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(0))

    str := strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)

    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d objects", 0)), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)

    // remove single exist object
    num := 2
    for i := 0; i < num; i++ {
        object := fmt.Sprintf("TestBatchRemoveProgress%d", i)
        s.putObject(bucket, object, uploadFileName, c)
    }

    err = s.initRemove(bucket, "TestBatchRemoveProgress", true, true, false)
    c.Assert(err, IsNil)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = removeCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, true)
    c.Assert(strings.Contains(pstr, fmt.Sprintf("total %d objects", num)), Equals, true)
    c.Assert(strings.Contains(pstr, "error"), Equals, false)

    snap = removeCommand.monitor.getSnapshot()
    c.Assert(snap.objectNum, Equals, int64(num)) 
    c.Assert(snap.uploadIdNum, Equals, int64(0))
    c.Assert(snap.errObjectNum, Equals, int64(0))
    c.Assert(snap.errUploadIdNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(num))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(removeCommand.monitor.getPrecent(snap), Equals, 100)

    str = strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", num)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, false)
    c.Assert(strings.Contains(str, "progress"), Equals, true)

    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("total %d objects", num)), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d objects", num)), Equals, true)
    c.Assert(strings.Contains(str, "err"), Equals, false)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)

    removeCommand.monitor.init() 
    removeCommand.updateObjectMonitor(0, 1)
    c.Assert(removeCommand.monitor.objectNum, Equals, int64(0))
    c.Assert(removeCommand.monitor.uploadIdNum, Equals, int64(0))
    c.Assert(removeCommand.monitor.errObjectNum, Equals, int64(1))
    c.Assert(removeCommand.monitor.errUploadIdNum, Equals, int64(0))

    str = strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.TrimSpace(str), Equals, "")

    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.TrimSpace(str), Equals, "")

    removeCommand.monitor.setOP(rmObject)
    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, false)
    c.Assert(strings.Contains(str, "when error happens"), Equals, true)
    c.Assert(strings.Contains(str, "scanned"), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d objects", 0)), Equals, true)

    str = strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("scanned %d objects", 1)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, true)
    c.Assert(strings.Contains(str, "progress"), Equals, false)
}

func (s *OssutilCommandSuite) TestRemoveUploadIdProgress(c *C) {
    removeCommand.monitor.init() 
    removeCommand.monitor.setOP(rmUploadId)
    removeCommand.monitor.updateUploadIdNum(2)
    removeCommand.monitor.updateErrUploadIdNum(1)
    c.Assert(removeCommand.monitor.objectNum, Equals, int64(0))
    c.Assert(removeCommand.monitor.uploadIdNum, Equals, int64(2))
    c.Assert(removeCommand.monitor.errObjectNum, Equals, int64(0))
    c.Assert(removeCommand.monitor.errUploadIdNum, Equals, int64(1))

    str := strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.Contains(str, fmt.Sprintf("scanned %d uploadids", 3)), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d uploadids", 2)), Equals, true)
    c.Assert(strings.Contains(str, "error"), Equals, true)
    c.Assert(strings.Contains(str, "progress"), Equals, false)
   
    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(str, "succeed:"), Equals, false)
    c.Assert(strings.Contains(str, "when error happens"), Equals, true)
    c.Assert(strings.Contains(str, "scanned"), Equals, true)
    c.Assert(strings.Contains(str, fmt.Sprintf("removed %d uploadids", 2)), Equals, true)
}

func (s *OssutilCommandSuite) TestRemoveBucketProgress(c *C) {
    // remove not exist bucket 
    err := s.initRemove(bucketNameNotExist, "", false, true, true)
    c.Assert(err, IsNil)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = removeCommand.RunCommand()
    c.Assert(err, NotNil)
    os.Stdout = out
    pstr := strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, "succeed"), Equals, false)

    bucket := bucketNamePrefix + "progress" 
    s.putBucket(bucket, c)
    //time.Sleep(10*time.Second)

    err = s.initRemove(bucket, "", false, true, true)
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out = os.Stdout
    os.Stdout = testResultFile
    err = removeCommand.RunCommand()
    c.Assert(err, IsNil)
    os.Stdout = out
    pstr = strings.ToLower(s.readFile(resultPath, c))
    c.Assert(strings.Contains(pstr, fmt.Sprintf("removed bucket: %s", bucket)), Equals, true)

    snap := removeCommand.monitor.getSnapshot()
    c.Assert(int64(removeCommand.monitor.op), Equals, int64(rmBucket))
    c.Assert(snap.objectNum, Equals, int64(0)) 
    c.Assert(snap.uploadIdNum, Equals, int64(0))
    c.Assert(snap.errObjectNum, Equals, int64(0))
    c.Assert(snap.errUploadIdNum, Equals, int64(0))
    c.Assert(snap.dealNum, Equals, int64(0))
    c.Assert(snap.errNum, Equals, int64(0))
    c.Assert(snap.removedBucket, Equals, bucket)

    str := strings.ToLower(removeCommand.monitor.getProgressBar())
    c.Assert(strings.TrimSpace(str), Equals, "")

    str = strings.ToLower(removeCommand.monitor.getFinishBar(normalExit))
    c.Assert(strings.Contains(pstr, fmt.Sprintf("removed bucket: %s", bucket)), Equals, true)
    c.Assert(strings.Contains(strings.TrimSpace(pstr), strings.TrimSpace(str)), Equals, true)
}