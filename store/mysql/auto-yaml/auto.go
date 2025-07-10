package auto_yaml

import (
	"github.com/mengri/utils-store/store"
	store_mysql "github.com/mengri/utils-store/store/mysql"
	"github.com/mengri/utils/autowire-v2"
	"github.com/mengri/utils/cftool"
)

func init() {
	cftool.Register[DBConfig]("mysql")

	autowire.Auto(func() store.IDB {
		return &mysqlInit{}
	})
}

type mysqlInit struct {
	store.IDB
	config *DBConfig `autowired:""`
}

func (m *mysqlInit) OnPreComplete() {
	m.InitDb()
}
func (m *mysqlInit) InitDb() {

	m.IDB = store_mysql.CreateDb(m.config.getDBNS())
}
