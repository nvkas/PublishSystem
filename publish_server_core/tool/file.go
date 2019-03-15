package tool

import (
	"fmt"
	"io/ioutil"
	"strings"
)
/*
filePath:文件路径 例如:/usr/local/test
fileName:文件名 例如:test.conf
fileContent:文件内容 例如:Hello World!
 */
func WriteFile(filePath string, fileName string,fileContent string) bool{
	if !strings.HasSuffix(filePath,"/") {
		filePath += "/"
	}
	//创建文件
	content:=[]byte(fileContent)
	err:=ioutil.WriteFile(filePath+fileName,content,0777)
	if err!=nil {
		fmt.Println(err)
		return false
	}
	return true
}
