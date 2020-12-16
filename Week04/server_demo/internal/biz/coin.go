package biz

import (
	"context"
	"server_demo/internal/data"
	"server_demo/pkg/cache"
	"server_demo/pkg/db"
)

type CoinDO struct {
	Uid uint64
	Num int32
}

type CoinRepo interface {
	GetCoinByUid(ctx context.Context, uid uint64) (*data.CoinPO, error)
}

type CoinBiz struct {
	repo CoinRepo
}

func NewCoinRepo(cacheStorage cache.Storage, db db.DB) CoinRepo {
	return data.NewCoinRepo(cacheStorage, db)
}

func NewCoinBiz(repo CoinRepo) *CoinBiz {
	return &CoinBiz{
		repo: repo,
	}
}

func (b *CoinBiz) GetUserCoin(ctx context.Context, uid uint64) (*CoinDO, error) {
	po, err := b.repo.GetCoinByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	return CoinPOToDo(po), nil
}

func CoinPOToDo(po *data.CoinPO) *CoinDO {
	return &CoinDO{
		Uid: po.Uid,
		Num: po.Num,
	}
}
