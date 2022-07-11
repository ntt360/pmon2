package logf

import (
	"bufio"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/spf13/cobra"
	"io"
	"os/exec"
	"sync"
)

var Cmd = &cobra.Command{
	Use:   "logf",
	Short: "display process log dynamic by id or name",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	if len(args) == 0 {
		app.Log.Fatal("please input start process id or name")
	}
	val := args[0]
	var m model.Process
	if err := app.Db().First(&m, "id = ? or name = ?", val, val).Error; err != nil {
		app.Log.Fatal(fmt.Sprintf("the process %s not exist", val))
	}
	displayLog(m.Log)
}

func displayLog(log string) {
	c := exec.Command("bash", "-c", "tail -f "+log)
	stdout, _ := c.StdoutPipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
		}
	}()
	if err := c.Start(); err != nil {
		app.Log.Error(err)
	}
	wg.Wait()
}
