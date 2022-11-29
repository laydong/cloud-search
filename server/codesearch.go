package server

import (
	"codesearch/model/gitlab"
	"codesearch/model/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
)

func ProjectTag(c *gin.Context, project gitlab.Projects) {
	projectName := strings.Replace(project.Name, "_", "-", -1)
	var data []mysql.DataModel
	data, _ = new(mysql.ProjectModel).QueryByCode(c, projectName)
	var wait sync.WaitGroup
	if len(data) > 0 {
		for _, v := range data {
			wait.Add(1)
			if v.Tag != "" {

				fmt.Println("添加项目" + project.Name + "---" + v.Tag)
				ProjectsTag <- gitlab.ProjectsTag{
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
			ProjectsTag <- gitlab.ProjectsTag{
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

func ProjectList(c *gin.Context, project gitlab.ProjectsTag) {
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
				//ProjectsFileChan <- v
				fmt.Println("投递1", v)
				ProjectTree(c, strconv.Itoa(project.Id), project.Tag, v.Path, project.Code, 1, project.EnvID)
			} else {
				v.Content = gitlab.GetFileRaw(c, strconv.Itoa(project.Id), v.Path, project.Tag)
				if v.Content != "" {
					ProjectsFileListChan <- v
					fmt.Println("投递2", v.Path)
				}
			}
			wait.Done()
		}
		wait.Wait()
		if len(resp) == 100 {
			ProjectTree(c, strconv.Itoa(project.Id), project.Tag, "", project.Code, page+1, project.EnvID)
		}
	}

	return
}

func ProjectTree(c *gin.Context, projectsId, ref, filePath, projectsName string, page, envID int) {
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
				//ProjectsFileChan <- v
				fmt.Println("投递3", v)
				ProjectTree(c, v.Id, v.Tag, v.Path, v.ProjectsName, 1, v.EnvID)
			} else {
				v.Content = gitlab.GetFileRaw(c, projectsId, v.Path, ref)
				if v.Content != "" {
					ProjectsFileListChan <- v
					fmt.Println("投递4", v.Path)
				}
			}
			wait.Done()
		}
		wait.Wait()
		if len(resp) == 100 {
			ProjectTree(c, projectsId, ref, filePath, projectsName, page+1, envID)
		}
	}

	return
}
