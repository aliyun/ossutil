package lib

import (
	"fmt"
	"strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var signURLHeaderOptionMap = map[string]interface{}{
	oss.HTTPHeaderContentType:             oss.ContentType,
	oss.HTTPHeaderOssServerSideEncryption: oss.ServerSideEncryption,
	oss.HTTPHeaderOssObjectACL:            oss.ObjectACL,
	oss.HTTPHeaderContentMD5:              oss.ContentMD5,
}

var specChineseSignurl = SpecText{

	synopsisText: "签名URL",

	paramText: "cloud_url [meta] [options]",

	syntaxText: ` 
    ossutil signurl cloud_url [header:value#header:value...] [--method m] [--timeout t] [--encoding-type url] [-c file] 
`,

	detailHelpText: ` 
    该命令签名用户指定的cloud_url，生成经过签名的url可供第三方用户访问object，其中cloud_url
    必须为形如：oss://bucket/object的cloud_url，bucket和object不可缺少。通过--method指定签名
    的method，默认为GET。通过--timeout选项指定url的过期时间，默认为60s。

    如果签名url需要指定content-type或自定义header等，可以通过header:value#header:value...参
    数指定。

Headers:

    可选的header列表如下：
        ` + formatHeaderString(signURLHeaderOptionMap, "\n        ") + `
        以及以` + oss.HTTPHeaderOssMetaPrefix + `开头的header

    注意：header不区分大小写，但value区分大小写。


用法：

    ossutil signurl oss://bucket/object [header:value#header:value...] [--method m] [--timeout t] [--encoding-type url]
`,

	sampleText: ` 
    ossutil signurl oss://bucket1/object1
        生成oss://bucket1/object1的GET签名url。

    ossutil signurl oss://bucket1/object1 --method PUT
        生成oss://bucket1/object1的PUT签名url。
`,
}

var specEnglishSignurl = SpecText{

	synopsisText: "Sign URL",

	paramText: "cloud_url [options]",

	syntaxText: ` 
    ossutil signurl cloud_url [--method m] [--timeout t] [--encoding-type url] [-c file]
`,

	detailHelpText: ` 
    The command create symlink of object in oss, the target object must be object in the 
    same bucket of symlink object, and the file type of target object must not be symlink. 
    So, cloud_url must be in format: oss://bucket/object, and target_object is the object 
    name of target object.  

Headers:

    ossutil supports following headers:
        ` + formatHeaderString(signURLHeaderOptionMap, "\n        ") + `
        and headers starts with: ` + oss.HTTPHeaderOssMetaPrefix + `

    Warning: headers are case-insensitive, but value are case-sensitive.


Usage:

    ossutil signurl oss://bucket/object [--method m] [--timeout t] [--encoding-type url]
`,

	sampleText: ` 
    ossutil signurl oss://bucket1/object1
        生成oss://bucket1/object1的GET签名url。

    ossutil signurl oss://bucket1/object1 --method PUT
        生成oss://bucket1/object1的PUT签名url。
`,
}

// SignurlCommand is the command list buckets or objects
type SignurlCommand struct {
	command Command
}

var signURLCommand = SignurlCommand{
	command: Command{
		name:        "signurl",
		nameAlias:   []string{},
		minArgc:     1,
		maxArgc:     2,
		specChinese: specChineseSignurl,
		specEnglish: specEnglishSignurl,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionMethod,
			OptionTimeout,
			OptionEncodingType,
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
func (sc *SignurlCommand) formatHelpForWhole() string {
	return sc.command.formatHelpForWhole()
}

func (sc *SignurlCommand) formatIndependHelp() string {
	return sc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (sc *SignurlCommand) Init(args []string, options OptionMapType) error {
	return sc.command.Init(args, options, sc)
}

// RunCommand simulate inheritance, and polymorphism
func (sc *SignurlCommand) RunCommand() error {
	encodingType, _ := GetString(OptionEncodingType, sc.command.options)
	cloudURL, err := ObjectURLFromString(sc.command.args[0], encodingType)
	if err != nil {
		return err
	}

	method, _ := GetString(OptionMethod, sc.command.options)
	timeout, _ := GetInt(OptionTimeout, sc.command.options)

	str := ""
	if len(sc.command.args) > 1 {
		str = strings.TrimSpace(sc.command.args[1])
	}

	headers, err := sc.parseHeaders(str)
	if err != nil {
		return err
	}

	options, err := sc.command.getOSSOptions(signURLHeaderOptionMap, headers)
	if err != nil {
		return err
	}

	bucket, err := sc.command.ossBucket(cloudURL.bucket)
	if err != nil {
		return err
	}

	str, err = sc.ossSignurlRetry(bucket, cloudURL.object, method, timeout, options...)
	if err != nil {
		return err
	}

	fmt.Println(str)
	return nil
}

func (sc *SignurlCommand) parseHeaders(str string) (map[string]string, error) {
	if str == "" {
		return nil, nil
	}

	headers := map[string]string{}
	sli := strings.Split(str, "#")
	for _, s := range sli {
		pair := strings.SplitN(s, ":", 2)
		name := pair[0]
		value := ""
		if len(pair) > 1 {
			value = pair[1]
		}
		if _, err := fetchHeaderOptionMap(signURLHeaderOptionMap, name); err != nil && !strings.HasPrefix(strings.ToLower(name), strings.ToLower(oss.HTTPHeaderOssMetaPrefix)) {
			return nil, fmt.Errorf("unsupported header:%s, please try \"help %s\" to see supported headers", name, sc.command.name)
		}
		headers[name] = value
	}
	return headers, nil
}

func (sc *SignurlCommand) ossSignurlRetry(bucket *oss.Bucket, object, method string, timeout int64, options ...oss.Option) (string, error) {
	signURLConfig := oss.SignURLConfiguration{Expires: timeout, Method: oss.HTTPMethod(method)}
	retryTimes, _ := GetInt(OptionRetryTimes, sc.command.options)
	for i := 1; ; i++ {
		str, err := bucket.SignURL(object, signURLConfig, options...)
		if err == nil {
			return str, err
		}
		if int64(i) >= retryTimes {
			return str, ObjectError{err, bucket.BucketName, object}
		}
	}
}
