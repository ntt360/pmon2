package cmd

import (
	"github.com/ntt360/pmon2/client/cmd/del"
	"github.com/ntt360/pmon2/client/cmd/desc"
	"github.com/ntt360/pmon2/client/cmd/exec"
	"github.com/ntt360/pmon2/client/cmd/list"
	"github.com/ntt360/pmon2/client/cmd/reload"
	"github.com/ntt360/pmon2/client/cmd/start"
	"github.com/ntt360/pmon2/client/cmd/stop"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "pmon2",
	Short:   "pmon2 client cli",
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
	)
	return rootCmd.Execute()
}
