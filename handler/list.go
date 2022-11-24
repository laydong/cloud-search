package handler

import (
	"codesearch/server"
	"codesearch/utils"
	"github.com/gin-gonic/gin"
)

func CodeInit(c *gin.Context) {

}

func CodeList(c *gin.Context) {
	server.UpProjects(c)
	//server.ProjectCodeUp(c, "devops")
	utils.OkWithData(nil, c)
}
