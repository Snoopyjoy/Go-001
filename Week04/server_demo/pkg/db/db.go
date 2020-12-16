package db

type DBConf struct {
	Addr   string `yaml:"addr"`
	Passwd string `yaml:"passwd"`
}

type DB interface {
}

type db struct {
}

func NewDB(cfg *DBConf) (DB, error) {
	return &db{}, nil
}
