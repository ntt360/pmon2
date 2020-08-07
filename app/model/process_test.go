package model

import (
	"fmt"
	"github.com/struCoder/pidusage"
	"log"
	"testing"
)

func TestProcess_RenderTable(t *testing.T) {
	info, err := pidusage.GetStat(4120)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info.CPU, info.Memory / 1024)
}
