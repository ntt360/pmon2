package conf

import (
	"os"
	"path"
)

type Tpl struct {
	AppDir string
	Data   string `yaml:"data"`
	Logs   string `yaml:"logs"`
	Sock   string `yaml:"sock"`
}

func (c *Tpl) GetSockFile() string {
	if len(c.Sock) <= 0 {
		panic("run sock file is empty")
	}

	sockDir := path.Dir(c.Sock)
	_, err := os.Stat(sockDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(sockDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// if sock file exist, try to clear
	_, err = os.Stat(c.Sock)
	if !os.IsNotExist(err) {
		err = os.Remove(c.Sock)
		if err != nil {
			panic(err)
		}
	}

	return c.Sock
}

func (c *Tpl) GetDataDir() string {
	if path.IsAbs(c.Data) {
		return c.Data
	}

	return c.AppDir + "/" + c.Data
}

func (c *Tpl) GetLogsDir() string {
	if path.IsAbs(c.Logs) {
		return c.Logs
	}

	return c.AppDir + "/" + c.Logs
}
