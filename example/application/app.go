package application

import (
	rclient "config-platform/client"
)

type AppForTest struct {
	configPlatform *rclient.Client
}

func NewAppForTest(c *rclient.Client) *AppForTest {
	return &AppForTest{
		configPlatform: c,
	}
}

func (a *AppForTest) Start() {
	a.configPlatform.Start()
}
