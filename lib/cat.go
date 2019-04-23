package lib

import (
	"fmt"
	"io"
	"os"
)

var specChineseCat = SpecText{
	synopsisText: "将文件内容输出到标准输出",

	paramText: "object [options]",

	syntaxText: ` 
	ossutil cat oss://bucket/object 
`,
	detailHelpText: ` 
    cat命令可以将oss的object内容输出到标准输出,object内容最好是文本格式

用法:
    该命令仅有一种用法:
	
    1) ossutil cat oss://bucket/object
       将object内容输出到标准输出
`,
	sampleText: ` 
    1) 将object内容输出到标准输出
       ossutil cat oss://bucket/object
`,
}

var specEnglishCat = SpecText{
	synopsisText: "Output object content to standard output",

	paramText: "object [options]",

	syntaxText: ` 
	ossutil cat oss://bucket/object 
`,
	detailHelpText: ` 
	The cat command can output the object content of oss to standard output
    The object content is preferably text format

Usage:
    There is only one usage for this command:
	
    1) ossutil cat oss://bucket/object
       The command output object content to standard output
`,
	sampleText: ` 
    1) output object content to standard output
       ossutil cat oss://bucket/object
`,
}

type catOptionType struct {
	bucketName   string
	objectName   string
	encodingType string
}

type CatCommand struct {
	command   Command
	catOption catOptionType
}

var catCommand = CatCommand{
	command: Command{
		name:        "cat",
		nameAlias:   []string{"cat"},
		minArgc:     1,
		maxArgc:     1,
		specChinese: specChineseCat,
		specEnglish: specEnglishCat,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionEncodingType,
			OptionLogLevel,
		},
	},
}

// function for FormatHelper interface
func (catc *CatCommand) formatHelpForWhole() string {
	return catc.command.formatHelpForWhole()
}

func (catc *CatCommand) formatIndependHelp() string {
	return catc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (catc *CatCommand) Init(args []string, options OptionMapType) error {
	return catc.command.Init(args, options, catc)
}

// RunCommand simulate inheritance, and polymorphism
func (catc *CatCommand) RunCommand() error {
	catc.catOption.encodingType, _ = GetString(OptionEncodingType, catc.command.options)
	srcBucketUrL, err := GetCloudUrl(catc.command.args[0], catc.catOption.encodingType)
	if err != nil {
		return err
	}

	if srcBucketUrL.object == "" {
		return fmt.Errorf("object key is empty")
	}

	catc.catOption.bucketName = srcBucketUrL.bucket
	catc.catOption.objectName = srcBucketUrL.object

	// check object exist or not
	client, err := catc.command.ossClient(catc.catOption.bucketName)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(catc.catOption.bucketName)
	if err != nil {
		return err
	}

	isExist, err := bucket.IsObjectExist(catc.catOption.objectName)
	if err != nil {
		return err
	}

	if !isExist {
		return fmt.Errorf("oss object is not exist")
	}

	body, err := bucket.GetObject(catc.catOption.objectName)
	if err != nil {
		return err
	}

	defer body.Close()
	io.Copy(os.Stdout, body)
	fmt.Printf("\n")

	return err
}
