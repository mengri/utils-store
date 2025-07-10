package store_mysql

import (
	"github.com/mengri/utils-store/store"
	store_mysql2 "github.com/mengri/utils-store/store/mysql"
	"github.com/mengri/utils-store/store/mysql/auto-yaml"
	"github.com/mengri/utils/autowire-v2"
	"github.com/mengri/utils/cftool"
)

func AutoYaml() {
	cftool.Register[auto_yaml.DBConfig]("mysql")

	autowire.Auto(func() store.IDB {
		return &store_mysql2.mysqlInit{}
	})
}
