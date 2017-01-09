package lib

import (
	"fmt"
	"strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseList = SpecText{

	synopsisText: "列举Buckets或者Objects",

	paramText: "[url] [options]",

	syntaxText: ` 
    ossutil ls [oss://bucket[/prefix]] [-s] [-d] [-c file] 
`,

	detailHelpText: ` 
    该命令列举指定身份凭证下的buckets，或该身份凭证下对应endpoint的objects。默认显示
    长格式，ossutil在列举buckets或者objects的同时展示它们的一些附加信息。如果指定了
    --short-format选项，则显示精简格式。

    对于用户使用multipart方式上传且未complete的object，ossutil在显示objects或者目录时，
    不会显示这些objects。（关于multipart的更多信息请查看oss官网API文档。）

用法：

    该命令有两种用法：

    1) ossutil ls [oss://] [-s] [-m] [-a]
        如果用户列举时缺失url参数，则ossutil获取用户的身份凭证信息（从配置文件中读取），
    并列举该身份凭证下的所有buckets，并显示每个bucket的最新更新时间和位置信息。如果指定
    了--short-format选项则只输出bucket名称。该用法不支持--directory选项。--multipart选
    项只输出碎片名称(未completed的multipart)及其对应uploadID。--all-type选项显示普通的
    object和碎片。

    2) ossutil ls oss://bucket[/prefix] [-s] [-d] [-m] [-a]
        该用法列举指定bucket下的objects（如果指定了前缀，则列举拥有该前缀的objects），同时
    展示了object大小，最新更新时间和etag，但是如果指定了--short-format选项则只输出object名
    称。如果指定了--directory选项，则返回指定bucket下以指定前缀开头的第一层目录下的文件和子
    目录，但是不递归显示所有子目录，此时默认为精简格式。如果指定了--multipart选项，则返回指
    定bucket下以指定前缀开头的第一层目录下的碎片名称(未completed的multipart)和及其对应uploadID。
    --all-type选项显示普通的object和碎片。

`,

	sampleText: ` 
    1)ossutil ls -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    2)ossutil ls oss:// -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    3)ossutil ls oss://bucket1 -s
        oss://bucket1/obj1
        oss://bucket1/dir1/obj11
        Object Number is: 2

    4)ossutil ls oss://bucket1
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2016-04-08 14:50:47 +0000 UTC 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/obj1
        2015-06-05 14:06:29 +0000 UTC  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
        Object Number is: 2

    5)ossutil ls oss://bucket1 -d
        oss://bucket1/obj1
        oss://bucket1/dir1
        Object or Directory Number is: 2

    6)ossutil ls oss://bucket1 -m 
        UploadID                            MultipartName
        15754AF7980C4DFB8193F190837520BB    oss://bucket1/obj1
        3998971ACAF94AD9AC48EAC1988BE863    oss://bucket2/obj2
        UploadId Number is: 2
    
    7)ossutil ls oss://bucket1 -a 
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2016-04-08 14:50:47 +0000 UTC 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/sample.txt
        2015-06-05 14:06:29 +0000 UTC  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
        Object Number is: 2
        UploadID                            MultipartName
        15754AF7980C4DFB8193F190837520BB    oss://bucket1/obj1
        3998971ACAF94AD9AC48EAC1988BE863    oss://bucket2/obj2
        Multipart Number is: 2
         
`,
}

var specEnglishList = SpecText{

	synopsisText: "List Buckets or Objects",

	paramText: "[url] [options]",

	syntaxText: ` 
    ossutil ls [oss://bucket[/prefix]] [-s] [-d] [-c file] 
`,

	detailHelpText: ` 
    The command list buckets of the specified credentials. or objects of the specified 
    endpoint and credentials, with simple additional information, about each matching 
    provider, bucket, subdirectory, or object. If --short-format option is specified, 
    ossutil will show by short format. 

Usage:

    There are two usages:

    1) ossutil ls [oss://] [-s]
        If you list without a url, ossutil lists all the buckets using the credentials
    in config file with last modified time and location in addition. --show_format option 
    will ignore last modified time and location. The usage do not support --directory 
    option.

    2) ossutil ls oss://bucket[/prefix] [-s] [-d]
        The usage list objects under the specified bucket(with the prefix if you specified), 
    with object size, last modified time and etag in addition, --short-format option ignores 
    all the additional information. --directory option returns top-level subdirectory names 
    instead of contents of the subdirectory, which in default show by short format. 
`,

	sampleText: ` 
    1)ossutil ls -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    2)ossutil ls oss:// -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    3)ossutil ls oss://bucket1 -s
        oss://bucket1/obj1
        oss://bucket1/dir1/obj11
        Object Number is: 2

    4)ossutil ls oss://bucket1
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2016-04-08 14:50:47 +0000 UTC 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/obj1
        2015-06-05 14:06:29 +0000 UTC  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
        Object Number is: 2

    5)ossutil ls oss://bucket1 -d
        oss://bucket1/obj1
        oss://bucket1/dir1
        Object or Directory Number is: 2

    6)ossutil ls oss://bucket1 -m 
             ObjectName                                UploadID
        oss://bucket1/obj1                15754AF7980C4DFB8193F190837520BB
        oss://bucket2/obj2                3998971ACAF94AD9AC48EAC1988BE863 
        Object Number is: 2
    
    7)ossutil ls oss://bucket1 -a 
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2016-04-08 14:50:47 +0000 UTC 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/sample.txt
        2015-06-05 14:06:29 +0000 UTC  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
         ObjectName                                UploadID
        oss://bucket1/obj1                15754AF7980C4DFB8193F190837520BB
        oss://bucket2/obj2                3998971ACAF94AD9AC48EAC1988BE863
        Object Number is: 4

`,
}

// ListCommand is the command list buckets or objects
type ListCommand struct {
	command Command
}

var listCommand = ListCommand{
	command: Command{
		name:        "ls",
		nameAlias:   []string{"list"},
		minArgc:     0,
		maxArgc:     1,
		specChinese: specChineseList,
		specEnglish: specEnglishList,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionShortFormat,
			OptionDirectory,
			OptionMultipart,
			OptionAllType,
			OptionConfigFile,
            OptionEndpoint,
            OptionAccessKeyID,
            OptionAccessKeySecret,
            OptionSTSToken,
			OptionRetryTimes,
		},
	},
}

// function for FormatHelper interface
func (lc *ListCommand) formatHelpForWhole() string {
	return lc.command.formatHelpForWhole()
}

func (lc *ListCommand) formatIndependHelp() string {
	return lc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (lc *ListCommand) Init(args []string, options OptionMapType) error {
	return lc.command.Init(args, options, lc)
}

// RunCommand simulate inheritance, and polymorphism
func (lc *ListCommand) RunCommand() error {
	if len(lc.command.args) == 0 {
		return lc.listBuckets("")
	}

	cloudURL, err := CloudURLFromString(lc.command.args[0])
	if err != nil {
		return err
	}

	if cloudURL.bucket == "" {
		return lc.listBuckets("")
	}

	return lc.listFiles(cloudURL)
}

func (lc *ListCommand) listBuckets(prefix string) error {
	if err := lc.lbCheckArgOptions(); err != nil {
		return err
	}

	shortFormat, _ := GetBool(OptionShortFormat, lc.command.options)
	num := 0

	client, err := lc.command.ossClient("")
	if err != nil {
		return err
	}

	// list all buckets
	pre := oss.Prefix(prefix)
	marker := oss.Marker("")
	for {
		lbr, err := lc.ossListBucketsRetry(client, pre, marker)
		if err != nil {
			return err
		}
		pre = oss.Prefix(lbr.Prefix)
		marker = oss.Marker(lbr.NextMarker)
        if num == 0 && !shortFormat && len(lbr.Buckets) > 0 {
            fmt.Printf("%-30s %20s%s%s\n", "CreationTime", "Region", FormatTAB, "BucketName")
        }
		for _, bucket := range lbr.Buckets {
			if !shortFormat {
				fmt.Printf("%-30s %20s%s%s\n", utcToLocalTime(bucket.CreationDate), bucket.Location, FormatTAB, CloudURLToString(bucket.Name, ""))
			} else {
				fmt.Println(CloudURLToString(bucket.Name, ""))
			}
		}
		num += len(lbr.Buckets)
		if !lbr.IsTruncated {
			break
		}
	}
	fmt.Printf("Bucket Number is: %d\n", num)
	return nil
}

func (lc *ListCommand) lbCheckArgOptions() error {
	if ok, _ := GetBool(OptionDirectory, lc.command.options); ok {
		return fmt.Errorf("ListBucket does not support option: \"%s\"", OptionDirectory)
	}
	return nil
}

func (lc *ListCommand) ossListBucketsRetry(client *oss.Client, options ...oss.Option) (oss.ListBucketsResult, error) {
	retryTimes, _ := GetInt(OptionRetryTimes, lc.command.options)
	for i := 1; ; i++ {
		lbr, err := client.ListBuckets(options...)
		if err == nil || int64(i) >= retryTimes {
			return lbr, err
		}
	}
}

func (lc *ListCommand) listFiles(cloudURL CloudURL) error {
	bucket, err := lc.command.ossBucket(cloudURL.bucket)
	if err != nil {
		return err
	}

	shortFormat, _ := GetBool(OptionShortFormat, lc.command.options)
	directory, _ := GetBool(OptionDirectory, lc.command.options)
    isMultipart, _ := GetBool(OptionMultipart, lc.command.options)
    isAllType, _ := GetBool(OptionAllType, lc.command.options) 
    isObject := true

    if isMultipart {
        isObject = false
    }
    if isAllType {
        isObject = true
        isMultipart = true
    }

	return lc.listObjects(bucket, cloudURL, shortFormat, directory, isObject, isMultipart)
}

func (lc *ListCommand) listObjects(bucket *oss.Bucket, cloudURL CloudURL, shortFormat bool, directory bool, isObject, isMultipart bool) error {
	objectNum := 0
	multipartNum := 0
	pre := oss.Prefix(cloudURL.object)
	marker := oss.Marker("")
	del := oss.Delimiter("")
	if directory {
		del = oss.Delimiter("/")
	}

    if isObject {
	    for i := 0; ; i++ {
	    	lor, err := lc.command.ossListObjectsRetry(bucket, marker, pre, del)
    		if err != nil {
			    return err
		    }
		    pre = oss.Prefix(lor.Prefix)
		    marker = oss.Marker(lor.NextMarker)
	    	objectNum += lc.displayResult(lor, cloudURL.bucket, shortFormat, directory, i)
    		if !lor.IsTruncated {
			    break
		    }
	    }
    
        if !directory {
		    fmt.Printf("Object Number is: %d\n", objectNum)
	    } else {
		    fmt.Printf("Object or Directory Number is: %d\n", objectNum)
	    }
    }

    if isMultipart {
    	for i := 0; ; i++ {
	    	lmr, err := lc.command.ossListMultipartUploadsRetry(bucket, marker, pre, del)
		    if err != nil {
			    return err
		    }
		    pre = oss.Prefix(lmr.Prefix)
		    marker = oss.Marker(lmr.NextKeyMarker)
		    multipartNum += lc.displayMultipartUploadsResult(lmr, cloudURL.bucket, shortFormat, directory, i)
		    if !lmr.IsTruncated {
			    break
		    }
	    }
        fmt.Printf("Multipart Number is: %d\n", multipartNum)
    }
	return nil
}

func (lc *ListCommand) displayResult(lor oss.ListObjectsResult, bucket string, shortFormat bool, directory bool, i int) int {
	if i == 0 && !shortFormat && !directory && len(lor.Objects) > 0 {
		fmt.Printf("%-30s %12s%s%-38s%s%s\n", "LastModifiedTime", "Size(B)", "   ", "ETAG", "  ", "ObjectName")
	}

	var output string
	var num int
	if !directory {
		output, num = lc.showObjects(lor, bucket, shortFormat)
	} else {
		output, num = lc.showObjects(lor, bucket, true)
		output1, num1 := lc.showDirectories(lor, bucket)
        output += output1
        num += num1 
	}
	fmt.Printf(output)
	return num
}

func (lc *ListCommand) showObjects(lor oss.ListObjectsResult, bucket string, shortFormat bool) (string, int) {
	var output string
	for _, object := range lor.Objects {
		if !shortFormat {
			output += fmt.Sprintf(
				"%-30s %12d%s%-38s%s%s\n", utcToLocalTime(object.LastModified), object.Size, "   ",
				strings.Trim(object.ETag, "\""), "  ", CloudURLToString(bucket, object.Key),
			)
		} else {
			output += CloudURLToString(bucket, object.Key) + "\n"
		}
	}
	return output, len(lor.Objects)
}

func (lc *ListCommand) showDirectories(lor oss.ListObjectsResult, bucket string) (string, int) {
	var output string
	for _, prefix := range lor.CommonPrefixes {
		output += CloudURLToString(bucket, prefix) + "\n"
	}
	return output, len(lor.CommonPrefixes)
}


// multipart objects display
func (lc *ListCommand) displayMultipartUploadsResult(lmr oss.ListMultipartUploadResult, bucket string, shortFormat bool, directory bool, i int) int {
	if i == 0 && !shortFormat && len(lmr.Uploads) > 0 {
		fmt.Printf("%s%-28s%s\n", "UploadId", " ", "MultipartName")
	}

	var output string
    output = lc.showMultipartUploads(lmr, bucket, shortFormat)	
    fmt.Printf(output)
	return len(lmr.Uploads)
}

func (lc *ListCommand) showMultipartUploads(lmr oss.ListMultipartUploadResult, bucket string, shortFormat bool)(string) {
	var output string
	for _, upload := range lmr.Uploads {
		if !shortFormat {
			output += fmt.Sprintf(
				"%s%-4s%s\n", upload.UploadID, " ", CloudURLToString(bucket, upload.Key),
			)
		} else {
			output += CloudURLToString(bucket, upload.Key) + "\n"
		}
	}
	return output
}

