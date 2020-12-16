package data

import (
	"context"
	"fmt"
	"server_demo/pkg/cache"
	"server_demo/pkg/db"
	"time"
)

type UserPO struct {
	Id       uint64
	Nickname string
}

type UserRepoImpl struct {
	cacheStorage cache.Storage
	db           db.DB
}

func (r *UserRepoImpl) FindById(ctx context.Context, uid uint64) (*UserPO, error) {
	rst, err := cache.Cache(r.cacheStorage, fmt.Sprintf("coin:%d", uid), time.Minute*3, func() (interface{}, error) {
		// 调用r.db查询数据
		return &UserPO{
			Id:       1,
			Nickname: "藤原拓海",
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return rst.(*UserPO), nil
}

func NewUserRepo(cacheStorage cache.Storage, db db.DB) *UserRepoImpl {
	return &UserRepoImpl{
		cacheStorage: cacheStorage,
		db:           db,
	}
}
