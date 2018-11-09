package config

import (
	"sync"

	"github.com/spf13/viper"
)

var init_once sync.Once
var over_once sync.Once

type rdoConf struct {
}

func NewRdoConf() *rdoConf {
	return &rdoConf{}
}

func init() {
	init_once.Do(func() {
		set_rdo()
	})
}

func set_rdo() {
	viper.SetDefault("rdo.auth_url", "http://127.0.0.1:35357/v3")
	viper.SetDefault("rdo.os_region_name", "regione")
	viper.SetDefault("rdo.project_domain_name", "default")
	viper.SetDefault("rdo.project_name", "admin")
	viper.SetDefault("rdo.user_domain_name", "default")
	viper.SetDefault("rdo.username", "admin")
	viper.SetDefault("rdo.password", "admin")
	viper.SetDefault("rdo.api_version", "2.6")
}

func (this *rdoConf) OverWriteConfig() {
	over_once.Do(func() {
		AUTH_URL = viper.GetString("rdo.auth_url")
		OS_REGION_NAME = viper.GetString("rdo.os_region_name")
		PROJECT_DOMAIN_NAME = viper.GetString("rdo.project_domain_name")
		PROJECT_NAME = viper.GetString("rdo.project_name")
		USER_DOMAIN_NAME = viper.GetString("rdo.user_domain_name")
		USERNAME = viper.GetString("rdo.username")
		PASSWORD = viper.GetString("rdo.password")
		API_VERSION = viper.GetString("rdo.api_version")
	})
}
