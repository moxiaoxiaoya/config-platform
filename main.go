package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"config-platform/config"
	"config-platform/resolver"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	ctx := context.Background()
	conf := initConfig()
	ipPort := conf.IP + ":" + conf.Port
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ipPort},
		DialTimeout: time.Duration(conf.TimeOut) * time.Second,
	})
	if err != nil {
		panic(err)
	}
	configCenter := NewConfigCenterServer(resolver.NewETCDRegister(etcdCli))

	if err = configCenter.Start(ctx); err != nil {
		panic(err)
	}

	// wait for shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	if err = configCenter.Stop(ctx); err != nil {
		// TODO:log here
	}
}

func initConfig() config.ETCDConfig {
	var conf config.ETCDConfig
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return conf
}
