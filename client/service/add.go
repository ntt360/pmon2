package service

import (
	"encoding/json"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
)

func AddData(process *model.Process) (string, error) {
	// save to db
	var originOne model.Process
	err := app.Db().First(&originOne, "process_file = ?", process.ProcessFile).Error
	if err == nil && originOne.ID > 0 { // process already exist
		process.ID = originOne.ID
	}

	err = app.Db().Save(&process).Error
	if err != nil {
		return "", fmt.Errorf("pmon2 run err: %w", err)
	}

	output, err := json.Marshal(process.RenderTable())
	if err != nil {
		return "", err
	}

	return string(output), nil
}
