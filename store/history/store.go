package history

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mengri/utils-store/store"
	"github.com/mengri/utils/utils"
	"gorm.io/gorm/schema"
	"time"
)

var _ ICommitStore[any, any] = (*CommitStore[any, any])(nil)

type ICommitStore[T any, D any] interface {
	Latest(ctx context.Context, id int64) (*History[D], error)
	LatestCommits(ctx context.Context, ids ...int64) (map[int64]string, error)
	FindLatest(ctx context.Context, id ...int64) ([]*History[D], error)
	AddHistory(ctx context.Context, id int64, userId string, data *D) (*History[D], error)
	SetLatest(ctx context.Context, id int64, latest int64, commit string) error
	ListHistory(ctx context.Context, id int64, page, size int) ([]*Base, int64, error)
	DeleteHistory(ctx context.Context, id int64) error
}

type CommitStore[T any, D any] struct {
	store.Store[T]
	historyName string
	latestName  string
}

func (s *CommitStore[T, D]) SetLatest(ctx context.Context, id int64, latest int64, commit string) error {
	return s.DB(ctx).Table(s.latestName).Save(&Latest{
		Id:     id,
		Latest: latest,
		Commit: commit,
	}).Error
}
func (s *CommitStore[T, D]) DeleteHistory(ctx context.Context, id int64) error {
	return s.Transaction(ctx, func(ctx context.Context) error {
		err := s.DB(ctx).Table(s.historyName).Where("target=?", id).Delete(&History[D]{}).Error
		if err != nil {
			return err
		}
		return s.DB(ctx).Table(s.latestName).Where("id=?", id).Delete(&Latest{}).Error
	})
}
func (s *CommitStore[T, D]) AddHistory(ctx context.Context, id int64, userId string, data *D) (*History[D], error) {

	v := &History[D]{
		Id:     0,
		UUID:   uuid.NewString(),
		User:   userId,
		Target: id,
		Time:   time.Now(),
		Data:   data,
	}
	err := s.Transaction(ctx, func(ctx context.Context) error {
		if err := s.Store.DB(ctx).Table(s.historyName).Create(v).Error; err != nil {
			return err
		}
		return s.SetLatest(ctx, id, v.Id, v.UUID)

	})
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *CommitStore[T, D]) ListHistory(ctx context.Context, id int64, page, size int) ([]*Base, int64, error) {
	var count int64
	var rs []*Base
	s.DB(ctx).Table(s.historyName).Where("target=?", id).Order("id desc").Count(&count).Offset(page * size).Limit(size).Find(&rs)
	return rs, count, nil
}

func (s *CommitStore[T, D]) OnComplete() {

	s.Store.OnComplete()
	ctx := context.Background()

	var mi interface{} = new(T)
	if tn, ok := mi.(schema.Tabler); ok {
		s.historyName = fmt.Sprintf("%s_history", tn.TableName())
		s.latestName = fmt.Sprintf("%s_latest", tn.TableName())
	} else {
		panic("not support")
	}
	err := s.DB(ctx).Table(s.historyName).AutoMigrate(&History[D]{})

	if err != nil {
		panic(err)
	}
	err = s.DB(ctx).Table(s.latestName).AutoMigrate(&Latest{})
	if err != nil {
		panic(err)
	}
}

func (s *CommitStore[T, D]) FindLatest(ctx context.Context, ids ...int64) ([]*History[D], error) {
	if len(ids) == 0 {
		return nil, nil
	}
	list := make([]*History[D], 0, len(ids))
	err := s.Transaction(ctx, func(ctx context.Context) error {
		latestList := make([]*Latest, 0, len(ids))
		err := s.DB(ctx).Table(s.latestName).Where(map[string]any{"id": ids}).Find(&latestList).Error
		if err != nil {
			return err
		}
		hids := utils.SliceToSlice(latestList, func(s *Latest) int64 {
			return s.Latest
		})
		return s.DB(ctx).Table(s.historyName).Where(map[string]any{"id": hids}).Find(&list).Error
	})
	if err != nil {
		return nil, err
	}
	return list, nil

}
func (s *CommitStore[T, D]) Latest(ctx context.Context, id int64) (*History[D], error) {

	db := s.DB(ctx)
	latest := new(Latest)
	err := db.Table(s.latestName).First(latest, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	ho := new(History[D])
	err = db.Table(s.historyName).First(ho, "id=?", latest.Latest).Error
	if err != nil {
		return nil, err
	}
	return ho, nil
}
func (s *CommitStore[T, D]) LatestCommits(ctx context.Context, ids ...int64) (map[int64]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	latestList := make([]*Latest, 0, len(ids))
	err := s.DB(ctx).Table(s.latestName).Where(map[string]any{"id": ids}).Find(&latestList).Error
	if err != nil {
		return nil, err
	}
	return utils.SliceToMapO(latestList, func(s *Latest) (int64, string) {
		return s.Id, s.Commit
	}), nil

}
