package desc

import (
	"github.com/jinzhu/gorm"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
)

func Run(args []string) {
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

}
