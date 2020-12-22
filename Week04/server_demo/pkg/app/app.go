package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Addr             string `yaml:"addr"`
	ReflectionEnable bool   `yaml:"reflectionEnable"`
}

func ReadConf(cfgFile string) ([]byte, error) {
	content, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	return content, nil
}

type ProfileApp struct {
	Conf       *Conf
	ConfRaw    []byte
	GrpcServer *grpc.Server
	SignalChan chan os.Signal
}

func NewProfileApp(confRaw []byte) (*ProfileApp, error) {
	cfg := &Conf{}
	err := yaml.Unmarshal(confRaw, cfg)
	if err != nil {
		return nil, err
	}
	return &ProfileApp{
		Conf:       cfg,
		ConfRaw:    confRaw,
		GrpcServer: grpc.NewServer(),
	}, nil
}

func (app *ProfileApp) Start() error {
	group, ctx := errgroup.WithContext(context.Background())

	if app.Conf.ReflectionEnable {
		// 开启grpc反射，方便调试
		reflection.Register(app.GrpcServer)
	}

	app.SignalChan = make(chan os.Signal, 10)
	group.Go(func() error {
		return app.listenStopSignal(ctx)
	})
	group.Go(func() error {
		return app.serve()
	})

	err := group.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (app *ProfileApp) serve() error {
	lis, err := net.Listen("tcp", app.Conf.Addr)
	if err != nil {
		return err
	}
	fmt.Printf("service serve at %s\n", app.Conf.Addr)
	return app.GrpcServer.Serve(lis)
}
func (app *ProfileApp) listenStopSignal(ctx context.Context) error {
	signal.Notify(app.SignalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-ctx.Done():
	case <-app.SignalChan:
		fmt.Println("receive close signal!")
	}
	app.GrpcServer.GracefulStop()
	return nil
}
