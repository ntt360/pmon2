package desc

import (
	"github.com/jinzhu/gorm"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	"github.com/spf13/cobra"
	"strconv"
)

var Cmd = &cobra.Command{
	Use:     "desc",
	Aliases: []string{"show"},
	Short:   "print the process detail message",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	val := args[0]

	var process model.Process
	err := app.Db().Find(&process, "name = ? or id = ?", val, val).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			app.Log.Fatal("pmon2 run err: %v", err)
		}

		// not found
		app.Log.Errorf("process %s not exist", val)
		return
	}

	rel := [][]string{
		{"status", process.Status},
		{"id", strconv.Itoa(int(process.ID))},
		{"name", process.Name},
		{"pid", strconv.Itoa(process.Pid)},
		{"process", process.ProcessFile},
		{"args", process.Args},
		{"user", process.Username},
		{"log", process.Log},
		{"no-autorestart", process.NoAutoRestartStr()},
		{"created_at", process.CreatedAt.Format("2006-01-02 15:04:05")},
		{"updated_at", process.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	output.DescTable(rel)
}
