package search

// Index 索引, 用于快速查询
type Index struct {
	Id     int64  `gorm:"type:BIGINT(20);size:20;not null;auto_increment;primary_key;column:id;comment:主键ID;"`
	Target int64  `gorm:"type:BIGINT(20);size:20;not null;column:target;comment:target id;index:tid;"`
	Label  string `gorm:"type:varchar(255);not null;column:label;comment:标签"`
}

func (i *Index) IdValue() int64 {
	return i.Id
}
