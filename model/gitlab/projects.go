package gitlab

import (
	"codesearch/global/glogs"
	"codesearch/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
	"time"
)

type Projects struct {
	Id                int                `json:"id"`
	EnvID             int                `json:"env_id"`
	Description       string             `json:"description"`
	Name              string             `json:"name"`
	NameWithNamespace string             `json:"name_with_namespace"`
	Path              string             `json:"path"`
	PathWithNamespace string             `json:"path_with_namespace"`
	CreatedAt         time.Time          `json:"created_at"`
	DefaultBranch     string             `json:"default_branch"`
	TagList           []interface{}      `json:"tag_list"`
	SshUrlToRepo      string             `json:"ssh_url_to_repo"`
	HttpUrlToRepo     string             `json:"http_url_to_repo"`
	WebUrl            string             `json:"web_url"`
	ReadmeUrl         string             `json:"readme_url"`
	AvatarUrl         interface{}        `json:"avatar_url"`
	StarCount         int                `json:"star_count"`
	ForksCount        int                `json:"forks_count"`
	CodeFile          []ProjectsFileList `json:"code_file"`
	Tag               string             `json:"tag"`
}

type ProjectsTag struct {
	Id    int    `json:"id"`
	EnvID int    `json:"env_id"`
	Code  string `json:"code"`
	Tag   string `json:"tag"`
}

type ProjectsFileList struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Path         string `json:"path"`
	Mode         string `json:"mode"`
	Content      string `json:"content"`
	Tag          string `json:"tag"`
	ProjectsName string `json:"projects_name"`
	EnvID        int    `json:"env_id"`
	//ProjectsID   string `json:"projects_id"`
}

// GetPrivateToken 获取git密钥
func GetPrivateToken() string {
	key := viper.GetString("git.Key")
	return "&private_token=" + key
}

// QueryByID 获取项目详情
func QueryByID(c *gin.Context, id int) (resp Projects, err error) {
	urls := viper.GetString("git.url") + "/api/v4/projects/" + strconv.Itoa(id) + GetPrivateToken()
	body, err := utils.HttpGet(c, urls, nil)
	if err != nil {
		glogs.Error("项目信息获取错误", err.Error())
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}

// QueryByName 获取项目信息ID
func QueryByName(c *gin.Context, project string) (resp Projects, err error) {
	urls := "https://gitlab.xthktech.cn/api/v4/projects?search=" + project + GetPrivateToken()
	body, err := utils.HttpGet(c, urls, nil)
	if err != nil {
		glogs.ErrorF(c, "项目信息获取错误", err.Error())
		return
	}
	var data []Projects
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	if len(data) > 0 {
		for _, v := range data {
			if v.Name == project {
				return v, err
			}
		}
	}
	return
}

// GetProjectsList 分页获取项目列表信息
func GetProjectsList(c *gin.Context, page, perPage int) (resp []Projects, err error) {
	urls := viper.GetString("git.url") + "/api/v4/projects?pagination=keyset&page=" + strconv.Itoa(page) + "&per_page=" + strconv.Itoa(perPage) + GetPrivateToken()
	body, err := utils.HttpGet(c, urls, nil)
	if err != nil {
		glogs.Error("项目信息获取错误", err.Error())
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}

/**
@ProjectFileList
@ projectsId 项目ID
@ ref 分支名字
@ recursive 是否递归获取文件 用于获取递归树的布尔值（默认为 false）
@ path 存储库内的路径。用于获取子目录的内容
@ page 页码
*/
func ProjectFileList(c *gin.Context, projectsId, ref, recursive string, page int, path string) (data []ProjectsFileList, err error) {
	urls := viper.GetString("git.url") + "/api/v4/projects/" + projectsId + "/repository/tree?recursive=true&per_page=100&page=" + strconv.Itoa(page) + GetPrivateToken()
	if ref != "" {
		urls = urls + "&ref=" + ref
	}
	if path != "" {
		urls = urls + "&path=" + path
	}
	body, err := utils.HttpGet(c, urls, nil)
	if err != nil {
		glogs.Error(err.Error())
		return
	}
	data = []ProjectsFileList{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		glogs.Error(err.Error())
	}
	return
}

/**
@ projectsId 项目ID
@ ref 分支名字
@ recursive 是否递归获取文件 用于获取递归树的布尔值（默认为 false）
@ path 存储库内的路径。用于获取子目录的内容
@ page 页码
@ envID 环境标识
@ projectsName 项目标识
*/

func ProjectTree(c *gin.Context, projectsId, ref, recursive, path string, page, envID int, projectsName string, req []interface{}) (resp []interface{}) {
	resp = req
	urls := viper.GetString("git.url") + "/api/v4/projects/" + projectsId + "/repository/tree?per_page=100" + "&page=" + strconv.Itoa(page) + GetPrivateToken()
	if ref != "" {
		urls = urls + "&ref=" + ref
	}
	if recursive != "true" {
		urls = urls + "&recursive=" + recursive
	}
	if path != "" {
		urls = urls + "&path=" + path
	}
	body, err := utils.HttpGet(c, urls, nil)
	if err != nil {
		glogs.Error(err.Error())
		return
	}
	var data []ProjectsFileList
	err = json.Unmarshal(body, &data)
	if err != nil {
		glogs.Error(err.Error())
		return
	}
	if len(data) > 0 {
		for _, v := range data {
			if v.Type == "tree" {
				v.Id = projectsId
				v.Tag = ref
				v.ProjectsName = projectsName
				v.EnvID = envID
				resp = ProjectTree(c, projectsId, ref, recursive, v.Path, 1, envID, projectsName, resp)
			} else {
				v.Content = GetFileRaw(c, projectsId, v.Path, ref)
				if v.Content != "" {
					v.Id = projectsId
					v.ProjectsName = projectsName
					v.EnvID = envID
					resp = append(resp, v)
				}
			}

		}
		if len(data) == 100 {
			resp = ProjectTree(c, projectsId, ref, recursive, path, page+1, envID, projectsName, resp)
		}
	}
	return
}

// GetFileRaw 获取单个文件内容 filePath 文件路径 ref 分支名
func GetFileRaw(c *gin.Context, projectsId, filePath, ref string) (resp string) {
	ext := utils.FileExt(filePath)
	if ext == false {
		urls := viper.GetString("git.url") + "/api/v4/projects/" + projectsId + "/repository/files/" + url.QueryEscape(filePath) + "/raw?ref=" + ref + GetPrivateToken()
		res, err := utils.HttpGet(c, urls, nil)
		if err != nil {
			glogs.Error(err.Error())
			return
		}
		return string(res)
	}
	return
}

// GetProjectsList 获取项目列表
//func GetProjectsList(c *gin.Context, envID int, projectsId, ref, recursive, path, projectsName string, page int, task *TaskEsPool) {
//	url := viper.GetString("git.url") + "/api/v4/projects/" + projectsId + "/repository/tree?per_page=100" + "&page=" + strconv.Itoa(page) + GetPrivateToken()
//	if ref != "" {
//		url = url + "&ref=" + ref
//	}
//	if recursive != "true" {
//		url = url + "&recursive=" + recursive
//	}
//	if path != "" {
//		url = url + "&path=" + path
//	}
//	body, err := utils.HttpGet(c, url, nil)
//	if err != nil {
//		glogs.ErrorF(c, err.Error())
//		return
//	}
//	var resp []ProjectsFileList
//	err = json.Unmarshal(body, &resp)
//	if err != nil {
//		glogs.ErrorF(c, err.Error())
//		return
//	}
//	if len(resp) > 0 {
//		for _, v := range resp {
//			v.Id = projectsId
//			v.Tag = ref
//			v.EnvID = envID
//			v.ProjectsName = projectsName
//			if v.Type == "tree" {
//				task.ProjectsFileChan <- v
//			} else {
//				task.ProjectsFileListChan <- v
//			}
//		}
//	}
//	if len(resp) == 100 {
//		GetProjectsList(c, envID, projectsId, ref, recursive, path, projectsName, page+1, task)
//	}
//	return
//}
