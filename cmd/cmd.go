package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tinyhole/ap/config"
	"github.com/tinyhole/ap/logger"
	"github.com/tinyhole/ap/server"
)

var (
	APSrv   server.Server
	RootCmd = &cobra.Command{
		Use:   "",
		Short: "run server",
		Long:  "run server",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err error
			)
			//加载配置文件
			err = config.Init("config.yaml")
			if err != nil {
				//return errors.Wrap(err, "config Init failed ")
			}
			//配置日志
			logger.Init(config.SrvConfig.LogLevel, config.SrvConfig.LogPath, config.SrvConfig.LogFileName)
			//启动ap服务
			APSrv = server.NewAPServer()
			APSrv.Init(server.WithLocalAddr(fmt.Sprintf(":%d", config.SrvConfig.ApPort)),
				server.WithSrvID(config.SrvConfig.SrvID))
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("version: %d.%d.%d.%d", MAJOR, MINOR, PATCH, BUILD))
		},
	}
)

func Execute() error {
	RootCmd.Execute()
	return nil
}
