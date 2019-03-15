package tool

import (
	"fmt"
	"os"
)

//Clone项目
//repoDir string	项目保存地址(Git本地仓库)
//branch  string	项目分支(Git分支)
//gitAdd  string	项目git地址(Git地址)
func CloneRepo(repoDir string,branch string,gitAdd string) bool {
	flag := MkDirs(repoDir)
	if flag {
		return SplitCommandAndExec(fmt.Sprintf("cd %v||git clone -b %v %v ",repoDir,branch,gitAdd))
	}
	return false
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建文件夹
func MkDirs(dir string) bool {
	exist, _ := PathExists(dir)

	if exist {
		fmt.Printf("has dir![%v]\n", dir)
		return true
	} else {
		fmt.Printf("no dir![%v]\n", dir)
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			return false
		} else {
			fmt.Printf("mkdir success!\n")
			return true
		}
	}
	return false
}

