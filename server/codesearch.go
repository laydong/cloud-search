package server

import (
	"cloud-search/model/gitlab"
	"cloud-search/model/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
)

func ProjectTag(c *gin.Context, project gitlab.Projects, task *TaskEsPool) {
	projectName := strings.Replace(project.Name, "_", "-", -1)
	var data []mysql.DataModel
	data, _ = new(mysql.ProjectModel).QueryByCode(c, projectName)
	var wait sync.WaitGroup
	if len(data) > 0 {
		for _, v := range data {
			wait.Add(1)
			if v.Tag != "" {
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
			v.ProjectName = project.Code
			v.EnvID = project.EnvID
			if v.Type == "tree" {
				fmt.Println("投递1", v.Path)
				task.ProjectsFileChan <- v
			} else {
				v.Content = gitlab.GetFileRaw(c, strconv.Itoa(project.Id), v.Path, project.Tag)
				if v.Content != "" {
					fmt.Println("投递2", v.Path)
					task.ProjectsPathChan <- v
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

func ProjectTree(c *gin.Context, projectsId, ref, filePath, projectName string, page, envID int, task *TaskEsPool) {
	resp, _ := gitlab.ProjectFileList(c, projectsId, ref, page, filePath)
	if len(resp) > 0 {
		var wait sync.WaitGroup
		for _, v := range resp {
			v.Id = projectsId
			v.Tag = ref
			v.ProjectName = projectName
			v.EnvID = envID
			wait.Add(1)
			if v.Type == "tree" {
				fmt.Println("投递3", v.Path)
				task.ProjectsFileChan <- v
			} else {
				v.Content = gitlab.GetFileRaw(c, projectsId, v.Path, ref)
				if v.Content != "" {
					fmt.Println("投递4", v.Path)
					task.ProjectsPathChan <- v
				}
			}
			wait.Done()
		}
		wait.Wait()
		if len(resp) == 100 {
			ProjectTree(c, projectsId, ref, filePath, projectName, page+1, envID, task)
		}
	}
	return
}
