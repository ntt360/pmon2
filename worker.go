package main

import (
	"os"

	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/utils/array"
	w "github.com/ntt360/pmon2/client/worker"
)

var cmdTypes = []string{"start", "restart"}

func worker() {
	args := os.Args

	if len(args) <= 2 {
		_, _ = os.Stderr.WriteString("process params not valid")
		os.Exit(2)
	}

	err := app.Instance(os.Getenv("PMON2_CONF"))
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(2)
	}

	// check run type param
	typeCli := args[0]
	if !array.In(cmdTypes, typeCli) {
		_, _ = os.Stderr.WriteString("run params not illegal")
		os.Exit(2)
	}

	var output string

	flagModel, err := model.ExecFlags{}.Parse(args[2])
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(2)
	}

	switch typeCli {
	case "start":
		output, err = w.Start(args[1], flagModel)
		break
	case "restart":
		output, err = w.Restart(args[1], flagModel)
		break
	}

	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(2)
	}

	_, _ = os.Stdout.WriteString(output)
}
