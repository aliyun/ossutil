package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var specChineseReservedCapacity = SpecText{
	synopsisText: "创建、更新、查询或列举预留空间的配置及列举预留空间下的bucket列表",

	paramText: "[local_xml_file|ReservedCapacityId] [options]",

	syntaxText: ` 
    ossutil reserved-capacity --method create local_xml_file [options]
    ossutil reserved-capacity --method update id local_xml_file [options]
    ossuitl reserved-capacity --method get id [local_xml_file] [options]
    ossuitl reserved-capacity --method list [local_xml_file] [options]
    ossuitl reserved-capacity --method list-bucket id [local_xml_file] [options]
`,
	detailHelpText: ` 
    reserved-capacity命令通过设置method选项值为create、update、get、list、list-bucket,可以创建、更新、查询或者列举预留空间的配置以及列举预留空间下的bucket列表

用法:
    该命令有五种用法:
	
    1) ossutil reserved-capacity --method create local_xml_file [options]
        这个命令从配置文件local_xml_file中读取预留空间配置,然后添加一个预留空间
        配置文件是一个xml格式的文件
        下面是一个配置文件例子
   
        <?xml version="1.0" encoding="UTF-8"?>
        <ReservedCapacityConfiguration>
		  <Name>your-rc-name</Name>
		  <DataRedundancyType>LRS</DataRedundancyType>
		  <ReservedCapacity>10240</ReservedCapacity>
		</ReservedCapacityConfiguration>
	
    2) ossutil reserved-capacity --method update id local_xml_file [options]
        这个命令从配置文件local_xml_file中读取预留空间配置,然后更新标识为id预留空间
        配置文件是一个xml格式的文件
        下面是一个配置文件例子
        <?xml version="1.0" encoding="UTF-8"?>
        <ReservedCapacityConfiguration>
		  <Status>Enabled</Status>
		  <ReservedCapacity>10240</ReservedCapacity>
		  <AutoExpansionSize>100</AutoExpansionSize>
		  <AutoExpansionMaxSize>20480</AutoExpansionMaxSize>
		</ReservedCapacityConfiguration>
	
    3) ossuitl reserved-capacity --method get id [local_xml_file] [options]
        这个命令查询标识为id的预留空间配置
        如果输入参数local_xml_file，清单配置将输出到该文件，否则输出到屏幕上

    4) ossuitl reserved-capacity --method list [local_xml_file] [options]
        这个命令列举预留空间配置
        如果输入参数local_xml_file，清单配置将输出到该文件，否则输出到屏幕上

    5) ossuitl reserved-capacity --method list-bucket id [local_xml_file] [options]
        这个命令列举标识为id下面的bucket列表
        如果输入参数local_xml_file，清单配置将输出到该文件，否则输出到屏幕上

`,
	sampleText: ` 
    1) 创建预留空间
       ossutil reserved-capacity --method create local_xml_file
    
    2) 更新预留空间
       ossutil reserved-capacity --method update id local_xml_file
	
    3) 查询预留空间配置，结果输出到标准输出
       ossutil reserved-capacity --method get id

    4) 查询预留空间配置，结果输出到本地文件
       ossutil reserved-capacity --method get id local_xml_file
	
    5) 列举预留空间配置，结果输出到标准输出
       ossuitl reserved-capacity --method list

    6) 列举预留空间配置，结果输出到本地文件
       ossuitl reserved-capacity --method list local_xml_file

    7) 列举标识为id下面的bucket列表，结果输出到标准输出
       ossuitl reserved-capacity --method list-bucket

    8) 列举标识为id下面的bucket列表，结果输出到本地文件
       ossuitl reserved-capacity --method list-bucket local_xml_file
`,
}

var specEnglishReservedCapacity = SpecText{
	synopsisText: "Create, update, query, or list the configuration of reserved capacity and list buckets under the reserved capacity",

	paramText: "[local_xml_file|ReservedCapacityId] [options]",

	syntaxText: ` 
    ossutil reserved-capacity --method create local_xml_file
    ossutil reserved-capacity --method update id local_xml_file
    ossuitl reserved-capacity --method get id [local_xml_file]
    ossuitl reserved-capacity --method list [local_xml_file]
    ossuitl reserved-capacity --method list-bucket id [local_xml_file]
`,
	detailHelpText: ` 
    reserved-capacity command  command can create, update, get, list the reserved capacity configuration and list buckets under the the reserved capacity by set method option value to create, update, get, list, list-bucket

Usage:
	
    1) ossutil reserved-capacity --method create local_xml_file [options]
       The command creates the reserved capacity configuration from local file local_xml_file
       the local_xml_file is xml format
       The following is an example of configure:
   
       <?xml version="1.0" encoding="UTF-8"?>
       <ReservedCapacityConfiguration>
		  <Name>your-rc-name</Name>
		  <DataRedundancyType>LRS</DataRedundancyType>
		  <ReservedCapacity>10240</ReservedCapacity>
       </ReservedCapacityConfiguration>
	
    2) ossutil reserved-capacity --method update id local_xml_file [options]
       The command updates the reserved capacity configuration from local file local_xml_file
       the local_xml_file is xml format
       The following is an example of configure:
        
       <?xml version="1.0" encoding="UTF-8"?>
       <ReservedCapacityConfiguration>
		  <Status>Enabled</Status>
		  <ReservedCapacity>10240</ReservedCapacity>
		  <AutoExpansionSize>100</AutoExpansionSize>
		  <AutoExpansionMaxSize>20480</AutoExpansionMaxSize>
       </ReservedCapacityConfiguration>
	
    3) ossuitl reserved-capacity --method get id [local_xml_file] [options]
        The command gets the reserved capacity configuration, The identifier of the reserved capacity is id
        If you input parameter local_xml_file,the configuration will be output to local_xml_file
        If you don't input parameter local_xml_file,the configuration will be output to stdout

    4) ossuitl reserved-capacity --method list [local_xml_file] [options]
       List all the reserved capacity configuration
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout

    5) ossuitl reserved-capacity --method list-bucket id [local_xml_file] [options]
       List all the buckets under the reserved capacity,The identifier of the reserved capacity is id
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout
`,
	sampleText: ` 
    1) create reserved capacity
       ossutil reserved-capacity --method create local_xml_file
    
    2) update reserved capacity
       ossutil reserved-capacity --method update id local_xml_file
	
    3) get reserved capacity configuration to stdout, The identifier of the reserved capacity is id
       ossutil reserved-capacity --method get id

    4) get reserved capacity configuration to local file, The identifier of the reserved capacity is id
       ossutil reserved-capacity --method get id local_xml_file
	
    5) list reserved capacity configuration to stdout
       ossuitl reserved-capacity --method list

    6) list reserved capacity configuration to local file
       ossuitl reserved-capacity --method list local_xml_file

    7) list buckets under the reserved capacity to stdout, The identifier of the reserved capacity is id
       ossuitl reserved-capacity --method list-bucket

    8) list buckets under the reserved capacity to local file, The identifier of the reserved capacity is id
       ossuitl reserved-capacity --method list-bucket local_xml_file
`,
}

type ReservedCommand struct {
	command    Command
	bucketName string
}

var reservedCommand = ReservedCommand{
	command: Command{
		name:        "reserved-capacity",
		nameAlias:   []string{"reserved-capacity"},
		minArgc:     0,
		maxArgc:     2,
		specChinese: specChineseReservedCapacity,
		specEnglish: specEnglishReservedCapacity,
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
			OptionMode,
			OptionECSRoleName,
			OptionTokenTimeout,
			OptionRamRoleArn,
			OptionRoleSessionName,
			OptionReadTimeout,
			OptionConnectTimeout,
			OptionSTSRegion,
			OptionMethod,
			OptionItem,
			OptionSkipVerifyCert,
			OptionUserAgent,
			OptionSignVersion,
			OptionRegion,
			OptionCloudBoxID,
		},
	},
}

// function for FormatHelper interface
func (reservedc *ReservedCommand) formatHelpForWhole() string {
	return reservedc.command.formatHelpForWhole()
}

func (reservedc *ReservedCommand) formatIndependHelp() string {
	return reservedc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (reservedc *ReservedCommand) Init(args []string, options OptionMapType) error {
	return reservedc.command.Init(args, options, reservedc)
}

// RunCommand simulate inheritance, and polymorphism
func (reservedc *ReservedCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, reservedc.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}
	strMethod = strings.ToLower(strMethod)
	if strMethod != "create" && strMethod != "update" && strMethod != "get" && strMethod != "list" && strMethod != "list-bucket" {
		return fmt.Errorf("--method value is not in the optional value:create|update|get|list|list-bucket")
	}
	var err error
	switch strMethod {
	case "create":
		err = reservedc.CreateReservedCapacity(true)
	case "get":
		err = reservedc.GetReservedCapacity(false)
	case "update":
		err = reservedc.CreateReservedCapacity(false)
	case "list":
		err = reservedc.GetReservedCapacity(true)
	case "list-bucket":
		err = reservedc.ListBucketWithReservedCapacity()
	}
	return err
}

func (reservedc *ReservedCommand) CreateReservedCapacity(isCreate bool) error {
	xmlFile := ""
	id := ""
	if isCreate {
		if len(reservedc.command.args) < 1 {
			return fmt.Errorf("create reserved capacity need at least 1 parameters,the local xml file is empty")
		}
		xmlFile = reservedc.command.args[0]
	} else {
		if len(reservedc.command.args) < 2 {
			return fmt.Errorf("update reserved capacity need at least 2 parameters,the local xml file is empty")
		}
		id = reservedc.command.args[0]
		xmlFile = reservedc.command.args[1]
	}

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
	client, err := reservedc.command.ossClient(reservedc.bucketName)
	if err != nil {
		return err
	}
	return client.CreateReservedCapacityXml(id, string(text))
}

func (reservedc *ReservedCommand) GetReservedCapacity(isList bool) error {

	fileName := ""
	id := ""
	if isList {
		if len(reservedc.command.args) == 1 {
			fileName = reservedc.command.args[0]
		}
	} else {
		if len(reservedc.command.args) < 1 {
			return fmt.Errorf("get reserved capacity need at least 1 parameters,the id is empty")
		}
		id = reservedc.command.args[0]
		if len(reservedc.command.args) == 2 {
			fileName = reservedc.command.args[1]
		}
	}
	client, err := reservedc.command.ossClient(reservedc.bucketName)
	if err != nil {
		return err
	}
	output, err := client.GetReservedCapacityXml(id)
	if err != nil {
		return err
	}

	var outFile *os.File
	if fileName == "" {
		outFile = os.Stdout
	} else {
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := reservedc.confirm(fileName)
			if !bConitnue {
				return nil
			}
		}
		outFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
		if err != nil {
			return err
		}
		defer outFile.Close()
	}
	outFile.Write([]byte(output))
	fmt.Printf("\n\n")
	return nil
}
func (reservedc *ReservedCommand) ListBucketWithReservedCapacity() error {
	if len(reservedc.command.args) < 1 {
		return fmt.Errorf("list bucket under reserved capacity need at least 1 parameters,the id is empty")
	}
	fileName := ""
	id := reservedc.command.args[0]
	if len(reservedc.command.args) == 2 {
		fileName = reservedc.command.args[1]
	}
	client, err := reservedc.command.ossClient(reservedc.bucketName)
	if err != nil {
		return err
	}
	output, err := client.ListBucketWithReservedCapacityXml(id)
	if err != nil {
		return err
	}

	var outFile *os.File
	if fileName == "" {
		outFile = os.Stdout
	} else {
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := reservedc.confirm(fileName)
			if !bConitnue {
				return nil
			}
		}

		outFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
		if err != nil {
			return err
		}
		defer outFile.Close()
	}
	outFile.Write([]byte(output))
	fmt.Printf("\n\n")
	return nil
}
func (reservedc *ReservedCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("reserved capacity: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}
