package tool

import (
	"github.com/prometheus/common/log"
	"io/ioutil"
	"os/exec"
)

/**
根据名字查询PID
*/
func GetPidByProcessName(name string) string {
	cmd := exec.Command("bash", "-c", "pgrep -f "+name)
	log.Info("打印CMD", cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("Error:can not obtain stdout pipe for command:%s\n", err)
		return ""
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		log.Error("Error:The command is err,", err)
		return ""
	}
	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Error("ReadAll Stdout:", err.Error())
		return ""
	}
	if err := cmd.Wait(); err != nil {
		log.Error("wait:", err.Error())
		return ""
	}
	log.Info("stdout:\n\n %s", bytes)
	pid := string(bytes[:])
	return pid
}
