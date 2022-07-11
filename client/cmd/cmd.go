package cmd

import (
	"fmt"
	"github.com/ntt360/pmon2/app/conf"
	"github.com/ntt360/pmon2/client/cmd/completion"
	"github.com/ntt360/pmon2/client/cmd/del"
	"github.com/ntt360/pmon2/client/cmd/desc"
	"github.com/ntt360/pmon2/client/cmd/exec"
	"github.com/ntt360/pmon2/client/cmd/list"
	"github.com/ntt360/pmon2/client/cmd/log"
	"github.com/ntt360/pmon2/client/cmd/logf"
	"github.com/ntt360/pmon2/client/cmd/reload"
	"github.com/ntt360/pmon2/client/cmd/restart"
	"github.com/ntt360/pmon2/client/cmd/start"
	"github.com/ntt360/pmon2/client/cmd/stop"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "pmon2",
	Short: "pmon2 client cli",
}

var verCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pmon2: %s \n", conf.Version)
	},
}

func Exec() error {
	rootCmd.AddCommand(
		del.Cmd,
		desc.Cmd,
		list.Cmd,
		exec.Cmd,
		stop.Cmd,
		reload.Cmd,
		start.Cmd,
		restart.Cmd,
		completion.Cmd,
		log.Cmd,
		logf.Cmd,
		verCmd,
	)

	return rootCmd.Execute()
}
