package service

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"

	"server_demo/api/profile"
	"server_demo/internal/biz"
	"server_demo/pkg/cache"
	"server_demo/pkg/db"
)

type Conf struct {
	DBConf    *db.DBConf       `yaml:"db"`
	CacheConf *cache.CacheConf `yaml:"cache"`
}

type ProfileService struct {
	Conf *Conf

	userBiz *biz.UserBiz
	coinBiz *biz.CoinBiz

	profileApp *ProfileApp

	profile.UnimplementedServiceServer
}

func NewProfileService(cfg *Conf, profileApp *ProfileApp) *ProfileService {
	return &ProfileService{
		Conf:       cfg,
		profileApp: profileApp,
	}
}

func ParseConf(cfgRaw []byte) (*Conf, error) {
	cfg := &Conf{}
	err := yaml.Unmarshal(cfgRaw, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *ProfileService) Profile(ctx context.Context, in *profile.ProfileArgs) (*profile.ProfileResp, error) {
	ret, err := s.profileApp.GetUserProfile(ctx, in.Uid)
	if err != nil {
		fmt.Printf("GetUserProfile err %+v", err)
		return nil, err
	}
	return ret, nil
}
