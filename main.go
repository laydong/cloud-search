package main

import (
	"cloud-search/conf"
	"cloud-search/global"
	"cloud-search/router"
	"fmt"
	"github.com/laydong/toolpkg"
	"github.com/laydong/toolpkg/db"
)

func main() {
	//初始化配置
	err := conf.InitDoAfter()
	if err != nil {
		panic(err)
	}
	//初始化 日志服务
	toolpkg.InitLog(toolpkg.AppConf{
		AppName: conf.ConfInfo.AppConf.Name,
		AppMode: conf.ConfInfo.AppConf.Mode,
	})
	global.DB, err = db.InitDB(conf.ConfInfo.DBConf.Dsn)
	if err != nil {
		panic(err)
	}
	global.Rdb, err = db.InitRdb(conf.ConfInfo.RDConf.Addr, conf.ConfInfo.RDConf.Password, conf.ConfInfo.RDConf.DB)
	if err != nil {
		panic(err)
	}
	global.Mdb, err = db.InitMongoDb(conf.ConfInfo.MGConf.Dsn, conf.ConfInfo.MGConf.ConnMaxPoolSize, conf.ConfInfo.MGConf.ConnTimeOut)
	if err != nil {
		panic(err)
	}
	//初始化路由
	routers := router.Routers()
	address := fmt.Sprintf("0.0.0.0:%v", "80")
	_ = routers.Run(address)
}
