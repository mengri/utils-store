package search

import (
	"context"
	"fmt"
)

var _ ISearchStore[any] = (*StoreSoftDelete[any])(nil)

type StoreSoftDelete[T any] struct {
	Store[T]
}

func (s *StoreSoftDelete[T]) Delete(ctx context.Context, id ...int64) (int, error) {

	r := s.Store.DB(ctx).Where(map[string]interface{}{
		"id": id,
	}).Update("is_delete", true)
	return int(r.RowsAffected), r.Error
}

func (s *StoreSoftDelete[T]) DeleteWhere(ctx context.Context, m map[string]interface{}) (int64, error) {

	r := s.Store.DB(ctx).Where(m).Update("is_delete", true)
	return r.RowsAffected, r.Error
}

func (s *StoreSoftDelete[T]) DeleteUUID(ctx context.Context, uuid string) error {
	r := s.Store.DB(ctx).Where(map[string]interface{}{
		"uuid": uuid,
	}).Update("is_delete", true)
	return r.Error
}

func (s *StoreSoftDelete[T]) DeleteQuery(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	r := s.Store.DB(ctx).Where(sql, args...).Update("is_delete", true)
	return r.RowsAffected, r.Error
}

func (s *StoreSoftDelete[T]) CountWhere(ctx context.Context, m map[string]interface{}) (int64, error) {
	vm := m
	if vm == nil {
		vm = map[string]interface{}{}
	}
	vm["is_delete"] = false
	return s.Store.CountWhere(ctx, vm)
}

func (s *StoreSoftDelete[T]) CountQuery(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return s.Store.CountQuery(ctx, "is_delete = false", false)
	}
	return s.Store.CountQuery(ctx, fmt.Sprintf("(%s) AND is_delete = false", sql), args...)
}

func (s *StoreSoftDelete[T]) List(ctx context.Context, m map[string]interface{}, order ...string) ([]*T, error) {
	vm := m
	if vm == nil {
		vm = map[string]interface{}{}
	}
	return s.Store.List(ctx, vm, order...)
}

func (s *StoreSoftDelete[T]) ListQuery(ctx context.Context, sql string, args []interface{}, order string) ([]*T, error) {
	if sql != "" {
		sql = fmt.Sprintf("(%s) AND is_delete = false", sql)
	} else {
		sql = "is_delete = false"
	}
	return s.Store.ListQuery(ctx, sql, args, order)
}

func (s *StoreSoftDelete[T]) First(ctx context.Context, m map[string]interface{}, order ...string) (*T, error) {
	if m == nil {
		m = map[string]interface{}{}
	}
	m["is_delete"] = false
	return s.Store.First(ctx, m, order...)
}

func (s *StoreSoftDelete[T]) FirstQuery(ctx context.Context, sql string, args []interface{}, order string) (*T, error) {
	if sql != "" {
		sql = fmt.Sprintf("(%s) AND is_delete = false", sql)
	} else {
		sql = "is_delete = false"
	}
	return s.Store.FirstQuery(ctx, sql, args, order)
}
func (s *StoreSoftDelete[T]) ListPageWhere(ctx context.Context, where map[string]any, pageNum, pageSize int, order string) ([]*T, int64, error) {
	if where == nil {
		where = map[string]any{}
	}
	where["is_delete"] = false
	return s.Store.ListPageWhere(ctx, where, pageNum, pageSize, order)

}
func (s *StoreSoftDelete[T]) ListPage(ctx context.Context, sql string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int64, error) {
	if sql != "" {
		sql = fmt.Sprintf("(%s) AND is_delete = false", sql)
	} else {
		sql = "is_delete = false"
	}

	return s.Store.ListPage(ctx, sql, pageNum, pageSize, args, order)

}

func (s *StoreSoftDelete[T]) Search(ctx context.Context, keyword string, condition map[string]interface{}, sortRule ...string) ([]*T, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.Store.Search(ctx, keyword, condition, sortRule...)
}

func (s *StoreSoftDelete[T]) Count(ctx context.Context, keyword string, condition map[string]interface{}) (int64, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.Store.Count(ctx, keyword, condition)
}

func (s *StoreSoftDelete[T]) CountByGroup(ctx context.Context, condition map[string]interface{}, groupBy string) (map[string]int64, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.Store.CountByGroup(ctx, condition, groupBy)
}

func (s *StoreSoftDelete[T]) SearchByPage(ctx context.Context, keyword string, condition map[string]interface{}, page int, pageSize int, sortRule ...string) ([]*T, int64, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.Store.SearchByPage(ctx, keyword, condition, page, pageSize, sortRule...)
}
