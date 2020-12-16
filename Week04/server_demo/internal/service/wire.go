//+build wireinject

package service

import (
	"server_demo/internal/biz"
	"server_demo/pkg/cache"
	"server_demo/pkg/db"

	"github.com/google/wire"
)

func InitializeService(cfgRaw []byte) (*ProfileService, error) {
	wire.Build(NewProfileService,
		biz.NewCoinBiz,
		biz.NewUserBiz,
		biz.NewUserRepo,
		biz.NewCoinRepo,
		db.NewDB,
		cache.NewCacheClient,
		wire.FieldsOf(new(*Conf), "DBConf", "CacheConf"),
		ParseConf)
	return &ProfileService{}, nil
}
