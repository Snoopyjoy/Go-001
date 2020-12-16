package biz

import (
	"context"
	"server_demo/internal/data"
	"server_demo/pkg/cache"
	"server_demo/pkg/db"
)

type UserDO struct {
	Id       uint64
	Nickname string
}

type UserRepo interface {
	FindById(ctx context.Context, uid uint64) (*data.UserPO, error)
}

func NewUserRepo(cacheStorage cache.Storage, db db.DB) UserRepo {
	return data.NewUserRepo(cacheStorage, db)
}

type UserBiz struct {
	repo UserRepo
}

func NewUserBiz(repo UserRepo) *UserBiz {
	return &UserBiz{
		repo: repo,
	}
}

func (b *UserBiz) GetUserById(ctx context.Context, uid uint64) (*UserDO, error) {
	po, err := b.repo.FindById(ctx, uid)
	if err != nil {
		return nil, err
	}

	return UserPOToDO(po), nil
}

func UserPOToDO(po *data.UserPO) *UserDO {
	return &UserDO{
		Id:       po.Id,
		Nickname: po.Nickname,
	}
}
