package lib 

import (
    "log"
    "os"
    "os/user"
    "time"
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    "path/filepath"
    "math/rand"
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
    ConfigFile          = "ossutil_test.boto" + randStr(5)
    configFile          = ConfigFile 
    testLogFile, _      = os.OpenFile(logPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    testLogger          = log.New(testLogFile, "", log.Ldate|log.Ltime|log.Lshortfile)
    resultPath          = "ossutil_test.result" + randStr(5)
    testResultFile, _   = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    uploadFileName      = "ossutil_test.upload_file" + randStr(5)
    downloadFileName    = "ossutil_test.download_file" + randStr(5)
    downloadDir         = "ossutil_test.download_dir" + randStr(5)
    inputFileName       = "ossutil_test.input_file" + randStr(5)
    content             = "abc"
    cm                  = CommandManager{}
    out                 = os.Stdout
    errout              = os.Stderr
    sleepTime           = time.Second
)

var (
    bucketNamePrefix    = "ossutil-test-" + randLowStr(6)
    bucketNameExist     = "special-" + bucketNamePrefix + "existbucket" 
    bucketNameDest      = "special-" + bucketNamePrefix + "destbucket" 
    bucketNameNotExist  = "nodelete-ossutil-test-notexist"
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
    s.SetUpBucketEnv(c)
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

func (s *OssutilCommandSuite) SetUpBucketEnv(c *C) {
    s.removeBuckets(bucketNamePrefix, c)
    s.putBucket(bucketNameExist, c)
    s.putBucket(bucketNameDest, c)
}

// Run before each test or benchmark starts running
func (s *OssutilCommandSuite) TearDownSuite(c *C) {
    s.removeBuckets(bucketNamePrefix, c)
    s.removeBucket(bucketNameExist, true, c)
    s.removeBucket(bucketNameDest, true, c)
    testLogger.Println("test command completed")
    _ = os.Remove(configFile)
    _ = os.Remove(configFile+".bak")
    _ = os.Remove(resultPath)
    _ = os.Remove(uploadFileName)
    _ = os.Remove(downloadFileName)
    _ = os.RemoveAll(downloadDir)
    _ = os.RemoveAll(DefaultOutputDir)
    os.Stdout = out
    os.Stderr = errout
}

// Run after each test or benchmark runs
func (s *OssutilCommandSuite) SetUpTest(c *C) {
    configFile = ConfigFile 
}

// Run once after all tests or benchmarks have finished running
func (s *OssutilCommandSuite) TearDownTest(c *C) {
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
    b := make([]rune, n)
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := range b {
        b[i] = letters[r.Intn(len(letters))]
    }
    return string(b)
}

func randLowStr(n int) string {
    return strings.ToLower(randStr(n))
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

func (s *OssutilCommandSuite) rawList(args []string, cmdline string) (bool, error) {
    array := strings.Split(cmdline, " ")
    if len(array) < 2 {
        return false, fmt.Errorf("ls test wrong cmdline given")
    }

    parameter := strings.Split(array[1], "-")
    if len(parameter) < 2 {
        return false, fmt.Errorf("ls test wrong cmdline given")
    }

    command := array[0]
    sf := strings.Contains(parameter[1], "s")
    d := strings.Contains(parameter[1], "d")
    m := strings.Contains(parameter[1], "m")
    a := strings.Contains(parameter[1], "a")

    str := ""
    options := OptionMapType{
        "endpoint":        &str,
        "accessKeyID":     &str,
        "accessKeySecret": &str,
        "stsToken":        &str,
        "configFile":      &configFile,
        "shortFormat":     &sf,
        "directory":       &d,
        "multipart":       &m,
        "allType":         &a,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) listBuckets(shortFormat bool, c *C) []string {
    var args []string
    out := os.Stdout
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    os.Stdout = testResultFile
    showElapse, err := s.rawList(args, "ls -s")
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
    return s.getFileResult(resultPath, c)
}

func (s *OssutilCommandSuite) getFileResult(fileName string, c *C) ([]string) {
    str := s.readFile(fileName, c)
    sli := strings.Split(str, "\n")
    result := []string{}
    for _, str := range sli {
        if str != ""{
            result = append(result, str)
        }
    }
    return result 
}

func (s *OssutilCommandSuite) getReportResult(fileName string, c *C) ([]string) {
    result := s.getFileResult(fileName, c)
    c.Assert(len(result) >= 1, Equals, true)
    c.Assert(strings.HasPrefix(result[0], "#"), Equals, true)
    result = result[1:]
    for _, r := range result {
        c.Assert(strings.HasPrefix(r, "[Error]"), Equals, true)
    }
    return result 
}

func (s *OssutilCommandSuite) removeBucket(bucket string, clearObjects bool, c *C) {
    args := []string{CloudURLToString(bucket, "")}
    var showElapse bool
    var err error
    if !clearObjects {
        showElapse, err = s.rawRemove(args, false, true, true)
    } else {
        showElapse, err = s.removeWrapper("rm -arfb", bucket, "", c)
    }
    if err != nil {
        verr := err.(BucketError).err
        c.Assert(verr.(oss.ServiceError).Code == "NoSuchBucket" || verr.(oss.ServiceError).Code == "BucketNotEmpty", Equals, true)
        c.Assert(showElapse, Equals, false)
    } else {
        c.Assert(showElapse, Equals, true)
    }
}

func (s *OssutilCommandSuite) rawRemove(args []string, recursive, force, bucket bool) (bool, error) {
    command := "rm"
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "recursive": &recursive,
        "force": &force,
        "bucket": &bucket,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    time.Sleep(sleepTime)
    return showElapse, err
}

func (s *OssutilCommandSuite) removeWrapper(cmdline string, bucket string, object string, c *C) (bool, error) {
    array := strings.Split(cmdline, " ")
    if len(array) < 2 {
        return false, fmt.Errorf("rm test wrong cmdline given")
    }

    parameter := strings.Split(array[1], "-")
    if len(parameter) < 2 {
        return false, fmt.Errorf("rm test wrong cmdline given")
    }

    command := array[0]
    a := strings.Contains(parameter[1], "a")
    m := strings.Contains(parameter[1], "m")
    b := strings.Contains(parameter[1], "b")
    r := strings.Contains(parameter[1], "r")
    f := strings.Contains(parameter[1], "f")

    args := []string{CloudURLToString(bucket, object)}
    str := ""
    options := OptionMapType{
        "endpoint":        &str,
        "accessKeyID":     &str,
        "accessKeySecret": &str,
        "stsToken":        &str,
        "configFile":      &configFile,
        "bucket":          &b,
        "allType":         &a,
        "multipart":       &m,
        "recursive":       &r,
        "force":           &f,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    time.Sleep(sleepTime)
    return showElapse, err
}

func (s *OssutilCommandSuite) initRemove(bucket string, object string, cmdline string) error {
    array := strings.Split(cmdline, " ")
    if len(array) < 2 {
        return fmt.Errorf("rm test wrong cmdline given")
    }

    parameter := strings.Split(array[1], "-")
    if len(parameter) < 2 {
        return fmt.Errorf("rm test wrong cmdline given")
    }

    a := strings.Contains(parameter[1], "a")
    m := strings.Contains(parameter[1], "m")
    b := strings.Contains(parameter[1], "b")
    r := strings.Contains(parameter[1], "r")
    f := strings.Contains(parameter[1], "f")

    args := []string{CloudURLToString(bucket, object)}
    str := ""
    options := OptionMapType{
        "endpoint":        &str,
        "accessKeyID":     &str,
        "accessKeySecret": &str,
        "stsToken":        &str,
        "configFile":      &configFile,
        "bucket":          &b,
        "allType":         &a,
        "multipart":       &m,
        "recursive":       &r,
        "force":           &f,
    }
    err := removeCommand.Init(args, options)
    return err
}

func (s *OssutilCommandSuite) removeObjects(bucket, prefix string, recursive, force bool, c *C) {
    args := []string{CloudURLToString(bucket, prefix)}
    showElapse, err := s.rawRemove(args, recursive, force, false)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) clearObjects(bucket, prefix string, c *C) {
    showElapse, err := s.removeWrapper("rm -afr", bucket, prefix, c)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) listObjects(bucket, prefix string, cmdline string, c *C) []string {
    args := []string{CloudURLToString(bucket, prefix)}
    out := os.Stdout
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    os.Stdout = testResultFile
    showElapse, err := s.rawList(args, cmdline)
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
    time.Sleep(sleepTime)
}

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

func (s *OssutilCommandSuite) rawCPWithOutputDir(srcURL, destURL string, recursive, force, update bool, threshold int64, outputDir string) (bool, error) {
    command := "cp"
    str := ""
    args := []string{srcURL, destURL}
    thre := strconv.FormatInt(threshold, 10)
    routines := strconv.Itoa(Routines)
    cpDir := CheckpointDir
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
        "outputDir": &outputDir,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) initCopyCommand(srcURL, destURL string, recursive, force, update bool, threshold int64, cpDir, outputDir string) error {
    str := ""
    args := []string{srcURL, destURL}
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
        "outputDir": &outputDir,
    }
    err := copyCommand.Init(args, options)
    return err
}

func (s *OssutilCommandSuite) initCopyWithSnapshot(srcURL, destURL string, recursive, force, update bool, threshold int64, snapshotPath string) error {
    str := ""
    args := []string{srcURL, destURL}
    thre := strconv.FormatInt(threshold, 10)
    routines := strconv.Itoa(Routines)
    cpDir := CheckpointDir
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
        "snapshotPath": &snapshotPath,
    }
    err := copyCommand.Init(args, options)
    return err
}

func (s *OssutilCommandSuite) putObject(bucket, object, fileName string, c *C) {
    args := []string{fileName, CloudURLToString(bucket, object)}
    showElapse, err := s.rawCPWithArgs(args, false, true, false, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    time.Sleep(sleepTime)
}

func (s *OssutilCommandSuite) getObject(bucket, object, fileName string, c *C) {
    args := []string{CloudURLToString(bucket, object), fileName}
    showElapse, err := s.rawCPWithArgs(args, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) copyObject(srcBucket, srcObject, destBucket, destObject string, c *C) {
    args := []string{CloudURLToString(srcBucket, srcObject), CloudURLToString(destBucket, destObject)}
    showElapse, err := s.rawCPWithArgs(args, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) rawGetStat(bucket, object string) (bool, error) {
    args := []string{CloudURLToString(bucket, object)}
    showElapse, err := s.rawGetStatWithArgs(args)
    return showElapse, err 
}

func (s *OssutilCommandSuite) rawGetStatWithArgs(args []string) (bool, error) {
    command := "stat"
    str := ""
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err 
}

func (s *OssutilCommandSuite) getStat(bucket, object string, c *C) (map[string]string) {
    args := []string{CloudURLToString(bucket, object)}
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
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
    return showElapse, err
}

func (s *OssutilCommandSuite) initSetACL(bucket, object, acl string, recursive, tobucket, force bool) error {
    args := []string{CloudURLToString(bucket, object), acl}
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
        "bucket": &tobucket,
        "force": &force,
    }
    err := setACLCommand.Init(args, options)
    return err
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
    return showElapse, err
}

func (s *OssutilCommandSuite) setObjectMeta(bucket, object, meta string, update, delete, recursive, force bool, c *C) {
    showElapse, err := s.rawSetMeta(bucket, object, meta, update, delete, recursive, force, DefaultLanguage) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
}

func (s *OssutilCommandSuite) initSetMeta(bucket, object, meta string, update, delete, recursive, force bool, language string) error {
    args := []string{CloudURLToString(bucket, object), meta}
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
    err := setMetaCommand.Init(args, options)
    return err
}

func (s *OssutilCommandSuite) getFileList(dpath string) ([]string, error) {
    fileList := []string{}
    err := filepath.Walk(dpath, func(fpath string, f os.FileInfo, err error) error {
        if f == nil {
            return err
        }

        dpath = filepath.Clean(dpath)
        fpath = filepath.Clean(fpath)
        fileName, err := filepath.Rel(dpath, fpath) 
        if err != nil {
            return fmt.Errorf("list file error: %s, info: %s", fpath, err.Error())
        }

        if f.IsDir(){
            if fpath != dpath {
                fileList = append(fileList, fileName + string(os.PathSeparator))
            }
            return nil
        }
        fileList = append(fileList, fileName)
        return nil
    })
    return fileList, err
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
    err = copyCommand.ossDownloadFileRetry(bucket, object, object)
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestUserAgent(c *C) {
    userAgent := getUserAgent()
    c.Assert(userAgent != "", Equals, true)

    client, err := listCommand.command.ossClient("")
    c.Assert(err, IsNil)
    c.Assert(client, NotNil)
}

func (s *OssutilCommandSuite) TestParseAndRunCommand(c *C) {
    args := []string{}
    options := OptionMapType{}
    showElapse, err := RunCommand(args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestGetSizeString(c *C) {
    c.Assert(getSizeString(0), Equals, "0")
    c.Assert(getSizeString(1), Equals, "1")
    c.Assert(getSizeString(12), Equals, "12")
    c.Assert(getSizeString(123), Equals, "123")
    c.Assert(getSizeString(1234), Equals, "1,234")
    c.Assert(getSizeString(12345), Equals, "12,345")
    c.Assert(getSizeString(123456), Equals, "123,456")
    c.Assert(getSizeString(1234567), Equals, "1,234,567")
    c.Assert(getSizeString(123456789012), Equals, "123,456,789,012")
    c.Assert(getSizeString(1234567890123), Equals, "1,234,567,890,123")
    c.Assert(getSizeString(-0), Equals, "0")
    c.Assert(getSizeString(-1), Equals, "-1")
    c.Assert(getSizeString(-12), Equals, "-12")
    c.Assert(getSizeString(-123), Equals, "-123")
    c.Assert(getSizeString(-1234), Equals, "-1,234")
    c.Assert(getSizeString(-12345), Equals, "-12,345")
    c.Assert(getSizeString(-123456), Equals, "-123,456")
    c.Assert(getSizeString(-1234567), Equals, "-1,234,567")
    c.Assert(getSizeString(-123456789012), Equals, "-123,456,789,012")
    c.Assert(getSizeString(-1234567890123), Equals, "-1,234,567,890,123")
}
