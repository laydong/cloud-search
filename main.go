package main

import (
	"codesearch/conf"
	"codesearch/global"
	"codesearch/router"
	"fmt"
)

func main() {
	err := conf.InitDoAfter()
	if err != nil {
		return
	}

	global.InitMongoDb(conf.ConfInfo.MGConf.Dsn, conf.ConfInfo.MGConf.ConnMaxPoolSize, conf.ConfInfo.MGConf.ConnTimeOut)
	routers := router.Routers()
	address := fmt.Sprintf("0.0.0.0:%v", conf.ConfInfo.AppConf.HttpListen)
	_ = routers.Run(address)
}
