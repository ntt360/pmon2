package process

import (
	"github.com/ntt360/pmon2/app/boot"
	"github.com/ntt360/pmon2/app/model"
)

func FindByProcessFile(pFile string) *model.Process{
	var rel model.Process
	err := boot.Db().First(&rel, "process_file = ?", pFile).Error
	if err != nil {
		return nil
	}

	return &rel
}
