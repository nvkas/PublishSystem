package tool

import (
	"crypto/md5"
	"fmt"
)

func MD5Encryption(str string){
	h := md5.New()
	h.Write([]byte(str)) // 需要加密的字符串为 123456
	strs := h.Sum(nil)
	fmt.Println(strs)
}