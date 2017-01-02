package lib

import (
	"fmt"
    "strings"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var specChineseRemove = SpecText{

	synopsisText: "删除Bucket或Objects",

	paramText: "url [options]",

	syntaxText: ` 
    ossutil rm oss://bucket[/prefix] [-r] [-b] [-f] [-c file] 
`,

	detailHelpText: ` 
    该命令删除Bucket或objects，在某些情况下可一并删除二者。请小心使用该命令！！
    在删除objects前确定objects可以删除，在删除bucket前确定整个bucket连同其下的所有
    objects都可以删除！

        （1）删除单个object，参考用法1)
        （2）删除bucket，不删除objects，参考用法2)
        （3）删除objects，不删除bucket，参考用法3)
        （4）删除bucket和objects，参考用法4)

        对bucket进行删除，都需要添加--bucket选项。
        如果指定了--force选项，则删除前不会进行询问提示。
        
        结果：显示命令耗时前未报错，则表示成功删除。

用法：

    该命令有四种用法：

    1) ossutil rm oss://bucket/object [-m] 
        （删除单个object）
        如果未指定--recursive和--bucket选项，删除指定的单个object，此时请确保url精确指
    定了待删除的object，ossutil不会进行前缀匹配。无论是否指定--force选项，ossutil都不会
    进行询问提示。如果此时指定了--bucket选项，将会报错，单独删除bucket参考用法4)。
        如果指定--multipart选项, 删除指定的object下对应的所有uploadId，即删除这个multipart

    2) ossutil rm oss://bucket -b [-f]
        （删除bucket，不删除objects）
        如果指定了--bucket选项，未指定--recursive选项，ossutil删除指定的bucket，但并不
    去删除该bucket下的objects。此时请确保url精确匹配待删除的bucket，并且指定的bucket内
    容为空，否则会报错。如果指定了--force选项，则删除前不会进行询问提示。

    3) ossutil rm oss://bucket[/prefix] -r [-m] [-a] [-f]
        （删除objects，不删除bucket）
        如果指定了--recursive选项，未指定--bucket选项。则可以进行objects的批量删除。该
    用法查找与指定url前缀匹配的所有objects（prefix为空代表bucket下的所有objects），删除
    这些objects。由于未指定--bucket选项，则ossutil保留bucket。如果指定了--force选项，则
    删除前不会进行询问提示。
        如果指定--multipart选项, 该用法查找与指定url前缀匹配的所有multipart object（prefix
    为空代表bucket下的所有multipart object），并删除对应的所有uploadId。即删除所有符合这
    个前缀的multipart。
        如果指定--all-type, 该操作不会区分multipart和普通的object，执行删除上述multipart
    和普通object的操作。

    4) ossutil rm oss://bucket[/prefix] -r -b [-a] [-f]
        （删除bucket和objects）
        如果同时指定了--bucket和--recursive选项，ossutil进行批量删除后会尝试去一并删除
    bucket。当用户想要删除某个bucket连同其中的所有objects时，可采用该操作。如果指定了
    --force选项，则删除前不会进行询问提示。
         如果指定--all-type, 该操作不会区分multipart和普通的object，执行上述删除bucket
    和multipart object及普通object操作。

`,

	sampleText: ` 
    ossutil rm oss://bucket1/obj1
    ossutil rm oss://bucket1/dir -r 
    ossutil rm oss://bucket1 -b
    ossutil rm oss://bucket2 -r -b -f
`,
}

var specEnglishRemove = SpecText{

	synopsisText: "Remove Bucket or Objects",

	paramText: "url [options]",

	syntaxText: ` 
    ossutil rm oss://bucket[/prefix] [-r] [-b] [-f] [-c file]
`,

	detailHelpText: ` 
    The command remove bucket or objects, in some case remove both. Please use the 
    command carefully!! 
    Make sure the objects can be removed before useing the command to remove objects! 
    Make sure the bucket and objects inside can be removed before useing the command 
    to remove bucket!

    (1) Remove single object, see usage 1)
    (2) Remove bucket, don't remove objects inside, see usage 2)
    (3) Batch remove many objects, reserve bucket, see usage 3)
    (4) Remove bucket and objects inside, see usage 4)

    When remove bucket, the --bucket option must be specified.
    If --force option is specified, remove silently without asking user to confirm the 
    operation.  

    Result: if no error displayed before show elasped time, then the target is removed 
    successfully.

Usage:

    There are four usages:

    1) ossutil rm oss://bucket/object
        (Remove single object)
        If you remove without --recursive and --bucket option, ossutil remove the single 
    object specified in url. In the usage, please make sure url exactly specified the 
    object you want to remove, ossutil will not treat object as prefix and remove prefix 
    matching objects. No matter --force is specified or not, ossutil will not show prompt 
    question.

    2) ossutil rm oss://bucket -b [-f]
        (Remove bucket, don't remove objects inside)
        If you remove with --bucket option, without --recursive option, ossutil try to 
    remove the bucket, if the bucket is not empty, error occurs. In the usage, please make 
    sure url exactly specified the bucket you want to remove, or error occurs. If --force 
    option is specified, ossutil will not show prompt question. 

    3) ossutil rm oss://bucket[/prefix] -r [-f]
        (Remove objects, reserve bucket)
        If you remove with --recursive option, without --bucket option, ossutil remove all 
    the objects that prefix-matching the url you specified(empty prefix means all objects in 
    the bucket), bucket will be reserved because of missing --bucket option.

    4) ossutil rm oss://bucket[/prefix] -r -b [-f] 
        (Remove bucket and objects inside)
        If you remove with both --recursive and --bucket option, after ossutil removed all 
    the prefix-matching objects, ossutil will try to remove the bucket together. If user want 
    to remove bucket and objects inside, the usage is recommended. If --force option is 
    specified, ossutil will not show prompt question. 
`,

	sampleText: ` 
    ossutil rm oss://bucket1/obj1
    ossutil rm oss://bucket1/dir -r 
    ossutil rm oss://bucket1 -b
    ossutil rm oss://bucket2 -r -b -f
`,
}

// RemoveCommand is the command remove bucket or objects 
type RemoveCommand struct {
	command Command
}

var removeCommand = RemoveCommand{
	command: Command{
		name:        "rm",
		nameAlias:   []string{"remove", "delete", "del"},
		minArgc:     1,
		maxArgc:     1,
		specChinese: specChineseRemove,
		specEnglish: specEnglishRemove,
		group:       GroupTypeNormalCommand,
		validOptionNames: []string{
			OptionRecursion,
			OptionBucket,
			OptionForce,
            OptionMultipart,
            OptionAllType,
			OptionConfigFile,
            OptionEndpoint,
            OptionAccessKeyID,
            OptionAccessKeySecret,
            OptionSTSToken,
			OptionRetryTimes,
		},
	},
}

// function for FormatHelper interface
func (rc *RemoveCommand) formatHelpForWhole() string {
	return rc.command.formatHelpForWhole()
}

func (rc *RemoveCommand) formatIndependHelp() string {
	return rc.command.formatIndependHelp()
}


// Init simulate inheritance, and polymorphism 
func (rc *RemoveCommand) Init(args []string, options OptionMapType) error {
	return rc.command.Init(args, options, rc)
}

// RunCommand simulate inheritance, and polymorphism
func (rc *RemoveCommand) RunCommand() error {
    err := PreCheck()
    if err != nil {
        return err
    }
    
    if (isObject) {
		err = rc.removeObjectEntry(bucket, cloudURL, recursive)
        if err != nil {
            return err
        }
    }
    if (isMultipart) {
        err = rc.removeMultipartObjectEntry(bucket, cloudURL, recursive)
        if err != nil {
            return err
        }
    }
    if (toBucket) {
	    return rc.removeBucket(bucket, cloudURL, true)
    }
}

func (rc *RemoveCommand) PreCheck() error {
	recursive, _ := GetBool(OptionRecursion, rc.command.options)
	toBucket, _ := GetBool(OptionBucket, rc.command.options)
	force, _ := GetBool(OptionForce, rc.command.options)
    isMultipart, _ := GetBool(OptionMultipart, rc.command.options)
    isAllType, _ := GetBool(OptionAllType, rc.command.options)
    isObject := true

    if isMultipart {
        isObject = false
    }
    if isAllType {
        isMultipart = true
        isObject = true
    }

    if !recuresive && toBucket && isMultipart {
		return fmt.Errorf("invalid remove args: %s, miss bucket", rc.command.args[0])
    } 
    if !recuresive && toBucket && isAllType {
		return fmt.Errorf("invalid remove args: %s, miss bucket", rc.command.args[0])
    } 

	cloudURL, err := CloudURLFromString(rc.command.args[0])
	if err != nil {
		return err
	}
	if cloudURL.bucket == "" {
		return fmt.Errorf("invalid cloud url: %s, miss bucket", rc.command.args[0])
	}

   return nil; 
}


func (rc *RemoveCommand) removeObjectEntry(bucket *oss.Bucket, cloudURL CloudURL, recursive bool) error {
    if recursive {
        err := removeObject(bucket, cloudURL);
        if err != nil {
            return err
        }
    }
    
    return rc.recursiveRemoveObject(bucket, cloudURL, force)
}

func (rc *RemoveCommand) removeObject(bucket *oss.Bucket, cloudURL CloudURL) error {
	if cloudURL.object == "" {
		return fmt.Errorf("remove bucket, miss --bucket option, if you mean remove object, invalid url: %s, miss object", rc.command.args[0])
	}

	return rc.ossDeleteObjectRetry(bucket, cloudURL.object)
}

func (rc *RemoveCommand) ossDeleteObjectRetry(bucket *oss.Bucket, object string) error {
	retryTimes, _ := GetInt(OptionRetryTimes, rc.command.options)
	for i := 1; ; i++ {
		err := bucket.DeleteObject(object)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return ObjectError{err, object}
		}
	}
}

func (rc *RemoveCommand) removeMultipartObjectEntry(bucket *oss.Bucket, cloudURL CloudURL, recursive bool) error {
	if cloudURL.object == "" {
		return fmt.Errorf("remove bucket, miss --bucket option, if you mean remove multipart object, invalid url: %s, miss object", rc.command.args[0])
	}

	return rc.ossDeleteMultipartObjectRetry(bucket, cloudURL.object, recursive)
}

func (rc *RemoveCommand) ossDeleteMultipartObjectRetry(bucket *oss.Bucket, object string, recursive bool) error {
	retryTimes, _ := GetInt(OptionRetryTimes, rc.command.options)
	num := 0
	pre := oss.Prefix(object)
	marker := oss.Marker("")
	del := oss.Delimiter("")

	for i := 0; ; i++ {
	    lmr, err := rc.command.ossListMultipartObjectsRetry(bucket, marker, pre, del)
		if err != nil {
			 return err
		}
		pre = oss.Prefix(lmr.Prefix)
		marker = oss.Marker(lmr.NextKeyMarker)

        for _, upload := range lmr.Uploads {
            if !recursive {
                if object != upload.Key {
                    break
                }
            } 
            var imur = oss.InitiateMultipartUploadResult{Bucket: bucket.BucketName,
                Key: upload.Key, UploadID: upload.UploadID}
            err = bucket.AbortMultipartUpload(imur)
        }
		if err == nil {
			return err
		}

		if !lmr.IsTruncated {
			break
		}

		num += 1
		if int64(i) >= retryTimes {
			return ObjectError{err, object}
		}
	}
    return nil
}

func (rc *RemoveCommand) removeBucket(bucket *oss.Bucket, cloudURL CloudURL, force bool) error {
	if cloudURL.object != "" {
		return fmt.Errorf("remove bucket invalid url: %s, object not empty, if you mean remove object, you should not use --bucket option", rc.command.args[0])
	}

	if !force {
		var val string
		fmt.Printf("Do you really mean to remove the bucket:%s(y or N)? ", cloudURL.bucket)
		if _, err := fmt.Scanln(&val); err != nil || (strings.ToLower(val) != "yes" && strings.ToLower(val) != "y") {
			fmt.Println("operation is canceled.")
			return nil
		}
	}

	err := rc.ossDeleteBucketRetry(&bucket.Client, cloudURL.bucket)
    if err == nil {
		fmt.Printf("removed bucket: %s.\n", cloudURL.bucket)
    }
    return err
}

func (rc *RemoveCommand) ossDeleteBucketRetry(client *oss.Client, bucket string) error {
	retryTimes, _ := GetInt(OptionRetryTimes, rc.command.options)
	for i := 1; ; i++ {
		err := client.DeleteBucket(bucket)
		if err == nil {
			return err
		}
		if int64(i) >= retryTimes {
			return BucketError{err, bucket}
		}
	}
}

func (rc *RemoveCommand) recursiveRemoveObject(bucket *oss.Bucket, cloudURL CloudURL, force bool) error {
	if !force {
		var val string
		fmt.Printf("Do you really mean to recursivlly remove %s? ", rc.command.args[0])
		if _, err := fmt.Scanln(&val); err != nil || (val != "yes" && val != "y") {
			fmt.Println("operation is canceled.")
			return nil
		}
	}

	// batch delete objects
	num, err := rc.batchDeleteObjects(bucket, cloudURL)
    if err != nil {
		fmt.Printf("removed %d objects, when error happens.\n", num)
		return err
	}
	fmt.Printf("scaned %d objects, removed %d objects.\n", num, num)
	return nil
}

func (rc *RemoveCommand) batchDeleteObjects(bucket *oss.Bucket, cloudURL CloudURL) (int, error) {
	// list objects
	num := 0
	pre := oss.Prefix(cloudURL.object)
	marker := oss.Marker("")
	for i := 0; ; i++ {
		lor, err := rc.command.ossListObjectsRetry(bucket, marker, pre)
		if err != nil {
			return num, BucketError{err, bucket.BucketName}
		}

		// batch delete
		delNum, err := rc.ossBatchDeleteObjectsRetry(bucket, rc.getObjectsFromListResult(lor))
		num += delNum
		if err != nil {
			return num, BucketError{err, bucket.BucketName}
		}
		pre = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
	return num, nil
}

func (rc *RemoveCommand) ossBatchDeleteObjectsRetry(bucket *oss.Bucket, objects []string) (int, error) {
	retryTimes, _ := GetInt(OptionRetryTimes, rc.command.options)
	num := len(objects)
    if num <= 0 {
        return 0, nil
    }

	for i := 1; ; i++ {
		delRes, err := bucket.DeleteObjects(objects, oss.DeleteObjectsQuiet(true))
		if err == nil && len(delRes.DeletedObjects) == 0 {
			return num, nil
		}
		if int64(i) >= retryTimes {
			if err != nil {
				return num - len(objects), err
			}
			return num - len(delRes.DeletedObjects), fmt.Errorf("delete objects: %s failed", delRes.DeletedObjects)
		}
		objects = delRes.DeletedObjects
	}
}

func (rc *RemoveCommand) getObjectsFromListResult(lor oss.ListObjectsResult) []string {
	objects := []string{}
	for _, object := range lor.Objects {
		objects = append(objects, object.Key)
	}
	return objects
}

