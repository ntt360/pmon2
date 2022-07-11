package proxy

import (
	"github.com/ntt360/errors"
	"github.com/ntt360/pmon2/app"
	conf2 "github.com/ntt360/pmon2/app/conf"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/utils/array"
	"github.com/ntt360/pmon2/client/worker"
)

var cmdTypes = []string{"start", "restart"}

func RunProcess(args []string) ([]byte, error) {
	if len(args) <= 2 {
		return nil, errors.New("process params not valid")
	}
	conf := conf2.GetDefaultConf()
	err := app.Instance(conf)

	if err != nil {
		return nil, err
	}
	// check run type param
	typeCli := args[0]

	if !array.In(cmdTypes, typeCli) {
		return nil, errors.WithStack(err)
	}

	var output string

	flagModel, err := model.ExecFlags{}.Parse(args[2])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch typeCli {
	case "start":
		output, err = worker.Start(args[1], flagModel)
	case "restart":
		output, err = worker.Restart(args[1], flagModel)
	}

	if err != nil {
		return []byte(err.Error()), err
	}

	return []byte(output), nil
}
