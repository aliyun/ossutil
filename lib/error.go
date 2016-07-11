package lib

import (
	"fmt"
)

// CommandError happens when use command in invalid way
type CommandError struct {
	command string
	reason  string
}

func (e CommandError) Error() string {
	return fmt.Sprintf("invalid useage of \"%s\" command, reason: %s, please try \"help %s\" or \"%s --man\" for more information", e.command, e.reason, e.command, e.command)
}

// BucketError happens when access bucket error
type BucketError struct {
	err    error
	bucket string
}

func (e BucketError) Error() string {
	return fmt.Sprintf("%s, Bucket=%s", e.err.Error(), e.bucket)
}

// ObjectError happens when access object error
type ObjectError struct {
	err    error
	object string
}

func (e ObjectError) Error() string {
	return fmt.Sprintf("%s, Object=%s", e.err.Error(), e.object)
}

// FileError happens when access file error
type FileError struct {
	err  error
	file string
}

func (e FileError) Error() string {
	return fmt.Sprintf("%s, File=%s", e.err.Error(), e.file)
}
