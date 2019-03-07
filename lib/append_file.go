package lib

import (
	"fmt"
	"os"
	"strconv"
	"time"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseAppendFile = SpecText{}

var specEnglishAppendFile = SpecText{}

type AppendProgressListener struct {
	lastMs   int64
	lastSize int64
	currSize int64
}

// ProgressChanged handle progress event
func (l *AppendProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	if event.EventType == oss.TransferDataEvent || event.EventType == oss.TransferCompletedEvent {
		if l.lastMs == 0 {
			l.lastSize = l.currSize
			l.currSize = event.ConsumedBytes
			l.lastMs = time.Now().UnixNano() / 1000 / 1000
		} else {
			now := time.Now()
			cost := now.UnixNano()/1000/1000 - l.lastMs
			if cost > 1000 || event.EventType == oss.TransferCompletedEvent {
				l.lastSize = l.currSize
				l.currSize = event.ConsumedBytes
				l.lastMs = now.UnixNano() / 1000 / 1000

				speed := float64(l.currSize-l.lastSize) / float64(cost)
				rate := float64(l.currSize) * 100 / float64(event.TotalBytes)
				fmt.Printf("\rtotal append %d(%.2f%%) byte,speed is %.2f(KB/s)", event.ConsumedBytes, rate, speed)
			}
		}
	}
}

type appendFileOptionType struct {
	bucketName   string
	objectName   string
	encodingType string
	fileName     string
	fileSize     int64
	ossMeta      string
}

type AppendFileCommand struct {
	command  Command
	afOption appendFileOptionType
}

var appendFileCommand = AppendFileCommand{
	command: Command{
		name:        "appendfromfile",
		nameAlias:   []string{"appendfromfile"},
		minArgc:     2,
		maxArgc:     2,
		specChinese: specChineseAppendFile,
		specEnglish: specEnglishAppendFile,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionConfigFile,
			OptionEndpoint,
			OptionAccessKeyID,
			OptionAccessKeySecret,
			OptionSTSToken,
			OptionEncodingType,
			OptionMeta,
			OptionMaxUpSpeed,
			OptionLogLevel,
		},
	},
}

// function for FormatHelper interface
func (afc *AppendFileCommand) formatHelpForWhole() string {
	return afc.command.formatHelpForWhole()
}

func (afc *AppendFileCommand) formatIndependHelp() string {
	return afc.command.formatIndependHelp()
}

// Init simulate inheritance, and polymorphism
func (afc *AppendFileCommand) Init(args []string, options OptionMapType) error {
	return afc.command.Init(args, options, afc)
}

// RunCommand simulate inheritance, and polymorphism
func (afc *AppendFileCommand) RunCommand() error {
	afc.afOption.encodingType, _ = GetString(OptionEncodingType, afc.command.options)
	afc.afOption.ossMeta, _ = GetString(OptionMeta, afc.command.options)

	srcBucketUrL, err := afc.CheckBucketUrl(afc.command.args[0], afc.afOption.encodingType)
	if err != nil {
		return err
	}

	if srcBucketUrL.object == "" {
		return fmt.Errorf("object key is empty")
	}

	afc.afOption.bucketName = srcBucketUrL.bucket
	afc.afOption.objectName = srcBucketUrL.object

	// check input file
	fileName := afc.command.args[1]
	stat, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("%s is dir", fileName)
	}

	if stat.Size() > MaxAppendObjectSize {
		return fmt.Errorf("locafile:%s is bigger than %d, it is not support by append", fileName, MaxAppendObjectSize)
	}

	afc.afOption.fileName = fileName
	afc.afOption.fileSize = stat.Size()

	// check object exist or not
	client, err := afc.command.ossClient(afc.afOption.bucketName)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(afc.afOption.bucketName)
	if err != nil {
		return err
	}

	isExist, err := bucket.IsObjectExist(afc.afOption.objectName)
	if err != nil {
		return err
	}

	if isExist && afc.afOption.ossMeta != "" {
		return fmt.Errorf("setting meta on existing append object is not supported")
	}

	position := int64(0)
	if isExist {
		//get object size
		props, err := bucket.GetObjectMeta(afc.afOption.objectName)
		if err != nil {
			return err
		}

		position, err = strconv.ParseInt(props.Get("Content-Length"), 10, 64)
		if err != nil {
			return err
		}
	}

	err = afc.AppendFromFile(bucket, position)

	return err
}

func (afc *AppendFileCommand) CheckBucketUrl(strlUrl, encodingType string) (*CloudURL, error) {
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

func (afc *AppendFileCommand) AppendFromFile(bucket *oss.Bucket, position int64) error {
	file, err := os.OpenFile(afc.afOption.fileName, os.O_RDONLY, 0660)
	if err != nil {
		return err
	}
	defer file.Close()

	var options []oss.Option
	if afc.afOption.ossMeta != "" {
		metas, err := afc.command.parseHeaders(afc.afOption.ossMeta, false)
		if err != nil {
			return err
		}

		options, err = afc.command.getOSSOptions(headerOptionMap, metas)
		if err != nil {
			return err
		}
	}

	var listener *AppendProgressListener = &AppendProgressListener{}
	options = append(options, oss.Progress(listener))

	startT := time.Now()
	newPosition, err := bucket.AppendObject(afc.afOption.objectName, file, position, options...)
	endT := time.Now()
	if err != nil {
		return err
	} else {
		cost := endT.UnixNano()/1000/1000 - startT.UnixNano()/1000/1000
		speed := float64(afc.afOption.fileSize) / float64(cost)
		fmt.Printf("\nlocal file size is %d,the object new size is %d,average speed is %.2f(KB/s)\n\n", afc.afOption.fileSize, newPosition, speed)
		return nil
	}
}
