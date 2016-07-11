package lib

import (
	"fmt"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
	"strconv"
	"strings"
)

// group spec text of all commands
const (
	GroupTypeNormalCommand string = "\nCommands:\n"

	GroupTypeAdditionalCommand string = "\nAdditional Commands:\n"

	GroupTypeDeprecatedCommand string = "\nDeprecated Commands:\n"
)

// CommandGroups is the array of all group types
var CommandGroups = []string{
	GroupTypeNormalCommand,
	GroupTypeAdditionalCommand,
	GroupTypeDeprecatedCommand,
}

// SpecText is the spec text of a command
type SpecText struct {
	synopsisText   string
	paramText      string
	syntaxText     string
	detailHelpText string
	sampleText     string
}

// Command contains all elements of a command, it's the base class of all commands
type Command struct {
	name             string
	nameAlias        []string
	minArgc          int
	maxArgc          int
	specChinese      SpecText
	specEnglish      SpecText
	validOptionNames []string
	group            string
	args             []string
	options          OptionMapType
	configOptions    OptionMapType
}

// Commander is the interface of all commands
type Commander interface {
	RunCommand() error
	Init(args []string, options OptionMapType) error
}

// RewriteLoadConfiger is the interface for those commands, which do not need to load config, or have other action
type RewriteLoadConfiger interface {
	rewriteLoadConfig(string) error
}

// RewriteAssembleOptioner is the interface for those commands, which do not need to assemble options
type RewriteAssembleOptioner interface {
	rewriteAssembleOptions()
}

// Init is the common functions for all commands, they use Init to initialize itself
func (cmd *Command) Init(args []string, options OptionMapType, cmder interface{}) error {
	cmd.args = args
	cmd.options = options
	cmd.configOptions = OptionMapType{}

	if err := cmd.checkArgs(); err != nil {
		return err
	}

	val, _ := GetString(OptionConfigFile, cmd.options)
	if err := cmd.loadConfig(val, cmder); err != nil {
		return err
	}

	if err := cmd.checkOptions(); err != nil {
		return err
	}

	cmd.assembleOptions(cmder)
	return nil
}

func (cmd *Command) checkArgs() error {
	str := ""
	if cmd.minArgc > 1 {
		str = "s"
	}
	if len(cmd.args) < cmd.minArgc {
		msg := fmt.Sprintf("the command needs at least %d argument%s", cmd.minArgc, str)
		return CommandError{cmd.name, msg}
	}

	str = ""
	if cmd.minArgc > 1 {
		str = "s"
	}

	if len(cmd.args) > cmd.maxArgc {
		msg := fmt.Sprintf("the command needs at most %d argument%s", cmd.maxArgc, str)
		return CommandError{cmd.name, msg}
	}
	return nil
}

func (cmd *Command) loadConfig(configFile string, cmder interface{}) error {
	if cmdder, ok := cmder.(RewriteLoadConfiger); ok {
		return cmdder.rewriteLoadConfig(configFile)
	}
	var err error
	if cmd.configOptions, err = LoadConfig(configFile); err != nil {
		return err
	}
	return nil
}

func (cmd *Command) checkOptions() error {
	for name := range cmd.options {
		msg := fmt.Sprintf("the command does not support option: \"%s\"", name)
		switch OptionMap[name].optionType {
		case OptionTypeFlagTrue:
			if val, _ := GetBool(name, cmd.options); val {
				if FindPos(name, cmd.validOptionNames) == -1 {
					return CommandError{cmd.name, msg}
				}
			}
		default:
			if val, _ := GetString(name, cmd.options); val != "" {
				if FindPos(name, cmd.validOptionNames) == -1 {
					return CommandError{cmd.name, msg}
				}
			}
		}
	}
	return nil
}

func (cmd *Command) assembleOptions(cmder interface{}) {
	if cmdder, ok := cmder.(RewriteAssembleOptioner); ok {
		cmdder.rewriteAssembleOptions()
		return
	}

	for name, option := range cmd.configOptions {
		if _, ok := cmd.options[name]; ok {
			if OptionMap[name].optionType != OptionTypeFlagTrue {
				if val, _ := GetString(name, cmd.options); val == "" {
					opval := option.(string)
					cmd.options[name] = &opval
					delete(cmd.configOptions, name)
				}
			}
		}
	}

	for name := range cmd.options {
		if OptionMap[name].def != "" {
			switch OptionMap[name].optionType {
			case OptionTypeInt64:
				if val, _ := GetString(name, cmd.options); val == "" {
					def, _ := strconv.ParseInt(OptionMap[name].def, 10, 64)
					cmd.options[name] = &def
				}
			case OptionTypeString:
				if val, _ := GetString(name, cmd.options); val == "" {
					def := OptionMap[name].def
					cmd.options[name] = &def
				}
			}
		}
	}
}

// FormatHelper is the interface for all commands to format spec information
type FormatHelper interface {
	formatHelpForWhole() string
	formatIndependHelp() string
}

func (cmd *Command) formatHelpForWhole() string {
	formatStr := "  %-" + strconv.Itoa(MaxCommandNameLen) + "s%s\n%s%s%s\n"
	spec := cmd.getSpecText()
	return fmt.Sprintf(formatStr, cmd.name, spec.paramText, FormatTAB, FormatTAB, spec.synopsisText)
}

func (cmd *Command) getSpecText() SpecText {
	switch Language {
	case "中文":
		return cmd.specChinese
	default:
		return cmd.specEnglish
	}
}

func (cmd *Command) formatIndependHelp() string {
	spec := cmd.getSpecText()
	text := fmt.Sprintf("SYNOPSIS\n\n%s%s\n", FormatTAB, strings.TrimSpace(spec.synopsisText))
	if spec.syntaxText != "" {
		text += fmt.Sprintf("\nSYNTAX\n\n%s%s\n", FormatTAB, strings.TrimSpace(spec.syntaxText))
	}
	if spec.detailHelpText != "" {
		text += fmt.Sprintf("\nDETAIL DESCRIPTION\n\n%s%s\n", FormatTAB, strings.TrimSpace(spec.detailHelpText))
	}
	if spec.sampleText != "" {
		text += fmt.Sprintf("\nSAMPLE\n\n%s%s\n", FormatTAB, strings.TrimSpace(spec.sampleText))
	}
	if len(cmd.validOptionNames) != 0 {
		text += fmt.Sprintf("\nOPTIONS\n\n%s\n", cmd.formatOptionsHelp(cmd.validOptionNames))
	}
	return text
}

func (cmd *Command) formatOptionsHelp(validOptionNames []string) string {
	text := ""
	for _, optionName := range validOptionNames {
		if option, ok := OptionMap[optionName]; ok {
			text += cmd.formatOption(option)
		}
	}
	return text
}

func (cmd *Command) formatOption(option Option) string {
	text := FormatTAB
	if option.name != "" {
		text += option.name
		if option.def != "" {
			text += " " + option.def
		}
	}

	if option.name != "" && option.nameAlias != "" {
		text += ", "
	}

	if option.nameAlias != "" {
		text += option.nameAlias
		if option.def != "" {
			text += fmt.Sprintf("=%s", option.def)
		}
	}

	if option.getHelp() != "" {
		text += fmt.Sprintf("\n%s%s%s\n\n", FormatTAB, FormatTAB, option.getHelp())
	}

	return text
}

// OSS common function
// get oss client according to bucket(if bucket not empty)
func (cmd *Command) ossClient(bucket string) (*oss.Client, error) {
	endpoint, isCname := cmd.getEndpoint(bucket)
	accessKeyID, _ := GetString(OptionAccessKeyID, cmd.options)
	accessKeySecret, _ := GetString(OptionAccessKeySecret, cmd.options)
	stsToken, _ := GetString(OptionSTSToken, cmd.options)
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret, oss.UseCname(isCname), oss.SecurityToken(stsToken))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (cmd *Command) getEndpoint(bucket string) (string, bool) {
	if cnameMap, ok := cmd.configOptions[BucketCnameSection]; ok {
		if endpoint, ok := cnameMap.(map[string]string)[bucket]; ok {
			return endpoint, true
		}
	}

	if eMap, ok := cmd.configOptions[BucketEndpointSection]; ok {
		if endpoint, ok := eMap.(map[string]string)[bucket]; ok {
			return endpoint, false
		}
	}

	endpoint, _ := GetString(OptionEndpoint, cmd.options)
	return endpoint, false
}

// get oss operable bucket
func (cmd *Command) ossBucket(bucketName string) (*oss.Bucket, error) {
	client, err := cmd.ossClient(bucketName)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func (cmd *Command) ossListObjectsRetry(bucket *oss.Bucket, options ...oss.Option) (oss.ListObjectsResult, error) {
	retryTimes, _ := GetInt(OptionRetryTimes, cmd.options)
	for i := 1; ; i++ {
		lor, err := bucket.ListObjects(options...)
		if err == nil || int64(i) >= retryTimes {
			return lor, err
		}
	}
}

func (cmd *Command) ossGetObjectStatRetry(bucket *oss.Bucket, object string) (http.Header, error) {
	retryTimes, _ := GetInt(OptionRetryTimes, cmd.options)
	for i := 1; ; i++ {
		props, err := bucket.GetObjectDetailedMeta(object)
		if err == nil {
			return props, err
		}
		if int64(i) >= retryTimes {
			return props, ObjectError{err, object}
		}
	}
}

func (cmd *Command) objectProducer(bucket *oss.Bucket, cloudURL CloudURL, chObjects chan<- string, chError chan<- error) {
	pre := oss.Prefix(cloudURL.object)
	marker := oss.Marker("")
	for i := 0; ; i++ {
		lor, err := cmd.ossListObjectsRetry(bucket, marker, pre)
		if err != nil {
			chError <- err
			break
		}

		for _, object := range lor.Objects {
			chObjects <- object.Key
		}

		pre = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
	defer close(chObjects)
	chError <- nil
}

// GetAllCommands returns all commands list
func GetAllCommands() []interface{} {
	return []interface{}{
		&helpCommand,
		&configCommand,
        &makeBucketCommand,
		&listCommand,
		&removeCommand,
		&statCommand,
		&setACLCommand,
		&setMetaCommand,
		&copyCommand,
	}
}
