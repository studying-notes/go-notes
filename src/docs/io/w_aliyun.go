package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

const (
	accessKeyId     = "LTAI4GBGw5ndeWaHtUQTZGJz"
	accessKeySecret = "UhKrxBOvA6Er51tA8dvx0Ivh0hzF2L"

	endPointExternal = "http://oss-cn-hangzhou.aliyuncs.com"          // 外网访问
	endPointInternal = "http://oss-cn-hangzhou-internal.aliyuncs.com" //内网访问

	bucketName = "jixinwulian"
)

func HandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func main() {
	client, err := oss.New(endPointExternal, accessKeyId, accessKeySecret)
	//client, err := oss.New(endPointInternal, accessKeyId, accessKeySecret)
	if err != nil {
		HandleError(err)
	}

	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// 上传文件
	err = bucket.PutObjectFromFile("go.mod", "go.mod")
	err = bucket.PutObjectFromFile("LICENSE", "LICENSE")
	if err != nil {
		HandleError(err)
	}

	// 下载文件
	err = bucket.GetObjectToFile("go.mod", "go.mod2")
	if err != nil {
		HandleError(err)
	}

	// 删除文件
	err = bucket.DeleteObject("go.mod")
	if err != nil {
		HandleError(err)
	}
}
