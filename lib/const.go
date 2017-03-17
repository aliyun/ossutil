package lib

import (
	"os"
)

// all supported options of ossutil
const (
	OptionConfigFile       string = "configFile"
	OptionEndpoint                = "endpoint"
	OptionAccessKeyID             = "accessKeyID"
	OptionAccessKeySecret         = "accessKeySecret"
	OptionSTSToken                = "stsToken"
	OptionACL                     = "acl"
	OptionShortFormat             = "shortFormat"
	OptionLimitedNum              = "limitedNum"
	OptionMarker                  = "marker"
	OptionUploadIDMarker          = "uploadIDMarker"
	OptionDirectory               = "directory"
	OptionMultipart               = "multipart"
	OptionAllType                 = "allType"
	OptionRecursion               = "recursive"
	OptionBucket                  = "bucket"
	OptionForce                   = "force"
	OptionUpdate                  = "update"
	OptionDelete                  = "delete"
	OptionContinue                = "continue"
	OptionOutputDir               = "outputDir"
	OptionBigFileThreshold        = "bigfileThreshold"
	OptionCheckpointDir           = "checkpointDir"
	OptionSnapshotPath            = "snapshotPath"
	OptionRetryTimes              = "retryTimes"
	OptionRoutines                = "routines"
	OptionParallel                = "parallel"
    OptionRange                   = "range"
	OptionLanguage                = "language"
	OptionHashType                = "hashType"
	OptionVersion                 = "version"
)

// the elements show in stat object
const (
	StatName             string = "Name"
	StatLocation                = "Location"
	StatCreationDate            = "CreationDate"
	StatExtranetEndpoint        = "ExtranetEndpoint"
	StatIntranetEndpoint        = "IntranetEndpoint"
	StatACL                     = "ACL"
	StatOwner                   = "Owner"
	StatLastModified            = "Last-Modified"
	StatContentMD5              = "Content-Md5"
	StatCRC64                   = "X-Oss-Hash-Crc64ecma"
)

// the elements show in hash file
const (
	HashCRC64      = "CRC64-ECMA"
	HashMD5        = "MD5"
	HashContentMD5 = "Content-MD5"
)

const (
	updateEndpoint       string = "oss-cn-hangzhou.aliyuncs.com"
	updateBucket                = "ossutil-version-update"
	updateVersionObject         = "ossutilversion"
	updateBinaryLinux           = "ossutil"
	updateBinaryWindow32        = "ossutil32.exe"
	updateBinaryWindow64        = "ossutil64.exe"
	updateBinaryMac64           = "ossutilmac64"
	updateTmpVersionFile        = ".ossutil_tmp_vsersion"
)

// global public variable
const (
	Package                 string = "ossutil"
	ChannelBuf              int    = 1000
	Version                 string = "1.0.0.Beta2"
	DefaultEndpoint         string = "oss.aliyuncs.com"
	ChineseLanguage                = "CH"
	EnglishLanguage                = "EN"
	DefaultLanguage                = ChineseLanguage
	Scheme                  string = "oss"
	DefaultConfigFile              = "~" + string(os.PathSeparator) + ".ossutilconfig"
	MaxUint                        = ^uint(0)
	MaxInt                         = int(MaxUint >> 1)
	MaxUint64                      = ^uint64(0)
	MaxInt64                       = int64(MaxUint64 >> 1)
	ReportPrefix                   = "ossutil_report_"
	ReportSuffix                   = ".report"
	DefaultOutputDir               = "ossutil_output"
	CheckpointDir                  = ".ossutil_checkpoint"
	CheckpointSep                  = "---"
	SnapshotConnector              = "==>"
	SnapshotSep                    = "#"
	MaxPartNum                     = 10000
	MaxIdealPartNum                = MaxPartNum / 10
	MinIdealPartNum                = MaxPartNum / 500
	MaxIdealPartSize               = 524288000
	MinIdealPartSize               = 1048576
	DefaultBigFileThreshold        = 104857600
	MaxBigFileThreshold            = MaxInt64
	MinBigFileThreshold            = 0
	DefaultLimitedNum              = -1
	MinLimitedNum                  = 0
	RetryTimes              int    = 3
	MaxRetryTimes           int64  = 500
	MinRetryTimes           int64  = 1
	Routines                int    = 5
	MaxRoutines             int64  = 10000
	MinRoutines             int64  = 1
	MaxParallel             int64  = 10000
	MinParallel             int64  = 1
	DefaultHashType         string = "crc64"
	MD5HashType             string = "md5"
	LogFilePrefix                  = "ossutil_log_"
)

const (
	objectType    = 0x00000001
	multipartType = 0x00000010
	allType       = objectType | multipartType // marker for objects
	bucketType    = 0x10000000
)
