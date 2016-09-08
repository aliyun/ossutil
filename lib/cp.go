package lib

import (
	"fmt"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type operationType int

const (
	operationTypePut operationType = iota
	operationTypeGet
	operationTypeCopy
)

type copyOptionType struct {
	recursive bool
	force     bool
	update    bool
	threshold int64
	cpDir     string
	routines  int64
}

type fileInfoType struct {
	filePath string
	dir      string
}

type objectInfoType struct {
	key          string
	size         int64
	lastModified time.Time
}

var mu sync.Mutex

var specChineseCopy = SpecText{

	synopsisText: "上传，下载或拷贝Objects",

	paramText: "src_url dest_url [options]",

	syntaxText: ` 
    ossutil cp file_url cloud_url  [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file] [-c file] 
    ossutil cp cloud_url file_url  [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file] [-c file] 
    ossutil cp cloud_url cloud_url [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file] [-c file] 
`,

	detailHelpText: ` 
    该命令允许：从本地文件系统上传文件到oss，从oss下载object到本地文件系统，在oss
    上进行object拷贝。分别对应下述三种操作：
        ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]

    其中file_url代表本地文件系统中的文件路径，支持相对路径或绝对路径，请遵循本地文
    件系统的使用格式；
    oss://bucket[/prefix]代表oss上的object，支持前缀匹配，不支持通配符。
    ossutil通过oss://前缀区分本地文件系统的文件和oss文件。


--recursive选项

    （1）上传文件到oss时，如果file_url为目录，则必须指定--recursive选项，否则无需指
        定--recursive选项。

    （2）从oss下载或在oss间拷贝文件时：
        如果未指定--recursive选项，则认为拷贝单个object，此时请确保src_url精确指定待
        拷贝的object，如果object不存在，则报错。

        如果指定了--recursive选项，ossutil会对src_url进行prefix匹配查找，对这些objects
        批量拷贝，如果拷贝失败，会终止操作，已经执行的拷贝不会回退。

--update选项
    
    如果指定了该选项，ossutil只有当目标文件（或object）不存在，或源文件（或object）新于
    目标文件（或object）时，才执行拷贝。当指定了该选项时，无论--force选项是否指定了，在
    目标文件存在时，ossutil都不会提示，直接采取上述策略。
    该选项可用于当批量拷贝失败时，重传时跳过已经成功的文件。

--force选项

    如果dest_url指定的文件或objects已经存在，并且未指定--update选项，ossutil会询问是否进
    行替换操作（输入非法时默认不替换），如果指定了--force选项，则不询问，强制替换。该选项
    只有在未指定--update选项时有效，否则按--update选项操作。


大文件断点续传：

    如果源文件大小超过--bigfile_threshold选项指定的大小（默认为500M），ossutil会认为该文件
    为大文件，并自动使用断点续传策略，策略如下：
    （1）上传到oss时：ossutil会对大文件自动分片，进行multipart分片上传，如果上传失败，会
        在本地的.ossutil_checkpoint目录记录失败信息，下次重传时会读取.ossutil_checkpoint目
        录中的信息进行断点续传，当上传成功时会删除.ossutil_checkpoint目录。
    （2）从oss下载时：ossutil会自动对大文件分片下载，组装成一个文件，如果下载失败，同样会
        在.ossutil_checkpoint目录记录失败信息，重试成功后会删除.ossutil_checkpoint目录。
    （3）在oss间拷贝：ossutil会自动对大文件分片，使用Upload Part Copy方式拷贝，同样会在
        .ossutil_checkpoint目录记录失败信息，重试成功后会删除.ossutil_checkpoint目录。

    注意：
    1）小文件不会采用断点续传策略，失败后下次直接重传。
    2）在操作（1）和（3）中，如果操作失败，oss上可能会产生未complete的uploadId，但是只要最
    终操作成功，就不会存在未complete的uploadId（被组装成object）。
    3）上传到oss时，如果.ossutil_checkpoint目录包含在file_url中，.ossutil_checkpoint目录不会
    被上传到oss。该目录路径可以用--checkpoint_dir选项指定，如果指定了该选项，请确保指定的目录
    可以被删除。


用法：

    该命令有三种用法：

    1) ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        该用法上传本地文件系统中文件或目录到oss。file_url可以为文件或目录。当file_url为文件
    时，无论是否指定--recursive选项都不会影响结果。当file_url为目录时，即使目录为空或者只含
    有一个文件，也必须使用--recursive选项，注意，此时ossutil会将file_url下的文件或子目录上传
    到oss，但不同于shell拷贝，file_url所代表的首层目录不会被创建。
    object命名规则：
        当file_url为文件时，如果prefix为空或以/结尾，object名为：dest_url+文件名。
                            否则，object名为：dest_url。
        当file_url为目录时，如果prefix为空或以/结尾，object名为：dest_url+文件或子目录相对
                            file_url的路径。
                            否则，object名为：dest_url+/+文件或子目录相对file_url的路径。

    2) ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        该用法下载oss上的单个或多个Object到本地文件系统。如果未指定--recursive选项，则ossutil
    认为src_url精确指定了待拷贝的单个object，此时不支持prefix匹配，如果object不存在则报错。如
    果指定了--recursive选项，ossutil会搜索prefix匹配的objects，批量拷贝这些objects，此时file_url
    必须为目录，如果该目录不存在，ossutil自动创建该目录。
    文件命名规则：
        当file_url为文件时，下载到file_url指定的文件，文件名与file_url保持一致。
        当file_url为目录时，下载到file_url指定的目录中，文件名为：object名称去除prefix。
    注意：对于以/结尾且大小为0的object，会在本地文件系统创建一个目录，而不是尝试创建一个文件。
    对于其他object会尝试创建文件。

    3) ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        该用法在oss间进行object的拷贝。其中src_bucket与dest_bucket可以相同，注意，当src_url与
    dest_url完全相同时，ossutil不会做任何事情，直接提示退出。设置meta请使用setmeta命令。如果未
    指定--recursive选项，则认为src_url精确指定了待拷贝的单个object，此时不支持prefix匹配，如果
    object不存在则报错。如果指定了--recursive选项，ossutil会搜索prefix匹配的objects，批量拷贝这
    些objects。
    注意：批量拷贝时，src_url包含dest_url，或dest_url包含src_url是不允许的（dest_url以src_url为
    前缀时，会产生递归拷贝，src_url以dest_url为前缀时，会覆盖待拷贝文件）。单个拷贝时，该情况是
    允许的。
    object命名规则：
        当src_url为单个文件时，如果dest_url的prefix为空或以/结尾，object名为：dest_url+object名去除所在父目录的路径。
                               否则，object名为：dest_url。
        当src_url为多个文件时，object名为：dest_url+源object名去除src_prefix。
`,

	sampleText: ` 
    1) 上传文件到oss
    假设本地local_dir目录中有文件a，b，目录c和d，目录c为空，目录d中包含文件dd。
    
    ossutil cp local_dir/a oss://bucket1
    生成：
        oss://bucket1/a

    ossutil cp local_dir/a oss://bucket1/b
    生成：
        oss://bucket1/b

    ossutil cp local_dir/a oss://bucket1/b/
    生成：
        oss://bucket1/b/a

    ossutil cp local_dir oss://bucket1/b/
    报错

    ossutil cp local_dir oss://bucket1/b -r
    生成：
        oss://bucket1/b/a
        oss://bucket1/b/b
        oss://bucket1/b/c/
        oss://bucket1/b/d/
        oss://bucket1/b/d/dd

    2) 从oss下载object
    假设oss上有下列objects：
        oss://bucket/abcdir1/a
        oss://bucket/abcdir1/b
        oss://bucket/abcdir1/c
        oss://bucket/abcdir2/a/
        oss://bucket/abcdir2/b/e
    其中oss://bucket/abcdir2/a/的size为0。

    ossutil cp oss://bucket/abcdir1/a b
    生成文件b

    ossutil cp oss://bucket/abcdir1/a b/
    在目录b下生成文件a

    ossutil cp oss://bucket/abcdir2/a/ b
    如果b为已存在文件，报错。
    如果b为已存在目录，在目录b下生成目录a。

    ossutil cp oss://bucket/abc b
    报错，object不存在。

    ossutil cp oss://bucket/abc b -r
    如果b为已存在文件，报错    
    否则在目录b下生成目录dir1和dir2，
        目录dir1中生成文件a，b，c
        目录dir2中生成目录a和b，目录b中生成文件e
        
    3) 在oss间拷贝
    假设oss上有下列objects：
        oss://bucket/abcdir1/a
        oss://bucket/abcdir1/b
        oss://bucket/abcdir1/c
        oss://bucket/abcdir2/a/
        oss://bucket/abcdir2/b/e

    ossutil cp oss://bucket/abcdir1/a oss://bucket1
    生成：
        oss://bucket1/a

    ossutil cp oss://bucket/abcdir1/a oss://bucket1/b
    生成:
        oss://bucket1/b

    ossutil cp oss://bucket/abcdir1/a oss://bucket/abcdir1/a/ 
    生成:
        oss://bucket/abcdir1/a/a

    ossutil cp oss://bucket/abcdir1/a/ oss://bucket/abcdir1/b/ 
    生成：
        oss://bucket/abcdir1/b/a/

    ossutil cp oss://bucket/abcdir1/a oss://bucket/abcdir1/a/ -r 
    报错，递归拷贝

    ossutil cp oss://bucket/abcdir1/a oss://bucket1/b/
    生成：
        oss://bucket1/b/a

    ossutil cp oss://bucket/abc oss://bucket1/b/
    报错，object不存在

    ossutil cp oss://bucket/abc oss://bucket1/123 -r
    生成：
        oss://bucket1/123dir1/a
        oss://bucket1/123dir1/b
        oss://bucket1/123dir1/c
        oss://bucket1/123dir2/a/
        oss://bucket1/123dir2/b/e

    ossutil cp oss://bucket/abc oss://bucket1/123/ -r
    生成：
        oss://bucket1/123/dir1/a
        oss://bucket1/123/dir1/b
        oss://bucket1/123/dir1/c
        oss://bucket1/123/dir2/a/
        oss://bucket1/123/dir2/b/e
`,
}

var specEnglishCopy = SpecText{

	synopsisText: "Upload, Download or Copy Objects",

	paramText: "src_url dest_url [options]",

	syntaxText: ` 
    ossutil cp file_url cloud_url  [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file] [-c file] 
    ossutil cp cloud_url file_url  [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file] [-c file] 
    ossutil cp cloud_url cloud_url [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file] [-c file] 
`,

	detailHelpText: ` 
    The command allows: 
    1. Upload file from local file system to oss 
    2. Download object from oss to local file system
    3. Copy objects between oss
    Which matches with the following three kinds of operations:
        ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]

    file_url means the file in local file system, it supports relative path and absolute 
    path, the usage of file_url is same with your local file system. oss://bucket[/prefix] 
    means object in oss, it supports prefix matching, but not support wildcard.

    ossutil sperate file of local system and oss objects by the prefix of oss://, which means 
    if the url starts with oss://, ossutil considers it as object, else, ossutil considers it 
    as file in local system. 

--recursive option:

    (1) Upload file to oss: if file_url is directory, the --recursive option must be specified. 

    (2) When download objects or copy objects between oss:
        If --recursive option is not specified, ossutil download or copy the specified single 
    object, in the usage, please make sure url exactly specified the object you want to set meta 
    on, if object not exist, error occurs. 
        If --recursive option is specified, ossutil will search for prefix-matching objects and 
    download or copy these objects. If error occurs, the operation will be terminated, objects 
    which has been download or copyed will not rollback. 

--update option

    Use the --update option to copy only when the source file is newer than the destination file 
    when the destination file is missing. If --update option is specified, when the destionation 
    file is existed, ossutil will not prompt and copy when newer, no matter if --force option is 
    specified or not.
    The option can be used when batch copy failed, skip the succeed files in retry.

--force option

    If the file dest_url specified is existed, and --update option is not specified, ossutil will 
    ask if replace the file(if the input is invalid, the file will not be replaced). If --force 
    option is specified here, ossutil will not prompt, replace by force. The option is useful only 
    when --update not specified. 


Resume copy of big file:

    If the size of source file is bigger than what --bigfile_threshold option specified(default: 
    500M), ossutil will consider the file as a big file, and use resume copy policy to these files:
    (1) Upload file to oss: ossutil will split the big file to many parts, use multipart upload. If 
        upload is failed, ossutil will record failure information in .ossutil_checkpoint directory 
        in local file system. When retry, ossutil will read the checkpoint information and resume 
        upload, if the upload is succeed, ossutil will remove the .ossutil_checkpoint directory. 
    (2) Download object from oss: ossutil will split the big file to many parts, range get each part. 
        If download is failed, ossutil wll record failure information in .ossutil_checkpoint directory 
        in local file system. If success, ossutil will remove the directory.
    (3) Copy between oss: ossutil will split the big file to many parts, use Upload Part Copy, and 
        record failure information in .ossutil_checkpoint directory in local file system. If success, 
        ossutil will remove the directory.

    Warning:
    1) Resume copy will not be implemented on small file, if failure happens, ossutil will copy the 
        whole file the next time.
    2) In operation (1) and (3), if failure happens, uploadId that has not been completed may appear in 
        oss. If the operation success after retry, these uploadId will be completed automatically. 
    3) When upload file to oss, if .ossutil_checkpoint directory is included in file_url, .ossutil_checkpoint 
        will not be uploaded to oss. The path of checkpoint directory can be specified by --checkpoint_dir 
        option, please make sure the directory you specified can be removed.


Usage:

    There are three usages:

    1) ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        The usage upload file in local system to oss. file_url can be file or directory. If file_url 
    is file, no matter --recursive option is specified or not will not affect the result. If file_url 
    is directory, even if the directory is empty or only contains one file, we must specify --recursive 
    option. Mind that, ossutil will upload all sub files and directories(include empty directory) inside 
    file_url to oss, but differe from shell cp, the first level directory specified by file_url will not 
    be upload to oss. 
    Object Naming Rules:
        If file_url is file: if prefix is empty or end with "/", object name is: dest_url + file name.
                             else, object name is: dest_url.
        If file_url is directory: if prefix is empty or end with "/", object name is: dest_url + file path relative to file_url.
        
    2) ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        The usage download one or many objects to local system. If --recursive option is not specified, 
    ossutil considers src_url exactly specified the single object you want to download, prefix-matching 
    is not supported now, if the object not exists, error occurs. If --recursive option is specified, 
    ossutil will search for prefix-matching objects and batch download those objects, at this time file_url 
    must be directory, if the directory not exists, ossutil will create the directory automatically.
    File Naming Rules:
        If file_url is file, ossutil download file to the path of file_url, and the file name is got from file_url.
        If file_url is directory, ossutil download file to the directory, and the file name is: object name exclude prefix.
    Warning: If the object name is end with / and size is zero, ossutil will create a directory in local 
    system, instead of creating a file.

    3) ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [--update] [--bigfile_threshold=size] [--checkpoint_dir=file]
        The usage copy objects between oss. The src_bucket can be same with dest_bucket. Pay attention 
    please, if src_url is the same with dest_url, ossutil will do nothing but exit after prompt. Set meta 
    please use "setmeta" command. If --recursive option is not specified, ossutil considers src_url exactly 
    specified the single object you want to copy. If --recursive option is specified, ossutil will search 
    for prefix-matching objects and batch copy those objects. 

    Warning: when batch copy, it's not allowed that src_url is the prefix of dest_url, because recursivlly 
    copy will happen under the situation. dest_ur is the prefix of src_url is not allowed too, because of 
    covering source file. But they are allowed in single file copy.

    Object Naming Rules:
        If src_url is one object: if prefix of dest_object is empty or end with "/", object name is: dest_url + object name exclude parenet directory path. 
                                  else, object name is: dest_url.
        If src_url means multiple objects: object name is: dest_url+ source object name exclude src_prefix.
`,

	sampleText: ` 
    1) Upload to oss
    Suppose there are directory local_dir in local system, 
        local_dir contains file a, b directory c, d, 
        c is empty, d contains file dd.
    
    ossutil cp local_dir/a oss://bucket1
    Generate:
        oss://bucket1/a

    ossutil cp local_dir/a oss://bucket1/b
    Generate:
        oss://bucket1/b

    ossutil cp local_dir/a oss://bucket1/b/
    Generate:
        oss://bucket1/b/a

    ossutil cp local_dir oss://bucket1/b/
    Error

    ossutil cp local_dir oss://bucket1/b -r
    Generate:
        oss://bucket1/b/a
        oss://bucket1/b/b
        oss://bucket1/b/c/
        oss://bucket1/b/d/
        oss://bucket1/b/d/dd

    2) download from oss
    Suppose there are following objects in oss:
        oss://bucket/abcdir1/a
        oss://bucket/abcdir1/b
        oss://bucket/abcdir1/c
        oss://bucket/abcdir2/a/
        oss://bucket/abcdir2/b/e
    And size of oss://bucket/abcdir2/a/ is zero. 

    ossutil cp oss://bucket/abcdir1/a b
    Generate file b

    ossutil cp oss://bucket/abcdir1/a b/
    Generate file a under directory b

    ossutil cp oss://bucket/abcdir2/a/ b
    If b exists and is a file, error occurs.
    If b exists and is a directory, generate directory a under directory b.

    ossutil cp oss://bucket/abc b
    Error: object not exist

    ossutil cp oss://bucket/abc b -r
    If b exists and is a file, error occurs.
    Else generate directory dir1, dir2,
        generate file a, b, c in dir1,
        generate directory a, b in dir2, generate file e in directory b.
        
    3) Copy between oss 
    Suppose there are following objects in oss:
        oss://bucket/abcdir1/a
        oss://bucket/abcdir1/b
        oss://bucket/abcdir1/c
        oss://bucket/abcdir2/a/
        oss://bucket/abcdir2/b/e

    ossutil cp oss://bucket/abcdir1/a oss://bucket1
    Generate:
        oss://bucket1/a

    ossutil cp oss://bucket/abcdir1/a oss://bucket1/b
    Generate:
        oss://bucket1/b

    ossutil cp oss://bucket/abcdir1/a oss://bucket/abcdir1/a/ 
    Generate:
        oss://bucket/abcdir1/a/a

    ossutil cp oss://bucket/abcdir1/a/ oss://bucket/abcdir1/b/ 
    Generate:
        oss://bucket/abcdir1/b/a/

    ossutil cp oss://bucket/abcdir1/a oss://bucket/abcdir1/a/ -r 
    Error, recursivlly copy

    ossutil cp oss://bucket/abcdir1/a oss://bucket1/b/
    Generate:
        oss://bucket1/b/a

    ossutil cp oss://bucket/abc oss://bucket1/b/
    Error: object not exist

    ossutil cp oss://bucket/abc oss://bucket1/123 -r
    Generate:
        oss://bucket1/123dir1/a
        oss://bucket1/123dir1/b
        oss://bucket1/123dir1/c
        oss://bucket1/123dir2/a/
        oss://bucket1/123dir2/b/e

    ossutil cp oss://bucket/abc oss://bucket1/123/ -r
    Generate:
        oss://bucket1/123/dir1/a
        oss://bucket1/123/dir1/b
        oss://bucket1/123/dir1/c
        oss://bucket1/123/dir2/a/
        oss://bucket1/123/dir2/b/e
`,
}

// CopyCommand is the command upload, download and copy objects
type CopyCommand struct {
	command Command
}

var copyCommand = CopyCommand{
	command: Command{
		name:        "cp",
		nameAlias:   []string{"copy"},
		minArgc:     2,
		maxArgc:     MaxInt,
		specChinese: specChineseCopy,
		specEnglish: specEnglishCopy,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionConfigFile,
			OptionRecursion,
			OptionForce,
			OptionUpdate,
			OptionBigFileThreshold,
			OptionCheckpointDir,
			OptionRetryTimes,
			OptionRoutines,
		},
	},
}

// function for FormatHelper interface
func (cc *CopyCommand) formatHelpForWhole() string {
	return cc.command.formatHelpForWhole()
}

func (cc *CopyCommand) formatIndependHelp() string {
	return cc.command.formatIndependHelp()
}


// Init simulate inheritance, and polymorphism 
func (cc *CopyCommand) Init(args []string, options OptionMapType) error {
	return cc.command.Init(args, options, cc)
}

// RunCommand simulate inheritance, and polymorphism
func (cc *CopyCommand) RunCommand() error {
	cpOption := copyOptionType{}
	cpOption.recursive, _ = GetBool(OptionRecursion, cc.command.options)
	cpOption.force, _ = GetBool(OptionForce, cc.command.options)
	cpOption.update, _ = GetBool(OptionUpdate, cc.command.options)
	cpOption.threshold, _ = GetInt(OptionBigFileThreshold, cc.command.options)
	cpOption.cpDir, _ = GetString(OptionCheckpointDir, cc.command.options)
	cpOption.routines, _ = GetInt(OptionRoutines, cc.command.options)

	//get file list
	srcURLList, err := cc.getStorageURLs(cc.command.args[0 : len(cc.command.args)-1])
	if err != nil {
		return err
	}

	destURL, err := StorageURLFromString(cc.command.args[len(cc.command.args)-1])
	if err != nil {
		return err
	}

	opType := cc.getCommandType(srcURLList, destURL)
	if err := cc.checkCopyArgs(srcURLList, destURL, opType); err != nil {
		return err
	}

	//create ckeckpoint dir
	if err := os.MkdirAll(cpOption.cpDir, 0777); err != nil {
		return err
	}

	switch opType {
	case operationTypePut:
		err = cc.uploadFiles(srcURLList, destURL.(CloudURL), cpOption)
	case operationTypeGet:
		err = cc.downloadFiles(srcURLList[0].(CloudURL), destURL.(FileURL), cpOption)
	default:
		err = cc.copyFiles(srcURLList[0].(CloudURL), destURL.(CloudURL), cpOption)
	}

	if err == nil {
		return os.RemoveAll(cpOption.cpDir)
	}
	return err
}

func (cc *CopyCommand) getStorageURLs(urls []string) ([]StorageURLer, error) {
	urlList := []StorageURLer{}
	for _, url := range urls {
		storageURL, err := StorageURLFromString(url)
		if err != nil {
			return nil, err
		}
		if storageURL.IsCloudURL() && storageURL.(CloudURL).bucket == "" {
			return nil, fmt.Errorf("invalid cloud url: %s, miss bucket", url)
		}
		urlList = append(urlList, storageURL)
	}
	return urlList, nil
}

func (cc *CopyCommand) getCommandType(srcURLList []StorageURLer, destURL StorageURLer) operationType {
	if srcURLList[0].IsCloudURL() {
		if destURL.IsFileURL() {
			return operationTypeGet
		}
		return operationTypeCopy
	}
	return operationTypePut
}

func (cc *CopyCommand) checkCopyArgs(srcURLList []StorageURLer, destURL StorageURLer, opType operationType) error {
	for _, url := range srcURLList {
		if url.IsCloudURL() && url.(CloudURL).bucket == "" {
			return fmt.Errorf("invalid cloud url: %s, miss bucket", url.ToString())
		}
	}
	if destURL.IsCloudURL() && destURL.(CloudURL).bucket == "" {
		return fmt.Errorf("invalid cloud url: %s, miss bucket", destURL.ToString())
	}

	switch opType {
	case operationTypePut:
		if destURL.IsFileURL() {
			return fmt.Errorf("Copy files between local file system is not allowed in ossutil, if you want to upload to oss, please make sure dest_url starts with \"%s\", which is: %s now", SchemePrefix, destURL.ToString())
		}
		for _, url := range srcURLList {
			if url.IsCloudURL() {
				return fmt.Errorf("invalid url: %s, copy between oss operation appear in upload operation, multi-type operations is not supported in one command", url.ToString())
			}
		}
		if len(srcURLList) > 1 {
			return fmt.Errorf("invalid url: %s, multiple source url in upload operation", srcURLList[1].ToString())
		}
    case operationTypeGet:
		if len(srcURLList) > 1 {
			return fmt.Errorf("invalid url: %s, multiple source url in download operation", srcURLList[1].ToString())
		}
	default:
		if len(srcURLList) > 1 {
			return fmt.Errorf("invalid url: %s, multiple source url in copy operation", srcURLList[1].ToString())
		}
	}
	return nil
}

//function for upload files
func (cc *CopyCommand) uploadFiles(srcURLList []StorageURLer, destURL CloudURL, cpOption copyOptionType) error {
	bucket, err := cc.command.ossBucket(destURL.bucket)
	if err != nil {
		return err
	}

	//adjust oss prefix name
	destURL, err = cc.adjustDestURLForUpload(srcURLList, destURL)
	if err != nil {
		return err
	}

	// producer list files
	// consumer set acl
	chFiles := make(chan fileInfoType, ChannelBuf)
	chFinishFiles := make(chan fileInfoType, ChannelBuf)
	chSkipFiles := make(chan fileInfoType, ChannelBuf)
	chError := make(chan error, cpOption.routines+1)
	go cc.fileProducer(srcURLList, cpOption, chFiles, chError)
	for i := 0; int64(i) < cpOption.routines; i++ {
		go cc.uploadConsumer(bucket, destURL, cpOption, chFiles, chFinishFiles, chSkipFiles, chError)
	}

	completed := 0
	fnum := 0
	dnum := 0
	snum := 0
	for int64(completed) <= cpOption.routines {
		select {
		case file := <-chFinishFiles:
			if strings.HasSuffix(file.filePath, "/") || strings.HasSuffix(file.filePath, "\\") {
				dnum++
			} else {
				fnum++
			}
			cc.schedule(cpOption, fmt.Sprintf("\rdealed %d files or directories(upload %d files, %d directories, skip %d files)...", fnum+dnum+snum, fnum, dnum, snum))
		case <-chSkipFiles:
			snum++
			cc.schedule(cpOption, fmt.Sprintf("\rdealed %d files or directories(upload %d files, %d directories, skip %d files)...", fnum+dnum+snum, fnum, dnum, snum))
		case err := <-chError:
			if err != nil {
				fmt.Printf("\rdealed %d files or directories(upload %d files, %d directories, skip %d files), when error happens.\n", fnum+dnum+snum, fnum, dnum, snum)
				return err
			}
			completed++
		}
	}
	fmt.Printf("\rSucceed: scaned %d files or directories, dealed %d files or directories(upload %d files, %d directories, skip %d files).\n", fnum+dnum+snum, fnum+dnum+snum, fnum, dnum, snum)
	return nil
}

func (cc *CopyCommand) schedule(cpOption copyOptionType, str string) {
	if !cpOption.update && !cpOption.force {
		return
	}
	fmt.Printf(str)
}

func (cc *CopyCommand) adjustDestURLForUpload(srcURLList []StorageURLer, destURL CloudURL) (CloudURL, error) {
	if len(srcURLList) == 1 {
		f, err := os.Stat(srcURLList[0].ToString())
		if err != nil {
			return destURL, err
		}
		if !f.IsDir() {
			return destURL, nil
		}
	}

	if destURL.object != "" && !strings.HasSuffix(destURL.object, "/") && !strings.HasSuffix(destURL.object, "\\") {
		destURL.object += "/"
	}
	return destURL, nil
}

func (cc *CopyCommand) fileProducer(srcURLList []StorageURLer, cpOption copyOptionType, chFiles chan<- fileInfoType, chError chan<- error) {
	for _, url := range srcURLList {
		name := url.ToString()
		f, err := os.Stat(name)
		if err != nil {
			chError <- err
			return
		}
		if f.IsDir() {
			if !cpOption.recursive {
				chError <- fmt.Errorf("omitting directory \"%s\", please use --recursive option", name)
				return
			}
            fl, err := cc.getFileList(name)
            if err != nil {
                chError <- err
            }
			for _, fname := range fl {
				chFiles <- fileInfoType{fname, name}
			}
		} else {	
            dir, fname := filepath.Split(name) 
		    chFiles <- fileInfoType{fname, dir}
        }
    }

	defer close(chFiles)
	chError <- nil
}

func (cc *CopyCommand) getFileList(dpath string) ([]string, error) {
    fileList := []string{}
    err := filepath.Walk(dpath, func(fpath string, f os.FileInfo, err error) error {
        if f == nil {
            return err
        }

        dpath = filepath.Clean(dpath)
        fpath = filepath.Clean(fpath)
        fileName, err := filepath.Rel(dpath, fpath) 
        if err != nil {
            return fmt.Errorf("list file error: %s, info: %s", fpath, err.Error())
        }

        if f.IsDir(){
            if fpath != dpath {
                fileList = append(fileList, fileName + string(os.PathSeparator))
            }
            return nil
        }
        fileList = append(fileList, fileName)
        return nil
    })
    return fileList, err
}

func (cc *CopyCommand) uploadConsumer(bucket *oss.Bucket, destURL CloudURL, cpOption copyOptionType, chFiles <-chan fileInfoType, chFinishFiles, chSkipFiles chan<- fileInfoType, chError chan<- error) {
	for file := range chFiles {
		if cc.filterFile(file, cpOption.cpDir) {
			skip, err := cc.uploadFile(bucket, destURL, cpOption, file)
			if err != nil {
				chError <- err
				return
			}
			if skip {
				chSkipFiles <- file
			} else {
				chFinishFiles <- file
			}
		}
	}

	chError <- nil
}

func (cc *CopyCommand) filterFile(file fileInfoType, cpDir string) bool {
	filePath := file.filePath
	if file.dir != "" {
		filePath = file.dir + string(os.PathSeparator) + file.filePath
	}
	if !strings.Contains(filePath, cpDir) {
		return true
	}
	absFile, _ := filepath.Abs(filePath)
	absCPDir, _ := filepath.Abs(cpDir)
	return !strings.Contains(absFile, absCPDir)
}

func (cc *CopyCommand) uploadFile(bucket *oss.Bucket, destURL CloudURL, cpOption copyOptionType, file fileInfoType) (bool, error) {
	//first make object name
	objectName := cc.makeObjectName(destURL, file)

	filePath := file.filePath
	if file.dir != "" {
		filePath = file.dir + string(os.PathSeparator) + file.filePath
	}

	//get file size and last modify time
	f, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	if skip, err := cc.skipUpload(bucket, objectName, destURL, cpOption, f.ModTime()); err != nil || skip {
		return skip, err
	}

	if f.IsDir() {
		return false, cc.ossPutObjectRetry(bucket, objectName, "")
	}

	//decide whether to use resume upload
	if f.Size() < cpOption.threshold {
		return false, cc.ossUploadFileRetry(bucket, objectName, filePath)
	}

	//make options for resume multipart upload
	//part size
	partSize, rt := cc.preparePartOption(f.Size())
	//checkpoint file
	cp := oss.Checkpoint(true, cc.formatCPFileName(cpOption.cpDir, filePath, objectName))
	return false, cc.ossResumeUploadRetry(bucket, objectName, filePath, partSize, oss.Routines(rt), cp)
}

func (cc *CopyCommand) makeObjectName(destURL CloudURL, file fileInfoType) string {
	if destURL.object == "" || strings.HasSuffix(destURL.object, "/") || strings.HasSuffix(destURL.object, "\\") {
		return destURL.object + file.filePath
	}
	return destURL.object
}

func (cc *CopyCommand) skipUpload(bucket *oss.Bucket, objectName string, destURL CloudURL, cpOption copyOptionType, srct time.Time) (bool, error) {
	if cpOption.update {
		if props, err := cc.command.ossGetObjectStatRetry(bucket, objectName); err == nil {
			destt, err := time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified))
			if err != nil {
				return false, err
			}
			if destt.Unix() >= srct.Unix() {
				return true, nil
			}
		}
	} else {
		if !cpOption.force {
			if _, err := cc.command.ossGetObjectStatRetry(bucket, objectName); err == nil {
				if !cc.confirm(CloudURLToString(destURL.bucket, objectName)) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (cc *CopyCommand) confirm(str string) bool {
	mu.Lock()
	defer mu.Unlock()

	var val string
	fmt.Printf("\rcp: overwrite \"%s\"(y or n)? ", str)
	if _, err := fmt.Scanln(&val); err != nil || (val != "yes" && val != "y") {
		return false
	}
	return true
}

func (cc *CopyCommand) ossPutObjectRetry(bucket *oss.Bucket, objectName string, content string) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.PutObject(objectName, strings.NewReader(content))
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, objectName}
		}
	}
}

func (cc *CopyCommand) ossUploadFileRetry(bucket *oss.Bucket, objectName string, filePath string) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.PutObjectFromFile(objectName, filePath)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return FileError{err, filePath}
		}
	}
}

func (cc *CopyCommand) preparePartOption(fileSize int64) (int64, int) {
	partSize := int64(math.Ceil(float64(fileSize) / float64(MaxPartNum)))
	if partSize < oss.MinPartSize {
		partSize = oss.MinPartSize
	}
	partNum := (fileSize - 1) / partSize + 1

	for partNum > MaxIdealPartNum && partSize < MaxIdealPartSize {
		partNum /= 5
		partSize = int64(math.Ceil(float64(fileSize) / float64(partNum)))
	}

	for partSize < MinIdealPartSize && partNum > MinIdealPartNum {
		partSize *= 5
		partNum = (fileSize-1)/partSize + 1
	}

	if partNum < 3 {
		return partSize, 1
	}
	if partNum <= 10 {
		return partSize, 2
	}
	if partNum <= 500 {
		return partSize, 5
	}
	return partSize, 10
}

func (cc *CopyCommand) formatCPFileName(cpDir, srcf, destf string) string {
	return cpDir + string(os.PathSeparator) + srcf + CheckpointSep + destf + ".cp"
}

func (cc *CopyCommand) ossResumeUploadRetry(bucket *oss.Bucket, objectName string, filePath string, partSize int64, options ...oss.Option) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.UploadFile(objectName, filePath, partSize, options...)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return FileError{err, filePath}
		}
	}
}

//function for download files
func (cc *CopyCommand) downloadFiles(srcURL CloudURL, destURL FileURL, cpOption copyOptionType) error {
	bucket, err := cc.command.ossBucket(srcURL.bucket)
	if err != nil {
		return err
	}

	filePath, err := cc.adjustDestURLForDownload(destURL, cpOption)
	if err != nil {
		return err
	}

	if !cpOption.recursive {
		if srcURL.object == "" {
			return fmt.Errorf("copy object invalid url: %s, object empty. If you mean batch copy objects, please use --recursive option", srcURL.ToString())
		}

		_, err := cc.downloadSingleFile(bucket, objectInfoType{srcURL.object, -1, time.Now()}, filePath, cpOption)
		return err
	}
	return cc.batchDownloadFiles(bucket, srcURL, filePath, cpOption)
}

func (cc *CopyCommand) adjustDestURLForDownload(destURL FileURL, cpOption copyOptionType) (string, error) {
	filePath := destURL.ToString()

	isDir := false
	if f, err := os.Stat(filePath); err == nil {
		isDir = f.IsDir()
	}

	if cpOption.recursive || isDir {
		if !strings.HasSuffix(filePath, "/") && !strings.HasSuffix(filePath, "\\") {
			filePath += "/"
		}
	}
	if strings.HasSuffix(filePath, "/") || strings.HasSuffix(filePath, "\\") {
		if err := os.MkdirAll(filePath, 0777); err != nil {
			return filePath, err
		}
	}
	return filePath, nil
}

func (cc *CopyCommand) downloadSingleFile(bucket *oss.Bucket, objectInfo objectInfoType, filePath string, cpOption copyOptionType) (bool, error) {
	//make file name
	fileName := cc.makeFileName(objectInfo.key, filePath)

	//get object size and last modify time
	object := objectInfo.key
	size := objectInfo.size
	srct := objectInfo.lastModified
	if size < 0 {
		props, err := cc.command.ossGetObjectStatRetry(bucket, object)
		if err != nil {
			return false, fmt.Errorf("%s, object: %s", err.Error(), object)
		}
		size, err = strconv.ParseInt(props.Get(oss.HTTPHeaderContentLength), 10, 64)
		if err != nil {
			return false, err
		}
		if srct, err = time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified)); err != nil {
			return false, err
		}
	}

	if cc.skipDownload(fileName, cpOption, srct) {
		return true, nil
	}

	if size == 0 && (strings.HasSuffix(object, "/") || strings.HasSuffix(object, "\\")) {
		return false, os.MkdirAll(fileName, 0777)
	}

	//create parent directory
	if err := cc.createParentDirectory(fileName); err != nil {
		return false, err
	}

	if size < cpOption.threshold {
		return false, cc.command.ossDownloadFileRetry(bucket, object, fileName)
	}

	partSize, rt := cc.preparePartOption(size)
	cp := oss.Checkpoint(true, cc.formatCPFileName(cpOption.cpDir, object, filePath))
	return false, cc.ossResumeDownloadRetry(bucket, object, fileName, size, partSize, oss.Routines(rt), cp)
}

func (cc *CopyCommand) makeFileName(object, filePath string) string {
	if strings.HasSuffix(filePath, "/") || strings.HasSuffix(filePath, "\\") {
		return filePath + object
	}
	return filePath
}

func (cc *CopyCommand) skipDownload(fileName string, cpOption copyOptionType, srct time.Time) bool {
	if cpOption.update {
		if f, err := os.Stat(fileName); err == nil {
			destt := f.ModTime()
			if destt.Unix() >= srct.Unix() {
				return true
			}
		}
	} else {
		if !cpOption.force {
			if _, err := os.Stat(fileName); err == nil {
				if !cc.confirm(fileName) {
					return true
				}
			}
		}
	}
	return false
}

func (cc *CopyCommand) createParentDirectory(fileName string) error {
	dir, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		return err
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	return os.MkdirAll(dir, 0777)
}

func (cc *CopyCommand) ossResumeDownloadRetry(bucket *oss.Bucket, objectName string, filePath string, size, partSize int64, options ...oss.Option) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.DownloadFile(objectName, filePath, partSize, options...)
		if err == nil {
			return cc.truncateFile(filePath, size) 
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, objectName}
		}
	}
}

func (cc *CopyCommand) truncateFile(filePath string, size int64) error {
    f, err := os.Stat(filePath)
    if err != nil {
        return err
    }
    if f.Size() > size {
        return os.Truncate(filePath, size)
    }
    return nil
}

func (cc *CopyCommand) batchDownloadFiles(bucket *oss.Bucket, srcURL CloudURL, filePath string, cpOption copyOptionType) error {
	chObjects := make(chan objectInfoType, ChannelBuf)
	chFinishObjects := make(chan string, ChannelBuf)
	chSkipObjects := make(chan string, ChannelBuf)
	chError := make(chan error, cpOption.routines+1)
	go cc.objectProducer(bucket, srcURL, chObjects, chError)
	for i := 0; int64(i) < cpOption.routines; i++ {
		go cc.downloadConsumer(bucket, filePath, cpOption, chObjects, chFinishObjects, chSkipObjects, chError)
	}

	completed := 0
	num := 0
	snum := 0
	for int64(completed) <= cpOption.routines {
		select {
		case <-chFinishObjects:
			num++
			cc.schedule(cpOption, fmt.Sprintf("\rdownload %d objects, skip %d objects...", num, snum))
		case <-chSkipObjects:
			snum++
			cc.schedule(cpOption, fmt.Sprintf("\rdownload %d objects, skip %d objects...", num, snum))
		case err := <-chError:
			if err != nil {
				fmt.Printf("\rdownload %d objects, skip %d objects, when error happens.\n", num, snum)
				return err
			}
			completed++
		}
	}
	fmt.Printf("\rSucceed: scaned %d objects, download %d objects, skip %d objects.\n", num+snum, num, snum)
	return nil
}

func (cc *CopyCommand) objectProducer(bucket *oss.Bucket, cloudURL CloudURL, chObjects chan<- objectInfoType, chError chan<- error) {
	pre := oss.Prefix(cloudURL.object)
	marker := oss.Marker("")
	for i := 0; ; i++ {
		lor, err := cc.command.ossListObjectsRetry(bucket, marker, pre)
		if err != nil {
			chError <- err
			break
		}

		for _, object := range lor.Objects {
			chObjects <- objectInfoType{object.Key, int64(object.Size), object.LastModified}
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

func (cc *CopyCommand) downloadConsumer(bucket *oss.Bucket, filePath string, cpOption copyOptionType, chObjects <-chan objectInfoType, chFinishObjects, chSkipObjects chan<- string, chError chan<- error) {
	for objectInfo := range chObjects {
		skip, err := cc.downloadSingleFile(bucket, objectInfo, filePath, cpOption)
		if err != nil {
			chError <- err
			return
		}
		if skip {
			chSkipObjects <- objectInfo.key
		} else {
			chFinishObjects <- objectInfo.key
		}
	}

	chError <- nil
}

func (cc *CopyCommand) copyFiles(srcURL, destURL CloudURL, cpOption copyOptionType) error {
	bucket, err := cc.command.ossBucket(srcURL.bucket)
	if err != nil {
		return err
	}

	if err := cc.checkCopyFileArgs(srcURL, destURL, cpOption); err != nil {
		return err
	}

	if !cpOption.recursive {
		if srcURL.object == "" {
			return fmt.Errorf("copy object invalid url: %s, object empty. If you mean batch copy objects, please use --recursive option", srcURL.ToString())
		}

		_, err := cc.copySingleFile(bucket, objectInfoType{srcURL.object, -1, time.Now()}, srcURL, destURL, cpOption)
		return err
	}
	return cc.batchCopyFiles(bucket, srcURL, destURL, cpOption)
}

func (cc *CopyCommand) checkCopyFileArgs(srcURL, destURL CloudURL, cpOption copyOptionType) error {
	if srcURL.bucket != destURL.bucket {
		return nil
	}
	srcPrefix := srcURL.object
	destPrefix := destURL.object
	if srcPrefix == destPrefix {
		return fmt.Errorf("\"%s\" and \"%s\" are the same, copy self will do nothing, set meta please use setmeta command", srcURL.ToString(), srcURL.ToString())
	}
	if cpOption.recursive {
		if strings.HasPrefix(destPrefix, srcPrefix) {
			return fmt.Errorf("\"%s\" include \"%s\", it's not allowed, recursivlly copy should be avoided", destURL.ToString(), srcURL.ToString())
		}
		if strings.HasPrefix(srcPrefix, destPrefix) {
			return fmt.Errorf("\"%s\" include \"%s\", it's not allowed, recover source object should be avoided", srcURL.ToString(), destURL.ToString())
		}
	}
	return nil
}

func (cc *CopyCommand) copySingleFile(bucket *oss.Bucket, objectInfo objectInfoType, srcURL, destURL CloudURL, cpOption copyOptionType) (bool, error) {
	//make object name
	srcObject := objectInfo.key
	destObject := cc.makeCopyObjectName(objectInfo.key, srcURL.object, destURL, cpOption)

	if srcURL.bucket == destURL.bucket && srcObject == destObject {
		return false, fmt.Errorf("\"%s\" and \"%s\" are the same, copy self will do nothing, set meta please use setmeta command", CloudURLToString(srcURL.bucket, srcObject), CloudURLToString(srcURL.bucket, srcObject))
	}

	//get object size
	size := objectInfo.size
	srct := objectInfo.lastModified
	if size < 0 {
		props, err := cc.command.ossGetObjectStatRetry(bucket, srcObject)
		if err != nil {
			return false, fmt.Errorf("%s, object: %s", err.Error(), srcObject)
		}
		size, err = strconv.ParseInt(props.Get(oss.HTTPHeaderContentLength), 10, 64)
		if err != nil {
			return false, err
		}
		if srct, err = time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified)); err != nil {
			return false, err
		}
	}

	if skip, err := cc.skipCopy(destURL, destObject, cpOption, srct); err != nil || skip {
		return skip, err
	}

	if size < cpOption.threshold {
		return false, cc.ossCopyObjectRetry(bucket, srcObject, destURL.bucket, destObject)
	}

	partSize, rt := cc.preparePartOption(size)
	cp := oss.Checkpoint(true, cc.formatCPFileName(cpOption.cpDir, srcURL.bucket + "-" + srcObject, destURL.bucket + "-" + destObject))
	return false, cc.ossResumeCopyRetry(srcURL.bucket, srcObject, destURL.bucket, destObject, partSize, oss.Routines(rt), cp)
}

func (cc *CopyCommand) makeCopyObjectName(srcObject, srcPrefix string, destURL CloudURL, cpOption copyOptionType) string {
	if !cpOption.recursive {
		if destURL.object == "" || strings.HasSuffix(destURL.object, "/") || strings.HasSuffix(destURL.object, "\\") {
			pos := strings.LastIndex(srcObject, "/")
			pos1 := strings.LastIndex(srcObject, "\\")
			pos = int(math.Max(float64(pos), float64(pos1)))
			if pos > 0 {
				srcObject = srcObject[pos+1:]
			}
			return destURL.object + srcObject
		}
		return destURL.object
	}
	return destURL.object + srcObject[len(srcPrefix):]
}

func (cc *CopyCommand) skipCopy(destURL CloudURL, destObject string, cpOption copyOptionType, srct time.Time) (bool, error) {
	destBucket, err := cc.command.ossBucket(destURL.bucket)
	if err != nil {
		return false, err
	}

	if cpOption.update {
		if props, err := cc.command.ossGetObjectStatRetry(destBucket, destObject); err == nil {
			destt, err := time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified))
			if err != nil {
				return false, err
			}
			if destt.Unix() >= srct.Unix() {
				return true, nil
			}
		}

	} else {
		if !cpOption.force {
			if _, err := cc.command.ossGetObjectStatRetry(destBucket, destObject); err == nil {
				if !cc.confirm(CloudURLToString(destURL.bucket, destObject)) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (cc *CopyCommand) ossCopyObjectRetry(bucket *oss.Bucket, objectName, destBucketName, destObjectName string) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		_, err := bucket.CopyObjectTo(destBucketName, destObjectName, objectName)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, objectName}
		}
	}
}

func (cc *CopyCommand) ossResumeCopyRetry(bucketName, objectName, destBucketName, destObjectName string, partSize int64, options ...oss.Option) error {
	bucket, err := cc.command.ossBucket(destBucketName)
	if err != nil {
		return err
	}
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
        err := bucket.CopyFile(bucketName, objectName, destObjectName, partSize, options...)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, objectName}
		}
	}
}

func (cc *CopyCommand) batchCopyFiles(bucket *oss.Bucket, srcURL, destURL CloudURL, cpOption copyOptionType) error {
	chObjects := make(chan objectInfoType, ChannelBuf)
	chFinishObjects := make(chan string, ChannelBuf)
	chSkipObjects := make(chan string, ChannelBuf)
	chError := make(chan error, cpOption.routines+1)
	go cc.objectProducer(bucket, srcURL, chObjects, chError)
	for i := 0; int64(i) < cpOption.routines; i++ {
		go cc.copyConsumer(bucket, srcURL, destURL, cpOption, chObjects, chFinishObjects, chSkipObjects, chError)
	}

	completed := 0
	num := 0
	snum := 0
	for int64(completed) <= cpOption.routines {
		select {
		case <-chFinishObjects:
			num++
			cc.schedule(cpOption, fmt.Sprintf("\rcopy %d objects, skip %d objects...", num, snum))
		case <-chSkipObjects:
			snum++
			cc.schedule(cpOption, fmt.Sprintf("\rcopy %d objects, skip %d objects...", num, snum))
		case err := <-chError:
			if err != nil {
				fmt.Printf("\rcopy %d objects, skip %d objects, when error happens.\n", num, snum)
				return err
			}
			completed++
		}
	}
	fmt.Printf("\rSucceed: scaned %d objects, copy %d objects, skip %d objects.\n", num+snum, num, snum)
	return nil
}

func (cc *CopyCommand) copyConsumer(bucket *oss.Bucket, srcURL, destURL CloudURL, cpOption copyOptionType, chObjects <-chan objectInfoType, chFinishObjects, chSkipObjects chan<- string, chError chan<- error) {
	for objectInfo := range chObjects {
		skip, err := cc.copySingleFile(bucket, objectInfo, srcURL, destURL, cpOption)
		if err != nil {
			chError <- err
			return
		}
		if skip {
			chSkipObjects <- objectInfo.key
		} else {
			chFinishObjects <- objectInfo.key
		}
	}

	chError <- nil
}
