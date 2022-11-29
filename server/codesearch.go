package server

import (
	"codesearch/model/gitlab"
	"codesearch/model/mysql"
	"codesearch/model/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"sync"
)

var fileKey = "project_file_list_"
var pathKey = "project_path_list_"

func ProjectTag(c *gin.Context, project gitlab.Projects, task *TaskEsPool) {
	projectName := strings.Replace(project.Name, "_", "-", -1)
	var data []mysql.DataModel
	data, _ = new(mysql.ProjectModel).QueryByCode(c, projectName)
	var wait sync.WaitGroup
	if len(data) > 0 {
		for _, v := range data {
			wait.Add(1)
			if v.Tag != "" {
				redis.Expire(c, fileKey+projectName+"_"+v.Tag, 0)
				redis.Expire(c, pathKey+projectName+"_"+v.Tag, 0)
				fmt.Println("添加项目" + project.Name + "---" + v.Tag)
				task.ProjectsTag <- gitlab.ProjectsTag{
					Id:    project.Id,
					Code:  projectName,
					EnvID: v.EnvId,
					Tag:   v.Tag,
				}
			}
			wait.Done()
		}
	} else {
		if project.Tag != "" {
			redis.Del(c, fileKey+projectName+"_"+project.Tag)
			redis.Del(c, pathKey+projectName+"_"+project.Tag)
			wait.Add(1)
			fmt.Println("添加项目2" + project.Name + "---" + project.Tag)
			task.ProjectsTag <- gitlab.ProjectsTag{
				Id:    project.Id,
				Code:  projectName,
				EnvID: 3,
				Tag:   project.Tag,
			}
			wait.Done()
		}
	}
	wait.Wait()
	return
}

func ProjectList(c *gin.Context, project gitlab.ProjectsTag, task *TaskEsPool) {
	page := 1
	resp, _ := gitlab.ProjectFileList(c, strconv.Itoa(project.Id), project.Tag, page, "")
	if len(resp) > 0 {
		var wait sync.WaitGroup
		for _, v := range resp {
			wait.Add(1)
			v.Id = strconv.Itoa(project.Id)
			v.Tag = project.Tag
			v.ProjectsName = project.Code
			v.EnvID = project.EnvID
			if v.Type == "tree" {
				if redis.Sismember(c, fileKey+v.ProjectsName+"_"+v.Tag, v.Path) == false {
					fmt.Println("投递1", v.Path)
					task.ProjectsFileChan <- v
					redis.SAdd(c, fileKey+v.ProjectsName+"_"+v.Tag, v.Path)
					redis.Expire(c, fileKey+v.ProjectsName+"_"+v.Tag, viper.GetInt64("git.key_end"))
				}
			} else {
				v.Content = gitlab.GetFileRaw(c, strconv.Itoa(project.Id), v.Path, project.Tag)
				if v.Content != "" {
					if redis.Sismember(c, pathKey+v.ProjectsName+"_"+v.Tag, v.Path) == false {
						fmt.Println("投递2", v.Path)
						task.ProjectsFileListChan <- v
						redis.SAdd(c, pathKey+v.ProjectsName+"_"+v.Tag, v.Path)
						redis.Expire(c, pathKey+v.ProjectsName+"_"+v.Tag, viper.GetInt64("git.key_end"))
					}
				}
			}
			wait.Done()
		}
		wait.Wait()
		if len(resp) == 100 {
			ProjectTree(c, strconv.Itoa(project.Id), project.Tag, "", project.Code, page+1, project.EnvID, task)
		}
	}

	return
}

func ProjectTree(c *gin.Context, projectsId, ref, filePath, projectsName string, page, envID int, task *TaskEsPool) {
	resp, _ := gitlab.ProjectFileList(c, projectsId, ref, page, filePath)
	if len(resp) > 0 {
		var wait sync.WaitGroup
		for _, v := range resp {
			v.Id = projectsId
			v.Tag = ref
			v.ProjectsName = projectsName
			v.EnvID = envID
			wait.Add(1)
			if v.Type == "tree" {
				if redis.Sismember(c, fileKey+v.ProjectsName+"_"+v.Tag, v.Path) == false {
					fmt.Println("投递3", v.Path)
					task.ProjectsFileChan <- v
					redis.SAdd(c, fileKey+v.ProjectsName+"_"+v.Tag, v.Path)
					redis.Expire(c, fileKey+v.ProjectsName+"_"+v.Tag, viper.GetInt64("git.key_end"))
				}
				////ProjectsFileChan <- v
				//fmt.Println("投递3", v)
				//ProjectTree(c, v.Id, v.Tag, v.Path, v.ProjectsName, 1, v.EnvID)
			} else {
				v.Content = gitlab.GetFileRaw(c, projectsId, v.Path, ref)
				if v.Content != "" {
					if redis.Sismember(c, pathKey+v.ProjectsName+"_"+v.Tag, v.Path) == false {
						fmt.Println("投递4", v.Path)
						task.ProjectsFileListChan <- v
						redis.SAdd(c, pathKey+v.ProjectsName+"_"+v.Tag, v.Path)
						redis.Expire(c, pathKey+v.ProjectsName+"_"+v.Tag, viper.GetInt64("git.key_end"))
					}
				}
			}
			wait.Done()
		}
		wait.Wait()
		if len(resp) == 100 {
			ProjectTree(c, projectsId, ref, filePath, projectsName, page+1, envID, task)
		}
	}

	return
}
