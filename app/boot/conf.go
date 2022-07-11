package boot

import (
	"github.com/ntt360/errors"
	"github.com/ntt360/pmon2/app/conf"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Conf(confFile string) (*conf.Tpl, error) {
	d, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var c conf.Tpl
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	c.Conf = confFile
	
	return &c,nil
}
