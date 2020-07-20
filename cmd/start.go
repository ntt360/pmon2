package cmd

import (
	"fmt"
	"log"
	"os"
)

func Start(args []string) {
	// 获取进程文件
	processFile := args[0]
	file, err := os.Stat(processFile)
	if os.IsNotExist(err) || file.IsDir() {
		// TODO 进程文件不存在
		log.Fatal(fmt.Sprintf("%s not exist", processFile))
		return
	}

	//syscall.ForkExec()
	fmt.Println(args)
}
