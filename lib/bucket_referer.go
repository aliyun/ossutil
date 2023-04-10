package lib

import (
	"encoding/xml"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"os"
	"strings"
)

var specChineseBucketReferer = SpecText{
	synopsisText: "设置、查询或者删除bucket的referer配置",

	paramText: "bucket_url referer [options]",

	syntaxText: ` 
	ossutil referer --method put oss://bucket referer [options]
    ossutil referer --method put --item raw oss://bucket local_xml_file [options]
    ossutil referer --method get oss://bucket [local_file]
    ossuitl referer --method delete oss://bucket
`,
	detailHelpText: ` 
    referer命令通过设置method选项值为put、get、delete,可以设置、查询或者删除bucket的referer配置

用法:
    该命令有三种用法:
	
    1) ossutil referer --method put oss://bucket referer [options]
       这个命令将bucket的referer设置成后面的referer值
       referer参数可以连续输入多个

    2) ossutil referer --method put --item raw oss://bucket local_xml_file [options]
       这个命令从配置文件local_xml_file中读取referer配置,然后设置bucket的referer
       配置文件是一个xml格式的文件,举例如下
       
       <?xml version="1.0" encoding="UTF-8"?>
       <RefererConfiguration>
		  <AllowEmptyReferer>false</AllowEmptyReferer>
		  <AllowTruncateQueryString>false</AllowTruncateQueryString>
		  <RefererList>
		        <Referer>http://www.aliyun.com</Referer>
		        <Referer>https://www.aliyun.com</Referer>
		        <Referer>http://www.*.com</Referer>
		        <Referer>https://www.?.aliyuncs.com</Referer>
		  </RefererList>
		  <RefererBlacklist>
		        <Referer>http://www.refuse.com</Referer>
		        <Referer>https://*.hack.com</Referer>
		        <Referer>http://ban.*.com</Referer>
		        <Referer>https://www.?.deny.com</Referer>
		  </RefererBlacklist>
       </RefererConfiguration>
        
    3) ossutil referer --method get oss://bucket [local_xml_file]
        这个命令查询bucket的referer配置
        如果输入参数local_xml_file，referer配置将输出到该文件，否则输出到屏幕上
	
    4)  ossutil referer --method delete oss://bucket
        这个命令删除bucket的referer配置
`,
	sampleText: ` 
    1) 设置bucket的referer配置
       ossutil referer --method put oss://bucket www.test1.com www.test2.com
	
    2) 设置bucket的referer配置，且不允许referer为空
       ossutil referer --method put oss://bucket www.test1.com www.test2.com  --disable-empty-referer

    3) 设置bucket的referer配置，使用local file 的方式
       ossutil referer --method put --item raw oss://bucket local_xml_file

    4) 查询bucket的referer配置，结果输出到标准输出
       ossutil referer --method get oss://bucket
	
    5) 查询bucket的referer配置，结果输出到本地文件
       ossutil referer --method get oss://bucket local_xml_file
	
    6) 删除bucket的referer配置
       ossutil referer --method delete oss://bucket
`,
}

var specEnglishBucketReferer = SpecText{
	synopsisText: "Set、get or delete bucket referer configuration",

	paramText: "bucket_url referer [options]",

	syntaxText: ` 
	ossutil referer --method put oss://bucket referer [options]
    ossutil referer --method put --item raw oss://bucket local_xml_file [options]
    ossutil referer --method get oss://bucket [local_file]
    ossuitl referer --method delete oss://bucket
`,
	detailHelpText: ` 
    referer command can set、get and delete the referer configuration of the oss bucket by
    set method option value to put, get,delete

Usage:
    There are three usages for this command:
	
    1) ossutil referer --method put oss://bucket referer [options]
       This command sets the referer of the bucket to the following referer value.
       You can input many referer parameter.

    2) ossutil referer --method put --item raw oss://bucket local_xml_file [options]
       The command sets the referer configuration of bucket from local file local_xml_file
        the local_xml_file is xml format,for example
       
       <?xml version="1.0" encoding="UTF-8"?>
       <RefererConfiguration>
		  <AllowEmptyReferer>false</AllowEmptyReferer>
		  <AllowTruncateQueryString>false</AllowTruncateQueryString>
		  <RefererList>
		        <Referer>http://www.aliyun.com</Referer>
		        <Referer>https://www.aliyun.com</Referer>
		        <Referer>http://www.*.com</Referer>
		        <Referer>https://www.?.aliyuncs.com</Referer>
		  </RefererList>
		  <RefererBlacklist>
		        <Referer>http://www.refuse.com</Referer>
		        <Referer>https://*.hack.com</Referer>
		        <Referer>http://ban.*.com</Referer>
		        <Referer>https://www.?.deny.com</Referer>
		  </RefererBlacklist>
       </RefererConfiguration>
        
    3) ossutil referer --method get oss://bucket  [local_xml_file]
       The command gets the referer configuration of bucket
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout
	
    4)  ossutil referer --method delete oss://bucket
       The command deletes the referer configuration of bucket
`,
	sampleText: ` 
    1) put bucket referer
       ossutil referer --method put oss://bucket www.test1.com www.test2.com
	
    2) put bucket referer, empty referer is forbidden  
       ossutil referer --method put oss://bucket www.test1.com www.test2.com --disable-empty-referer

    3) put bucket referer, read referer from local xml file 
       ossutil referer --method put --item raw oss://bucket local_xml_file

    4) get referer configuration to stdout
       ossutil referer --method get oss://bucket
	
    5) get referer configuration to local file
       ossutil referer --method get oss://bucket local_xml_file
	
    6) delete referer configuration
       ossutil referer --method delete oss://bucket
`,
}

type bucketReferOptionType struct {
	bucketName        string
	disableEmptyRefer bool
}

type BucketRefererCommand struct {
	command  Command
	brOption bucketReferOptionType
}

var bucketRefererCommand = BucketRefererCommand{
	command: Command{
		name:        "referer",
		nameAlias:   []string{"referer"},
		minArgc:     1,
		maxArgc:     MaxInt,
		specChinese: specChineseBucketReferer,
		specEnglish: specEnglishBucketReferer,
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
			OptionDisableEmptyReferer,
			OptionMethod,
			OptionItem,
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
			OptionSignVersion,
			OptionRegion,
			OptionCloudBoxID,
		},
	},
}

// function for FormatHelper interface
func (brc *BucketRefererCommand) formatHelpForWhole() string {
	return brc.command.formatHelpForWhole()
}

func (brc *BucketRefererCommand) formatIndependHelp() string {
	return brc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (brc *BucketRefererCommand) Init(args []string, options OptionMapType) error {
	return brc.command.Init(args, options, brc)
}

// RunCommand simulate inheritance, and polymorphism
func (brc *BucketRefererCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, brc.command.options)
	strItem, _ := GetString(OptionItem, brc.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}

	strMethod = strings.ToLower(strMethod)
	if strMethod != "put" && strMethod != "get" && strMethod != "delete" {
		return fmt.Errorf("--method value is not in the optional value:put|get|delete")
	}

	strItem = strings.ToLower(strItem)
	if strMethod == "put" && strItem != "" && strItem != "raw" {
		return fmt.Errorf("--item value is not in the optional value:raw")
	}

	srcBucketUrL, err := GetCloudUrl(brc.command.args[0], "")
	if err != nil {
		return err
	}

	brc.brOption.bucketName = srcBucketUrL.bucket
	brc.brOption.disableEmptyRefer, _ = GetBool(OptionDisableEmptyReferer, brc.command.options)

	switch strMethod {
	case "put":
		switch strItem {
		case "":
			err = brc.PutBucketRefer()
		case "raw":
			err = brc.SetBucketRefer()
		}
	case "get":
		err = brc.GetBucketRefer()
	case "delete":
		err = brc.DeleteBucketRefer()
	}
	return err
}

func (brc *BucketRefererCommand) PutBucketRefer() error {
	if len(brc.command.args) < 2 {
		return fmt.Errorf("put bucket referer need at least 2 parameters,the refer is empty")
	}

	referers := brc.command.args[1:len(brc.command.args)]

	// put bucket refer
	client, err := brc.command.ossClient(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	return client.SetBucketReferer(brc.brOption.bucketName, referers, !brc.brOption.disableEmptyRefer)
}

func (brc *BucketRefererCommand) SetBucketRefer() error {
	if len(brc.command.args) < 2 {
		return fmt.Errorf("put bucket referer need at least 2 parameters,the local xml file is empty")
	}

	xmlFile := brc.command.args[1]
	fileInfo, err := os.Stat(xmlFile)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s is dir,not the expected file", xmlFile)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("%s is empty file", xmlFile)
	}

	// parsing the xml file
	file, err := os.Open(xmlFile)
	if err != nil {
		return err
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	refererConfig := oss.SetBucketReferer{}
	err = xml.Unmarshal(text, &refererConfig)
	if err != nil {
		return err
	}

	// put bucket refer
	client, err := brc.command.ossClient(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	return client.SetBucketRefererV2(brc.brOption.bucketName, refererConfig)
}

func (brc *BucketRefererCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("bucket referer: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (brc *BucketRefererCommand) GetBucketRefer() error {
	client, err := brc.command.ossClient(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	referRes, err := client.GetBucketReferer(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	output, err := xml.MarshalIndent(referRes, "  ", "    ")
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(brc.command.args) >= 2 {
		fileName := brc.command.args[1]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := brc.confirm(fileName)
			if !bConitnue {
				return nil
			}
		}

		outFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
		if err != nil {
			return err
		}
		defer outFile.Close()
	} else {
		outFile = os.Stdout
	}

	outFile.Write([]byte(xml.Header))
	outFile.Write(output)

	fmt.Printf("\n\n")

	return nil
}

func (brc *BucketRefererCommand) DeleteBucketRefer() error {

	referers := []string{}

	// put bucket refer
	client, err := brc.command.ossClient(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	return client.SetBucketReferer(brc.brOption.bucketName, referers, true)
}
