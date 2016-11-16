package lib 

import (
    "log"
    "os"
    "os/user"
    "time"
    "fmt"
    "io/ioutil"
    "strings"
    "testing"
    oss "github.com/aliyun/aliyun-oss-go-sdk/oss"

    . "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { 
    TestingT(t) 
}

type OssutilCommandSuite struct{}

var _ = Suite(&OssutilCommandSuite{})

var (
    // Update before running test
    endpoint         = "<testEndpoint>"
    accessKeyID      = "<testAccessKeyID>"
    accessKeySecret  = "<testAccessKeySecret>"
    stsToken         = "<testSTSToken>"
)

var (
    logPath             = "ossutil_test_" + time.Now().Format("20060102_150405") + ".log"
    configFile          = "ossutil_test.boto"
    testLogFile, _      = os.OpenFile(logPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    testLogger          = log.New(testLogFile, "", log.Ldate|log.Ltime|log.Lshortfile)
    resultPath          = "ossutil_test.result"
    testResultFile, _   = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    bucketNamePrefix    = "ossutil-test-"
    bucketNameExist     = "nodelete-ossutil-test-normalcase"
    bucketNameDest      = "nodelete-ossutil-test-dest"
    bucketNameNotExist  = bucketNamePrefix + "notexistbucket"
    uploadFileName      = "ossutil_test.upload_file"
    downloadFileName    = "ossutil_test.download_file"
    inputFileName       = "ossutil_test.input_file"
    content             = "abc"
    cm                  = CommandManager{}
    out                 = os.Stdout
    errout              = os.Stderr
    sleepTime           = 7*time.Second
)

// Run once when the suite starts running
func (s *OssutilCommandSuite) SetUpSuite(c *C) {
    os.Stdout = testLogFile 
    os.Stderr = testLogFile 
    testLogger.Println("test command started")
    SetUpCredential()
    cm.Init()
    s.configNonInteractive(c)
    s.createFile(uploadFileName, content, c)
    s.removeBuckets(bucketNamePrefix, c)
    s.removeBucket(bucketNameExist, true, c)
    s.removeBucket(bucketNameDest, true, c)
    time.Sleep(2*sleepTime)
    s.putBucket(bucketNameExist, c)
    s.putBucket(bucketNameDest, c)
    time.Sleep(2*sleepTime)
}

func SetUpCredential() {
    if endpoint == "<testEndpoint>" {
        endpoint = os.Getenv("OSS_TEST_ENDPOINT") 
    }
    if strings.HasPrefix(endpoint, "https://") {
        endpoint = endpoint[8:]
    }
    if strings.HasPrefix(endpoint, "http://") {
        endpoint = endpoint[7:]
    }
    if accessKeyID == "<testAccessKeyID>" {
        accessKeyID = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
    }
    if accessKeySecret == "<testAccessKeySecret>" {
        accessKeySecret = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
    }
    if ue := os.Getenv("OSS_TEST_UPDATE_ENDPOINT"); ue != "" {
        vUpdateEndpoint = ue
    }
    if ub := os.Getenv("OSS_TEST_UPDATE_BUCKET"); ub != "" {
        vUpdateBucket = ub
    }
    if strings.HasPrefix(vUpdateEndpoint, "https://") {
        vUpdateEndpoint = vUpdateEndpoint[8:]
    }
    if strings.HasPrefix(vUpdateEndpoint, "http://") {
        vUpdateEndpoint = vUpdateEndpoint[7:]
    }
}

// Run before each test or benchmark starts running
func (s *OssutilCommandSuite) TearDownSuite(c *C) {
    testLogger.Println("test command completed")
    _ = os.Remove(configFile)
    _ = os.Remove(resultPath)
    _ = os.Remove(uploadFileName)
    _ = os.Remove(downloadFileName)
    os.Stdout = out
    os.Stderr = errout
}

// Run after each test or benchmark runs
func (s *OssutilCommandSuite) SetUpTest(c *C) {
}

// Run once after all tests or benchmarks have finished running
func (s *OssutilCommandSuite) TearDownTest(c *C) {
}

func (s *OssutilCommandSuite) configNonInteractive(c *C) {
    command := "config" 
    var args []string
    options := OptionMapType{
        "endpoint": &endpoint,
        "accessKeyID": &accessKeyID,
        "accessKeySecret": &accessKeySecret,
        "configFile": &configFile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(showElapse, Equals, false)
    c.Assert(err, IsNil)

    opts, err := LoadConfig(configFile) 
    c.Assert(err, IsNil)
    c.Assert(len(opts), Equals, 4)
    c.Assert(opts[OptionLanguage], Equals, DefaultLanguage)
    c.Assert(opts[OptionEndpoint], Equals, endpoint)
    c.Assert(opts[OptionAccessKeyID], Equals, accessKeyID)
    c.Assert(opts[OptionAccessKeySecret], Equals, accessKeySecret)
}

func (s *OssutilCommandSuite) createFile(fileName, content string, c *C) {
    fout, err := os.Create(fileName)
    defer fout.Close()
    c.Assert(err, IsNil)
    _, err = fout.WriteString(content)
    c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) readFile(fileName string, c *C) (string) {
    f, err := ioutil.ReadFile(fileName)
    c.Assert(err, IsNil)
    return string(f)
}

func (s *OssutilCommandSuite) removeBuckets(prefix string, c *C) {
    buckets := s.listBuckets(false, c)
    for _, bucket := range buckets {
        if strings.HasPrefix(bucket, prefix) {
            s.removeBucket(bucket, true, c)
        }
    }
}

func (s *OssutilCommandSuite) listBuckets(shortFormat bool, c *C) ([]string) {
    var args []string
    out := os.Stdout
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    os.Stdout = testResultFile 
    showElapse, err := s.rawList(args, shortFormat, false)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    os.Stdout = out

    // get result
    buckets := s.getBucketResults(c) 
    _ = os.Remove(resultPath)
    return buckets
}

func (s *OssutilCommandSuite) getBucketResults(c *C) ([]string) {
    result := s.getResult(c)
    c.Assert(len(result) >= 1, Equals, true)
    buckets := []string{}
    for _, str := range result {
        pos := strings.Index(str, SchemePrefix)
        if pos != -1 {
            buckets = append(buckets, str[pos + len(SchemePrefix):])
        }
    }
    return buckets 
}

func (s *OssutilCommandSuite) getResult(c *C) ([]string) {
    str := s.readFile(resultPath, c)
    sli := strings.Split(str, "\n")
    result := []string{}
    for _, str := range sli {
        if str != ""{
            result = append(result, str)
        }
    }
    return result 
}

func (s *OssutilCommandSuite) removeBucket(bucket string, clearObjects bool, c *C) {
    args := []string{CloudURLToString(bucket, "")}
    showElapse, err := s.rawRemove(args, clearObjects, true, true)
    if err != nil {
        os.Stdout = out 
        os.Stderr = errout 
        fmt.Println(bucket, err)
        os.Stdout = testLogFile 
        os.Stderr = testLogFile 
        c.Assert(err.(oss.ServiceError).Code == "NoSuchBucket" || err.(oss.ServiceError).Code == "BucketAlreadyExist", Equals, true)
        c.Assert(showElapse, Equals, false)
    } else {
        c.Assert(showElapse, Equals, true)
    }
}

func (s *OssutilCommandSuite) removeObjects(bucket, prefix string, recursive, force bool, c *C) {
    args := []string{CloudURLToString(bucket, prefix)}
    showElapse, err := s.rawRemove(args, recursive, force, false)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) listObjects(bucket, prefix string, shortFormat, directory bool, c *C) ([]string) {
    args := []string{CloudURLToString(bucket, prefix)}
    out := os.Stdout
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    os.Stdout = testResultFile 
    showElapse, err := s.rawList(args, shortFormat, directory)
    os.Stdout = out
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // get result
    objects := s.getObjectResults(c) 
    _ = os.Remove(resultPath)
    return objects 
}

func (s *OssutilCommandSuite) getObjectResults(c *C) ([]string) {
    result := s.getResult(c)
    c.Assert(len(result) >= 1, Equals, true)
    objects := []string{}
    for _, str := range result {
        pos := strings.Index(str, SchemePrefix)
        if pos != -1 {
            url := str[pos:] 
            cloudURL, err := CloudURLFromString(url)
            c.Assert(err, IsNil)
            c.Assert(cloudURL.object != "", Equals, true)
            objects = append(objects, cloudURL.object)
        }
    }
    return objects 
}

func (s *OssutilCommandSuite) putBucket(bucket string, c *C) {
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
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) putObject(bucket, object, fileName string, c *C) {
    args := []string{fileName, CloudURLToString(bucket, object)}
    showElapse, err := s.rawCPWithArgs(args, false, true, false, BigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) getObject(bucket, object, fileName string, c *C) {
    args := []string{CloudURLToString(bucket, object), fileName}
    showElapse, err := s.rawCPWithArgs(args, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) copyObject(srcBucket, srcObject, destBucket, destObject string, c *C) {
    args := []string{CloudURLToString(srcBucket, srcObject), CloudURLToString(destBucket, destObject)}
    showElapse, err := s.rawCPWithArgs(args, false, true, false, BigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) getStat(bucket, object string, c *C) (map[string]string) {
    args := []string{CloudURLToString(bucket, object)}
    out := os.Stdout
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    os.Stdout = testResultFile 
    showElapse, err := s.rawGetStatWithArgs(args)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    os.Stdout = out

    // get result
    stat := s.getStatResults(c)
    _ = os.Remove(resultPath)
    return stat 
}

func (s *OssutilCommandSuite) getStatResults(c *C) (map[string]string) {
    result := s.getResult(c)
    c.Assert(len(result) > 1, Equals, true)
    
    stat := map[string]string{}
    for _, str := range result {
        sli := strings.SplitN(str, ":", 2)
        if len(sli) == 2 {
            stat[strings.TrimSpace(sli[0])] = strings.TrimSpace(sli[1])
        }
    }
    return stat 
}

func (s *OssutilCommandSuite) getHashResults(c *C) (map[string]string) {
    result := s.getResult(c)
    c.Assert(len(result) >= 1, Equals, true)
    
    stat := map[string]string{}
    for _, str := range result {
        sli := strings.SplitN(str, ":", 2)
        if len(sli) == 2 {
            stat[strings.TrimSpace(sli[0])] = strings.TrimSpace(sli[1])
        }
    }
    return stat 
}

func (s *OssutilCommandSuite) setBucketACL(bucket, acl string, c *C) {
    args := []string{CloudURLToString(bucket, ""), acl}
    showElapse, err := s.rawSetACLWithArgs(args, false, true, false)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) setObjectACL(bucket, object, acl string, recursive, force bool, c *C) {
    args := []string{CloudURLToString(bucket, object), acl}
    showElapse, err := s.rawSetACLWithArgs(args, recursive, false, force)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) TestParseOptions(c *C) {
    bucket := bucketNameExist 
    s.putBucket(bucket, c)

    s.createFile(uploadFileName, content, c)

    // put object
    object := "中文" 
    v := []string{"", "cp", uploadFileName, CloudURLToString(bucket, object), "-f", "--update", "--bigfile-threshold=1", "--checkpoint-dir=checkpoint_dir", "-c", configFile}
    os.Args = v
    err := ParseAndRunCommand()
    c.Assert(err, IsNil)

    // get object
    s.getObject(bucket, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, content)
}

func (s *OssutilCommandSuite) TestNotExistCommand(c *C) {
    command := "notexistcmd"
    args := []string{}
    options := OptionMapType{}
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestDecideConfigFile(c *C) {
    usr, _ := user.Current()
    file := DecideConfigFile("")
    c.Assert(file, Equals, strings.Replace(DefaultConfigFile, "~", usr.HomeDir, 1))
    input := "~/a"
    file = DecideConfigFile(input)
    c.Assert(file, Equals, strings.Replace(input, "~", usr.HomeDir, 1))
}

func (s *OssutilCommandSuite) TestCheckConfig(c *C) {
    // config file not exist
    configMap := OptionMapType{OptionRetryTimes: "abc"}
    err := checkConfig(configMap)
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestOptions(c *C) {
    option := Option{"", "", "", OptionTypeString, "", "", "", ""}
    _, err := stringOption(option)
    c.Assert(err, NotNil)

    option = Option{"", "", "", OptionTypeFlagTrue, "", "", "", ""}
    _, err = flagTrueOption(option)
    c.Assert(err, NotNil)

    option = Option{"-a", "", "", OptionTypeFlagTrue, "", "", "", ""}
    _, err = flagTrueOption(option)
    c.Assert(err, IsNil)

    str := "abc"
    options := OptionMapType{OptionRetryTimes: &str}
    err = checkOption(options)
    c.Assert(err, NotNil)

    str = "-1"
    options = OptionMapType{OptionRetryTimes: &str}
    err = checkOption(options)
    c.Assert(err, NotNil)

    str = "1001"
    options = OptionMapType{OptionRetryTimes: &str}
    err = checkOption(options)
    c.Assert(err, NotNil)

    language := "unknown"
    options = OptionMapType{OptionLanguage: &language}
    err = checkOption(options)
    c.Assert(err, NotNil)

    options = OptionMapType{OptionConfigFile: &configFile}
    ok, err := GetBool(OptionConfigFile, options)
    c.Assert(err, NotNil)
    c.Assert(ok, Equals, false)
    
    i, err := GetInt(OptionConfigFile, options)
    c.Assert(err, NotNil)
    c.Assert(i, Equals, int64(0))
 
    str = ""
    options = OptionMapType{OptionConfigFile: &str}
    i, err = GetInt(OptionConfigFile, options)
    c.Assert(err, NotNil)
    c.Assert(i, Equals, int64(0))
 
    options = OptionMapType{OptionRetryTimes: &str}
    i, err = GetInt(OptionConfigFile, options)
    c.Assert(err, NotNil)
    c.Assert(i, Equals, int64(0))

    ok = true
    options = OptionMapType{OptionRetryTimes: &ok}
    i, err = GetInt(OptionConfigFile, options)
    c.Assert(err, NotNil)
    c.Assert(i, Equals, int64(0))

    options = OptionMapType{OptionConfigFile: &ok}
    i, err = GetInt(OptionConfigFile, options)
    c.Assert(err, NotNil)
    c.Assert(i, Equals, int64(0))

    options = OptionMapType{OptionConfigFile: "a"}
    val, err := GetString(OptionConfigFile, options)
    c.Assert(err, NotNil)
    c.Assert(val, Equals, "")
}

func (s *OssutilCommandSuite) TestErrors(c *C) {
    err := CommandError{"help", "abc"}
    c.Assert(err.Error(), Equals, "invalid usage of \"help\" command, reason: abc, please try \"help help\" for more information", )

    berr := BucketError{err, "b"}
    c.Assert(berr.Error(), Equals, fmt.Sprintf("%s, Bucket=%s", err.Error(), "b"))

    ferr := FileError{err, "f"}
    c.Assert(ferr.Error(), Equals, fmt.Sprintf("%s, File=%s", err.Error(), "f"))
}

func (s *OssutilCommandSuite) TestStorageURL(c *C) {
    var cloudURL CloudURL
    err := cloudURL.Init("/abc/d")
    c.Assert(err, IsNil)
    c.Assert(cloudURL.bucket, Equals, "abc")
    c.Assert(cloudURL.object, Equals, "d")

    usr, _ := user.Current()
    dir := usr.HomeDir
    url := "~/test"
    var fileURL FileURL
    fileURL.Init(url)
    c.Assert(fileURL.url, Equals, strings.Replace(url, "~", dir, 1))
    
    _, err = CloudURLFromString("oss:///object")
    c.Assert(err, NotNil)

    _, err = CloudURLFromString("./file")
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestErrOssDownloadFile(c *C) {
    bucketName := bucketNamePrefix + "b1"
	bucket, err := copyCommand.command.ossBucket(bucketName)
    c.Assert(err, IsNil)

    object := "object"
    err = copyCommand.command.ossDownloadFileRetry(bucket, object, object)
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestUserAgent(c *C) {
    userAgent := getUserAgent()
    c.Assert(userAgent != "", Equals, true)

    client, err := listCommand.command.ossClient("")
    c.Assert(err, IsNil)
    c.Assert(client, NotNil)
}
