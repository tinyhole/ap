package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tinyhole/ap/config"
	"github.com/tinyhole/ap/logger"
	"github.com/tinyhole/ap/server"
	"os"
)

var (
	APSrv   server.Server
	RootCmd = &cobra.Command{
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
			APSrv.Start()
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("version: %d.%d.%d.%d", MAJOR, MINOR, PATCH, BUILD))
			os.Exit(0)
		},
	}
)

func init(){
	RootCmd.AddCommand(versionCmd)
}

func Execute() error {
	RootCmd.Execute()
	return nil
}
