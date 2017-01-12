package lib

import (
	"fmt"
	"sort"
	"strings"
    "time"
    "net/http"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseStat = SpecText{

	synopsisText: "显示bucket或者object的描述信息",

	paramText: "url [options]",

	syntaxText: ` 
    ossutil stat oss://bucket[/object] [-c file] 
`,

	detailHelpText: ` 
    该命令获取指定bucket或者objects的描述信息，当指定--recursive选项时，获取所有与指定
    url匹配的objects，统计objects个数大小和最新更新时间并输出。否则，显示指定的bucket的
    信息或者单个object的元信息。

用法：

    该命令有两种用法：

    1) ossutil stat oss://bucket
        如果未指定--recursive选项，ossutil显示指定bucket的信息，包括创建时间，location，
    访问的外网域名，内网域名，拥有者，acl信息。

    2) ossutil stat oss://bucket/object
        如果未指定--recursive选项，ossutil显示指定object的元信息，包括文件大小，最新更新
    时间，etag，文件类型，acl，文件的自定义meta等信息。
`,

	sampleText: ` 
    ossutil stat oss://bucket1
    ossutil stat oss://bucket1/object  
`,
}

var specEnglishStat = SpecText{

	synopsisText: "Display status of bucket or objects",

	paramText: "url [options]",

	syntaxText: ` 
    ossutil stat oss://bucket[/prefix] [-c file] 
`,

	detailHelpText: ` 
    The command display status of bucket or objects. If --recursive option is specified, 
    ossutil access all the objects prefix matching the specified url, display total count 
    and size, last modify time. If not, ossutil will display the meta information of the 
    specified bucket or single object.

Usage：

    There are three usages:    

    1) ossutil stat oss://bucket
        If you use the command to bucket without --recursive option, ossutil will display 
    bucket meta info, include creation date, location, extranet endpoint, intranet endpoint, 
    Owner and acl info.

    2) ossutil stat oss://bucket/object
        If you use the command to object without --recursive option, ossutil will display
    object meta info, include file size, last modify time, etag, content-type, user meta etc.
`,

	sampleText: ` 
    ossutil stat oss://bucket1
    ossutil stat oss://bucket1/object  
`,
}

// StatCommand is the command get bucket's or objects' meta information 
type StatCommand struct {
	command Command
}

var statCommand = StatCommand{
	command: Command{
		name:        "stat",
		nameAlias:   []string{"meta", "info"},
		minArgc:     1,
		maxArgc:     1,
		specChinese: specChineseStat,
		specEnglish: specEnglishStat,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
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
func (sc *StatCommand) formatHelpForWhole() string {
	return sc.command.formatHelpForWhole()
}

func (sc *StatCommand) formatIndependHelp() string {
	return sc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (sc *StatCommand) Init(args []string, options OptionMapType) error {
	return sc.command.Init(args, options, sc)
}

// RunCommand simulate inheritance, and polymorphism
func (sc *StatCommand) RunCommand() error {
	cloudURL, err := CloudURLFromString(sc.command.args[0])
	if err != nil {
		return err
	}

	if cloudURL.bucket == "" {
		return fmt.Errorf("invalid cloud url: %s, miss bucket", sc.command.args[0])
	}

	bucket, err := sc.command.ossBucket(cloudURL.bucket)
	if err != nil {
		return err
	}

	if cloudURL.object == "" {
		return sc.bucketStat(bucket, cloudURL)
	}
	return sc.objectStat(bucket, cloudURL)
}

func (sc *StatCommand) bucketStat(bucket *oss.Bucket, cloudURL CloudURL) error {
	// TODO: go sdk should implement GetBucketInfo
    gbar, err := sc.ossGetBucketStatRetry(bucket)
    if err != nil {
        return err
    }

    attrMap := map[string]string{}
    attrMap["Location"] = gbar.BucketInfo.Location
	fmt.Printf("%-18s: %s\n", StatName, gbar.BucketInfo.Name)
	fmt.Printf("%-18s: %s\n", StatLocation, gbar.BucketInfo.Location)
	fmt.Printf("%-18s: %s\n", StatCreationDate, utcToLocalTime(gbar.BucketInfo.CreationDate))
	fmt.Printf("%-18s: %s\n", StatExtranetEndpoint, gbar.BucketInfo.ExtranetEndpoint)
	fmt.Printf("%-18s: %s\n", StatIntranetEndpoint, gbar.BucketInfo.IntranetEndpoint)
	fmt.Printf("%-18s: %s\n", StatACL, gbar.BucketInfo.ACL)
	fmt.Printf("%-18s: %s\n", StatOwner, gbar.BucketInfo.Owner.ID)
	return nil
}

func (sc *StatCommand) ossGetBucketStatRetry(bucket *oss.Bucket) (oss.GetBucketInfoResult, error) {
	retryTimes, _ := GetInt(OptionRetryTimes, sc.command.options)
	for i := 1; ; i++ {
        gbar, err := bucket.Client.GetBucketInfo(bucket.BucketName)
		if err == nil {
			return gbar, err
		}
		if int64(i) >= retryTimes {
			return gbar, BucketError{err, bucket.BucketName}
		}
	}
}

func (sc *StatCommand) objectStat(bucket *oss.Bucket, cloudURL CloudURL) error {
	// acl info
	goar, err := sc.ossGetObjectACLRetry(bucket, cloudURL.object)
	if err != nil {
		return err
	}

	// normal info
	props, err := sc.command.ossGetObjectStatRetry(bucket, cloudURL.object)
	if err != nil {
		return err
	}

	sortNames := []string{}
	attrMap := map[string]string{}
	for name := range props {
		ln := strings.ToLower(name)
		if ln != strings.ToLower(oss.HTTPHeaderDate) &&
			ln != strings.ToLower(oss.HTTPHeaderOssRequestID) &&
			ln != strings.ToLower(oss.HTTPHeaderServer) &&
			ln != "x-oss-server-time" &&
			ln != "connection" {
			sortNames = append(sortNames, name)
			attrMap[name] = props.Get(name)
		}
	}

	sortNames = append(sortNames, "Owner")
	sortNames = append(sortNames, "ACL")
	attrMap[StatOwner] = goar.Owner.ID
	attrMap[StatACL] = goar.ACL
    if lm, err := time.Parse(http.TimeFormat, attrMap[StatLastModified]); err == nil {
        attrMap[StatLastModified] = fmt.Sprintf("%s", utcToLocalTime(lm.UTC())) 
    }

	sort.Strings(sortNames)

	for _, name := range sortNames {
		if strings.ToLower(name) != "etag" {
			fmt.Printf("%-28s: %s\n", name, attrMap[name])
		} else {
			fmt.Printf("%-28s: %s\n", name, strings.Trim(attrMap[name], "\""))
		}
	}
	return nil
}

func (sc *StatCommand) ossGetObjectACLRetry(bucket *oss.Bucket, object string) (oss.GetObjectACLResult, error) {
	retryTimes, _ := GetInt(OptionRetryTimes, sc.command.options)
	for i := 1; ; i++ {
		goar, err := bucket.GetObjectACL(object)
		if err == nil {
			return goar, err
		}
		if int64(i) >= retryTimes {
			return goar, ObjectError{err, bucket.BucketName, object}
		}
	}
}
