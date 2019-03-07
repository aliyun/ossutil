package lib

import (
	"fmt"
	"io"
	"os"
)

var specChineseCat = SpecText{}

var specEnglishCat = SpecText{}

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
	srcBucketUrL, err := catc.CheckBucketUrl(catc.command.args[0], catc.catOption.encodingType)
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

func (catc *CatCommand) CheckBucketUrl(strlUrl, encodingType string) (*CloudURL, error) {
	bucketUrL, err := StorageURLFromString(strlUrl, encodingType)
	if err != nil {
		return nil, err
	}

	if !bucketUrL.IsCloudURL() {
		return nil, fmt.Errorf("parameter is not a cloud url,url is %s", bucketUrL.ToString())
	}

	cloudUrl := bucketUrL.(CloudURL)
	if cloudUrl.bucket == "" {
		return nil, fmt.Errorf("bucket name is empty,url is %s", bucketUrL.ToString())
	}
	return &cloudUrl, nil
}
