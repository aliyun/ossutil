package lib

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

var specChineseBucketLog = SpecText{}

var specEnglishBucketLog = SpecText{}

type bucketLogOptionType struct {
	srcBucketName  string
	destBucketName string
	destPrefix     string
}

type BucketLogCommand struct {
	command  Command
	blOption bucketLogOptionType
}

var bucketLogCommand = BucketLogCommand{
	command: Command{
		name:        "logging",
		nameAlias:   []string{"logging"},
		minArgc:     2,
		maxArgc:     3,
		specChinese: specChineseBucketLog,
		specEnglish: specEnglishBucketLog,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionLogLevel,
		},
	},
}

// function for FormatHelper interface
func (blc *BucketLogCommand) formatHelpForWhole() string {
	return blc.command.formatHelpForWhole()
}

func (blc *BucketLogCommand) formatIndependHelp() string {
	return blc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (blc *BucketLogCommand) Init(args []string, options OptionMapType) error {
	return blc.command.Init(args, options, blc)
}

// RunCommand simulate inheritance, and polymorphism
func (blc *BucketLogCommand) RunCommand() error {
	strMethod := blc.command.args[0]
	if strMethod != "put" && strMethod != "get" && strMethod != "delete" {
		return fmt.Errorf("%s is not in the optional value:put|get|delete", strMethod)
	}

	srcBucketUrL, err := blc.CheckBucketUrl(blc.command.args[1])
	if err != nil {
		return err
	}

	blc.blOption.srcBucketName = srcBucketUrL.bucket

	if strMethod == "put" {
		err = blc.PutBucketLog()
	} else if strMethod == "get" {
		err = blc.GetBucketLog()
	} else if strMethod == "delete" {
		err = blc.DeleteBucketLog()
	}
	return err
}

func (blc *BucketLogCommand) CheckBucketUrl(strlUrl string) (*CloudURL, error) {
	bucketUrL, err := StorageURLFromString(strlUrl, "")
	if err != nil {
		return nil, err
	}

	if !bucketUrL.IsCloudURL() {
		return nil, fmt.Errorf("the second parameter is not a cloud url,url is %s", bucketUrL.ToString())
	}

	cloudUrl := bucketUrL.(CloudURL)
	if cloudUrl.bucket == "" {
		return nil, fmt.Errorf("bucket name is empty,url is %s", bucketUrL.ToString())
	}
	return &cloudUrl, nil
}

func (blc *BucketLogCommand) PutBucketLog() error {
	if len(blc.command.args) < 3 {
		return fmt.Errorf("put bucket log need 3 parameters,the target bucket is empty")
	}

	destBucketUrL, err := blc.CheckBucketUrl(blc.command.args[2])
	if err != nil {
		return err
	}

	blc.blOption.destBucketName = destBucketUrL.bucket
	blc.blOption.destPrefix = destBucketUrL.object

	// put bucket log
	client, err := blc.command.ossClient(blc.blOption.srcBucketName)
	if err != nil {
		return err
	}

	return client.SetBucketLogging(blc.blOption.srcBucketName, blc.blOption.destBucketName, blc.blOption.destPrefix, true)
}

func (blc *BucketLogCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("bucket log: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (blc *BucketLogCommand) GetBucketLog() error {
	client, err := blc.command.ossClient(blc.blOption.srcBucketName)
	if err != nil {
		return err
	}

	logRes, err := client.GetBucketLogging(blc.blOption.srcBucketName)
	if err != nil {
		return err
	}

	output, err := xml.MarshalIndent(logRes, "  ", "    ")
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(blc.command.args) >= 3 {
		fileName := blc.command.args[2]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := blc.confirm(fileName)
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

func (blc *BucketLogCommand) DeleteBucketLog() error {
	client, err := blc.command.ossClient(blc.blOption.srcBucketName)
	if err != nil {
		return err
	}
	return client.DeleteBucketLogging(blc.blOption.srcBucketName)
}
