package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var specChineseBucketCallbackPolicy = SpecText{
	synopsisText: "设置、查询、删除bucket的callback policy配置",
	paramText:    "bucket_url local_xml_file [options]",

	syntaxText: ` 
	ossutil callback-policy --method put oss://bucket local_xml_file [options]
    ossutil callback-policy --method get oss://bucket [local_xml_file] [options]
    ossutil callback-policy --method delete oss://bucket [options]
`,
	detailHelpText: ` 
    callback-policy命令通过设置method选项值为put、get、delete,可以设置、查询、删除bucket的callback policy配置

用法:
    该命令有三种用法:
	
    1) ossutil callback-policy --method put oss://bucket local_xml_file [options]
        这个命令从配置文件local_xml_file中读取callback policy配置，然后设置bucket的callback policy规则
        配置文件是一个xml格式的文件，下面是一个例子
   
        <?xml version="1.0" encoding="UTF-8"?>
        <BucketCallbackPolicy>
		    <PolicyItem>
		        <PolicyName>first</PolicyName>
		        <Callback>e1wiY2Fs...R7YnU=</Callback>
		        <CallbackVar>Q2FsbG...mJcIn0=</CallbackVar>
		    </PolicyItem>
		    <PolicyItem>
		        <PolicyName>second</PolicyName>
		        <Callback>e1wiY2Fsb...9keVwiOlwiYnVja2V0PSR7YnU=</Callback>
		        <CallbackVar>Q2Fs...FcIiwgXCJ4OmJcIjpcImJcIn0=</CallbackVar>
		    </PolicyItem>
		</BucketCallbackPolicy>

    2) ossutil callback-policy --method get oss://bucket [local_xml_file] [options]
        这个命令查询bucket的callback policy配置
        如果输入参数local_xml_file，callback policy配置将输出到该文件，否则输出到屏幕上

    3) ossutil callback-policy --method delete oss://bucket [options]
        这个命令删除bucket的callback policy配置
`,
	sampleText: ` 
    1) 设置bucket的callback policy配置
       ossutil callback-policy --method put oss://bucket local_xml_file

    2) 查询bucket的callback policy配置，结果输出到标准输出
       ossutil callback-policy --method get oss://bucket
	
    3) 查询bucket的callback policy配置，结果输出到本地文件
       ossutil callback-policy --method get oss://bucket local_xml_file

    4) 删除bucket的callback policy配置
       ossutil callback-policy --method delete oss://bucket
`,
}

var specEnglishBucketCallbackPolicy = SpecText{
	synopsisText: "Set, get , delete bucket callback policy configuration",
	paramText:    "bucket_url local_xml_file [options]",

	syntaxText: ` 
	ossutil callback-policy --method put oss://bucket local_xml_file [options]
    ossutil callback-policy --method get oss://bucket [local_xml_file] [options]
    ossutil callback-policy --method delete oss://bucket [options]
`,

	detailHelpText: ` 
    callback-policy command can set, get ,delete the callback policy configuration of the oss bucket by
    set method option value to put, get, delete

Usage:
    1) ossutil callback-policy --method put oss://bucket local_xml_file [options]
	   The command sets the lifecycle configuration of bucket from local file local_xml_file
        the local_xml_file is xml format, The following is an example:

        <?xml version="1.0" encoding="UTF-8"?>
        <BucketCallbackPolicy>
		    <PolicyItem>
		        <PolicyName>first</PolicyName>
		        <Callback>e1wiY2Fs...R7YnU=</Callback>
		        <CallbackVar>Q2FsbG...mJcIn0=</CallbackVar>
		    </PolicyItem>
		    <PolicyItem>
		        <PolicyName>second</PolicyName>
		        <Callback>e1wiY2Fsb...9keVwiOlwiYnVja2V0PSR7YnU=</Callback>
		        <CallbackVar>Q2Fs...FcIiwgXCJ4OmJcIjpcImJcIn0=</CallbackVar>
		    </PolicyItem>
		</BucketCallbackPolicy>
	
	2) ossutil callback-policy --method get oss://bucket [local_xml_file] [options]
	   The command gets the callback policy configuration of bucket
       If you input parameter local_xml_file,the configuration will be output to local_xml_file
       If you don't input parameter local_xml_file,the configuration will be output to stdout

    3) ossutil callback-policy --method delete oss://bucket [options]
	   The command delete the callback policy configuration of bucket
`,

	sampleText: ` 
    1) put bucket callback policy
       ossutil callback-policy --method put oss://bucket local_xml_file

    2) get bucket callback policy configuration to stdout
       ossutil callback-policy --method get oss://bucket
	
    3) get bucket callback policy configuration to local file
       ossutil callback-policy --method get oss://bucket local_xml_file

    4) delete bucket callback policy configuration
       ossutil callback-policy --method delete oss://bucket
`,
}

type bucketCallbackPolicyOptionType struct {
	bucketName string
}

type BucketCallbackPolicyCommand struct {
	command   Command
	bcpOption bucketCallbackPolicyOptionType
}

var bucketCallbackPolicyCommand = BucketCallbackPolicyCommand{
	command: Command{
		name:        "callback-policy",
		nameAlias:   []string{"callback-policy"},
		minArgc:     1,
		maxArgc:     2,
		specChinese: specChineseBucketCallbackPolicy,
		specEnglish: specEnglishBucketCallbackPolicy,
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
func (bcpc *BucketCallbackPolicyCommand) formatHelpForWhole() string {
	return bcpc.command.formatHelpForWhole()
}

func (bcpc *BucketCallbackPolicyCommand) formatIndependHelp() string {
	return bcpc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (bcpc *BucketCallbackPolicyCommand) Init(args []string, options OptionMapType) error {
	return bcpc.command.Init(args, options, bcpc)
}

// RunCommand simulate inheritance, and polymorphism
func (bcpc *BucketCallbackPolicyCommand) RunCommand() error {
	strMethod, _ := GetString(OptionMethod, bcpc.command.options)
	if strMethod == "" {
		return fmt.Errorf("--method value is empty")
	}

	strMethod = strings.ToLower(strMethod)
	if strMethod != "put" && strMethod != "get" && strMethod != "delete" {
		return fmt.Errorf("--method value is not in the optional value:put|get|delete")
	}

	srcBucketUrL, err := GetCloudUrl(bcpc.command.args[0], "")
	if err != nil {
		return err
	}

	bcpc.bcpOption.bucketName = srcBucketUrL.bucket

	switch strMethod {
	case "put":
		err = bcpc.PutBucketCallbackPolicy()
	case "get":
		err = bcpc.GetBucketCallbackPolicy()
	case "delete":
		err = bcpc.DeleteBucketCallbackPolicy()
	}
	return err
}

func (bcpc *BucketCallbackPolicyCommand) PutBucketCallbackPolicy() error {
	if len(bcpc.command.args) < 2 {
		return fmt.Errorf("put bucket callback policy need at least 2 parameters,the local xml file is empty")
	}

	xmlFile := bcpc.command.args[1]
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

	client, err := bcpc.command.ossClient(bcpc.bcpOption.bucketName)
	if err != nil {
		return err
	}

	return client.PutBucketCallbackPolicyXml(bcpc.bcpOption.bucketName, string(xmlBody))
}

func (bcpc *BucketCallbackPolicyCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("bucket callback policy: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (bcpc *BucketCallbackPolicyCommand) GetBucketCallbackPolicy() error {
	client, err := bcpc.command.ossClient(bcpc.bcpOption.bucketName)
	if err != nil {
		return err
	}

	output, err := client.GetBucketCallbackPolicyXml(bcpc.bcpOption.bucketName)
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(bcpc.command.args) >= 2 {
		fileName := bcpc.command.args[1]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := bcpc.confirm(fileName)
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

func (bcpc *BucketCallbackPolicyCommand) DeleteBucketCallbackPolicy() error {
	client, err := bcpc.command.ossClient(bcpc.bcpOption.bucketName)
	if err != nil {
		return err
	}
	return client.DeleteBucketCallbackPolicy(bcpc.bcpOption.bucketName)
}
