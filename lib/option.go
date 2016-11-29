package lib

import (
	"errors"
	"fmt"
    "strings"
	"strconv"
	goopt "github.com/droundy/goopt"
)

type optionType int

// option types, only support three kinds now
const (
	OptionTypeString optionType = iota
	OptionTypeInt64
	OptionTypeFlagTrue
    OptionTypeAlternative
)

// Option describe the component of a option
type Option struct {
	name        string
	nameAlias   string
	def         string
	optionType  optionType
	minVal      string // empty means no check, for OptionTypeAlternative, minVal is the alternative values connected by '/', eg: CN|EN
	maxVal      string // empty means no check, for OptionTypeAlternative, maxVal is empty
	helpChinese string
	helpEnglish string
}

// LEnglishLanguage is the lower case of EnglishLanguage
var LEnglishLanguage = strings.ToLower(EnglishLanguage)

// OptionMap is a collection of ossutil supported options
var OptionMap = map[string]Option{
	OptionConfigFile: Option{"-c", "--config-file", "", OptionTypeString, "", "",
		"ossutil工具的配置文件路径，ossutil启动时从配置文件读取配置，在config命令中，ossutil将配置写入该文件。",
		"Path of ossutil configuration file, where to dump config in config command, or to load config in other commands that need credentials."},
	OptionEndpoint: Option{"-e", "--endpoint", "", OptionTypeString, "", "",
		fmt.Sprintf("ossutil工具的基本endpoint配置（该选项值会覆盖配置文件中的相应设置），注意其必须为一个二级域名。"),
		fmt.Sprintf("Base endpoint for oss endpoint(Notice that the value of the option will cover the value in config file). Take notice that it should be second-level domain(SLD).")},
	OptionAccessKeyID:      Option{"-i", "--access-key-id", "", OptionTypeString, "", "", "访问oss使用的AccessKeyID（该选项值会覆盖配置文件中的相应设置）。", "AccessKeyID while access oss(Notice that the value of the option will cover the value in config file)."},
	OptionAccessKeySecret:  Option{"-k", "--access-key-secret", "", OptionTypeString, "", "", "访问oss使用的AccessKeySecret（该选项值会覆盖配置文件中的相应设置）。", "AccessKeySecret while access oss(Notice that the value of the option will cover the value in config file)."},
	OptionSTSToken:         Option{"-t", "--sts-token", "", OptionTypeString, "", "", "访问oss使用的STSToken（该选项值会覆盖配置文件中的相应设置），非必须设置项。", "STSToken while access oss(Notice that the value of the option will cover the value in config file), not necessary."},
	OptionACL:              Option{"", "--acl", "", OptionTypeString, "", "", "acl信息的配置。", "acl information."},
	OptionShortFormat:      Option{"-s", "--short-format", "", OptionTypeFlagTrue, "", "", "显示精简格式，如果未指定该选项，默认显示长格式。", "Show by short format, if the option is not specified, show long format by default."},
	OptionDirectory:        Option{"-d", "--directory", "", OptionTypeFlagTrue, "", "", "返回当前目录下的文件和子目录，而非递归显示所有子目录下的所有object", "Return matching subdirectory names instead of contents of the subdirectory"},
	OptionRecursion:        Option{"-r", "--recursive", "", OptionTypeFlagTrue, "", "", "递归进行操作。对于支持该选项的命令，当指定该选项时，命令会对bucket下所有符合条件的objects进行操作，否则只对url中指定的单个object进行操作。", "operate recursively, for those commands which support the option, when use them, if the option is specified, the command will operate on all match objects under the bucket, else we will search the specified object and operate on the single object."},
	OptionBucket:           Option{"-b", "--bucket", "", OptionTypeFlagTrue, "", "", "对bucket进行操作，该选项用于确认操作作用于bucket", "the option used to make sure the operation will operate on bucket"},
	OptionForce:            Option{"-f", "--force", "", OptionTypeFlagTrue, "", "", "强制操作，不进行询问提示。", "operate silently without asking user to confirm the operation."},
	OptionUpdate:           Option{"", "--update", "", OptionTypeFlagTrue, "", "", "更新操作", "update"},
	OptionDelete:           Option{"", "--delete", "", OptionTypeFlagTrue, "", "", "删除操作", "delete"},
	OptionBigFileThreshold: Option{"", "--bigfile-threshold", strconv.Itoa(BigFileThreshold), OptionTypeInt64, strconv.FormatInt(MinBigFileThreshold, 10), strconv.FormatInt(MaxBigFileThreshold, 10), fmt.Sprintf("开启大文件断点续传的文件大小阀值，默认值:%dM，取值范围：%d-%d", BigFileThreshold/(1024*0124), MinBigFileThreshold, MaxBigFileThreshold), fmt.Sprintf("the threshold of file size, the file size larger than the threshold will use resume upload or download(default: %d), value range is: %d-%d", BigFileThreshold, MinBigFileThreshold, MaxBigFileThreshold)},
	OptionCheckpointDir:    Option{"", "--checkpoint-dir", CheckpointDir, OptionTypeString, "", "",
		fmt.Sprintf("checkpoint目录的路径(默认值为:%s)，断点续传时，操作失败ossutil会自动创建该目录，并在该目录下记录checkpoint信息，操作成功会删除该目录。如果指定了该选项，请确保所指定的目录可以被删除。", CheckpointDir),
		fmt.Sprintf("Path of checkpoint directory(default:%s), the directory is used in resume upload or download, when operate failed, ossutil will create the directory automatically, and record the checkpoint information in the directory, when the operation is succeed, the directory will be removed, so when specify the option, please make sure the directory can be removed.", CheckpointDir)},
	OptionRetryTimes:       Option{"", "--retry-times", strconv.Itoa(RetryTimes), OptionTypeInt64, strconv.FormatInt(MinRetryTimes, 10), strconv.FormatInt(MaxRetryTimes, 10), fmt.Sprintf("当错误发生时的重试次数，默认值：%d，取值范围：%d-%d", RetryTimes, MinRetryTimes, MaxRetryTimes), fmt.Sprintf("retry times when fail(default: %d), value range is: %d-%d", RetryTimes, MinRetryTimes, MaxRetryTimes)},
	OptionRoutines:         Option{"-j", "--jobs", strconv.Itoa(Routines), OptionTypeInt64, strconv.FormatInt(MinRoutines, 10), strconv.FormatInt(MaxRoutines, 10), fmt.Sprintf("多文件操作时的并发任务数，默认值：%d，取值范围：%d-%d", Routines, MinRoutines, MaxRoutines), fmt.Sprintf("amount of concurrency tasks between multi-files(default: %d), value range is: %d-%d", Routines, MinRoutines, MaxRoutines)},
	OptionParallel:         Option{"", "--parallel", "", OptionTypeInt64, strconv.FormatInt(MinParallel, 10), strconv.FormatInt(MaxParallel, 10), fmt.Sprintf("单文件内部操作的并发任务数，取值范围：%d-%d, 默认将由ossutil根据操作类型和文件大小自行决定。", MinRoutines, MaxRoutines), fmt.Sprintf("amount of concurrency tasks when work with a file, value range is: %d-%d, by default the value will be decided by ossutil intelligently.", MinRoutines, MaxRoutines)},
    OptionLanguage:         Option{"-L", "--language", DefaultLanguage, OptionTypeAlternative, fmt.Sprintf("%s/%s", DefaultLanguage, EnglishLanguage), "", fmt.Sprintf("设置ossutil工具的语言，默认值：%s，取值范围：%s/%s", DefaultLanguage, DefaultLanguage, EnglishLanguage), fmt.Sprintf("set the language of ossutil(default: %s), value range is: %s/%s", DefaultLanguage, DefaultLanguage, EnglishLanguage)}, 
    OptionHashType:         Option{"", "--type", DefaultHashType, OptionTypeAlternative, fmt.Sprintf("%s/%s", DefaultHashType, MD5HashType), "", fmt.Sprintf("计算的类型, 默认值：%s, 取值范围: %s/%s", DefaultHashType, DefaultHashType, MD5HashType),
        fmt.Sprintf("hash type, Default: %s, value range is: %s/%s", DefaultHashType, DefaultHashType, MD5HashType)},
	OptionVersion:          Option{"-v", "--version", "", OptionTypeFlagTrue, "", "", fmt.Sprintf("显示ossutil的版本（%s）并退出。", Version), fmt.Sprintf("Show ossutil version (%s) and exit.", Version)},
}

func (T *Option) getHelp(language string) string {
	switch strings.ToLower(language) {
	case LEnglishLanguage:
		return T.helpEnglish
    default:
		return T.helpChinese
	}
}

// OptionMapType is the type for ossutil got options
type OptionMapType map[string]interface{}

// ParseArgOptions parse command line and returns args and options
func ParseArgOptions() ([]string, OptionMapType, error) {
	options := initOption()
    goopt.Args = make([]string, 0, 4)
	goopt.Description = func() string {
		return "Simple tool for access OSS."
	}
	goopt.Parse(nil)
	if err := checkOption(options); err != nil {
		return nil, nil, err
	}
	return goopt.Args, options, nil
}

func initOption() OptionMapType {
	m := make(OptionMapType, len(OptionMap))
	for name, option := range OptionMap {
		switch option.optionType {
		case OptionTypeInt64:
			val, _ := stringOption(option)
			m[name] = val
		case OptionTypeFlagTrue:
			val, _ := flagTrueOption(option)
			m[name] = val
        case OptionTypeAlternative:
            val, _ := stringOption(option) 
            m[name] = val
		default:
			val, _ := stringOption(option)
			m[name] = val
		}
	}
	return m
}

func stringOption(option Option) (*string, error) {
	names, err := makeNames(option)
    if err == nil {
		// ignore option.def, set it to "", will assemble it after
		return goopt.String(names, "", option.getHelp(DefaultLanguage)), nil
	}
	return nil, err
}

func flagTrueOption(option Option) (*bool, error) {
	names, err := makeNames(option)
    if err == nil {
		return goopt.Flag(names, []string{}, option.getHelp(DefaultLanguage), ""), nil
	}
	return nil, err
}

func makeNames(option Option) ([]string, error) {
	if option.name == "" && option.nameAlias == "" {
		return nil, errors.New("Internal Error, invalid option whose name and nameAlias empty!")
	}

	var names []string
	if option.name == "" || option.nameAlias == "" {
		names = make([]string, 1)
		if option.name == "" {
			names[0] = option.nameAlias
		} else {
			names[0] = option.name
		}
	} else {
		names = make([]string, 2)
		names[0] = option.name
		names[1] = option.nameAlias
	}
	return names, nil
}

func checkOption(options OptionMapType) error {
	for name, optionInfo := range OptionMap {
		if option, ok := options[name]; ok {
		    if optionInfo.optionType == OptionTypeInt64 {
				if val, ook := option.(*string); ook && *val != "" {
					num, err := strconv.ParseInt(*val, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid option value of %s, the value: %s is not int64, please check", name, *val)
					}

					if optionInfo.minVal != "" {
						minv, _ := strconv.ParseInt(optionInfo.minVal, 10, 64)
						if num < minv {
							return fmt.Errorf("invalid option value of %s, the value: %d is smaller than the min value range: %d", name, num, minv)
						}
					}
					if optionInfo.maxVal != "" {
						maxv, _ := strconv.ParseInt(optionInfo.maxVal, 10, 64)
						if num > maxv {
							return fmt.Errorf("invalid option value of %s, the value: %d is bigger than the max value range: %d", name, num, maxv)
						}
					}
				}
			}
            if optionInfo.optionType == OptionTypeAlternative {
				if val, ook := option.(*string); ook && *val != "" {
                    vals := strings.Split(optionInfo.minVal, "/")
                    if FindPosCaseInsen(*val, vals) == -1 {
                        return fmt.Errorf("invalid option value of %s, the value: %s is not anyone of %s", name, *val, optionInfo.minVal)
                    }
                }
            }
		}
	}
	return nil
}

// GetBool is used to get bool option from option map parsed by ParseArgOptions
func GetBool(name string, options OptionMapType) (bool, error) {
	if option, ok := options[name]; ok {
		if val, ook := option.(*bool); ook {
			return *val, nil
		}
		return false, fmt.Errorf("Error: option value of %s is not bool", name)
	}
	return false, fmt.Errorf("Error: there is no option for %s", name)
}

// GetInt is used to get int option from option map parsed by ParseArgOptions
func GetInt(name string, options OptionMapType) (int64, error) {
	if option, ok := options[name]; ok {
		switch option.(type) {
		case *string:
			val, err := strconv.ParseInt(*(option.(*string)), 10, 64)
            if err == nil {
				return val, nil
			}
            if *(option.(*string)) == "" {
                return 0, fmt.Errorf("Option value of %s is empty", name)
            }
            return 0, err
		case *int64:
			return *(option.(*int64)), nil
		default:
			return 0, fmt.Errorf("Option value of %s is not int64", name)
		}
	} else {
		return 0, fmt.Errorf("There is no option for %s", name)
	}
	return 0, nil
}

// GetString is used to get string option from option map parsed by ParseArgOptions
func GetString(name string, options OptionMapType) (string, error) {
	if option, ok := options[name]; ok {
		if val, ook := option.(*string); ook {
			return *val, nil
		}
		return "", fmt.Errorf("Error: Option value of %s is not string", name)
	}
	return "", fmt.Errorf("Error: There is no option for %s", name)
}
