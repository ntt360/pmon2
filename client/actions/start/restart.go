package start

import (
	"encoding/json"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/proxy"
)

func restart(m *model.Process, dir string, args []string) ([]string, error) {
	// only stopped process or failed process allow run start
	if m.Status == model.StatusStopped || m.Status == model.StatusFailed {
		newData, err := reloadProcess(m, dir, args)
		if err != nil {
			return nil, err
		}

		return newData, nil
	}

	return m.RenderTable(), nil
}

func reloadProcess(m *model.Process, dir string, args []string) ([]string, error) {
	args = append([]string{"restart", m.ProcessFile}, args[1:]...)
	data, err := proxy.RunProcess(dir, args)

	if err != nil {
		return nil, err
	}

	var rel []string
	err = json.Unmarshal(data, &rel)
	if err != nil {
		return nil, err
	}

	return rel, nil
}
