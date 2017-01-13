package lib

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type operationType int

const (
	operationTypePut operationType = iota
	operationTypeGet
	operationTypeCopy
)

const (
    opUpload    string = "upload"
    opDownload         = "download"
    opCopy             = "copy"
)

type CopyOptionType struct {
	recursive bool
	force     bool
	update    bool
	threshold int64
	cpDir     string
	routines  int64
    ctnu      bool
    reporter  *Reporter
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

var (
    mu sync.RWMutex
    chProgressSignal chan bool
)

const SignalNum = 2

// OssProgressListener progress listener
type OssProgressListener struct {
    monitor     *CPMonitor
    lastSize    int64
    currSize    int64
}

// ProgressChanged handle progress event
func (l *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
    if event.EventType == oss.TransferDataEvent {
        l.lastSize = l.currSize
        l.currSize = event.ConsumedBytes
        l.monitor.updateTransferSize(l.currSize - l.lastSize)
        if len(chProgressSignal) <= SignalNum {
            chProgressSignal <- true
        }
    }
}


var specChineseCopy = SpecText{

	synopsisText: "上传，下载或拷贝Objects",

	paramText: "src_url dest_url [options]",

	syntaxText: ` 
    ossutil cp file_url cloud_url  [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=cdir] [-c file] 
    ossutil cp cloud_url file_url  [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=cdir] [-c file] 
    ossutil cp cloud_url cloud_url [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=cdir] [-c file] 
`,

	detailHelpText: ` 
    该命令允许：从本地文件系统上传文件到oss，从oss下载object到本地文件系统，在oss
    上进行object拷贝。分别对应下述三种操作：
        ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]

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
        批量拷贝，如果拷贝失败，已经执行的拷贝不会回退。

    在进行批量文件上传（或下载、拷贝）时，如果其中某个文件操作失败，ossutil不会退出，
    而是继续进行其他文件的上传（或下载、拷贝）动作，并将出错文件的错误信息记录到report
    文件中。成功上传（或下载、拷贝）的文件信息将不会被记录到report文件中。

    注意：批量操作出错时不继续运行，而是终止运行的情况：
    （1）如果未进入批量文件迭代过程，错误已经发生，则不会产生report文件，ossutil会终止
    运行，不继续迭代过程。如，用户输入cp命令出错时，不会产生report文件，而是屏幕输出错
    误并退出。
    （2）如果批量操作过程某文件发生的错误为：Bucket不存在、accessKeyID/accessKeySecret
    错误造成的权限验证非法等错误，ossutil会屏幕输出错误并退出。

    report文件名为：` + ReportPrefix + `日期_时间` + ReportSuffix + `。report文件是ossutil输出文件的一种，
    被放置在ossutil的输出目录下，该目录的路径可以用配置文件中的outputDir选项或命令行
    --output-dir选项指定，如果未指定，会使用默认的输出目录：当前目录下的` + DefaultOutputDir + `目录。

    注意：ossutil不做report文件的维护工作，请自行查看及清理您的report文件，避免产生过多的
    report文件。

--update选项（-u）
    
    如果指定了该选项，ossutil只有当目标文件（或object）不存在，或源文件（或object）新于
    目标文件（或object）时，才执行拷贝。当指定了该选项时，无论--force选项是否指定了，在
    目标文件存在时，ossutil都不会提示，直接采取上述策略。
    该选项可用于当批量拷贝失败时，重传时跳过已经成功的文件。

--output-dir选项
    
    该选项指定ossutil输出文件存放的目录，默认为：当前目录下的` + DefaultOutputDir + `目录。如果指定
    的目录不存在，ossutil会自动创建该目录，如果您指定的路径已存在并且不是目录，会报错。
    输出文件表示ossutil在运行过程中产生的输出文件，目前包含：在cp命令中ossutil运行出错时
    产生的report文件。
    
--force选项

    如果dest_url指定的文件或objects已经存在，并且未指定--update选项，ossutil会询问是否进
    行替换操作（输入非法时默认不替换），如果指定了--force选项，则不询问，强制替换。该选项
    只有在未指定--update选项时有效，否则按--update选项操作。


大文件断点续传：

    如果源文件大小超过--bigfile-threshold选项指定的大小（默认为100M），ossutil会认为该文件
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
    被上传到oss。该目录路径可以用--checkpoint-dir选项指定，如果指定了该选项，请确保指定的目录
    可以被删除。


批量文件迁移：

    ossutil支持通过本地文件系统中转的方式进行批量文件迁移，在这种场景下，通常的使用方式是：
    （1）批量上传：
        ossutil cp your_dir oss://your_bucket -r -u -f
    （2）批量下载：
        ossutil cp oss://your_bucket your_dir -r -u -f
    （3）同region的Bucket间迁移：
        ossutil cp oss://your_srcbucket oss://your_destbucket -r -u -f

    具体每个选项的意义，请见上文帮助。
    在运行完一轮文件迁移后，请根据屏幕提示查看report文件，处理出错文件。

用法：

    该命令有三种用法：

    1) ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [--update] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
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

    2) ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [--update] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        该用法下载oss上的单个或多个Object到本地文件系统。如果未指定--recursive选项，则ossutil
    认为src_url精确指定了待拷贝的单个object，此时不支持prefix匹配，如果object不存在则报错。如
    果指定了--recursive选项，ossutil会搜索prefix匹配的objects，批量拷贝这些objects，此时file_url
    必须为目录，如果该目录不存在，ossutil自动创建该目录。
    文件命名规则：
        当file_url为文件时，下载到file_url指定的文件，文件名与file_url保持一致。
        当file_url为目录时，下载到file_url指定的目录中，文件名为：object名称，当object名称中含有/或\\时，会创建相应子目录。
    注意：对于以/结尾且大小为0的object，会在本地文件系统创建一个目录，而不是尝试创建一个文件。
    对于其他object会尝试创建文件。

    3) ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [--update] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        该用法在oss间进行object的拷贝。其中src_bucket与dest_bucket可以相同，注意，当src_url与
    dest_url完全相同时，ossutil不会做任何事情，直接提示退出。设置meta请使用set-meta命令。如果未
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

    ossutil cp local_dir oss://bucket1/b -r
    如果某文件上传发生服务器内部错误等失败，会在当前目录下的ossutil_output目录中产生report文件记录错误信息，并尝试其他文件的上传操作。

    ossutil cp local_dir oss://bucket1/b -r --output-dir=your_dir 
    如果某文件上传发生服务器内部错误等失败，会在your_dir中产生report文件记录错误信息，并尝试其他文件的上传操作。

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

    ossutil cp oss://bucket/abcdir1/a b --update
    如果文件b已存在，且更新时间不晚于oss://bucket/abcdir1/a，则跳过本次操作。    

    ossutil cp oss://bucket/abcdir1/a b/
    在目录b下生成文件a

    ossutil cp oss://bucket/abcdir2/a/ b
    如果b为已存在文件，报错。
    如果b为已存在目录，在目录b下生成目录a。

    ossutil cp oss://bucket/abc b
    报错，object不存在。

    ossutil cp oss://bucket/abc b -r
    如果b为已存在文件，报错    
    否则在目录b下生成目录abcdir1和abcdir2，
        目录abcdir1中生成文件a，b，c
        目录abcdir2中生成目录a和b，目录b中生成文件e

    ossutil cp oss://bucket/ local_dir -r
    如果某文件下载发生服务器内部错误等失败，会在当前目录下的ossutil_output目录中产生report文件记录错误信息，并尝试其他文件的下载操作。

    ossutil cp oss://bucket/ local_dir -r --output-dir=your_dir 
    如果某文件下载发生服务器内部错误等失败，会在your_dir中产生report文件记录错误信息，并尝试其他文件的下载操作。
        
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

    ossutil cp oss://bucket/abcdir1/a oss://bucket1/ -r
    报错，因为此时目标object名称为空，非法

    ossutil cp oss://bucket/ oss://bucket1/ -r
    如果某文件拷贝发生服务器内部错误等失败，会在当前目录下的ossutil_output目录中产生report文件记录错误信息，并尝试其他文件的拷贝操作。

    ossutil cp oss://bucket/ oss://bucket1/ -r --output-dir=your_dir 
    如果某文件拷贝发生服务器内部错误等失败，会在your_dir中产生report文件记录错误文件的信息，并尝试其他文件的拷贝操作。
`,
}

var specEnglishCopy = SpecText{

	synopsisText: "Upload, Download or Copy Objects",

	paramText: "src_url dest_url [options]",

	syntaxText: ` 
    ossutil cp file_url cloud_url  [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=cdir] [-c file] 
    ossutil cp cloud_url file_url  [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=cdir] [-c file] 
    ossutil cp cloud_url cloud_url [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=cdir] [-c file] 
`,

	detailHelpText: ` 
    The command allows: 
    1. Upload file from local file system to oss 
    2. Download object from oss to local file system
    3. Copy objects between oss
    Which matches with the following three kinds of operations:
        ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [-u] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]

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
    download or copy these objects. If error occurs, objects which has been download or copyed 
    will not rollback. 

    By default, if an error occurs to a file in batch upload(/download/copy) files operation, 
    ossutil will continue to attempt to copy the remaining files, and ossutil will record the 
    error message to report file. The files succeed copied will not be recorded to report file.   

    Note: Ossutil will print error information and exit, instead of continue to run if an error 
    occurs in batch upload(/download/copy) files operation in several situations:
    (1) If the error occurs before of entering the upload(/download/copy) iteration, ossutil will 
        print error message and return, and the report file will not be generated. eg. user enter 
        an invalid cp command.
    (2) If the error occurs during upload(/download/copy) iteration is: NoSuchBucket, AccessDenied 
        caused by unauthorized authentication and other errors. ossutil will print error message 
        and return, the report file that has been generated will not be deleted.

    Report file name is: ` + ReportPrefix + `Date_Time` + ReportSuffix + `. Report file is one kind 
    of output files, and will be putted in output directory, the directory can be specified by 
    --output-dir option or outputDir option in config file. If it's not specified, ossutil will use 
    the default directory: ` + DefaultOutputDir + ` in current directory.

    Note: ossutil will not mainten the report file, please check and clear your output directory 
    regularlly to avoid too many report files in your output directory. 

--update option

    Use the --update option to copy only when the source file is newer than the destination file 
    when the destination file is missing. If --update option is specified, when the destionation 
    file is existed, ossutil will not prompt and copy when newer, no matter if --force option is 
    specified or not.
    The option can be used when batch copy failed, skip the succeed files in retry.

--output-dir option
    
    The option specify the directory to deposit output file generated by ossutil, the default value 
    is: ` + DefaultOutputDir + ` in current directory. If the directory specified not exist, ossutil will 
    create the directory automatically, if it exists but is not a directory, ossutil will return an 
    error.  

    Output file contains: report file which used to record error message generated by cp command.

--force option

    If the file dest_url specified is existed, and --update option is not specified, ossutil will 
    ask if replace the file(if the input is invalid, the file will not be replaced). If --force 
    option is specified here, ossutil will not prompt, replace by force. The option is useful only 
    when --update not specified. 


Resume copy of big file:

    If the size of source file is bigger than what --bigfile-threshold option specified(default: 
    100M), ossutil will consider the file as a big file, and use resume copy policy to these files:
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
        will not be uploaded to oss. The path of checkpoint directory can be specified by --checkpoint-dir 
        option, please make sure the directory you specified can be removed.


Batch file migration:

    ossutil support batch file migration by transfer files through local file system, the usual usage is: 
    (1) Batch file upload:
        ossutil cp your_dir oss://your_bucket -r -u -f
    (2) Batch file download:
        ossutil cp oss://your_bucket your_dir -r -u -f
    (3) File copy between buckets in the same region：
        ossutil cp oss://your_srcbucket oss://your_destbucket -r -u -f

    The meaning of every option, see help above.
    After each migration, please check your report file.


Usage:

    There are three usages:

    1) ossutil cp file_url oss://bucket[/prefix] [-r] [-f] [--update] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
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
        
    2) ossutil cp oss://bucket[/prefix] file_url [-r] [-f] [--update] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        The usage download one or many objects to local system. If --recursive option is not specified, 
    ossutil considers src_url exactly specified the single object you want to download, prefix-matching 
    is not supported now, if the object not exists, error occurs. If --recursive option is specified, 
    ossutil will search for prefix-matching objects and batch download those objects, at this time file_url 
    must be directory, if the directory not exists, ossutil will create the directory automatically.
    File Naming Rules:
        If file_url is file, ossutil download file to the path of file_url, and the file name is got from file_url.
        If file_url is directory, ossutil download file to the directory, and the file name is: object name.
    Warning: If the object name is end with / and size is zero, ossutil will create a directory in local 
    system, instead of creating a file.

    3) ossutil cp oss://src_bucket[/src_prefix] oss://dest_bucket[/dest_prefix] [-r] [-f] [--update] [--output-dir=odir] [--bigfile-threshold=size] [--checkpoint-dir=file]
        The usage copy objects between oss. The src_bucket can be same with dest_bucket. Pay attention 
    please, if src_url is the same with dest_url, ossutil will do nothing but exit after prompt. Set meta 
    please use "set-meta" command. If --recursive option is not specified, ossutil considers src_url exactly 
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

    ossutil cp local_dir oss://bucket1/b -r
    If an 5xx error occurs while upload a file, ossutil will generate a report file and record the error 
    information to the file, and store the file in ossutil_output directory in current path, and continue 
    to upload the remaining files.

    ossutil cp local_dir oss://bucket1/b -r --output-dir=your_dir 
    If an 5xx error occurs while upload a file, ossutil will generate a report file and record the error 
    information to the file, and store the file in your_dir, and continue to upload the remaining files.

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
        
    ossutil cp oss://bucket/ local_dir -r
    If an 5xx error occurs while download a file, ossutil will generate a report file and record the error 
    information to the file, and store the file in ossutil_output directory in current path, and continue 
    to download the remaining files.

    ossutil cp oss://bucket/ local_dir -r --output-dir=your_dir
    If an 5xx error occurs while download a file, ossutil will generate a report file and record the error 
    information to the file, and store the file in your_dir, and download to upload the remaining files.

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

    ossutil cp oss://bucket/ oss://bucket1/ -r
    If an 5xx error occurs while copy a file, ossutil will generate a report file and record the error 
    information to the file, and store the file in ossutil_output directory in current path, and continue 
    to copy the remaining files.

    ossutil cp oss://bucket/ oss://bucket1/ -r --output-dir=your_dir 
    If an 5xx error occurs while copy a file, ossutil will generate a report file and record the error 
    information to the file, and store the file in your_dir, and continue to copy the remaining files.
`,
}

// CopyCommand is the command upload, download and copy objects
type CopyCommand struct {
	command     Command
    cpOption    CopyOptionType
    monitor     CPMonitor
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
			OptionRecursion,
			OptionForce,
			OptionUpdate,
            OptionContinue,
            OptionOutputDir,
			OptionBigFileThreshold,
			OptionCheckpointDir,
			OptionConfigFile,
            OptionEndpoint,
            OptionAccessKeyID,
            OptionAccessKeySecret,
            OptionSTSToken,
			OptionRetryTimes,
			OptionRoutines,
            OptionParallel,
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
	cc.cpOption.recursive, _ = GetBool(OptionRecursion, cc.command.options)
	cc.cpOption.force, _ = GetBool(OptionForce, cc.command.options)
	cc.cpOption.update, _ = GetBool(OptionUpdate, cc.command.options)
	cc.cpOption.threshold, _ = GetInt(OptionBigFileThreshold, cc.command.options)
	cc.cpOption.cpDir, _ = GetString(OptionCheckpointDir, cc.command.options)
	cc.cpOption.routines, _ = GetInt(OptionRoutines, cc.command.options)
    cc.cpOption.ctnu = false 
    if cc.cpOption.recursive {
        cc.cpOption.ctnu = true
    }
    outputDir, _ := GetString(OptionOutputDir, cc.command.options)

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

	// create ckeckpoint dir
	if err := os.MkdirAll(cc.cpOption.cpDir, 0755); err != nil {
		return err
	}

    // init reporter 
    if cc.cpOption.reporter, err = GetReporter(cc.cpOption.ctnu, outputDir, commandLine); err != nil {
        return err
    }

    cc.monitor.init(opType)

    chProgressSignal = make(chan bool, 10)
    go cc.progressBar()

	switch opType {
	case operationTypePut:
		err = cc.uploadFiles(srcURLList, destURL.(CloudURL))
	case operationTypeGet:
		err = cc.downloadFiles(srcURLList[0].(CloudURL), destURL.(FileURL))
	default:
		err = cc.copyFiles(srcURLList[0].(CloudURL), destURL.(CloudURL))
	}

    cc.cpOption.reporter.Clear()

	if err == nil {
		return os.RemoveAll(cc.cpOption.cpDir)
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

func (cc *CopyCommand) progressBar() {
    // fetch all reveal
    for _ = range chProgressSignal {
        fmt.Printf(cc.monitor.progressBar(false, normalExit))
    }
}

//function for upload files
func (cc *CopyCommand) uploadFiles(srcURLList []StorageURLer, destURL CloudURL) error {
    if err := destURL.checkObjectPrefix(); err != nil {
        return err
    }

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
	chError := make(chan error, cc.cpOption.routines)
	chListError := make(chan error, 1)
    go cc.fileStatistic(srcURLList)
	go cc.fileProducer(srcURLList, chFiles, chListError)
	for i := 0; int64(i) < cc.cpOption.routines; i++ {
		go cc.uploadConsumer(bucket, destURL, chFiles, chError)
	}

	completed := 0
	for int64(completed) <= cc.cpOption.routines {
		select {
        case err := <-chListError:
            if err != nil {
                return err
            }
            completed++
		case err := <-chError:
            if err == nil {
                completed++
            } else {
                if !cc.cpOption.ctnu {
                    fmt.Printf(cc.monitor.progressBar(true, errExit))
                    return err
                }
            }
		}
	}
    fmt.Printf(cc.monitor.progressBar(true, normalExit))
	return nil
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

func (cc *CopyCommand) fileStatistic(srcURLList []StorageURLer) {
	for _, url := range srcURLList {
		name := url.ToString()
		f, err := os.Stat(name)
		if err != nil {
            cc.monitor.setScanError(err)
			return
		}
		if f.IsDir() {
			if !cc.cpOption.recursive {
                cc.monitor.setScanError(fmt.Errorf("omitting directory \"%s\", please use --recursive option", name))
				return
			}
            err := cc.getFileListStatistic(name)
            if err != nil {
                cc.monitor.setScanError(err)
                return
            }
		} else {	
            if cc.filterPath(name, cc.cpOption.cpDir) {
                cc.monitor.updateScanSizeNum(f.Size(), 1)
            }
        }
    }

    cc.monitor.setScanEnd()
}

func (cc *CopyCommand) getFileListStatistic(dpath string) error {
    err := filepath.Walk(dpath, func(fpath string, f os.FileInfo, err error) error {
        if f == nil {
            return err
        }

        if !cc.filterPath(fpath, cc.cpOption.cpDir) {
            return nil
        }

        dpath = filepath.Clean(dpath)
        fpath = filepath.Clean(fpath)
        _, err = filepath.Rel(dpath, fpath) 
        if err != nil {
            return fmt.Errorf("list file error: %s, info: %s", fpath, err.Error())
        }

        if f.IsDir(){
            if fpath != dpath {
                cc.monitor.updateScanNum(1)
            }
            return nil
        }
        cc.monitor.updateScanSizeNum(f.Size(), 1)
        return nil
    })
    return err
}

func (cc *CopyCommand) fileProducer(srcURLList []StorageURLer, chFiles chan<- fileInfoType, chListError chan<- error) {
	for _, url := range srcURLList {
		name := url.ToString()
		f, err := os.Stat(name)
		if err != nil {
			chListError <- err
			return
		}
		if f.IsDir() {
			if !cc.cpOption.recursive {
				chListError <- fmt.Errorf("omitting directory \"%s\", please use --recursive option", name)
				return
			}
            err := cc.getFileList(name, chFiles)
            if err != nil {
                chListError <- err
                return
            }
		} else {	
            dir, fname := filepath.Split(name) 
		    chFiles <- fileInfoType{fname, dir}
        }
    }

	defer close(chFiles)
	chListError <- nil
}

func (cc *CopyCommand) getFileList(dpath string, chFiles chan<- fileInfoType) error {
    name := dpath
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
                if strings.HasSuffix(fileName, "\\") || strings.HasSuffix(fileName, "/") { 
                    chFiles <- fileInfoType{fileName, name}
                } else {
                    chFiles <- fileInfoType{fileName + string(os.PathSeparator), name}
                }
            }
            return nil
        }
        chFiles <- fileInfoType{fileName, name}
        return nil
    })
    return err
}

func (cc *CopyCommand) uploadConsumer(bucket *oss.Bucket, destURL CloudURL, chFiles <-chan fileInfoType, chError chan<- error) {
	for file := range chFiles {
		if cc.filterFile(file, cc.cpOption.cpDir) {
			err := cc.uploadFileWithReport(bucket, destURL, file)
			if err != nil {
				chError <- err
                if !cc.cpOption.ctnu {
				    return
                }
                continue
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
    return cc.filterPath(filePath, cpDir)
}

func (cc *CopyCommand) filterPath(filePath string, cpDir string) bool {
	if !strings.Contains(filePath, cpDir) {
		return true
	}
	absFile, _ := filepath.Abs(filePath)
	absCPDir, _ := filepath.Abs(cpDir)
	return !strings.Contains(absFile, absCPDir)
}

func (cc *CopyCommand) uploadFileWithReport(bucket *oss.Bucket, destURL CloudURL, file fileInfoType) error {
    skip, err, isDir, size, msg := cc.uploadFile(bucket, destURL, file) 
    cc.report(msg, err)
    cc.updateMonitor(skip, err, isDir, size)
    return err
}

func (cc *CopyCommand) uploadFile(bucket *oss.Bucket, destURL CloudURL, file fileInfoType) (skip bool, rerr error, isDir bool, size int64, msg string) {
	//first make object name
	objectName := cc.makeObjectName(destURL, file)

	filePath := file.filePath
	if file.dir != "" {
        if strings.HasSuffix(file.dir, "/") || strings.HasSuffix(file.dir, "\\") { 
		    filePath = file.dir + file.filePath
        } else {
		    filePath = file.dir + string(os.PathSeparator) + file.filePath
        }
	}

    skip = false
    rerr = nil
    isDir = false
    size = 0        // the size update to monitor
    msg = fmt.Sprintf("%s %s to %s", opUpload, filePath, CloudURLToString(bucket.BucketName, objectName))

	//get file size and last modify time
	f, err := os.Stat(filePath)
	if err != nil {
        rerr = err
		return
	}

    if !f.IsDir() {
        size = f.Size()
    }

	if skip, rerr = cc.skipUpload(bucket, objectName, destURL, f.ModTime()); rerr != nil || skip {
		return
	}

    skip = false
	if f.IsDir() {
        rerr = cc.ossPutObjectRetry(bucket, objectName, "")
        isDir = true
		return
	}

    size = 0
    var listener *OssProgressListener = &OssProgressListener{&cc.monitor, 0, 0}
	//decide whether to use resume upload
	if f.Size() < cc.cpOption.threshold {
        rerr = cc.ossUploadFileRetry(bucket, objectName, filePath, oss.Progress(listener))
		return
	}

	//make options for resume multipart upload
	//part size
	partSize, rt := cc.preparePartOption(f.Size())
	//checkpoint file
	cp := oss.Checkpoint(true, cc.formatCPFileName(cc.cpOption.cpDir, filePath, objectName))
    rerr = cc.ossResumeUploadRetry(bucket, objectName, filePath, partSize, oss.Routines(rt), cp, oss.Progress(listener))
	return
}

func (cc *CopyCommand) makeObjectName(destURL CloudURL, file fileInfoType) string {
	if destURL.object == "" || strings.HasSuffix(destURL.object, "/") || strings.HasSuffix(destURL.object, "\\") || strings.HasSuffix(destURL.object, string(os.PathSeparator)){
        // replace "\" of file.filePath to "/"
        filePath := file.filePath
        filePath = strings.Replace(file.filePath, string(os.PathSeparator), "/", -1)
        filePath = strings.Replace(file.filePath, "\\", "/", -1)
		return destURL.object + filePath
	}
	return destURL.object
}

func (cc *CopyCommand) skipUpload(bucket *oss.Bucket, objectName string, destURL CloudURL, srct time.Time) (bool, error) {
	if cc.cpOption.update {
		if props, err := cc.command.ossGetObjectMetaRetry(bucket, objectName); err == nil {
			destt, err := time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified))
			if err != nil {
				return false, err
			}
			if destt.Unix() >= srct.Unix() {
				return true, nil
			}
		}
	} else {
		if !cc.cpOption.force {
			if _, err := cc.command.ossGetObjectMetaRetry(bucket, objectName); err == nil {
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
    fmt.Printf(getClearStr(fmt.Sprintf("cp: overwrite \"%s\"(y or N)? ", str)))
	if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
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
			return ObjectError{err, bucket.BucketName, objectName}
		}
	}
}

func (cc *CopyCommand) ossUploadFileRetry(bucket *oss.Bucket, objectName string, filePath string, options ...oss.Option) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.PutObjectFromFile(objectName, filePath, options...)
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

    if parallel, err := GetInt(OptionParallel, cc.command.options); err == nil {
        return partSize, int(parallel)
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

func (cc *CopyCommand) updateMonitor(skip bool, err error, isDir bool, size int64) {
    if err != nil {
        cc.monitor.updateErr(0, 1)
    } else if isDir {
        cc.monitor.updateDir(size, 1)
    } else if skip {
        cc.monitor.updateSkip(size, 1)
    } else {
        cc.monitor.updateFile(size, 1)
    }
    if len(chProgressSignal) <= SignalNum {
        chProgressSignal <- true 
    }
}

func (cc *CopyCommand) report(msg string, err error) {
    if cc.filterError(err) {
        cc.cpOption.reporter.ReportError(fmt.Sprintf("%s error, info: %s", msg, err.Error()))
        cc.cpOption.reporter.Prompt(err)
    }
}

func (cc *CopyCommand) filterError(err error) bool {
    if err == nil {
        return false
    }

    switch err.(type) {
    case FileError:
        err = err.(FileError).err
    case ObjectError:
        err = err.(ObjectError).err
    case BucketError:
        err = err.(BucketError).err
    }

    switch err.(type) {
    case oss.ServiceError:
        code := err.(oss.ServiceError).Code 
        if code == "NoSuchBucket" || code == "InvalidAccessKeyId" || code == "SignatureDoesNotMatch" || code == "AccessDenied" {
            cc.cpOption.ctnu = false
            return false
        }
    case CopyError:
        cc.cpOption.ctnu = false
        return false
    }
    return true
}

//function for download files
func (cc *CopyCommand) downloadFiles(srcURL CloudURL, destURL FileURL) error {
	bucket, err := cc.command.ossBucket(srcURL.bucket)
	if err != nil {
		return err
	}

	filePath, err := cc.adjustDestURLForDownload(destURL)
	if err != nil {
		return err
	}

	if !cc.cpOption.recursive {
		if srcURL.object == "" {
			return fmt.Errorf("copy object invalid url: %s, object empty. If you mean batch copy objects, please use --recursive option", srcURL.ToString())
		}

        go cc.objectStatistic(bucket, srcURL)
		err := cc.downloadSingleFileWithReport(bucket, objectInfoType{srcURL.object, -1, time.Now()}, filePath)
        return cc.formatResultPrompt(err)
	}
	return cc.batchDownloadFiles(bucket, srcURL, filePath)
}

func (cc *CopyCommand) formatResultPrompt(err error) error {
    fmt.Printf(cc.monitor.progressBar(true, normalExit)) 
    if err != nil && cc.cpOption.ctnu {
        return nil
    }
    return err
}

func (cc *CopyCommand) adjustDestURLForDownload(destURL FileURL) (string, error) {
	filePath := destURL.ToString()

	isDir := false
	if f, err := os.Stat(filePath); err == nil {
		isDir = f.IsDir()
	}

	if cc.cpOption.recursive || isDir {
		if !strings.HasSuffix(filePath, "/") && !strings.HasSuffix(filePath, "\\") {
			filePath += "/"
		}
	}
	if strings.HasSuffix(filePath, "/") || strings.HasSuffix(filePath, "\\") {
		if err := os.MkdirAll(filePath, 0755); err != nil {
			return filePath, err
		}
	}
	return filePath, nil
}

func (cc *CopyCommand) downloadSingleFileWithReport(bucket *oss.Bucket, objectInfo objectInfoType, filePath string) error {
    skip, err, size, msg := cc.downloadSingleFile(bucket, objectInfo, filePath)
    cc.report(msg, err)
    cc.updateMonitor(skip, err, false, size)
    return err
}

func (cc *CopyCommand) downloadSingleFile(bucket *oss.Bucket, objectInfo objectInfoType, filePath string) (bool, error, int64, string) {
	//make file name
	fileName := cc.makeFileName(objectInfo.key, filePath)

	//get object size and last modify time
	object := objectInfo.key
	size := objectInfo.size
	srct := objectInfo.lastModified

    msg := fmt.Sprintf("%s %s to %s", opDownload, CloudURLToString(bucket.BucketName, object), fileName)

	if size < 0 {
		props, err := cc.command.ossGetObjectMetaRetry(bucket, object)
		if err != nil {
			return false, err, size, msg
		}
		size, err = strconv.ParseInt(props.Get(oss.HTTPHeaderContentLength), 10, 64)
		if err != nil {
			return false, err, size, msg
		}
		if srct, err = time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified)); err != nil {
			return false, err, size, msg
		}
	}

	if cc.skipDownload(fileName, srct) {
		return true, nil, size, msg
	}

	if size == 0 && (strings.HasSuffix(object, "/") || strings.HasSuffix(object, "\\")) {
		return false, os.MkdirAll(fileName, 0755), size, msg
	}

	//create parent directory
	if err := cc.createParentDirectory(fileName); err != nil {
		return false, err, size, msg
	}

    var listener *OssProgressListener = &OssProgressListener{&cc.monitor, 0, 0}
	if size < cc.cpOption.threshold {
		return false, cc.ossDownloadFileRetry(bucket, object, fileName, oss.Progress(listener)), 0, msg
	}

	partSize, rt := cc.preparePartOption(size)
	cp := oss.Checkpoint(true, cc.formatCPFileName(cc.cpOption.cpDir, object, filePath))
	return false, cc.ossResumeDownloadRetry(bucket, object, fileName, size, partSize, oss.Routines(rt), cp, oss.Progress(listener)), 0, msg
}

func (cc *CopyCommand) makeFileName(object, filePath string) string {
	if strings.HasSuffix(filePath, "/") || strings.HasSuffix(filePath, "\\") {
		return filePath + object
	}
	return filePath
}

func (cc *CopyCommand) skipDownload(fileName string, srct time.Time) bool {
	if cc.cpOption.update {
		if f, err := os.Stat(fileName); err == nil {
			destt := f.ModTime()
			if destt.Unix() >= srct.Unix() {
				return true
			}
		}
	} else {
		if !cc.cpOption.force {
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
	return os.MkdirAll(dir, 0755)
}

func (cc *CopyCommand) ossDownloadFileRetry(bucket *oss.Bucket, objectName, fileName string, options ...oss.Option) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.GetObjectToFile(objectName, fileName, options...)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, bucket.BucketName, objectName}
		}
	}
}

func (cc *CopyCommand) ossResumeDownloadRetry(bucket *oss.Bucket, objectName string, filePath string, size, partSize int64, options ...oss.Option) error {
	retryTimes, _ := GetInt(OptionRetryTimes, cc.command.options)
	for i := 1; ; i++ {
		err := bucket.DownloadFile(objectName, filePath, partSize, options...)
		if err == nil {
			return cc.truncateFile(filePath, size) 
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, bucket.BucketName, objectName}
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

func (cc *CopyCommand) batchDownloadFiles(bucket *oss.Bucket, srcURL CloudURL, filePath string) error {
	chObjects := make(chan objectInfoType, ChannelBuf)
	chError := make(chan error, cc.cpOption.routines)
    chListError := make(chan error, 1)
    go cc.objectStatistic(bucket, srcURL)
	go cc.objectProducer(bucket, srcURL, chObjects, chListError)
	for i := 0; int64(i) < cc.cpOption.routines; i++ {
		go cc.downloadConsumer(bucket, filePath, chObjects, chError)
	}

    return cc.waitRoutinueComplete(chError, chListError, opDownload)
}

func (cc *CopyCommand) objectStatistic(bucket *oss.Bucket, cloudURL CloudURL) {
    if cc.cpOption.recursive {
        pre := oss.Prefix(cloudURL.object)
        marker := oss.Marker("")
        for i := 0; ; i++ {
            lor, err := cc.command.ossListObjectsRetry(bucket, marker, pre)
            if err != nil {
                cc.monitor.setScanError(err)
                return
            }

            for _, object := range lor.Objects {
                cc.monitor.updateScanSizeNum(object.Size, 1)
            }

            pre = oss.Prefix(lor.Prefix)
            marker = oss.Marker(lor.NextMarker)
            if !lor.IsTruncated {
                break
            }
        }
    } else {
        props, err := cc.command.ossGetObjectMetaRetry(bucket, cloudURL.object)
        if err != nil {
            cc.monitor.setScanError(err)
            return
        }
        size, err := strconv.ParseInt(props.Get(oss.HTTPHeaderContentLength), 10, 64)
		if err != nil {
            cc.monitor.setScanError(err)
            return
		}
        cc.monitor.updateScanSizeNum(size, 1)
    }

    cc.monitor.setScanEnd()
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

func (cc *CopyCommand) downloadConsumer(bucket *oss.Bucket, filePath string, chObjects <-chan objectInfoType, chError chan<- error) {
	for objectInfo := range chObjects {
		err := cc.downloadSingleFileWithReport(bucket, objectInfo, filePath)
		if err != nil {
			chError <- err
            if !cc.cpOption.ctnu {
			    return
            }
            continue
		}
	}

	chError <- nil
}

func (cc *CopyCommand) waitRoutinueComplete(chError, chListError <-chan error, opStr string) error {
	completed := 0
    var ferr error 
	for int64(completed) <= cc.cpOption.routines {
		select {
        case err := <-chListError:
            if err != nil {
                return err
            }
            completed++
        case err := <-chError:
            if err == nil {
                completed++
            } else {
                ferr = err
                if !cc.cpOption.ctnu {
                    fmt.Printf(cc.monitor.progressBar(true, errExit))
                    return err
                }
            }
		}
	}
    return cc.formatResultPrompt(ferr)
}

//function for copy objects 
func (cc *CopyCommand) copyFiles(srcURL, destURL CloudURL) error {
	bucket, err := cc.command.ossBucket(srcURL.bucket)
	if err != nil {
		return err
	}

	if err := cc.checkCopyFileArgs(srcURL, destURL); err != nil {
		return err
	}

	if !cc.cpOption.recursive {
		if srcURL.object == "" {
			return fmt.Errorf("copy object invalid url: %s, object empty. If you mean batch copy objects, please use --recursive option", srcURL.ToString())
		}

        go cc.objectStatistic(bucket, srcURL)
		err := cc.copySingleFileWithReport(bucket, objectInfoType{srcURL.object, -1, time.Now()}, srcURL, destURL)
        return cc.formatResultPrompt(err)
	}
	return cc.batchCopyFiles(bucket, srcURL, destURL)
}

func (cc *CopyCommand) checkCopyFileArgs(srcURL, destURL CloudURL) error {
    if err := destURL.checkObjectPrefix(); err != nil {
        return err
    }
	if srcURL.bucket != destURL.bucket {
		return nil
	}
	srcPrefix := srcURL.object
	destPrefix := destURL.object
	if srcPrefix == destPrefix {
		return fmt.Errorf("\"%s\" and \"%s\" are the same, copy self will do nothing, set meta please use set-meta command", srcURL.ToString(), srcURL.ToString())
	}
	if cc.cpOption.recursive {
		if strings.HasPrefix(destPrefix, srcPrefix) {
			return fmt.Errorf("\"%s\" include \"%s\", it's not allowed, recursivlly copy should be avoided", destURL.ToString(), srcURL.ToString())
		}
		if strings.HasPrefix(srcPrefix, destPrefix) {
			return fmt.Errorf("\"%s\" include \"%s\", it's not allowed, recover source object should be avoided", srcURL.ToString(), destURL.ToString())
		}
	}
	return nil
}

func (cc *CopyCommand) copySingleFileWithReport(bucket *oss.Bucket, objectInfo objectInfoType, srcURL, destURL CloudURL) error {
    skip, err, size, msg := cc.copySingleFile(bucket, objectInfo, srcURL, destURL)
    cc.report(msg, err)
    cc.updateMonitor(skip, err, false, size)
    return err
}

func (cc *CopyCommand) copySingleFile(bucket *oss.Bucket, objectInfo objectInfoType, srcURL, destURL CloudURL) (bool, error, int64, string) {
	//make object name
	srcObject := objectInfo.key
	destObject := cc.makeCopyObjectName(objectInfo.key, srcURL.object, destURL)
	size := objectInfo.size
	srct := objectInfo.lastModified

    msg := fmt.Sprintf("%s %s to %s", opCopy, CloudURLToString(srcURL.bucket, srcObject), CloudURLToString(destURL.bucket, destObject))

	if srcURL.bucket == destURL.bucket && srcObject == destObject {
		return false, fmt.Errorf("\"%s\" and \"%s\" are the same, copy self will do nothing, set meta please use set-meta command", CloudURLToString(srcURL.bucket, srcObject), CloudURLToString(srcURL.bucket, srcObject)), size, msg
	}

    if destObject == "" {
        return false, CopyError{fmt.Errorf("dest object name is empty, try add a prefix to dest_url ==> change dest_url to: oss://dest_bucket/prefix, see naming rules in \"help cp\"")}, size, msg
    }

	//get object size
	if size < 0 {
		props, err := cc.command.ossGetObjectMetaRetry(bucket, srcObject)
		if err != nil {
			return false, err, size, msg
		}
		size, err = strconv.ParseInt(props.Get(oss.HTTPHeaderContentLength), 10, 64)
		if err != nil {
			return false, err, size, msg
		}
		if srct, err = time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified)); err != nil {
			return false, err, size, msg
		}
	}

	if skip, err := cc.skipCopy(destURL, destObject, srct); err != nil || skip {
		return skip, err, size, msg
	}

	if size < cc.cpOption.threshold {
		return false, cc.ossCopyObjectRetry(bucket, srcObject, destURL.bucket, destObject), size, msg
	}

    var listener *OssProgressListener = &OssProgressListener{&cc.monitor, 0, 0}
	partSize, rt := cc.preparePartOption(size)
	cp := oss.Checkpoint(true, cc.formatCPFileName(cc.cpOption.cpDir, srcURL.bucket + "-" + srcObject, destURL.bucket + "-" + destObject))
	return false, cc.ossResumeCopyRetry(srcURL.bucket, srcObject, destURL.bucket, destObject, partSize, oss.Routines(rt), cp, oss.Progress(listener)), 0, msg
}

func (cc *CopyCommand) makeCopyObjectName(srcObject, srcPrefix string, destURL CloudURL) string {
	if !cc.cpOption.recursive {
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

func (cc *CopyCommand) skipCopy(destURL CloudURL, destObject string, srct time.Time) (bool, error) {
	destBucket, err := cc.command.ossBucket(destURL.bucket)
	if err != nil {
		return false, err
	}

	if cc.cpOption.update {
		if props, err := cc.command.ossGetObjectMetaRetry(destBucket, destObject); err == nil {
			destt, err := time.Parse(http.TimeFormat, props.Get(oss.HTTPHeaderLastModified))
			if err != nil {
				return false, err
			}
			if destt.Unix() >= srct.Unix() {
				return true, nil
			}
		}

	} else {
		if !cc.cpOption.force {
			if _, err := cc.command.ossGetObjectMetaRetry(destBucket, destObject); err == nil {
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
			return ObjectError{err, bucket.BucketName, objectName}
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
			return ObjectError{err, bucket.BucketName, objectName}
		}
	}
}

func (cc *CopyCommand) batchCopyFiles(bucket *oss.Bucket, srcURL, destURL CloudURL) error {
	chObjects := make(chan objectInfoType, ChannelBuf)
	chError := make(chan error, cc.cpOption.routines)
    chListError := make(chan error, 1)
    go cc.objectStatistic(bucket, srcURL)
	go cc.objectProducer(bucket, srcURL, chObjects, chListError)
	for i := 0; int64(i) < cc.cpOption.routines; i++ {
		go cc.copyConsumer(bucket, srcURL, destURL, chObjects, chError)
	}

    return cc.waitRoutinueComplete(chError, chListError, opDownload)
}

func (cc *CopyCommand) copyConsumer(bucket *oss.Bucket, srcURL, destURL CloudURL, chObjects <-chan objectInfoType, chError chan<- error) {
	for objectInfo := range chObjects {
		err := cc.copySingleFileWithReport(bucket, objectInfo, srcURL, destURL)
		if err != nil {
			chError <- err
            if !cc.cpOption.ctnu {
			    return
            }
            continue
		}
	}

	chError <- nil
}
