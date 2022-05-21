package lib

import (
	"encoding/xml"
	"fmt"
	"strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseBucketCname = SpecText{
	synopsisText: "管理bucket cname以及cname token配置",

	paramText: "bucket_url [options]",

	syntaxText: ` 
	ossutil bucket-cname --method put --item token oss://bucket  test-domain.com
    ossutil bucket-cname --method get --item token oss://bucket  test-domain.com
	ossutil bucket-cname --method put oss://bucket  test-domain.com
	ossutil bucket-cname --method delete oss://bucket  test-domain.com
	ossutil bucket-cname --method get oss://bucket  
`,
	detailHelpText: ` 
    cname命令通过设置method选项值可以创建、删除、查询bucket的cname配置

用法:
    1) ossutil bucket-cname --method put --item token oss://bucket test-domain.com
	    该命令会创建一个内部用的token, 设置bucket cname必须先创建这个token
	
	2) ossutil bucket-cname --method get --item token oss://bucket  test-domain.com
	    该命令查询token信息

	3) ossutil bucket-cname --method put oss://bucket test-domain.com
        这个命令设置bucket的cname配置
	    
    4) ossutil bucket-cname --method get oss://bucket
        这个命令查询bucket的cname配置

	5) ossutil bucket-cname --method delete oss://bucket test-domain.com
        这个命令删除bucket的cname配置
`,
	sampleText: ` 
    1) 查询bucket的cname配置，结果输出到标准输出
       ossutil bucket-cname --method get oss://bucket
`,
}

var specEnglishBucketCname = SpecText{
	synopsisText: "manage bucket canme and cname token configuration",

	paramText: "bucket_url [options]",

	syntaxText: ` 
	ossutil bucket-cname --method put --item token oss://bucket  test-domain.com
    ossutil bucket-cname --method get --item token oss://bucket  test-domain.com
	ossutil bucket-cname --method put oss://bucket  test-domain.com
	ossutil bucket-cname --method delete oss://bucket  test-domain.com
	ossutil bucket-cname --method get oss://bucket 
`,
	detailHelpText: ` 
    The command can create, delete and query the cname configuration of a bucket by setting the method option value

Usage:
    1) ossutil bucket-cname --method put --item token oss://bucket test-domain.com
	   This command will create an internal token, which must be created before setting bucket cname
	
	2) ossutil bucket-cname --method get --item token oss://bucket  test-domain.com
	   This command queries the token information

	3) ossutil bucket-cname --method put oss://bucket test-domain.com
	   This command sets the cname configuration of the bucket
	    
    4) ossutil bucket-cname --method get oss://bucket
       This command queries the cname configuration of the bucket

	5) ossutil bucket-cname --method delete oss://bucket test-domain.com
	   This command delete the cname configuration of the bucket
`,
	sampleText: ` 
    1) get cname configuration to stdout
       ossutil bucket-cname --method get oss://bucket
`,
}

type bucketCnameOptionType struct {
	bucketName string
	client     *oss.Client
}

type BucketCnameCommand struct {
	command  Command
	bwOption bucketCnameOptionType
}

var bucketCnameCommand = BucketCnameCommand{
	command: Command{
		name:        "bucket-cname",
		nameAlias:   []string{"bucket-cname"},
		minArgc:     1,
		maxArgc:     2,
		specChinese: specChineseBucketCname,
		specEnglish: specEnglishBucketCname,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionProxyHost,
			OptionProxyUser,
			OptionProxyPwd,
			OptionLogLevel,
			OptionMethod,
			OptionPassword,
			OptionMode,
			OptionECSRoleName,
			OptionTokenTimeout,
			OptionRamRoleArn,
			OptionRoleSessionName,
			OptionReadTimeout,
			OptionConnectTimeout,
			OptionSTSRegion,
			OptionSkipVerifyCert,
			OptionUserAgent,
			OptionItem,
			OptionSignVersion,
			OptionRegion,
			OptionCloudBoxID,
		},
	},
}

// function for FormatHelper interface
func (bwc *BucketCnameCommand) formatHelpForWhole() string {
	return bwc.command.formatHelpForWhole()
}

func (bwc *BucketCnameCommand) formatIndependHelp() string {
	return bwc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (bwc *BucketCnameCommand) Init(args []string, options OptionMapType) error {
	return bwc.command.Init(args, options, bwc)
}

// RunCommand simulate inheritance, and polymorphism
func (bwc *BucketCnameCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, bwc.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}

	strItem, _ := GetString(OptionItem, bwc.command.options)
	strMethod = strings.ToLower(strMethod)
	srcBucketUrL, err := GetCloudUrl(bwc.command.args[0], "")
	if err != nil {
		return err
	}

	bwc.bwOption.bucketName = srcBucketUrL.bucket
	bwc.bwOption.client, err = bwc.command.ossClient(bwc.bwOption.bucketName)
	if err != nil {
		return err
	}

	err = nil
	if strItem == "" {
		if strings.EqualFold(strMethod, "get") {
			err = bwc.GetBucketCname()
		} else if strings.EqualFold(strMethod, "put") {
			err = bwc.PutBucketCname()
		} else if strings.EqualFold(strMethod, "delete") {
			err = bwc.DeleteBucketCname()
		} else {
			err = fmt.Errorf("--method only support get,put,delete")
		}
	} else if strings.EqualFold(strItem, "token") {
		if strings.EqualFold(strMethod, "get") {
			err = bwc.GetBucketCnameToken()
		} else if strings.EqualFold(strMethod, "put") {
			err = bwc.PutBucketCnameToken()
		} else {
			err = fmt.Errorf("only support get bucket token or put bucket token")
		}
	} else {
		err = fmt.Errorf("--item only support token")
	}

	return err
}

func (bwc *BucketCnameCommand) GetBucketCname() error {
	client := bwc.bwOption.client
	output, err := client.GetBucketCname(bwc.bwOption.bucketName)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n\n", output)
	return nil
}

func (bwc *BucketCnameCommand) PutBucketCname() error {
	client := bwc.bwOption.client
	if len(bwc.command.args) < 2 {
		return fmt.Errorf("cname is emtpy")
	}
	cname := bwc.command.args[1]
	err := client.PutBucketCname(bwc.bwOption.bucketName, cname)
	return err
}

func (bwc *BucketCnameCommand) DeleteBucketCname() error {
	client := bwc.bwOption.client
	if len(bwc.command.args) < 2 {
		return fmt.Errorf("cname is emtpy")
	}
	cname := bwc.command.args[1]
	err := client.DeleteBucketCname(bwc.bwOption.bucketName, cname)
	return err
}

func (bwc *BucketCnameCommand) GetBucketCnameToken() error {
	client := bwc.bwOption.client
	if len(bwc.command.args) < 2 {
		return fmt.Errorf("cname is emtpy")
	}
	cname := bwc.command.args[1]
	out, err := client.GetBucketCnameToken(bwc.bwOption.bucketName, cname)
	if err == nil {
		var strXml []byte
		var xmlError error
		if strXml, xmlError = xml.MarshalIndent(out, "", " "); xmlError != nil {
			return xmlError
		}
		fmt.Println(string(strXml))
	}
	return err
}

func (bwc *BucketCnameCommand) PutBucketCnameToken() error {
	client := bwc.bwOption.client
	if len(bwc.command.args) < 2 {
		return fmt.Errorf("cname is emtpy")
	}
	cname := bwc.command.args[1]
	out, err := client.CreateBucketCnameToken(bwc.bwOption.bucketName, cname)
	if err == nil {
		var strXml []byte
		var xmlError error
		if strXml, xmlError = xml.MarshalIndent(out, "", " "); xmlError != nil {
			return xmlError
		}
		fmt.Println(string(strXml))
	}
	return err
}
