package util

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

func WritePid() error {
	if _, err := os.Stat(Rdir); err != nil && !os.IsExist(err) {
		if err = os.MkdirAll(Rdir, 0755); err != nil {
			return errors.New("write pid file fail")
		}
	}
	pid := os.Getpid()
	filename := path.Join(Rdir, App+".pid")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(strconv.Itoa(pid)); err != nil {
		return err
	}
	return nil
}

func CheckRun() *os.Process {
	pid := getPid()
	if pid == 0 {
		return nil
	}
	exe, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return nil
	}
	if !strings.Contains(exe, Dir+"/"+App) {
		return nil
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	return process
}

func getPid() int {
	filename := path.Join(Rdir, App+".pid")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer file.Close()
	data := make([]byte, 10)
	bytes, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	pid, err := strconv.Atoi(strings.Replace(string(data[:bytes]), "\n", "", -1))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return pid
}

func KillApp() error {
	process := CheckRun()
	if process == nil {
		return ErrorApp
	}
	return process.Kill()
}
