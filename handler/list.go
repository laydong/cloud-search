package handler

import (
	"codesearch/server"
	"codesearch/utils"
	"github.com/gin-gonic/gin"
)

func CodeInit(c *gin.Context) {

}

func CodeList(c *gin.Context) {
	resp, _ := server.GetProjects(c, 6, 100)
	utils.OkWithData(resp, c)
}
