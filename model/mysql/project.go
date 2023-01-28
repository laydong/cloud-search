package mysql

import (
	"cloud-search/global"
	"github.com/gin-gonic/gin"
)

type ProjectModel struct {
	global.ModelID
	Title       string `gorm:"column:title;NOT NULL" json:"title"`             // 项目名称
	Code        string `gorm:"column:code;NOT NULL" json:"code"`               // 项目标识码
	Des         string `gorm:"column:des" json:"des"`                          // 项目描述
	Status      int    `gorm:"column:status;default:2;NOT NULL" json:"status"` // 状态 1正常 2禁用
	Library     string `gorm:"column:library" json:"library"`                  // 仓库名称
	LibraryUrl  string `gorm:"column:library_url" json:"library_url"`          // 仓库地址
	ProjectType string `gorm:"column:project_type" json:"project_type"`        // 项目类型
	EnvId       int    `gorm:"column:env_id;default:0;NOT NULL" json:"env_id"` // 环境ID
	global.ModelTime
}

type DataModel struct {
	Code  string `json:"code"`   // 项目标识码
	Tag   string `json:"tag"`    // 版本
	EnvId int    `json:"env_id"` // 环境ID
}

func (m *ProjectModel) TableName() string {
	return "xthk_project"
}

func (m *ProjectModel) QueryByCode(c *gin.Context, code string) (data []DataModel, err error) {
	err = global.GetDB(c).Model(ProjectModel{}).
		Where("status", 1).
		Where("code", code).
		Select("distinct xthk_project.`code`,xthk_project.env_id,(SELECT tag FROM xthk_project_chart WHERE project_id=xthk_project.id ORDER BY updated_at DESC LIMIT 1) AS tag").
		Find(&data).Error
	return
}
