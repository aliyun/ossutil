package lib

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"strings"
)

var specChineseRegions = SpecText{
	synopsisText: "查询指定地域或所有支持地域的描述信息",

	paramText: "command_name [endpoint] [local_xml_file] [options]",

	syntaxText: ` 
    ossutil regions --method get endpoint [local_xml_file] [options]
    ossutil regions --method list [local_xml_file] [options]

`,
	detailHelpText: ` 
    regions命令通过设置method选项值为get、list,可以查询指定地域或者所有支持地域的Endpoint信息

用法:
    该命令有二种用法:
	
    1) ossutil regions --method get endpoint [local_xml_file]
       这个命令查询指定地域的Endpoint信息
       如果输入参数local_xml_file，结果将输出到该文件，否则输出到屏幕上
	
    2) ossutil regions --method list [local_xml_file]
       这个命令查询所有支持地域的Endpoint信息
       如果输入参数local_xml_file，结果将输出到该文件，否则输出到屏幕上
    
`,
	sampleText: ` 
    1) 查询指定地域的Endpoint信息，结果输出到标准输出
       ossutil regions --method get endpoint

    2) 查询指定地域的Endpoint信息，结果输出到本地文件
       ossutil regions --method get endpoint local_xml_file
	
    3) 查询所有支持地域的Endpoint信息，结果输出到标准输出
       ossutil regions --method list

    4) 查询所有支持地域的Endpoint信息，结果输出到本地文件
       ossutil regions --method list local_xml_file
`,
}

var specEnglishRegions = SpecText{
	synopsisText: "Query the description information of the specified region or all supported regions",

	paramText: "command_name [endpoint] [local_xml_file] [options]",

	syntaxText: ` 
    ossutil regions --method get endpoint [local_xml_file] [options]
    ossutil regions --method list [local_xml_file] [options]

`,
	detailHelpText: ` 
    The regions command can query endpoint information for a specified region or all supported regions by setting the method option values to get or list

Usage:
    There are 2 usages for this command:
	
    1) ossutil regions --method get endpoint [local_xml_file]
       This command query endpoint information for a specified region
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout
	
    2) ossutil regions --method list [local_xml_file]
       This command query all supported regions for a specified region
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout
    
`,
	sampleText: ` 
    1) query endpoint information for a specified region to stdout
       ossutil regions --method get endpoint

    2) query endpoint information for a specified region to local file
       ossutil regions --method get endpoint local_xml_file
	
    3) query all supported regions to stdout
       ossutil regions --method list

    4) query all supported regions to local file
       ossutil regions --method list local_xml_file
`,
}

type regionOptionType struct {
	bucketName string
}

type DescribeRegionCommand struct {
	command   Command
	regOption regionOptionType
}

var describeRegionCommand = DescribeRegionCommand{
	command: Command{
		name:        "regions",
		nameAlias:   []string{"regions"},
		minArgc:     0,
		maxArgc:     4,
		specChinese: specChineseRegions,
		specEnglish: specEnglishRegions,
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
			OptionMethod,
		},
	},
}

// function for FormatHelper interface
func (dere *DescribeRegionCommand) formatHelpForWhole() string {
	return dere.command.formatHelpForWhole()
}

func (dere *DescribeRegionCommand) formatIndependHelp() string {
	return dere.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (dere *DescribeRegionCommand) Init(args []string, options OptionMapType) error {
	return dere.command.Init(args, options, dere)
}

// RunCommand simulate inheritance, and polymorphism
func (dere *DescribeRegionCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, dere.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}

	strMethod = strings.ToLower(strMethod)
	if strMethod != "get" && strMethod != "list" {
		return fmt.Errorf("--method value is not in the optional value:get|list")
	}

	var err error
	switch strMethod {
	case "get":
		err = dere.DescribeRegions(false)
	case "list":
		err = dere.DescribeRegions(true)
	}
	return err
}

func (dere *DescribeRegionCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("describe regions: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (dere *DescribeRegionCommand) DescribeRegions(isList bool) error {
	client, err := dere.command.ossClient(dere.regOption.bucketName)
	if err != nil {
		return err
	}
	var output string
	endpoint := ""
	var fileName string
	if isList {
		if len(dere.command.args) >= 1 {
			fileName = dere.command.args[0]
		}
	} else {
		if len(dere.command.args) < 1 {
			return fmt.Errorf("get describe regions need at least 1 parameters,the endpoint is empty")
		}
		if len(dere.command.args) >= 2 {
			fileName = dere.command.args[1]
		}
		endpoint = dere.command.args[0]
	}
	output, err = client.DescribeRegionsXml(oss.AddParam("regions", endpoint))
	if err != nil {
		return err
	}
	var outFile *os.File
	if fileName != "" {
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := dere.confirm(fileName)
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
