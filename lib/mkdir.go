package lib

import (
	"fmt"
	"strings"
)

var specChineseMkdir = SpecText{}

var specEnglishMkdir = SpecText{}

type mkOptionType struct {
	encodingType string
}

type MkdirCommand struct {
	command  Command
	mkOption mkOptionType
}

var mkdirCommand = MkdirCommand{
	command: Command{
		name:        "mkdir",
		nameAlias:   []string{"mkdir"},
		minArgc:     1,
		maxArgc:     1,
		specChinese: specChineseMkdir,
		specEnglish: specEnglishMkdir,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionLogLevel,
			OptionEncodingType,
		},
	},
}

// function for FormatHelper interface
func (mkc *MkdirCommand) formatHelpForWhole() string {
	return mkc.command.formatHelpForWhole()
}

func (mkc *MkdirCommand) formatIndependHelp() string {
	return mkc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (mkc *MkdirCommand) Init(args []string, options OptionMapType) error {
	return mkc.command.Init(args, options, mkc)
}

// RunCommand simulate inheritance, and polymorphism
func (mkc *MkdirCommand) RunCommand() error {
	mkc.mkOption.encodingType, _ = GetString(OptionEncodingType, mkc.command.options)

	dirUrL, err := StorageURLFromString(mkc.command.args[0], mkc.mkOption.encodingType)
	if err != nil {
		return err
	}

	if !dirUrL.IsCloudURL() {
		return fmt.Errorf("parameter is not a cloud url,url is %s", dirUrL.ToString())
	}

	cloudUrl := dirUrL.(CloudURL)

	if cloudUrl.bucket == "" {
		return fmt.Errorf("bucket name is empty,url is %s", cloudUrl.ToString())
	}

	if cloudUrl.object == "" {
		return fmt.Errorf("object name is empty,url is %s", cloudUrl.ToString())
	}

	if !strings.HasSuffix(cloudUrl.object, "/") {
		cloudUrl.object += "/"
	}

	return mkc.MkBucketDir(cloudUrl)
}

func (mkc *MkdirCommand) MkBucketDir(dirUrl CloudURL) error {
	bucket, err := mkc.command.ossBucket(dirUrl.bucket)
	if err != nil {
		return err
	}

    bExist, err := bucket.IsObjectExist(dirUrl.object)
    if err != nil {
        return err
    }
    
	if bExist {
		return fmt.Errorf("%s already exists", dirUrl.object)
	}

	return bucket.PutObject(dirUrl.object, strings.NewReader(""))
}
