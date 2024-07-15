package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var specChineseBucketAccessPoint = SpecText{
	synopsisText: "设置、查询bucket的access point配置",
	paramText:    "bucket_url|ap-alias local_xml_file [options]",

	syntaxText: ` 
	ossutil access-point --method put oss://bucket local_xml_file [local_xml_result] [options]
    ossutil access-point --method get oss://bucket ap-name [local_xml_file] [options]
    ossutil access-point --method list oss://bucket [local_xml_file] [options]
    ossutil access-point --method delete oss://bucket ap-name [options]
    ossutil access-point --method put --item policy oss://bucket local_json_file [options]
    ossutil access-point --method put --item policy oss://ap-alias local_json_file [options]
    ossutil access-point --method get --item policy oss://bucket ap-name [local_json_file] [options]
    ossutil access-point --method get --item policy oss://ap-alias ap-name [local_json_file] [options]
    ossutil access-point --method delete --item policy oss://bucket ap-name [options]
    ossutil access-point --method delete --item policy oss://ap-alias ap-name [options]
`,
	detailHelpText: ` 
    access-point命令通过设置method选项值为put、get、list和delete,可以设置、查询、列举和删除bucket的access point配置

用法:
    该命令有七种用法:
	
    1) ossutil access-point --method put oss://bucket local_xml_file [local_xml_result] [options]
       这个命令从配置文件local_xml_file中读取access point配置，然后设置bucket的access point规则
       如果输入参数local_xml_result，执行结果将输出到该文件，否则输出到屏幕上
       配置文件是一个xml格式的文件，下面是一个例子
   
		<?xml version="1.0" encoding="UTF-8"?>
		<CreateAccessPointConfiguration>
			<AccessPointName>ap-01</AccessPointName>
			<NetworkOrigin>vpc</NetworkOrigin>
			<VpcConfiguration>
			  <VpcId>vpc-t4nlw426y44rd3iq4****</VpcId>
			</VpcConfiguration>
		</CreateAccessPointConfiguration>
		
		或
		
		<?xml version="1.0" encoding="UTF-8"?>
		<CreateAccessPointConfiguration>
			<AccessPointName>ap-01</AccessPointName>
			<NetworkOrigin>internet</NetworkOrigin>
		</CreateAccessPointConfiguration>

    2) ossutil access-point --method get oss://bucket ap-name [local_xml_file] [options]
       这个命令查询bucket的access point名字为ap-name的配置
       如果输入参数local_xml_file，access point配置将输出到该文件，否则输出到屏幕上

    3) ossutil access-point --method list oss://bucket [local_xml_file] [options]
       这个命令列举bucket的access point的配置
       如果输入参数local_xml_file，access point配置将输出到该文件，否则输出到屏幕上

    4) ossutil access-point --method delete oss://bucket ap-name [options]
       这个命令删除bucket的access point名字为ap-name的配置

    5) ossutil access-point --method put --item policy oss://bucket ap-name local_json_file [options]
       这个命令从配置文件local_json_file中读取policy配置,然后设置bucket的access point名字为ap-name的policy规则
	   配置文件是一个json格式的文件,举例如下：

        {
			"Version": "1",
			"Statement": [{
				"Effect": "Deny",
				"Action": [
					"oss:PutObject",
					"oss:GetObject"
				],
				"Principal": [
					"123456"
				],
				"Resource": [
					"acs:oss:oss-cn-hangzhou:123456:accesspoint/ap-name",
					"acs:oss:oss-cn-hangzhou:123456:accesspoint/ap-name/object/*"
				]
			}]
		}

    6) ossutil access-point --method get --item policy oss://bucket ap-name [local_json_file] [options]
       这个命令获取access point名字为ap-name的bucket的policy规则
	   如果输入参数local_json_file，access point配置将输出到该文件，否则输出到屏幕上

    7) ossutil access-point --method get --item policy oss://bucket ap-name [options]
       这个命令删除access point名字为ap-name的bucket的policy规则
`,
	sampleText: ` 
    1) 设置bucket的access point配置，结果输出到标准输出
       ossutil access-point --method put oss://bucket local_xml_file

    2) 设置bucket的access point配置，结果输出到本地文件
       ossutil access-point --method put oss://bucket local_xml_file local_xml_result

    3) 查询bucket的access point名字为ap-name配置，结果输出到标准输出
       ossutil access-point --method get oss://bucket ap-name

    4) 查询bucket的access point名字为ap-name配置，结果输出到本地文件
       ossutil access-point --method get oss://bucket ap-name local_xml_file
	
    5) 列举bucket的access point配置，结果输出到标准输出
       ossutil access-point --method list oss://bucket

    6) 列举bucket的access point配置，结果输出到本地文件
       ossutil access-point --method list oss://bucket local_xml_file

    7) 删除bucket的access point名字为ap-name配置
       ossutil access-point --method delete oss://bucket ap-name

    8) 设置bucket的access point名字为ap-name的policy规则
       ossutil access-point --method put --item policy oss://bucket ap-name local_json_file

    9) 查询bucket的access point名字为ap-name的policy规则，结果输出到标准输出
       ossutil access-point --method get --item policy oss://bucket ap-name

    10) 查询bucket的access point名字为ap-name的policy规则，结果输出到本地文件
       ossutil access-point --method get --item policy oss://bucket ap-name local_json_file

    11) 删除bucket的access point名字为ap-name的policy规则
       ossutil access-point --method delete --item policy oss://bucket ap-name
`,
}

var specEnglishBucketAccessPoint = SpecText{
	synopsisText: "Set, get, list, delete bucket access point configuration",
	paramText:    "bucket_url|ap-alias local_xml_file [options]",

	syntaxText: ` 
	ossutil access-point --method put oss://bucket local_xml_file [local_xml_result] [options]
    ossutil access-point --method get oss://bucket ap-name [local_xml_file] [options]
    ossutil access-point --method list oss://bucket [local_xml_file] [options]
    ossutil access-point --method delete oss://bucket ap-name [options]
    ossutil access-point --method put --item policy oss://bucket local_json_file [options]
    ossutil access-point --method put --item policy oss://ap-alias local_json_file [options]
    ossutil access-point --method get --item policy oss://bucket ap-name [local_json_file] [options]
    ossutil access-point --method get --item policy oss://ap-alias ap-name [local_json_file] [options]
    ossutil access-point --method delete --item policy oss://bucket ap-name [options]
    ossutil access-point --method delete --item policy oss://ap-alias ap-name [options]
`,

	detailHelpText: ` 
    access-point command can set, get, list, delete the access point configuration of the oss bucket by
    set method option value to put, get, list, delete,

Usage:
    1) ossutil access-point --method put oss://bucket local_xml_file [local_xml_result] [options]
	   The command sets the access point configuration of bucket from local file local_xml_file
       If you input parameter local_xml_result,the configuration will be output to local_xml_result
       If you don't input parameter local_xml_result,the configuration will be output to stdout
       the local_xml_file is xml format, The following is an example:

        <?xml version="1.0" encoding="UTF-8"?>
		<CreateAccessPointConfiguration>
			<AccessPointName>ap-01</AccessPointName>
			<NetworkOrigin>vpc</NetworkOrigin>
			<VpcConfiguration>
			  <VpcId>vpc-t4nlw426y44rd3iq4****</VpcId>
			</VpcConfiguration>
		</CreateAccessPointConfiguration>

		or

       <?xml version="1.0" encoding="UTF-8"?>
		<CreateAccessPointConfiguration>
			<AccessPointName>ap-01</AccessPointName>
			<NetworkOrigin>internet</NetworkOrigin>
		</CreateAccessPointConfiguration>
	
	2) ossutil access-point --method get oss://bucket ap-name [local_xml_file] [options]
	   The command gets the access point name of bucket as ap-name configuration
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout

    3) ossutil access-point --method list oss://bucket [local_xml_file] [options]
	   The command lists the access point of bucket configuration
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout

    4) ossutil access-point --method delete oss://bucket ap-name [options]
	   The command delete the access point name of bucket as ap-name configuration

    5) ossutil access-point --method put --item policy oss://bucket ap-name local_json_file [options]
       The command sets the policy for access point name of bucket as ap-name from local file local_json_file
       the local_json_file is json format,for example

        {
			"Version": "1",
			"Statement": [{
				"Effect": "Deny",
				"Action": [
					"oss:PutObject",
					"oss:GetObject"
				],
				"Principal": [
					"123456"
				],
				"Resource": [
					"acs:oss:oss-cn-hangzhou:123456:accesspoint/ap-name",
					"acs:oss:oss-cn-hangzhou:123456:accesspoint/ap-name/object/*"
				]
			}]
		}

    6) ossutil access-point --method get --item policy oss://bucket ap-name [local_json_file] [options]
       The command gets the policy for access point name of bucket as ap-name from local file local_json_file
       If you input parameter local_json_file,the configuration will be output to local_json_file
       If you don't input parameter local_json_file,the configuration will be output to stdout

    7) ossutil access-point --method get --item policy oss://bucket ap-name [options]
       The command deletes the policy for access point name of bucket as ap-name
`,

	sampleText: ` 
    1) put bucket access point to stdout
       ossutil access-point --method put oss://bucket local_xml_file

    2) put bucket access point to local file
       ossutil access-point --method put oss://bucket local_xml_file local_xml_result

    3) get the access point name of bucket as ap-name configuration to stdout
       ossutil access-point --method get oss://bucket ap-name

    4) get the access point name of bucket as ap-name configuration to local file
       ossutil access-point --method get oss://bucket ap-name local_xml_file

	5) list bucket access point configuration to stdout
       ossutil access-point --method list oss://bucket

	6) list bucket access point configuration to local file
       ossutil access-point --method list oss://bucket local_xml_file

    7) delete the access point name of bucket as ap-name configuration
       ossutil access-point --method delete oss://bucket ap-name

    8) put the policy for access point name of bucket as ap-name
       ossutil access-point --method put --item policy oss://bucket ap-name local_json_file

    9) get the policy for access point name of bucket as ap-name to stdout
       ossutil access-point --method get --item policy oss://bucket ap-name

    10) get the policy for access point name of bucket as ap-name to local file
       ossutil access-point --method get --item policy oss://bucket ap-name local_json_file

    11) delete the policy for access point name of bucket as ap-name
       ossutil access-point --method delete --item policy oss://bucket ap-name
`,
}

type bucketAccessPointOptionType struct {
	bucketName string
}

type BucketAccessPointCommand struct {
	command  Command
	blOption bucketAccessPointOptionType
}

var bucketAccessPointCommand = BucketAccessPointCommand{
	command: Command{
		name:        "access-point",
		nameAlias:   []string{"access-point"},
		minArgc:     1,
		maxArgc:     3,
		specChinese: specChineseBucketAccessPoint,
		specEnglish: specEnglishBucketAccessPoint,
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
			OptionItem,
		},
	},
}

// function for FormatHelper interface
func (bapc *BucketAccessPointCommand) formatHelpForWhole() string {
	return bapc.command.formatHelpForWhole()
}

func (bapc *BucketAccessPointCommand) formatIndependHelp() string {
	return bapc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (bapc *BucketAccessPointCommand) Init(args []string, options OptionMapType) error {
	return bapc.command.Init(args, options, bapc)
}

// RunCommand simulate inheritance, and polymorphism
func (bapc *BucketAccessPointCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, bapc.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}

	strMethod = strings.ToLower(strMethod)
	if strMethod != "put" && strMethod != "get" && strMethod != "list" && strMethod != "delete" {
		return fmt.Errorf("--method value is not in the optional value:put|get|list|delete")
	}

	strItem, _ := GetString(OptionItem, bapc.command.options)

	strItem = strings.ToLower(strItem)
	if strItem != "" && strItem != "policy" {
		return fmt.Errorf("--item value is not in the optional value:policy")
	}

	srcBucketUrL, err := GetCloudUrl(bapc.command.args[0], "")
	if err != nil {
		return err
	}

	bapc.blOption.bucketName = srcBucketUrL.bucket

	if strItem == "policy" {
		switch strMethod {
		case "put":
			err = bapc.PutAccessPointPolicy()
		case "get":
			err = bapc.GetAccessPointPolicy()
		case "delete":
			err = bapc.DeleteAccessPointPolicy()
		}
	} else {
		switch strMethod {
		case "put":
			err = bapc.PutBucketAccessPoint()
		case "get":
			err = bapc.GetBucketAccessPoint()
		case "list":
			err = bapc.ListBucketAccessPoint()
		case "delete":
			err = bapc.DeleteBucketAccessPoint()
		}
	}

	return err
}

func (bapc *BucketAccessPointCommand) PutBucketAccessPoint() error {
	if len(bapc.command.args) < 2 {
		return fmt.Errorf("put bucket access point need at least 2 parameters,the local xml file is empty")
	}

	xmlFile := bapc.command.args[1]
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

	client, err := bapc.command.ossClient(bapc.blOption.bucketName)
	if err != nil {
		return err
	}

	output, err := client.CreateBucketAccessPointXml(bapc.blOption.bucketName, string(xmlBody))
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(bapc.command.args) >= 3 {
		fileName := bapc.command.args[2]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := bapc.confirm(fileName)
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

func (bapc *BucketAccessPointCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("bucket access point: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (bapc *BucketAccessPointCommand) GetBucketAccessPoint() error {
	if len(bapc.command.args) < 2 {
		return fmt.Errorf("get bucket access point need at least 2 parameters,the parameter ap name is empty")
	}
	client, err := bapc.command.ossClient(bapc.blOption.bucketName)
	if err != nil {
		return err
	}
	apName := bapc.command.args[1]
	output, err := client.GetBucketAccessPointXml(bapc.blOption.bucketName, apName)
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(bapc.command.args) >= 3 {
		fileName := bapc.command.args[2]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := bapc.confirm(fileName)
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

func (bapc *BucketAccessPointCommand) ListBucketAccessPoint() error {
	client, err := bapc.command.ossClient(bapc.blOption.bucketName)
	if err != nil {
		return err
	}

	output, err := client.ListBucketAccessPointXml(bapc.blOption.bucketName)
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(bapc.command.args) >= 2 {
		fileName := bapc.command.args[1]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := bapc.confirm(fileName)
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

func (bapc *BucketAccessPointCommand) DeleteBucketAccessPoint() error {
	if len(bapc.command.args) < 2 {
		return fmt.Errorf("delete bucket access point need at least 2 parameters,the parameter ap name is empty")
	}
	client, err := bapc.command.ossClient(bapc.blOption.bucketName)
	if err != nil {
		return err
	}
	apName := bapc.command.args[1]
	return client.DeleteBucketAccessPoint(bapc.blOption.bucketName, apName)
}

func (bapc *BucketAccessPointCommand) PutAccessPointPolicy() error {
	if len(bapc.command.args) < 3 {
		return fmt.Errorf("put access point policy need at least 3 parameters,the local xml file is empty")
	}
	xmlFile := bapc.command.args[2]
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

	client, err := bapc.command.ossClient(bapc.blOption.bucketName)
	if err != nil {
		return err
	}

	apName := bapc.command.args[1]

	err = client.PutAccessPointPolicy(bapc.blOption.bucketName, apName, string(xmlBody))
	if err != nil {
		return err
	}
	return nil
}

func (bapc *BucketAccessPointCommand) GetAccessPointPolicy() error {
	if len(bapc.command.args) < 2 {
		return fmt.Errorf("get bucket access point policy need at least 2 parameters,the parameter ap name is empty")
	}
	client, err := bapc.command.ossClient(bapc.blOption.bucketName)
	if err != nil {
		return err
	}
	apName := bapc.command.args[1]
	output, err := client.GetAccessPointPolicy(bapc.blOption.bucketName, apName)
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(bapc.command.args) >= 3 {
		fileName := bapc.command.args[2]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := bapc.confirm(fileName)
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

	var prettyJson bytes.Buffer
	if err = json.Indent(&prettyJson, []byte(output), "", "  "); err != nil {
		return err
	}
	outFile.Write([]byte(prettyJson.String()))
	fmt.Printf("\n\n")
	return nil
}

func (bapc *BucketAccessPointCommand) DeleteAccessPointPolicy() error {
	if len(bapc.command.args) < 2 {
		return fmt.Errorf("delete bucket access point policy need at least 2 parameters,the parameter ap name is empty")
	}
	client, err := bapc.command.ossClient(bapc.blOption.bucketName)
	if err != nil {
		return err
	}
	apName := bapc.command.args[1]
	return client.DeleteAccessPointPolicy(bapc.blOption.bucketName, apName)
}
