package lib

import (
	"fmt"
	configparser "github.com/alyu/configparser"
	"os/user"
	"strconv"
	"strings"
)

// sections in config file
const (
	CREDSection string = "Credentials"

	BucketEndpointSection string = "Bucket-Endpoint"

	BucketCnameSection string = "Bucket-Cname"
)

type configOption struct {
	showNames []string
	reveal    bool
}

// CredOptionList is all options in Credentials section
var CredOptionList = []string{
    OptionLanguage,
    OptionEndpoint,
    OptionAccessKeyID,
    OptionAccessKeySecret,
    OptionSTSToken,
}

// CredOptionMap allows alias name for options in Credentials section 
var CredOptionMap = map[string]configOption{
    OptionLanguage:        configOption{[]string{"language", "Language"}, false},
	OptionEndpoint:        configOption{[]string{"endpoint", "host"}, true},
	OptionAccessKeyID:     configOption{[]string{"accessKeyID", "accessKeyId", "AccessKeyID", "AccessKeyId", "access_key_id", "access_id", "accessid"}, false},
	OptionAccessKeySecret: configOption{[]string{"accessKeySecret", "AccessKeySecret", "access_key_secret", "access_key", "accesskey"}, false},
	OptionSTSToken:        configOption{[]string{"stsToken", "ststoken", "sts_token"}, false},
}

// DecideConfigFile return the config file, if user not specified, return default one 
func DecideConfigFile(configFile string) string {
	if configFile == "" {
		configFile = DefaultConfigFile
	}
	usr, _ := user.Current()
	dir := usr.HomeDir
	if len(configFile) >=2 && configFile[:2] == "~/" {
		configFile = strings.Replace(configFile, "~", dir, 1)
	}
	return configFile
}

// LoadConfig load the specified config file 
func LoadConfig(configFile string) (OptionMapType, error) {
	var configMap OptionMapType
	var err error
	configMap, err = readConfigFromFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("Read config file error, please check, try \"help config\" to set configuration or use \"--config_file\" option, error info: %s", err) 
	}
	if err = checkConfig(configMap); err != nil {
		return nil, err
	}
	return configMap, nil
}

func readConfigFromFile(configFile string) (OptionMapType, error) {
	configFile = DecideConfigFile(configFile)

	config, err := configparser.Read(configFile)
	if err != nil {
		return nil, err 
	}

	configMap := OptionMapType{}

	// get options in cred section
	credSection, err := config.Section(CREDSection)
	if err != nil {
		return nil, err
	}

	credOptions := credSection.Options()
	for name, option := range credOptions {
		if opName, ok := getOptionNameByStr(name); ok {
			configMap[opName] = option
		}
	}

	// get options in pair sections
	for _, sec := range []string{BucketEndpointSection, BucketCnameSection} {
		if section, err := config.Section(sec); err == nil {
			configMap[sec] = map[string]string{}
			options := section.Options()
			for bucket, host := range options {
				(configMap[sec]).(map[string]string)[bucket] = host
			}
		}
	}
	return configMap, nil
}

func getOptionNameByStr(name string) (string, bool) {
	for optionName, option := range CredOptionMap {
		for _, val := range option.showNames {
			if name == val {
				return optionName, true
			}
		}
	}
	return "", false
}

func checkConfig(configMap OptionMapType) error {
	for name, opval := range configMap {
		if option, ok := OptionMap[name]; ok {
			if option.optionType == OptionTypeInt64 {
				if _, err := strconv.ParseInt(opval.(string), 10, 64); err != nil {
					return fmt.Errorf("error value of option \"%s\", the value is: %s in config file, which needs int64 type", name, opval)
				}
			}
            if option.optionType == OptionTypeAlternative {
                vals := strings.Split(option.minVal, "|") 
                if FindPosCaseInsen(opval.(string), vals) == -1 {
                    return fmt.Errorf("error value of option \"%s\", the value is: %s in config file, which is not anyone of %s", name, opval, option.minVal)
                }
            }
		}
	}
	return nil
}
