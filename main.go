package main

import (
	"codesearch/conf"
	"codesearch/global/glogs"
	"codesearch/global/gstore"
	"codesearch/router"
	"fmt"
)

func main() {
	err := conf.InitDoAfter()
	if err != nil {
		return
	}
	glogs.InitLog()
	gstore.InitMongoDb(conf.ConfInfo.MGConf.Dsn, conf.ConfInfo.MGConf.ConnMaxPoolSize, conf.ConfInfo.MGConf.ConnTimeOut)
	gstore.InitDB(conf.ConfInfo.DBConf.Dsn)
	gstore.InitRdb(conf.ConfInfo.RDConf.Addr, conf.ConfInfo.RDConf.Password, conf.ConfInfo.RDConf.DB)
	routers := router.Routers()
	address := fmt.Sprintf("0.0.0.0:%v", conf.ConfInfo.AppConf.HttpListen)
	_ = routers.Run(address)
}
