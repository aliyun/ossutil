package lib

// all supported options of ossutil 
const (
	OptionConfigFile       string = "configFile"
	OptionEndpoint                = "endpoint"
	OptionAccessKeyID             = "accessKeyID"
	OptionAccessKeySecret         = "accessKeySecret"
	OptionSTSToken                = "stsToken"
    OptionACL                     = "acl"
	OptionShortFormat             = "shortFormat"
	OptionDirectory               = "directory"
	OptionRecursion               = "recursive"
	OptionBucket                  = "bucket"
	OptionForce                   = "force"
	OptionUpdate                  = "update"
	OptionDelete                  = "delete"
	OptionBigFileThreshold        = "bigfileThreshold"
	OptionCheckpointDir           = "checkpointDir"
	OptionRetryTimes              = "retryTimes"
	OptionRoutines                = "routines"
    OptionLanguage                = "language"
	OptionVersion                 = "version"
	OptionMan                     = "man"
)

// the elements show in stat object
const (
    StatName                string = "Name"
    StatLocation                   = "Location"
    StatCreationDate               = "CreationDate"  
    StatExtranetEndpoint           = "ExtranetEndpoint"
    StatIntranetEndpoint           = "IntranetEndpoint"
    StatACL                        = "ACL"
    StatOwner                      = "Owner"
    StatLastModified               = "Last-Modified"
)

const (
    updateEndpoint      string = "oss-cn-hangzhou.aliyuncs.com"
    updateBucket               = "ossutil-version-update"
    updateVersionObject        = "ossutilversion"         
    updateBinaryLinux          = "ossutil"
    updateBinaryWindow32       = "ossutil32.exe"
    updateBinaryWindow64       = "ossutil64.exe"
    updateBinaryMac64          = "ossutilmac64"
    updateTmpVersionFile       = ".ossutil_tmp_vsersion"
)

// global public variable
const (
	Package             string = "ossutil"
	ChannelBuf          int    = 1000
	Version             string = "1.0.0"
	DefaultEndpoint     string = "oss-cn-hangzhou.aliyuncs.com"
	DefaultLanguage            = "中文"
    EnglishLanguage            = "English"
	Scheme              string = "oss"
    DefaultConfigFile          = "~/.ossutilconfig"
	MaxUint                    = ^uint(0)
	MaxInt                     = int(MaxUint >> 1)
	MaxUint64                  = ^uint64(0)
	MaxInt64                   = int64(MaxUint64 >> 1)
	CheckpointDir              = ".ossutil_checkpoint"
	CheckpointSep              = "---"
	MaxPartNum                 = 10000
	MaxIdealPartNum            = MaxPartNum / 20
	MinIdealPartNum            = MaxPartNum / 500
	MaxIdealPartSize           = 524288000
	MinIdealPartSize           = 10485760
	BigFileThreshold           = 524288000
	MaxBigFileThreshold        = MaxInt64
	MinBigFileThreshold        = 0
	RetryTimes          int    = 3
	MaxRetryTimes       int64  = 500
	MinRetryTimes       int64  = 1
	Routines            int    = 3
	MaxRoutines         int64  = 32
	MinRoutines         int64  = 1
)
