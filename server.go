package main

import (
	"context"
	"fmt"
	"os"

	"config-platform/resolver"
)

var prefix = "config_ceter_"

type ConfigCenterServer struct {
	register *resolver.ETCDRegister
}

func NewConfigCenterServer(r *resolver.ETCDRegister) *ConfigCenterServer {
	return &ConfigCenterServer{register: r}
}

func (s *ConfigCenterServer) Start(ctx context.Context) error {
	hostIp, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("get host ip err when register %w", err)
	}
	return s.register.Register(ctx, prefix+hostIp, hostIp)
}

func (s *ConfigCenterServer) Stop(ctx context.Context) error {
	hostIp, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("get host ip err when stop %w", err)
	}
	return s.register.LogOut(ctx, prefix+hostIp, hostIp)
}
