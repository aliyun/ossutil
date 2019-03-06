package lib

import (
	"fmt"
	"strconv"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseListPart = SpecText{}

var specEnglishListPart = SpecText{}

type listPartOptionType struct {
	cloudUrl     CloudURL
	uploadId     string
	encodingType string
}

type ListPartCommand struct {
	command  Command
	lpOption listPartOptionType
}

var listPartCommand = ListPartCommand{
	command: Command{
		name:        "listpart",
		nameAlias:   []string{"listpart"},
		minArgc:     2,
		maxArgc:     2,
		specChinese: specChineseListPart,
		specEnglish: specEnglishListPart,
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
func (lpc *ListPartCommand) formatHelpForWhole() string {
	return lpc.command.formatHelpForWhole()
}

func (lpc *ListPartCommand) formatIndependHelp() string {
	return lpc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (lpc *ListPartCommand) Init(args []string, options OptionMapType) error {
	return lpc.command.Init(args, options, lpc)
}

// RunCommand simulate inheritance, and polymorphism
func (lpc *ListPartCommand) RunCommand() error {
	lpc.lpOption.encodingType, _ = GetString(OptionEncodingType, lpc.command.options)
	srcBucketUrL, err := lpc.CheckBucketUrl(lpc.command.args[0], lpc.lpOption.encodingType)
	if err != nil {
		return err
	}

	if srcBucketUrL.object == "" {
		return fmt.Errorf("object name is empty")
	}

	lpc.lpOption.cloudUrl = *srcBucketUrL
	lpc.lpOption.uploadId = lpc.command.args[1]

	return lpc.ListPart()

}

func (lpc *ListPartCommand) CheckBucketUrl(strlUrl, encodingType string) (*CloudURL, error) {
	bucketUrL, err := StorageURLFromString(strlUrl, encodingType)
	if err != nil {
		return nil, err
	}

	if !bucketUrL.IsCloudURL() {
		return nil, fmt.Errorf("the first parameter is not a cloud url,url is %s", bucketUrL.ToString())
	}

	cloudUrl := bucketUrL.(CloudURL)
	if cloudUrl.bucket == "" {
		return nil, fmt.Errorf("bucket name is empty,url is %s", bucketUrL.ToString())
	}
	return &cloudUrl, nil
}

func (lpc *ListPartCommand) ListPart() error {
	client, err := lpc.command.ossClient(lpc.lpOption.cloudUrl.bucket)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(lpc.lpOption.cloudUrl.bucket)
	if err != nil {
		return err
	}

	var imur oss.InitiateMultipartUploadResult
	imur.Bucket = lpc.lpOption.cloudUrl.bucket
	imur.Key = lpc.lpOption.cloudUrl.object
	imur.UploadID = lpc.lpOption.uploadId

	partNumberMarker := 0
	totalPartCount := 0
	var totalPartSize int64 = 0
	for i := 0; ; i++ {
		lpOptions := []oss.Option{}
		lpOptions = append(lpOptions, oss.MaxParts(1000))
		lpOptions = append(lpOptions, oss.PartNumberMarker(partNumberMarker))

		lpRes, err := bucket.ListUploadedParts(imur, lpOptions...)
		if err != nil {
			return err
		} else {
			totalPartCount += len(lpRes.UploadedParts)
			if i == 0 && len(lpRes.UploadedParts) > 0 {
				fmt.Printf("%-10s\t%-32s\t%-10s\t%s\n", "PartNumber", "Etag", "Size(Byte)", "LastModifyTime")
			}
		}

		for _, v := range lpRes.UploadedParts {
			//PartNumber,ETag,Size,LastModified
			fmt.Printf("%-10d\t%-32s\t%-10d\t%s\n", v.PartNumber, v.ETag, v.Size, v.LastModified.Format("2006-01-02 15:04:05"))
			totalPartSize += int64(v.Size)
		}

		if lpRes.IsTruncated {
			partNumberMarker, err = strconv.Atoi(lpRes.NextPartNumberMarker)
			if err != nil {
				return err
			}
		} else {
			if totalPartCount > 0 {
				fmt.Printf("\ntotal part count:%d\ttotal part size(MB):%.2f\n\n", totalPartCount, float64(totalPartSize/1024)/1024)
			}
			break
		}
	}
	return nil
}
