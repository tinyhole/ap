package config

import (
	micro "github.com/micro/go-micro/config"
)

type srvConfig struct {
	SrvName     string `yaml:"srvName"`
	SrvID       uint32 `yaml:"srvID"`
	ApPort      int32  `yaml:"apPort"`
	LogPath     string `yaml:"logPath"`
	LogLevel    string `yaml:"logLevel"`
	LogFileName string `yaml:"logFileName"`
}

var (
	SrvConfig = srvConfig{}
)

func Init(path string) (err error) {
	err = micro.LoadFile("config.yaml")
	if err != nil {
		return err
	}
	err = micro.Scan(&SrvConfig)
	return err
}
