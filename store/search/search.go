package search

import (
	"context"
	"fmt"
	"github.com/mengri/utils-store/store"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _ ISearchStore[any] = (*Store[any])(nil)

type ISearchStore[M any] interface {
	Search(ctx context.Context, keyword string, condition map[string]interface{}, sortRule ...string) ([]*M, error)
	SetLabels(ctx context.Context, id int64, label ...string) error
	Count(ctx context.Context, keyword string, condition map[string]interface{}) (int64, error)
	SearchByPage(ctx context.Context, keyword string, condition map[string]interface{}, page int, pageSize int, sortRule ...string) ([]*M, int64, error)
}
type Store[M any] struct {
	store.Store[M]
	name string
}

func (s *Store[M]) OnComplete() {

	s.Store.OnComplete()
	ctx := context.Background()

	var mi interface{} = new(M)
	if tn, ok := mi.(schema.Tabler); ok {
		s.name = fmt.Sprintf("%s_index", tn.TableName())
	} else {
		panic("not support")
	}

	err := s.Store.DB(ctx).Table(s.name).AutoMigrate(&Index{})

	if err != nil {
		panic(err)
	}
}

func (s *Store[M]) Search(ctx context.Context, keyword string, condition map[string]interface{}, sortRule ...string) ([]*M, error) {
	db := s.Store.DB(ctx)
	order := ""
	if len(sortRule) > 0 {
		order = strings.Join(sortRule, ",")
	}
	wm := condition
	if wm == nil {
		wm = map[string]interface{}{}
	}
	if keyword == "" {
		list := make([]*M, 0)
		err := db.Model(s.Model).Where(wm).Order(order).Find(&list).Error
		if err != nil {
			return nil, err
		}
		return list, err
	}
	ids := make([]interface{}, 0)
	err := db.Table(s.name).Select("DISTINCT target").Where("label like ?", "%"+keyword+"%").Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	wm["id"] = ids

	rs := make([]*M, 0)
	err = db.Model(s.Model).Where(wm).Order(order).Scan(&rs).Error
	return rs, err
}

func (s *Store[M]) Count(ctx context.Context, keyword string, condition map[string]interface{}) (int64, error) {
	db := s.Store.DB(ctx)

	wm := condition
	if wm == nil {
		wm = map[string]interface{}{}
	}
	if keyword == "" {
		var count int64
		err := db.Model(s.Model).Where(wm).Count(&count).Error
		if err != nil {
			return 0, err
		}
		return count, err
	}
	ids := make([]interface{}, 0)
	err := db.Table(s.name).Select("DISTINCT target").Where("label like ?", "%"+keyword+"%").Scan(&ids).Error
	if err != nil {
		return 0, err
	}
	wm["id"] = ids

	var count int64
	err = db.Model(s.Model).Where(wm).Count(&count).Error
	return count, err
}

func (s *Store[M]) SearchByPage(ctx context.Context, keyword string, condition map[string]interface{}, page int, pageSize int, sortRule ...string) ([]*M, int64, error) {
	db := s.Store.DB(ctx)
	order := "name asc"
	if len(sortRule) > 0 {
		order = strings.Join(sortRule, ",")
	}

	wm := condition
	if wm == nil {
		wm = map[string]interface{}{}
	}

	var count int64
	list := make([]*M, 0, pageSize)
	if keyword != "" {
		ids := make([]interface{}, 0)
		err := db.Table(s.name).Select("DISTINCT target").Where("label like ?", "%"+keyword+"%").Scan(&ids).Error
		if err != nil {
			return nil, 0, err
		}
		wm["id"] = ids
	}

	err := db.Order(order).Model(s.Model).Where(wm).Count(&count).Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	return list, count, err
}

func (s *Store[M]) SetLabels(ctx context.Context, id int64, label ...string) error {
	labelValid := make([]string, 0, len(label))
	for _, v := range label {
		if v == "" {
			continue
		}
		labelValid = append(labelValid, v)
	}
	label = labelValid
	return s.Store.DB(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table(s.name).Where("target = ?", id).Delete(s.name).Error
		if err != nil {
			return err
		}
		if len(label) == 0 {
			return nil
		}
		txn := tx.Table(s.name)
		for _, lv := range label {
			txn.Create(&Index{
				Id:     0,
				Target: id,
				Label:  lv,
			})

		}
		return txn.Error
	})

}
