package cpu

import (
	"fmt"
	"github.com/struCoder/pidusage"
	"strconv"
)

func GetExtraInfo(pid int) (string, string) {
	cpuVal := "0"
	memVal := "0"

	info, err := pidusage.GetStat(pid)
	if err != nil {
		return cpuVal, memVal
	}

	cpuVal = strconv.Itoa(int(info.CPU))

	if info.Memory <= 1024 {
		memVal = strconv.Itoa(int(info.Memory))
	} else if info.Memory <= 1024*1024 {
		memVal = fmt.Sprintf("%.1f KB", info.Memory/float64(1024))
	} else if info.Memory <= 1024*1024*1024 {
		memVal = fmt.Sprintf("%.1f MB", info.Memory/float64(1024*1024))
	} else if info.Memory <= 1024*1024*1024*1024 {
		memVal = fmt.Sprintf("%.1f GB", info.Memory/float64(1024*1024*1024))
	} else {
		memVal = strconv.Itoa(int(info.Memory))
	}

	return cpuVal, memVal
}
