// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package service

import (
	"server_demo/internal/biz"
	"server_demo/pkg/cache"
	"server_demo/pkg/db"
)

// Injectors from wire.go:

func InitializeService(cfgRaw []byte) (*ProfileService, error) {
	conf, err := ParseConf(cfgRaw)
	if err != nil {
		return nil, err
	}
	cacheConf := conf.CacheConf
	storage, err := cache.NewCacheClient(cacheConf)
	if err != nil {
		return nil, err
	}
	dbConf := conf.DBConf
	dbDB, err := db.NewDB(dbConf)
	if err != nil {
		return nil, err
	}
	userRepo := biz.NewUserRepo(storage, dbDB)
	userBiz := biz.NewUserBiz(userRepo)
	coinRepo := biz.NewCoinRepo(storage, dbDB)
	coinBiz := biz.NewCoinBiz(coinRepo)
	profileApp := NewProfileApp(userBiz, coinBiz)
	profileService := NewProfileService(conf, profileApp)
	return profileService, nil
}
