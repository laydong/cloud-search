package server

import (
	"codesearch/global"
	"codesearch/global/glogs"
	"codesearch/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"
)

// GetPrivateToken 获取git密钥
func GetPrivateToken() string {
	key := viper.GetString("git.Key")
	return "&private_token=" + key
}

// GetProjectsID 获取项目信息ID
func GetProjectsID(c *gin.Context, project string) (resp Projects, err error) {
	url := "https://gitlab.xthktech.cn/api/v4/projects?search=" + project + GetPrivateToken()
	body, err := utils.HttpGet(c, url, nil)
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

// GetProjects 获取项目列表信息
func GetProjects(c *gin.Context, page, perPage int) (resp []Projects, err error) {
	url := viper.GetString("git.url") + "/api/v4/projects?pagination=keyset&page=" + strconv.Itoa(page) + "&per_page=" + strconv.Itoa(perPage) + GetPrivateToken()
	//res, err := http.Get(url)
	body, err := utils.HttpGet(c, url, nil)
	if err != nil {
		glogs.Error("项目信息获取错误", err.Error())
		return
	}
	//defer res.Body.Close()
	//data, err := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	return resp, err
}

//func GetProjectsFile(c *gin.Context, projectsId, ref, recursive, path string, projectsName string, task *TaskEsPool, wait *sync.WaitGroup) (resp []ProjectsFileList) {
//	defer wait.Done()
//	GetTree(c, projectsId, ref, recursive, path, 1, projectsName, task)
//	return
//}

//// UpDataFile 合并更新数据
//func UpDataFile(c *gin.Context, projectsId, project, path, version string, data []esdb.EsInfo) (res []esdb.EsInfo) {
//	raw := GetFileRaw(c, projectsId, path, version)
//	if raw != "" {
//		split := strings.Split(raw, "\n")
//		if len(split) > 0 {
//			for k, v1 := range split {
//				if v1 != "" {
//					line := strconv.Itoa(k + 1)
//					id, _ := utils.Md5Encrypt(projectsId + path + line)
//					esInfo := esdb.EsInfo{
//						ID:         id,
//						User:       "",
//						UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
//						Version:    version,
//						Project:    project,
//						File:       path,
//						Line:       line,
//						Content:    v1,
//					}
//					data = append(data, esInfo)
//				}
//			}
//		}
//	}
//	return data
//}
//
//// GetFileList 创建es储存
//func GetFileList(c *gin.Context, projectsId, ref, path, projects, namespace string, respData []esdb.EsInfo) (data []esdb.EsInfo, err error) {
//	var resData []esdb.ProjectsFileList
//	resp := GetTree(c, projectsId, ref, "false", path, 1, resData)
//	for _, v := range resp {
//		if v.Type == "tree" {
//			res, _ := GetFileList(c, projectsId, ref, v.Path, projects, namespace, respData)
//			for _, v1 := range res {
//				data = append(data, v1)
//			}
//		} else {
//			if v.Path != "" {
//				v.Content = GetFileRaw(c, projectsId, v.Path, ref)
//				if v.Content != "" {
//					split := strings.Split(v.Content, "\n")
//					if len(split) > 0 {
//						for k, v1 := range split {
//							line := strconv.Itoa(k + 1)
//							id, _ := utils.Md5Encrypt(projectsId + v.Path + line)
//							esInfo := esdb.EsInfo{
//								ID:         id,
//								User:       "",
//								UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
//								Version:    ref,
//								Project:    projects,
//								File:       v.Path,
//								Line:       line,
//								Content:    v1,
//							}
//							data = append(data, esInfo)
//						}
//					}
//				}
//			}
//		}
//	}
//	return data, err
//}
//
/**
@ GetTree 循环获取代码文件列表
@ projectsId 项目ID
@ ref 分支名字
@ recursive 是否递归获取文件 用于获取递归树的布尔值（默认为 false）
@ path 存储库内的路径。用于获取子目录的内容
@ page 页码
@ projectsName 项目标识
*/
//func GetTree(c *gin.Context, projectsId, ref, recursive, path string, page int, projectsName string, task *TaskEsPool) {
//	data, _ := GetProjectsList(c, projectsId, ref, recursive, path, page)
//	if len(data) > 0 {
//		for _, v := range data {
//			v.Id = projectsId
//			v.ProjectsName = projectsName
//			v.Tag = ref
//			if v.Type == "tree" {
//				task.ProjectsFileChan <- v
//			} else {
//				v.Content = GetFileRaw(c, projectsId, v.Path, ref)
//				if v.Content != "" {
//					mongo.AddOne(c, "code_list", v)
//				}
//			}
//		}
//		if len(data) == 100 {
//			GetTree(c, projectsId, ref, recursive, path, page+1, projectsName, task)
//		}
//	}
//	return
//}

/**
@ projectsId 项目ID
@ ref 分支名字
@ recursive 是否递归获取文件 用于获取递归树的布尔值（默认为 false）
@ path 存储库内的路径。用于获取子目录的内容
@ page 页码
@ projectsName 项目标识
*/

func ProjectTree(c *gin.Context, projectsId, ref, recursive, path string, page int, projectsName string, req []interface{}) (resp []interface{}) {
	resp = req
	url := viper.GetString("git.url") + "/api/v4/projects/" + projectsId + "/repository/tree?per_page=100" + "&page=" + strconv.Itoa(page) + GetPrivateToken()
	if ref != "" {
		url = url + "&ref=" + ref
	}
	if recursive != "true" {
		url = url + "&recursive=" + recursive
	}
	if path != "" {
		url = url + "&path=" + path
	}
	body, err := utils.HttpGet(c, url, nil)
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
				resp = ProjectTree(c, projectsId, ref, recursive, v.Path, page, projectsName, resp)
			} else {
				v.Content = GetFileRaw(c, projectsId, v.Path, ref)
				if v.Content != "" {
					v.Id = projectsId
					v.ProjectsName = projectsName
					resp = append(resp, v)
				}
			}

		}
		if len(data) == 100 {
			resp = ProjectTree(c, projectsId, ref, recursive, path, page+1, projectsName, resp)
		}
	}
	return
}

// GetFileRaw 获取单个文件内容 filePath 文件路径 ref 分支名
func GetFileRaw(c *gin.Context, projectsId, filePath, ref string) (resp string) {
	ext := utils.FileExt(filePath)
	if ext == false {
		url := viper.GetString("git.url") + "/api/v4/projects/" + projectsId + "/repository/files/" + url2.QueryEscape(filePath) + "/raw?ref=" + ref + GetPrivateToken()
		res, err := http.Get(url)
		if err != nil {
			glogs.Error(err.Error())
			return
		}
		defer res.Body.Close()
		var headType = res.Header.Get("Content-Type")
		if utils.SupeString(headType) == true {
			glogs.ErrorF(c, "编码错误")
			return
		}
		str, err := io.ReadAll(res.Body)
		if err != nil {
			glogs.Error(err.Error())
			return
		}
		return string(str)
	}
	return
}

// GetProjectsList 获取项目列表
func GetProjectsList(c *gin.Context, envID int, projectsId, ref, recursive, path, projectsName string, page int, task *TaskEsPool) {
	url := viper.GetString("git.url") + "/api/v4/projects/" + projectsId + "/repository/tree?per_page=100" + "&page=" + strconv.Itoa(page) + GetPrivateToken()
	if ref != "" {
		url = url + "&ref=" + ref
	}
	if recursive != "true" {
		url = url + "&recursive=" + recursive
	}
	if path != "" {
		url = url + "&path=" + path
	}
	body, err := utils.HttpGet(c, url, nil)
	if err != nil {
		glogs.ErrorF(c, err.Error())
		return
	}
	var resp []ProjectsFileList
	err = json.Unmarshal(body, &resp)
	if err != nil {
		glogs.ErrorF(c, err.Error())
		return
	}
	if len(resp) > 0 {
		for _, v := range resp {
			v.Id = projectsId
			v.Tag = ref
			v.EnvID = envID
			v.ProjectsName = projectsName
			if v.Type == "tree" {
				task.ProjectsFileChan <- v
			} else {
				task.ProjectsFileListChan <- v
			}
		}
	}
	if len(resp) == 100 {
		GetProjectsList(c, envID, projectsId, ref, recursive, path, projectsName, page+1, task)
	}
	return
}

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

func GetCodeComment(c *gin.Context, str string) (resp interface{}) {
	collection := global.Mdb.Database("app_gitlab").Collection("code_list")
	//match := bson.M{"$match": bson.M{"content": bson.M{"$regex": primitive.Regex{Pattern: str}}}}
	filter := bson.M{"content": bson.M{"$regex": primitive.Regex{Pattern: str}}}
	find, err := collection.Find(c, filter)
	if err != nil {
		return
	}
	var data []ProjectsFileList
	find.All(c, data)
	return data
}
