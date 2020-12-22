package service

import (
	"context"

	"server_demo/api/profile"
	"server_demo/internal/biz"
)

// ProfileApp profile业务应用
type ProfileApp struct {
	userBiz *biz.UserBiz
	coinBiz *biz.CoinBiz
}

func NewProfileApp(userBiz *biz.UserBiz, coinBiz *biz.CoinBiz) *ProfileApp {
	return &ProfileApp{
		userBiz: userBiz,
		coinBiz: coinBiz,
	}
}

func (p *ProfileApp) GetUserProfile(ctx context.Context, uid uint64) (*profile.ProfileResp, error) {
	userInfo, err := p.userBiz.GetUserById(ctx, uid)
	if err != nil {
		return nil, err
	}
	coinInfo, err := p.coinBiz.GetUserCoin(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &profile.ProfileResp{
		Nickname: userInfo.Nickname,
		Coins:    coinInfo.Num,
	}, nil
}
