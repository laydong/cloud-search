package server

import (
	"codesearch/model/gitlab"
	"codesearch/model/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func ProjectTag(c *gin.Context, project gitlab.Projects) (resp []gitlab.ProjectsTag) {
	projectName := strings.Replace(project.Name, "_", "-", -1)
	data, _ := new(mysql.ProjectModel).QueryByCode(c, projectName)
	if len(data) > 0 {
		for _, v := range data {
			if v.Tag != "" {
				//resp = append(resp, gitlab.ProjectsTag{
				//	Id:    project.Id,
				//	Code:  projectName,
				//	EnvID: v.EnvId,
				//	Tag:   v.Tag,
				//})
				fmt.Println("添加项目" + project.Name + "---" + v.Tag)
				ProjectsTag <- gitlab.ProjectsTag{
					Id:    project.Id,
					Code:  projectName,
					EnvID: v.EnvId,
					Tag:   v.Tag,
				}
			}
		}
	} else {
		if project.Tag != "" {
			//resp = append(resp, gitlab.ProjectsTag{
			//	Id:    project.Id,
			//	Code:  projectName,
			//	EnvID: 3,
			//	Tag:   project.Tag,
			//})
			fmt.Println("添加项目2" + project.Name + "---" + project.Tag)
			ProjectsTag <- gitlab.ProjectsTag{
				Id:    project.Id,
				Code:  projectName,
				EnvID: 3,
				Tag:   project.Tag,
			}
		}
	}
	return
}

func ProjectList(c *gin.Context, project gitlab.ProjectsTag, task *TaskEsPool) {
	var resp []gitlab.ProjectsFileList
	resp, _ = gitlab.ProjectFileList(c, strconv.Itoa(project.Id), project.Tag, "", 1, "")
	if len(resp) > 0 {
		for _, v := range resp {
			if v.Type == "tree" {
				//ProjectsFileChan <- gitlab.ProjectsFileList{
				//	Id:   v.Id,
				//	Name: v.Name,
				//	Type: v.Type,
				//	Path: v.Path,
				//	Mode: v.Mode,
				//	//Content:      v.Content,
				//	Tag:          project.Tag,
				//	ProjectsName: project.Code,
				//	EnvID:        project.EnvID,
				//}
				fmt.Println("投递1", v)
				ProjectTree(c, strconv.Itoa(project.Id), project.Tag, v.Path, project.Code, 1, project.EnvID, task)
			} else {
				v.Content = gitlab.GetFileRaw(c, strconv.Itoa(project.Id), v.Path, project.Tag)
				if v.Content != "" {
					task.ProjectsFileListChan <- gitlab.ProjectsFileList{
						Id:           v.Id,
						Name:         v.Name,
						Type:         v.Type,
						Path:         v.Path,
						Mode:         v.Mode,
						Content:      v.Content,
						Tag:          project.Tag,
						ProjectsName: project.Code,
						EnvID:        project.EnvID,
					}
					fmt.Println("投递2", v.Path)
				}
			}
		}
		if len(resp) == 100 {
			ProjectTree(c, strconv.Itoa(project.Id), project.Tag, "", project.Code, 2, project.EnvID, task)
		}
	}
	return
}

func ProjectTree(c *gin.Context, projectsId, ref, path, projectsName string, page, envID int, task *TaskEsPool) {
	var resp []gitlab.ProjectsFileList
	resp, _ = gitlab.ProjectFileList(c, projectsId, ref, "", page, path)
	if len(resp) > 0 {
		for _, v := range resp {
			if v.Type == "tree" {
				//ProjectsFileChan <- gitlab.ProjectsFileList{
				//	Id:   v.Id,
				//	Name: v.Name,
				//	Type: v.Type,
				//	Path: v.Path,
				//	Mode: v.Mode,
				//	//Content:      v.Content,
				//	Tag:          ref,
				//	ProjectsName: projectsName,
				//	EnvID:        envID,
				//}
				fmt.Println("投递3", v)
				ProjectTree(c, v.Id, v.Tag, v.Path, v.ProjectsName, 1, v.EnvID, task)
			} else {
				v.Content = gitlab.GetFileRaw(c, projectsId, v.Path, ref)
				if v.Content != "" {
					task.ProjectsFileListChan <- gitlab.ProjectsFileList{
						Id:           v.Id,
						Name:         v.Name,
						Type:         v.Type,
						Path:         v.Path,
						Mode:         v.Mode,
						Content:      v.Content,
						Tag:          ref,
						ProjectsName: projectsName,
						EnvID:        envID,
					}
					fmt.Println("投递4", v.Path)
				}
			}
		}
		if len(resp) == 100 {
			ProjectTree(c, projectsId, ref, path, projectsName, page+1, envID, task)
		}
	}
	return
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
