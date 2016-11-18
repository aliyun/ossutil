package lib 

import (
    "os"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) rawUpdate(force bool, language string) (bool, error) {
    command := "update" 
    var args []string
    options := OptionMapType{
        OptionForce: &force,
        OptionLanguage: &language, 
    }
    showElapse, err := cm.RunCommand(command, args, options)
    return showElapse, err
}

func (s *OssutilCommandSuite) TestUpdate(c *C) {
    showElapse, err := s.rawUpdate(false, "中文")
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawUpdate(false, "English")
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawUpdate(true, "中文")

    showElapse, err = s.rawUpdate(true, "English")

    err = updateCommand.updateVersion(Version, "中文")

    fileName := "ossutil_test_not_exist"
    err = updateCommand.rewriteLoadConfig(fileName)
    c.Assert(err, IsNil)
}

func (s *OssutilCommandSuite) TestRevertRename(c *C) {
    filePath := ".ossutil_tempf"
    renameFilePath := ".ossutil_tempr"

    s.createFile(filePath, filePath + "i", c)
    s.createFile(renameFilePath, renameFilePath + "i", c)

    updateCommand.revertRename(filePath, renameFilePath)
    _, err := os.Stat(renameFilePath) 
    c.Assert(err, NotNil)

    str := s.readFile(filePath, c) 
    c.Assert(str, Equals, renameFilePath + "i")

    _ = os.Remove(filePath)
    _ = os.Remove(renameFilePath)
}

func (s *OssutilCommandSuite) TestDownloadLastestBinary(c *C) {
    tempBinaryFile := ".ossutil_test_update.temp"  
    err := updateCommand.getBinary(tempBinaryFile, "1.0.0.Beta") 
    c.Assert(err, IsNil)

    _ = os.Remove(tempBinaryFile)
}

func (s *OssutilCommandSuite) TestAnonymousGetToFileError(c *C) {
    bucket := bucketNameNotExist 
    object := "TestAnonymousGetToFileError"
    err := updateCommand.anonymousGetToFileRetry(bucket, object, object)
    c.Assert(err, NotNil)

    bucket = bucketNameExist
    s.putObject(bucket, object, uploadFileName, c)
    fileName := "*"
    err = updateCommand.anonymousGetToFileRetry(bucket, object, fileName)
    c.Assert(err, NotNil)
}
