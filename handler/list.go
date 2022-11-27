package handler

import (
	"codesearch/server"
	"codesearch/utils"
	"github.com/gin-gonic/gin"
)

func CodeInit(c *gin.Context) {

}

func CodeList(c *gin.Context) {
	go server.UpProjects(c)
	//server.ProjectCodeUp(c, "gxe")
	utils.OkWithData(nil, c)
}
