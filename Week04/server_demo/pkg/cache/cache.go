package cache

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Storage interface {
	Set(key string, value interface{}, duration time.Duration) error
	Get(key string) (interface{}, error)
}

func Cache(storage Storage, key string, duration time.Duration, fetcher func() (interface{}, error)) (interface{}, error) {
	v, err := storage.Get(key)
	if err != nil && !IsNotFound(err) {
		return nil, errors.Wrapf(err, "fetch err %s", key)
	}
	if err != nil && IsNotFound(err) {
		v, err = fetcher()
		if err != nil {
			return nil, err
		}
		err = storage.Set(key, v, duration)
		if err != nil {
			fmt.Printf("set %s err %+v\n", key, err)
		}
	}

	return v, nil
}

func IsNotFound(err error) bool {
	return true
}

type CacheConf struct {
	Addr   string `yaml:"addr"`
	Passwd string `yaml:"passwd"`
}

type cacheClient struct {
}

func (c *cacheClient) Get(key string) (interface{}, error) {
	return nil, errors.New("not found")
}

func (c *cacheClient) Set(key string, value interface{}, duration time.Duration) error {
	return nil
}

func NewCacheClient(cfg *CacheConf) (Storage, error) {
	return &cacheClient{}, nil
}
