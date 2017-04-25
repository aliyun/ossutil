package lib

import (
	"fmt"
	"strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type batchOptionType struct {
	ctnu     bool
	reporter *Reporter
}

var specChineseRestore = SpecText{

	synopsisText: "恢复冷冻状态的Objects为可读状态",

	paramText: "cloud_url [options]",

	syntaxText: ` 
    ossutil restore cloud_url [--encoding-type url] [-r] [-f] [--output-dir=odir] [-c file] 
`,

	detailHelpText: ` 
    该命令恢复处于冷冻状态的归档类型object进入可读状态，即操作对象object必须为` + StorageArchive + `存储
    类型的object。

    如果是针对处于冷冻状态的归档类型object第一次调用restore接口，则返回成功。
    如果已经成功调用过restore接口，且restore没有完全完成，再次调用时同样成功。
    如果已经成功调用过restore接口，且restore已经完成，再次调用时返回成功，且会将object
    的可下载时间延长一天，最多延长7天。


用法：

    该命令有两种用法：

    1) ossutil restore oss://bucket/object [--encoding-type url] 
        该用法恢复单个冷冻状态object为可读状态，当指定object不存在时，ossutil会提示错
    误，此时请确保指定的url精确匹配需要设置acl的object，并且不要指定--recursive选项（
    否则ossutil会进行前缀匹配，恢复多个冷冻状态的objects为可读状态）。无论--force选项
    是否指定，都不会进行询问提示。

    2) ossutil restore oss://bucket[/prefix] -r [--encoding-type url] [-f] [--output-dir=odir]
        该用法可批量恢复多个冷冻状态的objects为可读状态，此时必须输入--recursive选项，
    ossutil会查找所有前缀匹配url的objects，恢复它们为可读状态。当一个object操作出现错
    误时，会将出错object的错误信息记录到report文件，并继续操作其他object，成功操作的
    object信息将不会被记录到report文件中（更多信息见cp命令的帮助）。如果--force选项被
    指定，则不会进行询问提示。
`,

	sampleText: ` 
    1) ossutil restore oss://bucket-restore/object-store
    2) ossutil restore oss://bucket-restore/object-prefix -r
    3) ossutil restore oss://bucket-restore/object-prefix -r -f
    4) ossutil restore oss://bucket-restore/%e4%b8%ad%e6%96%87 --encoding-type url
`,
}

var specEnglishRestore = SpecText{

	synopsisText: "Restore Frozen State Object to Read Ready Status",

	paramText: "cloud_url [options]",

	syntaxText: ` 
    ossutil restore cloud_url [--encoding-type url] [-r] [-f] [--output-dir=odir] [-c file] 
`,

	detailHelpText: ` 
    The command restore frozen state object to read ready status, the object must be in 
    the storage class of ` + StorageArchive + `. 

    If it's the first time to restore a frozen state object, the operation will success.
    If the object is in restoring, and the restoring is not finished, do the operation 
    again will success.
    If an object has been restored, do the operation again will success, and the time that 
    the object can be downloaded will extend one day, we can at most extend seven days. 


Usage:

    There are two usages:

    1) ossutil restore oss://bucket/object [--encoding-type url] 
        If --recursive option is not specified, ossutil restore the specified single frozen state 
    object to read ready status. In the usage, please make sure url exactly specified the object 
    you want to restore, if object not exist, error occurs. No matter --force option is specified 
    or not, ossutil will not show prompt question. 

    2) ossutil restore oss://bucket[/prefix] -r [--encoding-type url] [-f] [--output-dir=odir]
        The usage can restore many objects that in frozen state to read ready status, --recursive 
    option is required for the usage, ossutil will search for prefix-matching objects and restore 
    those objects. When an error occurs when restore an object, ossutil will record the error message 
    to report file, and ossutil will continue to attempt to set acl on the remaining objects(more 
    information see help of cp command). If --force option is specified, ossutil will not show 
    prompt question. 
`,

	sampleText: ` 
    1) ossutil restore oss://bucket-restore/object-store
    2) ossutil restore oss://bucket-restore/object-prefix -r
    3) ossutil restore oss://bucket-restore/object-prefix -r -f
    4) ossutil restore oss://bucket-restore/%e4%b8%ad%e6%96%87 --encoding-type url
`,
}

// RestoreCommand is the command list buckets or objects
type RestoreCommand struct {
	command  Command
	monitor  Monitor
	reOption batchOptionType
}

var restoreCommand = RestoreCommand{
	command: Command{
		name:        "restore",
		nameAlias:   []string{},
		minArgc:     1,
		maxArgc:     1,
		specChinese: specChineseRestore,
		specEnglish: specEnglishRestore,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionRecursion,
			OptionForce,
			OptionEncodingType,
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionRetryTimes,
			OptionRoutines,
			OptionOutputDir,
		},
	},
}

// function for FormatHelper interface
func (rc *RestoreCommand) formatHelpForWhole() string {
	return rc.command.formatHelpForWhole()
}

func (rc *RestoreCommand) formatIndependHelp() string {
	return rc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (rc *RestoreCommand) Init(args []string, options OptionMapType) error {
	return rc.command.Init(args, options, rc)
}

// RunCommand simulate inheritance, and polymorphism
func (rc *RestoreCommand) RunCommand() error {
	rc.monitor.init("Restored")

	encodingType, _ := GetString(OptionEncodingType, rc.command.options)
	recursive, _ := GetBool(OptionRecursion, rc.command.options)

	cloudURL, err := CloudURLFromString(rc.command.args[0], encodingType)
	if err != nil {
		return err
	}

	if err = rc.checkArgs(cloudURL, recursive); err != nil {
		return err
	}

	bucket, err := rc.command.ossBucket(cloudURL.bucket)
	if err != nil {
		return err
	}

	if !recursive {
		return rc.ossRestoreObject(bucket, cloudURL.object)
	}
	return rc.batchRestoreObjects(bucket, cloudURL)
}

func (rc *RestoreCommand) checkArgs(cloudURL CloudURL, recursive bool) error {
	if cloudURL.bucket == "" {
		return fmt.Errorf("invalid cloud url: %s, miss bucket", rc.command.args[0])
	}
	if !recursive && cloudURL.object == "" {
		return fmt.Errorf("restore object invalid cloud url: %s, object empty. Restore bucket is not supported, if you mean batch restore objects, please use --recursive", rc.command.args[0])
	}
	return nil
}

func (rc *RestoreCommand) ossRestoreObject(bucket *oss.Bucket, object string) error {
	retryTimes, _ := GetInt(OptionRetryTimes, rc.command.options)
	for i := 1; ; i++ {
		err := bucket.RestoreObject(object)
		if err == nil {
			return err
		}

		switch err.(type) {
		case oss.ServiceError:
			if err.(oss.ServiceError).StatusCode == 409 && err.(oss.ServiceError).Code == "RestoreAlreadyInProgress" {
				return nil
			}
		}

		if int64(i) >= retryTimes {
			return ObjectError{err, bucket.BucketName, object}
		}
	}
}

func (rc *RestoreCommand) batchRestoreObjects(bucket *oss.Bucket, cloudURL CloudURL) error {
	force, _ := GetBool(OptionForce, rc.command.options)
	if !force {
		var val string
		fmt.Printf("Do you really mean to recursivlly restore objects of %s(y or N)? ", rc.command.args[0])
		if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
			fmt.Println("operation is canceled.")
			return nil
		}
	}

	rc.reOption.ctnu = true
	outputDir, _ := GetString(OptionOutputDir, rc.command.options)

	// init reporter
	var err error
	if rc.reOption.reporter, err = GetReporter(rc.reOption.ctnu, outputDir, commandLine); err != nil {
		return err
	}
	defer rc.reOption.reporter.Clear()

	return rc.restoreObjects(bucket, cloudURL)
}

func (rc *RestoreCommand) restoreObjects(bucket *oss.Bucket, cloudURL CloudURL) error {
	routines, _ := GetInt(OptionRoutines, rc.command.options)

	chObjects := make(chan string, ChannelBuf)
	chError := make(chan error, routines+1)
	chListError := make(chan error, 1)
	go rc.command.objectStatistic(bucket, cloudURL, &rc.monitor)
	go rc.command.objectProducer(bucket, cloudURL, chObjects, chListError)
	for i := 0; int64(i) < routines; i++ {
		go rc.restoreConsumer(bucket, cloudURL, chObjects, chError)
	}

	return rc.waitRoutinueComplete(chError, chListError, routines)
}

func (rc *RestoreCommand) restoreConsumer(bucket *oss.Bucket, cloudURL CloudURL, chObjects <-chan string, chError chan<- error) {
	for object := range chObjects {
		err := rc.restoreObjectWithReport(bucket, object)
		if err != nil {
			chError <- err
			if !rc.reOption.ctnu {
				return
			}
			continue
		}
	}

	chError <- nil
}

func (rc *RestoreCommand) restoreObjectWithReport(bucket *oss.Bucket, object string) error {
	err := rc.ossRestoreObject(bucket, object)
	rc.command.updateMonitor(err, &rc.monitor)
	msg := fmt.Sprintf("restore %s", CloudURLToString(bucket.BucketName, object))
	rc.command.report(msg, err, &rc.reOption)
	return err
}

func (rc *RestoreCommand) waitRoutinueComplete(chError, chListError <-chan error, routines int64) error {
	completed := 0
	var ferr error
	for int64(completed) <= routines {
		select {
		case err := <-chListError:
			if err != nil {
				return err
			}
			completed++
		case err := <-chError:
			if err == nil {
				completed++
			} else {
				ferr = err
				if !rc.reOption.ctnu {
					fmt.Printf(rc.monitor.progressBar(true, errExit))
					return err
				}
			}
		}
	}
	return rc.formatResultPrompt(ferr)
}

func (rc *RestoreCommand) formatResultPrompt(err error) error {
	fmt.Printf(rc.monitor.progressBar(true, normalExit))
	if err != nil && rc.reOption.ctnu {
		return nil
	}
	return err
}