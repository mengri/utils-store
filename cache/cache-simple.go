package cache

import (
	"context"
	"fmt"
	"time"
)

type SimpleCache[K any] interface {
	Get(ctx context.Context, k K) ([]byte, error)
	Delete(ctx context.Context, k ...K) error
	Set(ctx context.Context, k K, v []byte) error
}
type simple[K any] struct {
	client        ICommonCache
	formatHandler func(k K) string
	expiration    time.Duration
}

func (r *simple[K]) Get(ctx context.Context, k K) ([]byte, error) {
	kv := r.formatHandler(k)

	bytes, err := r.client.Get(ctx, kv)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}
func (r *simple[K]) Set(ctx context.Context, k K, v []byte) error {

	kv := r.formatHandler(k)

	return r.client.Set(ctx, kv, v, r.expiration)
}

func (r *simple[K]) Delete(ctx context.Context, ks ...K) error {
	for _, k := range ks {
		key := r.formatHandler(k)
		if err := r.client.Del(ctx, key); err != nil {
			return err
		}
	}
	return nil
}
func CreateCacheSimple[K any](client ICommonCache, expiration time.Duration, formatHandler ...func(k K) string) SimpleCache[K] {

	s := &simple[K]{
		client:     client,
		expiration: expiration,
	}
	if len(formatHandler) == 0 {
		s.formatHandler = func(k K) string { return fmt.Sprint(k) }
	} else {
		s.formatHandler = formatHandler[0]
	}
	return s
}
