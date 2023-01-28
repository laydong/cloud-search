package handler

import (
	"cloud-search/server"
	"cloud-search/utils"
	"github.com/gin-gonic/gin"
)

func CodeInit(c *gin.Context) {

}

func CodeList(c *gin.Context) {
	go server.UpProjects(c)
	//go server.ProjectCodeUp(c, 3, "gxe", "master")
	utils.OkWithData(nil, c)
}
