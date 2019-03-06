package lib

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

var specChineseBucketRefer = SpecText{}

var specEnglishBucketRefer = SpecText{}

type bucketReferOptionType struct {
	bucketName        string
	disableEmptyRefer bool
}

type BucketRefererCommand struct {
	command  Command
	brOption bucketReferOptionType
}

var bucketRefererCommand = BucketRefererCommand{
	command: Command{
		name:        "referer",
		nameAlias:   []string{"referer"},
		minArgc:     2,
		maxArgc:     MaxInt,
		specChinese: specChineseBucketRefer,
		specEnglish: specEnglishBucketRefer,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionLogLevel,
			OptionDisableEmptyReferer,
		},
	},
}

// function for FormatHelper interface
func (brc *BucketRefererCommand) formatHelpForWhole() string {
	return brc.command.formatHelpForWhole()
}

func (brc *BucketRefererCommand) formatIndependHelp() string {
	return brc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (brc *BucketRefererCommand) Init(args []string, options OptionMapType) error {
	return brc.command.Init(args, options, brc)
}

// RunCommand simulate inheritance, and polymorphism
func (brc *BucketRefererCommand) RunCommand() error {
	strMethod := brc.command.args[0]
	if strMethod != "put" && strMethod != "get" && strMethod != "delete" {
		return fmt.Errorf("%s is not in the optional value:put|get|delete", strMethod)
	}

	srcBucketUrL, err := brc.CheckBucketUrl(brc.command.args[1])
	if err != nil {
		return err
	}

	brc.brOption.bucketName = srcBucketUrL.bucket
	brc.brOption.disableEmptyRefer, _ = GetBool(OptionDisableEmptyReferer, brc.command.options)

	if strMethod == "put" {
		err = brc.PutBucketRefer()
	} else if strMethod == "get" {
		err = brc.GetBucketRefer()
	} else if strMethod == "delete" {
		err = brc.DeleteBucketRefer()
	}
	return err
}

func (brc *BucketRefererCommand) CheckBucketUrl(strlUrl string) (*CloudURL, error) {
	bucketUrL, err := StorageURLFromString(strlUrl, "")
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

func (brc *BucketRefererCommand) PutBucketRefer() error {
	if len(brc.command.args) < 3 {
		return fmt.Errorf("put bucket referer need at least 3 parameters,the refer is empty")
	}

	referers := brc.command.args[2:len(brc.command.args)]

	// put bucket refer
	client, err := brc.command.ossClient(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	return client.SetBucketReferer(brc.brOption.bucketName, referers, !brc.brOption.disableEmptyRefer)
}

func (brc *BucketRefererCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("bucket referer: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (brc *BucketRefererCommand) GetBucketRefer() error {
	client, err := brc.command.ossClient(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	referRes, err := client.GetBucketReferer(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	output, err := xml.MarshalIndent(referRes, "  ", "    ")
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(brc.command.args) >= 3 {
		fileName := brc.command.args[2]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := brc.confirm(fileName)
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

func (brc *BucketRefererCommand) DeleteBucketRefer() error {

	referers := []string{}

	// put bucket refer
	client, err := brc.command.ossClient(brc.brOption.bucketName)
	if err != nil {
		return err
	}

	return client.SetBucketReferer(brc.brOption.bucketName, referers, true)
}
