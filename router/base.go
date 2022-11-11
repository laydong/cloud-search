package router

import (
	"codesearch/handler"
	"codesearch/middleware"
	"github.com/gin-gonic/gin"
)

type server interface {
	ListenAndServe() error
}

// 初始化总路由

func Routers() *gin.Engine {
	var Router = gin.Default()

	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面4行注释
	//Router.LoadHTMLGlob("./dist/*.html") // npm打包成dist的路径
	//Router.Static("/favicon.ico", "./dist/favicon.ico")
	//Router.Static("/static", "./dist/static")   // dist里面的静态资源
	//Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面

	//Router.StaticFS(viper.GetString("upload.file_url"), http.Dir("."+viper.GetString("upload.file_url"))) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了
	// 跨域
	Router.Use(middleware.Cors()) // 如需跨域可以打开
	Router.NoRoute(middleware.NotRouter())
	Router.GET("init", handler.CodeInit) //代码整体初始化
	Router.GET("list", handler.CodeList) //代码整体初始化
	//ApiRouter(Router)
	//启动rpc
	//router.ListenHttp()
	return Router
}
