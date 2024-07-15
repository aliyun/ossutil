package lib

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"os"
	"strings"
)

var specChineseBucketArchiveDirectRead = SpecText{
	synopsisText: "设置、查询bucket的archive direct read配置",
	paramText:    "bucket_url local_xml_file [options]",

	syntaxText: ` 
	ossutil archive-direct-read --method put oss://bucket local_xml_file [options]
    ossutil archive-direct-read --method get oss://bucket [local_xml_file] [options]
`,
	detailHelpText: ` 
    archive-direct-read命令通过设置method选项值为put、get,可以设置、查询bucket的archive direct read配置

用法:
    该命令有二种用法:
	
    1) ossutil archive-direct-read --method put oss://bucket local_xml_file [options]
        这个命令从配置文件local_xml_file中读取archive direct read配置，然后设置bucket的archive direct read规则
        配置文件是一个xml格式的文件，下面是一个所有规则的例子
   
        <?xml version="1.0" encoding="UTF-8"?>
        <ArchiveDirectReadConfiguration>
            <Enabled>true</Enabled>
        </ArchiveDirectReadConfiguration>

    2) ossutil archive-direct-read --method get oss://bucket [local_xml_file] [options]
        这个命令查询bucket的archive direct read配置
        如果输入参数local_xml_file，archive direct read配置将输出到该文件，否则输出到屏幕上

`,
	sampleText: ` 
    1) 设置bucket的archive direct read配置
       ossutil archive-direct-read --method put oss://bucket local_xml_file

    2) 查询bucket的archive direct read配置，结果输出到标准输出
       ossutil archive-direct-read --method get oss://bucket
	
    3) 查询bucket的archive direct read配置，结果输出到本地文件
       ossutil archive-direct-read --method get oss://bucket local_xml_file
`,
}

var specEnglishBucketArchiveDirectRead = SpecText{
	synopsisText: "Set, get bucket archive direct read configuration",
	paramText:    "bucket_url local_xml_file [options]",

	syntaxText: ` 
	ossutil archive-direct-read --method put oss://bucket local_xml_file [options]
    ossutil archive-direct-read --method get oss://bucket [local_xml_file] [options]
`,

	detailHelpText: ` 
    archive-direct-read command can set, get the archive direct read configuration of the oss bucket by
    set method option value to put, get

Usage:
    1) ossutil archive-direct-read --method put oss://bucket local_xml_file [options]
	   The command sets the lifecycle configuration of bucket from local file local_xml_file
        the local_xml_file is xml format, you can choose to configure only some rules
        The following is an example of all rules:

        <?xml version="1.0" encoding="UTF-8"?>
        <ArchiveDirectReadConfiguration>
            <Enabled>true</Enabled>
        </ArchiveDirectReadConfiguration>
	
	2) ossutil archive-direct-read --method get oss://bucket [local_xml_file] [options]
	   The command gets the archive direct read configuration of bucket
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout
`,

	sampleText: ` 
    1) put bucket archive direct read
       ossutil archive-direct-read --method put oss://bucket local_xml_file

    2) get bucket archive direct read configuration to stdout
       ossutil archive-direct-read --method get oss://bucket
	
    3) get bucket archive direct read configuration to local file
       ossutil archive-direct-read --method get oss://bucket local_xml_file
`,
}

type bucketArchiveDirectReadOptionType struct {
	bucketName string
}

type BucketArchiveDirectReadCommand struct {
	command  Command
	blOption bucketArchiveDirectReadOptionType
}

var bucketArchiveDirectReadCommand = BucketArchiveDirectReadCommand{
	command: Command{
		name:        "archive-direct-read",
		nameAlias:   []string{"archive-direct-read"},
		minArgc:     1,
		maxArgc:     2,
		specChinese: specChineseBucketArchiveDirectRead,
		specEnglish: specEnglishBucketArchiveDirectRead,
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
			OptionSignVersion,
			OptionRegion,
			OptionCloudBoxID,
		},
	},
}

// function for FormatHelper interface
func (badr *BucketArchiveDirectReadCommand) formatHelpForWhole() string {
	return badr.command.formatHelpForWhole()
}

func (badr *BucketArchiveDirectReadCommand) formatIndependHelp() string {
	return badr.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (badr *BucketArchiveDirectReadCommand) Init(args []string, options OptionMapType) error {
	return badr.command.Init(args, options, badr)
}

// RunCommand simulate inheritance, and polymorphism
func (badr *BucketArchiveDirectReadCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, badr.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}

	strMethod = strings.ToLower(strMethod)
	if strMethod != "put" && strMethod != "get" && strMethod != "delete" {
		return fmt.Errorf("--method value is not in the optional value:put|get|delete")
	}

	srcBucketUrL, err := GetCloudUrl(badr.command.args[0], "")
	if err != nil {
		return err
	}

	badr.blOption.bucketName = srcBucketUrL.bucket

	switch strMethod {
	case "put":
		err = badr.PutBucketArchiveDirectRead()
	case "get":
		err = badr.GetBucketArchiveDirectRead()
	}
	return err
}

func (badr *BucketArchiveDirectReadCommand) PutBucketArchiveDirectRead() error {
	if len(badr.command.args) < 2 {
		return fmt.Errorf("put bucket archive direct read need at least 2 parameters,the local xml file is empty")
	}

	xmlFile := badr.command.args[1]
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
	xmlBody, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// put bucket lifecycle
	client, err := badr.command.ossClient(badr.blOption.bucketName)
	if err != nil {
		return err
	}

	options := []oss.Option{oss.AllowSameActionOverLap(true)}
	return client.PutBucketArchiveDirectReadXml(badr.blOption.bucketName, string(xmlBody), options...)
}

func (badr *BucketArchiveDirectReadCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("bucket archive direct read: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (badr *BucketArchiveDirectReadCommand) GetBucketArchiveDirectRead() error {
	client, err := badr.command.ossClient(badr.blOption.bucketName)
	if err != nil {
		return err
	}

	output, err := client.GetBucketArchiveDirectReadXml(badr.blOption.bucketName)
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(badr.command.args) >= 2 {
		fileName := badr.command.args[1]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := badr.confirm(fileName)
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

	outFile.Write([]byte(output))

	fmt.Printf("\n\n")

	return nil
}
