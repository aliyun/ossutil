package lib

import (
	"fmt"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
)

var aclMap = map[oss.ACLType][]string{
	oss.ACLPublicReadWrite: []string{},
	oss.ACLPublicRead:      []string{},
	oss.ACLPrivate:         []string{},
	oss.ACLDefault:         []string{},
}

var bucketACLList = []oss.ACLType{
	oss.ACLPublicReadWrite,
	oss.ACLPublicRead,
	oss.ACLPrivate,
}

var objectACLList = []oss.ACLType{
	oss.ACLPublicReadWrite,
	oss.ACLPublicRead,
	oss.ACLPrivate,
	oss.ACLDefault,
}

type setACLType int

const (
	bucketACL setACLType = iota
	objectACL
)

func formatACLString(aclType setACLType, sep string) string {
	var list []oss.ACLType
	if aclType == bucketACL {
		list = bucketACLList
	} else {
		list = objectACLList
	}

	strList := []string{}
	for _, acl := range list {
		str := string(acl)
		if len(aclMap[acl]) != 0 {
			str += "(" + strings.Join(aclMap[acl][:], ",") + ")"
		}
		strList = append(strList, str)
	}
	return strings.Join(strList, sep)
}

var specChineseSetACL = SpecText{

	synopsisText: "设置bucket或者objects的acl",

	paramText: "url [acl] [options]",

	syntaxText: ` 
    ossutil set-acl oss://bucket[/prefix] [acl] [-r] [-b] [-f] [-c file] 
`,

	detailHelpText: ` 
    该命令设置指定bucket或者objects的acl。使用命令时若缺失了acl信息时，ossutil会询问用户acl信息。

        （1）设置bucket的acl，参考用法1)
        （2）设置单个object的acl，参考用法2)
        （3）批量设置objects的acl，不设置bucket的acl，参考用法3)

    对bucket设置acl，需要添加--bucket选项，否则视为对其中的objects设置acl。
    该命令不支持同时设置bucket和objects的acl，请分开操作。

    结果：显示命令耗时前未报错，则表示成功设置。
    查看bucket或者object的acl信息，请使用stat命令。

ACL：

    bucket的acl有三种，括号里为ossutil额外支持的简写模式：
        ` + formatACLString(bucketACL, "\n        ") + `

    object的acl有四种：
        ` + formatACLString(objectACL, "\n        ") + `

    acl的详细信息请参见：https://help.aliyun.com/document_detail/31867.html?spm=5176.doc31960.6.147.8dVwsh中的权限控制。

用法：

    该命令有三种用法：

    1) ossutil set-acl oss://bucket [acl] -b [-c file]
        当设置了--bucket选项时，ossutil会尝试设置bucket的acl，此时不支持--recursive选项，并且请
    确保输入的url精确匹配想要设置acl的bucket，无论--force选项是否指定，都不会进行询问提示。如果
    用户在命令行中缺失acl信息，会进入交互模式，询问用户的acl信息。 

    2) ossutil set-acl oss://bucket/object [acl] [-c file]
        该用法设置指定单个object的acl，当指定object不存在时，ossutil会提示错误，此时请确保指定的
    url精确匹配需要设置acl的object，并且不要指定--recursive选项（否则ossutil会进行前缀匹配，设置
    多个objects的acl），无论--force选项是否指定，都不会进行询问提示。如果用户在命令行中缺失acl信
    息，会进入交互模式，询问用户的acl信息。

    3) ossutil set-acl oss://bucket[/prefix] [acl] -r [-f] [-c file]
        该用法可批量设置objects的acl，此时必须输入--recursive选项，ossutil会查找所有前缀匹配url的
    objects，设置它们的acl，当错误出现时会终止命令。此时不支持--bucket选项，即ossutil不支持同时设
    置bucket和其中objects的acl，如有需要，请分开操作。如果--force选项被指定，则不会进行询问提示。
    如果用户在命令行中缺失acl信息，会进入交互模式，询问用户的acl信息。
`,

	sampleText: ` 
    (1)ossutil set-acl oss://bucket1 public-read-write -b 

    (2)ossutil set-acl oss://bucket1/obj1 private 

    (3)ossutil set-acl oss://bucket1/obj default -r

    (4)ossutil set-acl oss://bucket1/%e4%b8%ad%e6%96%87 default --encoding-type url
`,
}

var specEnglishSetACL = SpecText{

	synopsisText: "Set acl on bucket or objects",

	paramText: "url [acl] [options]",

	syntaxText: ` 
    ossutil set-acl oss://bucket[/prefix] [acl] [-r] [-b] [-f] [-c file] 
`,

	detailHelpText: ` 
    The command set acl on the specified bucket or objects. If you use the command 
    witout acl information, ossutil will ask user for it.

    (1) set acl on bucket, see usage 1)
    (2) set acl on single object, see usage 2)
    (3) batch set acl on many objects, see usage 3)

    When set acl on bucket, the --bucket option must be specified. 
    Set acl on bucket an objects inside simultaneously is not supported, please 
    operate independently.

    Result: if no error displayed before show elasped time, then the setting is completed successfully.
    User can use stat command to check the acl information of bucket or objects.

ACL:

    ossutil supports following bucket acls, shorthand versions in brackets:
        ` + formatACLString(bucketACL, "\n        ") + `

    ossutil support following objet acls:
        ` + formatACLString(objectACL, "\n        ") + `

    More information about acl see ACL Control in https://help.aliyun.com/document_detail/31867.html?spm=5176.doc31960.6.147.8dVwsh.

Usage：

    There are three usages:    

    1) ossutil set-acl oss://bucket [acl] -b [-c file]
        If --bucket option is specified, ossutil will try to set acl on bucket. In the 
    usage, please make sure url exactly specified the bucket you want to set acl on, 
    and --recursive option is not supported here. No matter --force option is specified 
    or not, ossutil will not show prompt question. If acl information is missed, ossutil 
    will enter interactive mode and ask you for it. 

    2) ossutil set-acl oss://bucket/object [acl] [-c file]
        The usage set acl on single object, if object not exist, error occurs. In the 
    usage, please make sure url exactly specified the object you want to set acl on, 
    and --recursive option is not specified(or ossutil will search for prefix-matching 
    objects and set acl on those objects). No matter --force option is specified or not, 
    ossutil will not show prompt question. If acl information is missed, ossutil will 
    enter interactive mode and ask you for it. 

    3) ossutil set-acl oss://bucket[/prefix] [acl] -r [-f] [-c file]
        The usage can set acl on many objects, --recursive option is required for the 
    usage, ossutil will search for prefix-matching objects and set acl on those objects, 
    if error occurs, the operation is terminated. In the usage, --bucket option is not 
    supported, which means set acl on bucket an objects inside simultaneously is not 
    supported. If --force option is specified, ossutil will not show prompt question. If 
    acl information is missed, ossutil will enter interactive mode and ask you for it. 
`,

	sampleText: ` 
    (1)ossutil set-acl oss://bucket1 public-read-write -b 

    (2)ossutil set-acl oss://bucket1/obj1 private 

    (3)ossutil set-acl oss://bucket1/obj default -r

    (4)ossutil set-acl oss://bucket1/%e4%b8%ad%e6%96%87 default --encoding-type url
`,
}

// SetACLCommand is the command set acl
type SetACLCommand struct {
	command Command
	monitor Monitor
}

var setACLCommand = SetACLCommand{
	command: Command{
		name:        "set-acl",
		nameAlias:   []string{"setacl", "set_acl"},
		minArgc:     1,
		maxArgc:     2,
		specChinese: specChineseSetACL,
		specEnglish: specEnglishSetACL,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionRecursion,
			OptionBucket,
			OptionForce,
            OptionEncodingType,
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionRetryTimes,
			OptionRoutines,
		},
	},
}

// function for FormatHelper interface
func (sc *SetACLCommand) formatHelpForWhole() string {
	return sc.command.formatHelpForWhole()
}

func (sc *SetACLCommand) formatIndependHelp() string {
	return sc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (sc *SetACLCommand) Init(args []string, options OptionMapType) error {
	return sc.command.Init(args, options, sc)
}

// RunCommand simulate inheritance, and polymorphism
func (sc *SetACLCommand) RunCommand() error {
	sc.monitor.init("Setted acl on")

	recursive, _ := GetBool(OptionRecursion, sc.command.options)
	toBucket, _ := GetBool(OptionBucket, sc.command.options)
	force, _ := GetBool(OptionForce, sc.command.options)
	routines, _ := GetInt(OptionRoutines, sc.command.options)

    encodingType, _ := GetString(OptionEncodingType, sc.command.options)
	cloudURL, err := CloudURLFromString(sc.command.args[0], encodingType)
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

	if toBucket {
		return sc.setBucketACL(&bucket.Client, cloudURL, recursive)
	}
	if !recursive {
		return sc.setObjectACL(bucket, cloudURL)
	}
	return sc.batchSetObjectACL(bucket, cloudURL, force, routines)
}

func (sc *SetACLCommand) setBucketACL(client *oss.Client, cloudURL CloudURL, recursive bool) error {
	if cloudURL.object != "" {
		return fmt.Errorf("set bucket acl invalid url: %s, object not empty, if you mean set object acl, you should not use --bucket option", sc.command.args[0])
	}

	if recursive {
		return fmt.Errorf("set bucket acl do not support --recursive option, if you mean set object acl recursivlly, you should not use --bucket option")
	}

	acl, err := sc.getACL(bucketACL, recursive)
	if err != nil {
		return err
	}

	return sc.ossSetBucketACLRetry(client, cloudURL.bucket, acl)
}

func (sc *SetACLCommand) getACL(aclType setACLType, recursive bool) (oss.ACLType, error) {
	var acl string
	if len(sc.command.args) == 2 {
		acl = sc.command.args[1]
	} else {
		str := "bucket"
		if aclType == objectACL {
			str = "object"
			if recursive {
				str = "objects"
			}
		}
		fmt.Printf("Please enter the acl you want to set on the %s(%s):", str, formatACLString(aclType, ", "))
		if _, err := fmt.Scanln(&acl); err != nil {
			return "", fmt.Errorf("invalid acl: %s, please check", acl)
		}
	}

	return sc.checkACL(acl, aclType)
}

func (sc *SetACLCommand) checkACL(acl string, aclType setACLType) (oss.ACLType, error) {
	var list []oss.ACLType
	if aclType == bucketACL {
		list = bucketACLList
	} else {
		list = objectACLList
	}

	for _, acll := range list {
		if acl == string(acll) {
			return acll, nil
		}
		for _, aclll := range aclMap[acll] {
			if acl == aclll {
				return acll, nil
			}
		}
	}
	return "", fmt.Errorf("invalid acl: %s, please check", acl)
}

func (sc *SetACLCommand) ossSetBucketACLRetry(client *oss.Client, bucket string, acl oss.ACLType) error {
	retryTimes, _ := GetInt(OptionRetryTimes, sc.command.options)
	for i := 1; ; i++ {
		err := client.SetBucketACL(bucket, acl)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return BucketError{err, bucket}
		}
	}
}

func (sc *SetACLCommand) setObjectACL(bucket *oss.Bucket, cloudURL CloudURL) error {
	if cloudURL.object == "" {
		return fmt.Errorf("set object acl invalid url: %s, object empty, if you mean set bucket acl, you should use --bucket option", sc.command.args[0])
	}

	acl, err := sc.getACL(objectACL, false)
	if err != nil {
		return err
	}

	return sc.ossSetObjectACLRetry(bucket, cloudURL.object, acl)
}

func (sc *SetACLCommand) ossSetObjectACLRetry(bucket *oss.Bucket, object string, acl oss.ACLType) error {
	retryTimes, _ := GetInt(OptionRetryTimes, sc.command.options)
	for i := 1; ; i++ {
		err := bucket.SetObjectACL(object, acl)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, bucket.BucketName, object}
		}
	}
}

func (sc *SetACLCommand) batchSetObjectACL(bucket *oss.Bucket, cloudURL CloudURL, force bool, routines int64) error {
	if !force {
		var val string
		fmt.Printf("Do you really mean to recursivlly set acl on objects of %s(y or N)? ", sc.command.args[0])
		if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
			fmt.Println("operation is canceled.")
			return nil
		}
	}

	acl, err := sc.getACL(objectACL, true)
	if err != nil {
		return err
	}

	// producer list objects
	// consumer set acl
	chObjects := make(chan string, ChannelBuf)
	chError := make(chan error, routines+1)
	chListError := make(chan error, 1)
	go sc.command.objectStatistic(bucket, cloudURL, &sc.monitor)
	go sc.command.objectProducer(bucket, cloudURL, chObjects, chListError)
	for i := 0; int64(i) < routines; i++ {
		go sc.setObjectACLConsumer(bucket, acl, chObjects, chError)
	}

	completed := 0
	for int64(completed) <= routines {
		select {
		case err := <-chListError:
			if err != nil {
				return err
			}
			completed++
		case err := <-chError:
			if err != nil {
				fmt.Printf(sc.monitor.progressBar(true))
				return err
			}
			completed++
		}
	}
	fmt.Printf(sc.monitor.progressBar(true))
	return nil
}

func (sc *SetACLCommand) setObjectACLConsumer(bucket *oss.Bucket, acl oss.ACLType, chObjects <-chan string, chError chan<- error) {
	for object := range chObjects {
		err := sc.ossSetObjectACLRetry(bucket, object, acl)
		sc.command.updateMonitor(err, &sc.monitor)
		if err != nil {
			chError <- err
			return
		}
	}

	chError <- nil
}
