package conf

import "os"

// current app version
var Version = "1.10.0"

func GetDefaultConf() string{
	conf := os.Getenv("PMON2_CONF")
	if len(conf) == 0 {
		conf = "/etc/pmon2/config/config.yml"
	}
	return conf

}
