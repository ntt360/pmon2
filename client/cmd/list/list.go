package list

import (
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:                        "ls",
	Aliases:                    []string{"list"},
	Short:                      "list all processes",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(nil)
	},
}

// show all process list
func runCmd(_ []string) {
	var all []model.Process
	err := app.Db().Find(&all).Error
	if err != nil {
		app.Log.Fatalf("pmon2 run err: %v", err)
	}

	var allProcess [][]string
	for _, process := range all {
		allProcess = append(allProcess, process.RenderTable())
	}

	output.Table(allProcess)
}
