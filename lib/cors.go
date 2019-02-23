package lib

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseCors = SpecText{}

var specEnglishCors = SpecText{}

type corsOptionType struct {
	bucketName string
}

type CorsCommand struct {
	command  Command
	csOption corsOptionType
}

var corsCommand = CorsCommand{
	command: Command{
		name:        "cors",
		nameAlias:   []string{"cors"},
		minArgc:     2,
		maxArgc:     3,
		specChinese: specChineseCors,
		specEnglish: specEnglishCors,
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
func (corsc *CorsCommand) formatHelpForWhole() string {
	return corsc.command.formatHelpForWhole()
}

func (corsc *CorsCommand) formatIndependHelp() string {
	return corsc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (corsc *CorsCommand) Init(args []string, options OptionMapType) error {
	return corsc.command.Init(args, options, corsc)
}

// RunCommand simulate inheritance, and polymorphism
func (corsc *CorsCommand) RunCommand() error {
	strMethod := corsc.command.args[0]
	if strMethod != "put" && strMethod != "get" && strMethod != "delete" {
		return fmt.Errorf("%s is not in the optional value:put|get|delete", strMethod)
	}

	bucketUrL, err := StorageURLFromString(corsc.command.args[1], "")
	if err != nil {
		return err
	}

	if !bucketUrL.IsCloudURL() {
		return fmt.Errorf("the second parameter is not a cloud url,url is %s", bucketUrL.ToString())
	}

	cloudUrl := bucketUrL.(CloudURL)
	if cloudUrl.bucket == "" {
		return fmt.Errorf("bucket name is empty,url is %s", bucketUrL.ToString())
	}

	corsc.csOption.bucketName = cloudUrl.bucket

	if strMethod == "put" {
		err = corsc.PutBucketCors()
	} else if strMethod == "get" {
		err = corsc.GetBucketCors()
	} else if strMethod == "delete" {
		err = corsc.DeleteBucketCors()
	}
	return err
}

func (corsc *CorsCommand) PutBucketCors() error {
	if len(corsc.command.args) < 3 {
		return fmt.Errorf("put bucket cors need 3 parameters,the cors config file is empty")
	}

	corsFile := corsc.command.args[2]
	fileInfo, err := os.Stat(corsFile)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s is dir,not the expected file", corsFile)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("%s is empty file", corsFile)
	}

	// parsing the xml file
	file, err := os.Open(corsFile)
	if err != nil {
		return err
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	rulesConfig := oss.CORSXML{}
	err = xml.Unmarshal(text, &rulesConfig)
	if err != nil {
		return err
	}

	// put bucket cors
	client, err := corsc.command.ossClient(corsc.csOption.bucketName)
	if err != nil {
		return err
	}

	return client.SetBucketCORS(corsc.csOption.bucketName, rulesConfig.CORSRules)
}

func (corsc *CorsCommand) confirm(str string) bool {
	var val string
	fmt.Printf(getClearStr(fmt.Sprintf("cors: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
		return false
	}
	return true
}

func (corsc *CorsCommand) GetBucketCors() error {
	client, err := corsc.command.ossClient(corsc.csOption.bucketName)
	if err != nil {
		return err
	}

	corsRes, err := client.GetBucketCORS(corsc.csOption.bucketName)
	if err != nil {
		return err
	}

	output, err := xml.MarshalIndent(corsRes, "  ", "    ")
	if err != nil {
		return err
	}

	var outFile *os.File
	if len(corsc.command.args) >= 3 {
		fileName := corsc.command.args[2]
		_, err = os.Stat(fileName)
		if err == nil {
			bConitnue := corsc.confirm(fileName)
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

func (corsc *CorsCommand) DeleteBucketCors() error {
	client, err := corsc.command.ossClient(corsc.csOption.bucketName)
	if err != nil {
		return err
	}
	return client.DeleteBucketCORS(corsc.csOption.bucketName)
}
