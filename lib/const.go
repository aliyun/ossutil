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
	OptionDirectory               = "directory"
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
    OptionLanguage                = "language"
    OptionHashType                = "hashType"
	OptionVersion                 = "version"
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
    StatContentMD5                 = "Content-Md5"
    StatCRC64                      = "X-Oss-Hash-Crc64ecma"
)

// the elements show in hash file
const (
    HashCRC64                   = "CRC64-ECMA"
    HashMD5                     = "MD5"  
    HashContentMD5              = "Content-MD5"
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
	Package                         string = "ossutil"
	ChannelBuf                      int    = 1000
	Version                         string = "1.0.0.Beta1"
	DefaultEndpoint                 string = "oss.aliyuncs.com"
	DefaultLanguage                        = "CH"
    EnglishLanguage                        = "EN"
	Scheme                          string = "oss"
    DefaultConfigFile                      = "~" + string(os.PathSeparator) +".ossutilconfig"
	MaxUint                                = ^uint(0)
	MaxInt                                 = int(MaxUint >> 1)
	MaxUint64                              = ^uint64(0)
	MaxInt64                               = int64(MaxUint64 >> 1)
    ReportPrefix                           = "ossutil_report_"
    ReportSuffix                           = ".report"
    DefaultOutputDir                       = "ossutil_output"
	CheckpointDir                          = ".ossutil_checkpoint"
	CheckpointSep                          = "---"
    SnapshotConnector                      = "==>"
    SnapshotSep                            = "#"
	MaxPartNum                             = 10000
	MaxIdealPartNum                        = MaxPartNum / 10
	MinIdealPartNum                        = MaxPartNum / 500
	MaxIdealPartSize                       = 524288000
	MinIdealPartSize                       = 1048576
	DefaultBigFileThreshold                = 104857600 
	MaxBigFileThreshold                    = MaxInt64
	MinBigFileThreshold                    = 0
	RetryTimes                      int    = 3
	MaxRetryTimes                   int64  = 500
	MinRetryTimes                   int64  = 1
	Routines                        int    = 3
	MaxRoutines                     int64  = 100
	MinRoutines                     int64  = 1
	MaxParallel                     int64  = 100
	MinParallel                     int64  = 1
	DefaultHashType                 string = "crc64"
	MD5HashType                     string = "md5"
    LogFilePrefix                          = "ossutil_log_"
)
