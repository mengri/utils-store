package cache_redis

import (
	"context"
	"github.com/mengri/utils/cftool"
	"log"
	"time"

	"github.com/mengri/utils-store/cache"
)

type RedisConfig struct {
	UserName   string   `yaml:"user_name"`
	Password   string   `yaml:"password"`
	Addr       []string `yaml:"addr"`
	Prefix     string   `yaml:"prefix"`
	MasterName string   `yaml:"master_name"`
	DB         int      `yaml:"db"`
}

type redisInit struct {
	cache.ICommonCache
	conf *RedisConfig `autowired:""`
}

func (r *redisInit) OnPreComplete() {

	client := SimpleCluster(Option{
		Addrs:      r.conf.Addr,
		MasterName: r.conf.MasterName,
		Username:   r.conf.UserName,
		Password:   r.conf.Password,
		DB:         r.conf.DB,
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	if err := client.Ping(timeout).Err(); err != nil {
		_ = client.Close()

		log.Fatalf("ping redis %v error:%s", r.conf.Addr, err.Error())
	}

	r.ICommonCache = newCommonCache(client, r.conf.Prefix)
}

func init() {
	cftool.Register[RedisConfig]("redis")
	autowire.Auto[cache.ICommonCache](func() cache.ICommonCache {

		return new(redisInit)
	})

}
