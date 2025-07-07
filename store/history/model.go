package history

import (
	"time"
)

type Base struct {
	Id     int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;"`
	UUID   string    `gorm:"column:uuid;type:varchar(36);NOT NULL;comment:uuid;"`
	User   string    `gorm:"column:user;type:varchar(36);NOT NULL;comment:user;"`
	Target int64     `gorm:"column:target;type:BIGINT(20);NOT NULL;comment:sid;"`
	Time   time.Time `gorm:"column:time;type:timestamp;NOT NULL;comment:time;"`
}
type History[T any] struct {
	Id     int64     `gorm:"column:id;type:BIGINT(20);AUTO_INCREMENT;NOT NULL;comment:id;primary_key;"`
	UUID   string    `gorm:"column:uuid;type:varchar(36);NOT NULL;comment:uuid;"`
	User   string    `gorm:"column:user;type:varchar(36);NOT NULL;comment:user;"`
	Target int64     `gorm:"column:target;type:BIGINT(20);NOT NULL;comment:sid;"`
	Time   time.Time `gorm:"column:time;type:timestamp;NOT NULL;comment:time;"`
	Data   *T        `gorm:"column:data;type:LONGTEXT;NOT NULL;comment:data; serializer:json"`
}

type Latest struct {
	Id     int64  `gorm:"column:id;type:BIGINT(20);NOT NULL;comment:id;primary_key;"`
	Latest int64  `gorm:"column:latest;type:BIGINT(20);NOT NULL;comment:last;"`
	Commit string `gorm:"column:commit;type:varchar(36);NOT NULL;comment:commit;"`
}
