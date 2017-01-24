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

用法：

    该命令有两种用法：

    1) ossutil ls [oss://] [-s]
        如果用户列举时缺失url参数，则ossutil获取用户的身份凭证信息（从配置文件中读取），
    并列举该身份凭证下的所有buckets，并显示每个bucket的最新更新时间和位置信息。如果指定
    了--short-format选项则只输出bucket名称。该用法不支持--directory选项。

    2) ossutil ls oss://bucket[/prefix] [-s] [-d] [-m] [-a]
        如果未指定--multipart和--all-type选项，则ossutil列举指定bucket下的objects（如果指定
    了前缀，则列举拥有该前缀的objects）。并同时展示object大小，最新更新时间和etag，但是如果
    指定了--short-format选项则只输出object名称。如果指定了--directory选项，则返回指定bucket
    下以指定前缀开头的第一层目录下的文件和子目录，但是不递归显示所有子目录，此时默认为精简
    格式。所有的目录均以/结尾。
        如果指定了--multipart选项，则显示指定URL(oss://bucket[/prefix])下未完成的上传任务，
    即，列举未complete的Multipart Upload事件的uploadId，这些Multipart Upload事件的object名
    称以指定的prefix为前缀。ossutil同时显示uploadId的init时间。该选项同样支持--short-format
    和--directory选项。（Multipart同样用于cp命令中大文件的断点续传，关于Multipart的更多信息
    见：https://help.aliyun.com/document_detail/31991.html?spm=5176.doc31992.6.880.VOSDk5）。
        如果指定了--all-type选项，则显示指定URL(oss://bucket[/prefix])下的object和未完成的
	上传任务（即，同时列举以prefix为前缀的object，和object名称以prefix为前缀的所有未complete
    的uploadId）。该选项同样支持--short-format和--directory选项。
`,

	sampleText: ` 
    1) ossutil ls -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    2) ossutil ls oss:// -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    3) ossutil ls oss://bucket1 -s
        oss://bucket1/dir1/obj11
        oss://bucket1/obj1
        oss://bucket1/sample.txt
        Object Number is: 3

    4) ossutil ls oss://bucket1
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2015-06-05 14:06:29 +0000 CST  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
        2015-06-05 14:36:21 +0000 CST  201933  6185CA2E8EB8510A61B3A845EAFE4174  oss://bucket1/obj1
        2016-04-08 14:50:47 +0000 CST 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/sample.txt
        Object Number is: 3

    5) ossutil ls oss://bucket1 -d
        oss://bucket1/obj1
        oss://bucket1/dir1
        oss://bucket1/sample.txt
        Object and Directory Number is: 3

    6) ossutil ls oss://bucket1 -m 
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        2017-01-20 11:16:21 +0800 CST  A20157A7B2FEC4670626DAE0F4C0073C  oss://bucket1/tobj
        UploadId Number is: 3
    
    7) ossutil ls oss://bucket1/obj -m 
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        UploadId Number is: 2
 
    8) ossutil ls oss://bucket1 -a 
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2015-06-05 14:06:29 +0000 CST  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
        2015-06-05 14:36:21 +0000 CST  201933  6185CA2E8EB8510A61B3A845EAFE4174  oss://bucket1/obj1
        2016-04-08 14:50:47 +0000 CST 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/sample.txt
        Object Number is: 3
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:43:13 +0000 CST  2A1F9B4A95E341BD9285CC42BB950EE0  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        2017-01-20 11:16:21 +0800 CST  A20157A7B2FEC4670626DAE0F4C0073C  oss://bucket1/tobj
        UploadId Number is: 4
         
    9) ossutil ls oss://bucket1/obj -a 
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2015-06-05 14:36:21 +0000 CST  201933  6185CA2E8EB8510A61B3A845EAFE4174  oss://bucket1/obj1
        Object Number is: 2
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:43:13 +0000 CST  2A1F9B4A95E341BD9285CC42BB950EE0  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        UploadId Number is: 3

    10) ossutil ls oss://bucket1/obj -a -s 
        oss://bucket1/obj1
        Object Number is: 2
        UploadID                          MultipartName
        15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2A1F9B4A95E341BD9285CC42BB950EE0  oss://bucket1/obj1
        3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        UploadId Number is: 3
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

    1) ossutil ls [oss://] [-s] [-m] [-a]
        If you list without a url, ossutil lists all the buckets using the credentials
    in config file with last modified time and location in addition. --show_format option 
    will ignore last modified time and location. The usage do not support --directory 
    option.

    2) ossutil ls oss://bucket[/prefix] [-s] [-d] [-m] [-a]
        If you list without --multipart and --all-type option, ossutil will list objects 
    in the specified bucket(with the prefix if you specified), with object size, last 
    modified time and etag in addition, --short-format option ignores all the additional 
    information. --directory option returns top-level subdirectory names instead of contents 
    of the subdirectory, which in default show by short format. the directory is end with /. 
        --multipart option will show multipart upload tasks under the url(oss://bucket[/prefix]), 
    which means, ossutil will show the uploadId of those uncompleted multipart, whose object 
    name starts with the specified prefix. ossutil will show the init time of uploadId meanwhile. 
    The usage also supports --short-format and --directory option. (Multipart upload is also used 
    in resume cp. More information about multipart see: https://help.aliyun.com/document_detail/31991.html?spm=5176.doc31992.6.880.VOSDk5). 
        --all-type option will show objects and multipart upload tasks under the url(oss://bucket[/prefix]),  
    which means, ossutil will both show the objects with the specified prefix and the uploadId of 
    those uncompleted multipart, whose object name starts with the specified prefix. The usage also 
    support --short-format and --directory option.
`,

	sampleText: ` 
    1) ossutil ls -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    2) ossutil ls oss:// -s
        oss://bucket1
        oss://bucket2
        oss://bucket3
        Bucket Number is: 3

    3) ossutil ls oss://bucket1 -s
        oss://bucket1/dir1/obj11
        oss://bucket1/obj1
        oss://bucket1/sample.txt
        Object Number is: 3

    4) ossutil ls oss://bucket1
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2015-06-05 14:06:29 +0000 CST  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
        2015-06-05 14:36:21 +0000 CST  201933  6185CA2E8EB8510A61B3A845EAFE4174  oss://bucket1/obj1
        2016-04-08 14:50:47 +0000 CST 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/sample.txt
        Object Number is: 3

    5) ossutil ls oss://bucket1 -d
        oss://bucket1/obj1
        oss://bucket1/dir1
        oss://bucket1/sample.txt
        Object and Directory Number is: 3

    6) ossutil ls oss://bucket1 -m 
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        2017-01-20 11:16:21 +0800 CST  A20157A7B2FEC4670626DAE0F4C0073C  oss://bucket1/tobj
        UploadId Number is: 3
    
    7) ossutil ls oss://bucket1/obj -m 
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        UploadId Number is: 2
 
    8) ossutil ls oss://bucket1 -a 
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2015-06-05 14:06:29 +0000 CST  201933  7E2F4A7F1AC9D2F0996E8332D5EA5B41  oss://bucket1/dir1/obj11
        2015-06-05 14:36:21 +0000 CST  201933  6185CA2E8EB8510A61B3A845EAFE4174  oss://bucket1/obj1
        2016-04-08 14:50:47 +0000 CST 6476984  4F16FDAE7AC404CEC8B727FCC67779D6  oss://bucket1/sample.txt
        Object Number is: 3
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:43:13 +0000 CST  2A1F9B4A95E341BD9285CC42BB950EE0  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        2017-01-20 11:16:21 +0800 CST  A20157A7B2FEC4670626DAE0F4C0073C  oss://bucket1/tobj
        UploadId Number is: 4
         
    9) ossutil ls oss://bucket1/obj -a 
        LastModifiedTime              Size(B)  ETAG                              ObjectName
        2015-06-05 14:36:21 +0000 CST  201933  6185CA2E8EB8510A61B3A845EAFE4174  oss://bucket1/obj1
        Object Number is: 2
        InitiatedTime                  UploadID                          MultipartName
        2017-01-13 03:45:26 +0000 CST  15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2017-01-13 03:43:13 +0000 CST  2A1F9B4A95E341BD9285CC42BB950EE0  oss://bucket1/obj1
        2017-01-13 03:45:25 +0000 CST  3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        UploadId Number is: 3

    10) ossutil ls oss://bucket1/obj -a -s 
        oss://bucket1/obj1
        Object Number is: 2
        UploadID                          MultipartName
        15754AF7980C4DFB8193F190837520BB  oss://bucket1/obj1
        2A1F9B4A95E341BD9285CC42BB950EE0  oss://bucket1/obj1
        3998971ACAF94AD9AC48EAC1988BE863  oss://bucket1/obj2
        UploadId Number is: 3
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

    typeSet := lc.getSubjectType()
    if typeSet & objectType != 0 {
	    if err := lc.listObjects(bucket, cloudURL, shortFormat, directory); err != nil {
            return err
        }
    }
    if typeSet & multipartType != 0 {
	    if err := lc.listMultipartUploads(bucket, cloudURL, shortFormat, directory); err != nil {
            return err
        }
    }
    return nil
}

func (lc *ListCommand) getSubjectType() int64 {
    var typeSet int64
    typeSet = 0 
    if isMultipart, _ := GetBool(OptionMultipart, lc.command.options); isMultipart {
        typeSet |= multipartType 
    }
	if isAllType, _ := GetBool(OptionAllType, lc.command.options); isAllType {
        typeSet |= allType
    }
    if typeSet & allType == 0 {
        typeSet = objectType 
    }
    return typeSet
}

func (lc *ListCommand) listObjects(bucket *oss.Bucket, cloudURL CloudURL, shortFormat bool, directory bool) error {
    //list all objects or directories
    var num int64
    num = 0
    pre := oss.Prefix(cloudURL.object)
    marker := oss.Marker("")
    del := oss.Delimiter("")
    if directory {
        del = oss.Delimiter("/")
    }

    var i int64
    for i = 0; ; i++ {
        lor, err := lc.command.ossListObjectsRetry(bucket, marker, pre, del)
        if err != nil {
            return err
        }
        pre = oss.Prefix(lor.Prefix)
        marker = oss.Marker(lor.NextMarker)
        num += lc.displayObjectsResult(lor, cloudURL.bucket, shortFormat, directory, i)
        if !lor.IsTruncated {
            break
        }
    }

    if !directory {
        fmt.Printf("Object Number is: %d\n", num)
    } else {
        fmt.Printf("Object and Directory Number is: %d\n", num)
    }

    return nil
}

func (lc *ListCommand) displayObjectsResult(lor oss.ListObjectsResult, bucket string, shortFormat bool, directory bool, i int64) int64 {
    if i == 0 && !shortFormat && !directory && len(lor.Objects) > 0 {
        fmt.Printf("%-30s%12s%s%-38s%s%s\n", "LastModifiedTime", "Size(B)", "   ", "ETAG", "  ", "ObjectName")
    }

    var num int64
    if !directory {
        num = lc.showObjects(lor, bucket, shortFormat)
    } else {
        num = lc.showObjects(lor, bucket, true)
        num1 := lc.showDirectories(lor, bucket)
        num += num1
    }
    return num
}

func (lc *ListCommand) showObjects(lor oss.ListObjectsResult, bucket string, shortFormat bool) int64 {
	for _, object := range lor.Objects {
		if !shortFormat {
			fmt.Printf("%-30s%12d%s%-38s%s%s\n", utcToLocalTime(object.LastModified), object.Size, "   ", strings.Trim(object.ETag, "\""), "  ", CloudURLToString(bucket, object.Key))
		} else {
            fmt.Printf("%s\n", CloudURLToString(bucket, object.Key))
		}
	}
	return int64(len(lor.Objects))
}

func (lc *ListCommand) showDirectories(lor oss.ListObjectsResult, bucket string) int64 {
	for _, prefix := range lor.CommonPrefixes {
        fmt.Printf("%s\n", CloudURLToString(bucket, prefix))
	}
	return int64(len(lor.CommonPrefixes))
}

func (lc *ListCommand) listMultipartUploads(bucket *oss.Bucket, cloudURL CloudURL, shortFormat bool, directory bool) error {
    var multipartNum int64
	multipartNum = 0
	pre := oss.Prefix(cloudURL.object)
	marker := oss.Marker("")
	del := oss.Delimiter("")
	if directory {
		del = oss.Delimiter("/")
	}

    var i int64
    for i = 0; ; i++ {
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
	return nil
}

func (lc *ListCommand) displayMultipartUploadsResult(lmr oss.ListMultipartUploadResult, bucket string, shortFormat bool, directory bool, i int64) int64 {
    if directory {
        shortFormat = true
    }

	if i == 0 && len(lmr.Uploads) > 0 {
        if shortFormat {
		    fmt.Printf("%-32s%s%s\n", "UploadID", FormatTAB, "MultipartName")
        } else {
		    fmt.Printf("%-30s%s%-32s%s%s\n", "InitiatedTime", FormatTAB, "UploadID", FormatTAB, "MultipartName")
        }
	}

    num := lc.showMultipartUploads(lmr, bucket, shortFormat)
    return num
}

func (lc *ListCommand) showMultipartUploads(lmr oss.ListMultipartUploadResult, bucket string, shortFormat bool) int64 {
	for _, upload := range lmr.Uploads {
		if shortFormat {
            fmt.Printf("%-32s%s%s\n", upload.UploadID, FormatTAB, CloudURLToString(bucket, upload.Key))
		} else {
			fmt.Printf("%-30s%s%-32s%s%s\n", utcToLocalTime(upload.Initiated), FormatTAB, upload.UploadID, FormatTAB, CloudURLToString(bucket, upload.Key))
		}
	}
	return int64(len(lmr.Uploads)) 
}
