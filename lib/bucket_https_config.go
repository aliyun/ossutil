package lib

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var specChineseHttpsConfig = SpecText{
	synopsisText: "设置、查询bucket的TLS版本设置",

	paramText: "bucket_url [local_xml_file] [options]",

	syntaxText: ` 
    ossutil https-config --method put oss://bucket local_xml_file
    ossutil https-config --method get oss://bucket [local_xml_file]
`,
	detailHelpText: ` 
    cors命令通过设置method选项值为put、get,可以设置、查询bucket的TLS版本设置

用法:
    该命令有两种用法:
	
    1) ossutil https-config --method put oss://bucket local_xml_file
        这个命令从配置文件local_xml_file中读取TLS版本设置，然后设置bucket的TLS版本设置
        配置文件是一个xml格式的文件，举例如下
	   
        <?xml version="1.0" encoding="UTF-8"?>
        <HttpsConfiguration>  
		  <TLS>
			<Enable>true</Enable>   
			<TLSVersion>TLSv1.2</TLSVersion>
			<TLSVersion>TLSv1.3</TLSVersion>
		  </TLS>
		</HttpsConfiguration>
	
    2) ossutil https-config --method get oss://bucket [local_xml_file]
        这个命令查询bucket的TLS版本设置
        如果输入参数local_xml_file，TLS版本设置将输出到该文件，否则输出到屏幕上
`,
	sampleText: ` 
    1) 设置bucket的TLS版本设置
       ossutil https-config --method put oss://bucket local_xml_file

    2) 查询bucket的TLS版本设置，结果输出到标准输出
       ossutil https-config --method get oss://bucket
	
    3) 查询bucket的TLS版本设置，结果输出到本地文件
       ossutil https-config --method get oss://bucket local_xml_file
`,
}

var specEnglishHttpsConfig = SpecText{
	synopsisText: "Set, get or delete the https configuration of the oss bucket",

	paramText: "bucket_url [local_xml_file] [options]",

	syntaxText: ` 
    ossutil https-config --method put oss://bucket local_xml_file
    ossutil https-config --method get oss://bucket [local_xml_file]
`,
	detailHelpText: ` 
    cors command can set、get the https configuration of the oss bucket by
    set method option value to put, get

Usage:
    There are two usages for this command:
	
    1) ossutil https-config --method put oss://bucket local_xml_file
	   
        The command sets the https configuration of bucket from local file local_xml_file
    the local_xml_file is xml format
        The following is an example of the contents of local_xml_file
	   
        <?xml version="1.0" encoding="UTF-8"?>
        <HttpsConfiguration>  
		  <TLS>
			<Enable>true</Enable>   
			<TLSVersion>TLSv1.2</TLSVersion>
			<TLSVersion>TLSv1.3</TLSVersion>
		  </TLS>
		</HttpsConfiguration>
	
    2) ossutil https-config --method get oss://bucket [local_xml_file]
        The command gets the https configuration of bucket
        if you input parameter local_xml_file,the configuration will be output to local_xml_file
        if you don't input parameter local_xml_file,the configuration will be output to stdout
`,
	sampleText: ` 
    1) put https configuration
       ossutil https-config --method put oss://bucket  local_xml_file

    2) get https configuration to stdout
       ossutil https-config --method get oss://bucket
	
    3) get https configuration to local file
       ossutil https-config --method get oss://bucket  local_xml_file
`,
}

type httpsConfigOptionType struct {
	bucketName string
}

type HttpsConfigCommand struct {
	command  Command
	hcOption httpsConfigOptionType
}

var bucketHttpsConfigCommand = HttpsConfigCommand{
	command: Command{
		name:        "https-config",
		nameAlias:   []string{"https-config"},
		minArgc:     1,
		maxArgc:     2,
		specChinese: specChineseHttpsConfig,
		specEnglish: specEnglishHttpsConfig,
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
			OptionMethod,
			OptionLogLevel,
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
			OptionForcePathStyle,
		},
	},
}

// function for FormatHelper interface
func (hcc *HttpsConfigCommand) formatHelpForWhole() string {
	return hcc.command.formatHelpForWhole()
}

func (hcc *HttpsConfigCommand) formatIndependHelp() string {
	return hcc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (hcc *HttpsConfigCommand) Init(args []string, options OptionMapType) error {
	return hcc.command.Init(args, options, hcc)
}

// RunCommand simulate inheritance, and polymorphism
func (hcc *HttpsConfigCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, hcc.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}
	strMethod = strings.ToLower(strMethod)
	if strMethod != "put" && strMethod != "get" {
		return fmt.Errorf("--method value is not in the optional value:put|get")
	}

	bucketUrL, err := StorageURLFromString(hcc.command.args[0], "")
	if err != nil {
		return err
	}

	if !bucketUrL.IsCloudURL() {
		return fmt.Errorf("parameter is not a cloud url,url is %s", bucketUrL.ToString())
	}

	cloudUrl := bucketUrL.(CloudURL)
	if cloudUrl.bucket == "" {
		return fmt.Errorf("bucket name is empty,url is %s", bucketUrL.ToString())
	}

	hcc.hcOption.bucketName = cloudUrl.bucket

	if strMethod == "put" {
		err = hcc.PutBucketHttpsConfig()
	} else if strMethod == "get" {
		err = hcc.GetBucketHttpsConfig()
	}
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
	return err
}

func (hcc *HttpsConfigCommand) PutBucketHttpsConfig() error {
	if len(hcc.command.args) < 2 {
		return fmt.Errorf("missing parameters,the local https config file is empty")
	}

	httpsFile := hcc.command.args[1]
	fileInfo, err := os.Stat(httpsFile)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s is dir,not the expected file", httpsFile)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("%s is empty file", httpsFile)
	}

	// parsing the xml file
	file, err := os.Open(httpsFile)
	if err != nil {
		return err
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// put bucket cors
	client, err := hcc.command.ossClient(hcc.hcOption.bucketName)
	if err != nil {
		return err
	}

	return client.PutBucketHttpsConfigXml(hcc.hcOption.bucketName, string(text))
}

func (hcc *HttpsConfigCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("https config: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (hcc *HttpsConfigCommand) GetBucketHttpsConfig() error {
	client, err := hcc.command.ossClient(hcc.hcOption.bucketName)
	if err != nil {
		return err
	}

	config, err := client.GetBucketHttpsConfig(hcc.hcOption.bucketName)
	if err != nil {
		return err
	}
	output, err := xml.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	var outFile *os.File
	if len(hcc.command.args) >= 2 {
		fileName := hcc.command.args[1]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := hcc.confirm(fileName)
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
