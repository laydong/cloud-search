package handler

import (
	"cloud-search/server"
	"cloud-search/utils"
	"github.com/gin-gonic/gin"
)

func CodeInit(c *gin.Context) {
	go server.UpProjects(c)
	utils.OkWithData(nil, c)
}

func CodeList(c *gin.Context) {

	go server.ProjectCodeUp(c, 3, "gxe", "master")
	utils.OkWithData(nil, c)
}
