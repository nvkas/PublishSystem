package tool

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

//切割命令行,执行命令
//text 命令行
func SplitCommandAndExec(text string) (success bool) {
	success = true
	cmdarr := strings.Split(text, "||")
	for _, val := range cmdarr {
		tmpval := strings.TrimSpace(val)

		//如果是以cd开头，那么是切换目录操作
		if tmpval != "" && strings.HasPrefix(tmpval, "cd "){
			tmpval = tmpval[3 : ]
			fmt.Println("cd到目录:",tmpval)
			err := os.Chdir(tmpval)
			if err != nil {
				fmt.Println(err)
				//success = false
				//break
			}
		} else if tmpval != "" {
			//分割命令
			cmdarr := strings.Split(tmpval, " ")
			fmt.Println("-----------执行命令-----------")
			fmt.Println(cmdarr)
			//命令名称
			command := cmdarr[0]
			//命令参数
			params := cmdarr[1:]
			//执行cmd命令
			flag := ExecCommand(command, params)
			if !flag {
				//success = false
			}
		}
	}
	return
}

//执行命令函数
//commandName 命名名称，如cat，ls，git等
//params 命令参数，如ls -l的-l，git log 的log等
func ExecCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetPidByProcessName(name string) string {
	cmd := exec.Command("bash","-c","pgrep -f "+name)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return ""
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return ""
	}
	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return ""
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return ""
	}
	fmt.Printf("stdout:\n\n %s", bytes)
	str2 := string(bytes[:])
	//str2为PID
	return str2
}