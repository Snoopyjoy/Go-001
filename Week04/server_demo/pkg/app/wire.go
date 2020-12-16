//+build wireinject

package app

import (
	"github.com/google/wire"
)

func InitializeApp(cfgFile string) (*ProfileApp, error) {
	wire.Build(NewProfileApp, ReadConf)
	return &ProfileApp{}, nil
}
