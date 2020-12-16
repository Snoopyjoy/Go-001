package data

import (
	"context"
	"fmt"
	"server_demo/pkg/cache"
	"server_demo/pkg/db"
	"time"
)

type CoinPO struct {
	Id  uint64
	Uid uint64
	Num int32
}

type CoinRepoImpl struct {
	cacheStorage cache.Storage
	db           db.DB
}

func (r *CoinRepoImpl) GetCoinByUid(ctx context.Context, uid uint64) (*CoinPO, error) {
	rst, err := cache.Cache(r.cacheStorage, fmt.Sprintf("coin:%d", uid), time.Minute*3, func() (interface{}, error) {
		return &CoinPO{
			Id:  1,
			Uid: uid,
			Num: 3,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return rst.(*CoinPO), nil
}

func NewCoinRepo(cacheStorage cache.Storage, db db.DB) *CoinRepoImpl {
	return &CoinRepoImpl{
		cacheStorage: cacheStorage,
		db:           db,
	}
}
