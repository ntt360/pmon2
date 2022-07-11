package conf

type Tpl struct {
	Data string `yaml:"data"`
	Logs string `yaml:"logs"`
	Conf string
}

func (c *Tpl) GetDataDir() string {
	return c.Data
}

func (c *Tpl) GetLogsDir() string {
	return c.Logs
}
