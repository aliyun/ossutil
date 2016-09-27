# Aliyun OSSUTIL 
## 关于
> - 此工具采用go语言，基于OSS[阿里云对象存储服务](http://www.aliyun.com/product/oss/)官方GO SDK 1.0.0 构建。
> - 阿里云对象存储（Object Storage Service，简称OSS），是阿里云对外提供的海量，安全，低成本，高可靠的云存储服务。
> - OSS适合存放任意文件类型，适合各种网站、开发企业及开发者使用。
> - 该工具旨在为用户提供一个方便的，以命令行方式管理OSS数据的途径。
> - 当前版本未提供Bucket管理功能和Multipart管理功能，相关功能会在后续版本中开发。

## 版本
> - 当前版本：1.0.0.Beta

## 运行环境
> - linux, windows 

## 依赖的库 
> - goopt (github.com/droundy/goopt) 
> - configparser (github.com/alyu/configparser)
> - oss (github.com/aliyun/aliyun-oss-go-sdk/oss)
> - gopkg.in/check.v1 (gopkg.in/check.v1)

## 快速使用
#### 获取命令列表
```go
    ./ossutil
    或 ./ossutil help
```

#### 查看某命令的帮助文档
```go
    ./ossutil help cmd 
    或 ./ossutil cmd --man
```
    
#### 配置ossutil 
```go
    ./ossutil config
```

#### 列举Buckets
```go
    ./ossutil ls
    或 ./ossutil ls oss://
```

#### 上传文件
```go
    ./ossutil cp localfile oss://bucket
```

#### 其它
请使用./ossutil help cmd来查看想要使用的命令的帮助文档。

## 注意事项
### 运行
> - 首先配置您的go工程目录。
> - go get该工程依赖的库。
> - go get github.com/aliyun/ossutil。
> - 进入go工程目录下的src目录，build代码生成ossutil工具，例如：在linux下可以运行go build github.com/aliyun/ossutil/ossutil.go。
> - 参考上面示例运行ossutil工具。

### 测试
> - 进入go工程目录下的src目录，修改github.com/aliyun/ossutil/lib/command_test.go里的endpoint、AccessKeyId、AccessKeySecret、STSToken等配置。
> - 请在lib目录下执行`go test`。

### 源码
> - https://github.com/aliyun/ossutil

## 联系我们
> - [阿里云OSS官方网站](http://oss.aliyun.com)
> - [阿里云OSS官方论坛](http://bbs.aliyun.com)
> - [阿里云OSS官方文档中心](http://www.aliyun.com/product/oss#Docs)

## 作者
> - Ting Zhang 

## License
> - MIT 
