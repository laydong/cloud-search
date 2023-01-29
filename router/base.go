package router

import (
	"cloud-search/handler"
	"github.com/gin-gonic/gin"
	"github.com/laydong/toolpkg/middleware"
)

type server interface {
	ListenAndServe() error
}

// 初始化总路由

// 初始化总路由

func Routers() *gin.Engine {
	var Router = gin.Default()
	// 跨域  如需跨域可以打开
	Router.Use(middleware.Cors())
	// 记录API日志
	Router.NoRoute(middleware.NotRouter())
	Router.NoMethod(middleware.NoMethodHandle())
	Router.GET("/init", handler.CodeInit)
	return Router
}
