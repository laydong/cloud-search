package server

import (
	"codesearch/global/glogs"
	"codesearch/model/mongo"
	"codesearch/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TaskEsPool struct {
	projectNumber            *atomic.Int64
	projectFileNumber        *atomic.Int64
	projectPathNumber        *atomic.Int64
	projectFileAllRowsNumber *atomic.Int64
	start                    time.Time
	wait                     *sync.WaitGroup
	TaskPool                 *ants.Pool
	ProjectChan              chan Projects
	ProjectsFileListChan     chan ProjectsFileList
	ProjectsFileChan         chan ProjectsFileList
}

// 批量更新项目数据
func UpProjects(c *gin.Context) (err error) {
	var task = &TaskEsPool{
		projectNumber:        atomic.NewInt64(0),
		start:                time.Now(),
		wait:                 &sync.WaitGroup{},
		ProjectChan:          make(chan Projects, 1000),
		ProjectsFileListChan: make(chan ProjectsFileList, 10000),
		ProjectsFileChan:     make(chan ProjectsFileList, 10000),
	}
	task.TaskPool, err = ants.NewPool(runtime.NumCPU() * 5)
	if err != nil {
		return err
	}

	defer task.TaskPool.Release()
	startProject(c, task)
	task.wait.Add(3)
	go startProjectFileListData(c, task)
	go startProjectFileData(c, task)
	go startEsAddBatch(c, task)

	task.wait.Wait()
	stop := time.Since(task.start)
	fmt.Printf("task over, 耗时：%v", stop)

	return nil
}

// ProjectCodeUp 更新指定项目
func ProjectCodeUp(c *gin.Context, envID int, code, tag string) (err error) {
	project, _ := GetProjectsID(c, ProjectReplace(c, code))
	if project.Id > 0 {
		var task = &TaskEsPool{
			projectNumber:        atomic.NewInt64(0),
			start:                time.Now(),
			wait:                 &sync.WaitGroup{},
			ProjectsFileListChan: make(chan ProjectsFileList, 10000),
			ProjectsFileChan:     make(chan ProjectsFileList, 10000),
		}
		task.TaskPool, err = ants.NewPool(runtime.NumCPU() * 5)
		if err != nil {
			return err
		}
		defer task.TaskPool.Release()
		//删除项目数据
		err = mongo.DelCodeAll(c, project.Name)
		if err != nil {
			glogs.ErrorF(c, err.Error())
		}
		GetProjectsList(c, envID, strconv.Itoa(project.Id), tag, "true", "", project.Name, 1, task)
		task.wait.Add(2)
		//go startProjectFileListData(c, task)
		go startProjectFileData(c, task)
		go startEsAddBatch(c, task)

		task.wait.Wait()
		stop := time.Since(task.start)
		fmt.Printf("task over, 耗时：%v", stop)

	} else {
		err = errors.New("更新失败-项目不存在")
	}
	return
}

//StartProject 获取项目信息
func startProject(c *gin.Context, task *TaskEsPool) {
	//all, err := new(mysql.ProjectModel).GetProjectAll(c)
	//if err != nil || len(all) == 0 {
	//	return
	//}
	//var projects []interface{}
	////获取项目信息
	//if len(all) > 0 {
	//	//删除原数据
	//	mongo.CodeInit(c)
	//	for _, v := range all {
	//		project, _ := GetProjectsID(c, ProjectReplace(c, v.Code))
	//		if project.Id > 0 {
	//			project.EnvID = int(v.EnvID)
	//			project.Tag = v.Tag
	//			projects = append(projects, project)
	//			//task.projectFileNumber.Add(1)
	//			task.ProjectChan <- project
	//		}
	//	}
	//}
	//if len(projects) > 0 {
	//	//写入数据库
	//	mongo.AddALL(c, viper.GetString("git.gitlab_project_name"), projects)
	//}
}

//ProjectReplace 项目下划线替换
func ProjectReplace(c *gin.Context, project string) (resp string) {
	str := viper.GetString("git.project_replace")
	if str != "" {
		strArr := strings.Split(str, `,`)
		if len(strArr) > 0 {
			if utils.InSliceString(project, strArr) == true {
				//执行替换
				project = strings.Replace(project, "-", "_", -1)
			}
		}
	}
	return project
}

//项目文件列表
func startProjectFileListData(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectChan)
	defer task.wait.Done()
	var wait sync.WaitGroup
	defer wait.Done()
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectChan {
		wait.Add(1)
		//err := task.TaskPool.Submit(func() {
		GetProjectsList(c, x.EnvID, strconv.Itoa(x.Id), x.Tag, "true", "", x.Name, 1, task)
		//})
		//if err != nil {
		//	glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
		//}
		wait.Done()
	}
	wait.Wait()
}

//func SStartProjectFileListData(c *gin.Context, task *TaskEsPool) {
//	//defer close(task.ProjectChan)
//	defer task.wait.Done()
//	var wait sync.WaitGroup
//	defer wait.Done()
//	// 从ch中接收值并赋值给变量x
//	for x := range task.ProjectChan {
//		wait.Add(1)
//		err := task.TaskPool.Submit(func() {
//			GetProjectsFile(c, strconv.Itoa(x.Id), x.DefaultBranch, "true", x.Path, x.Name, task, &wait)
//			fmt.Println("项目文件获取：", x)
//			//GetTree(c, strconv.Itoa(x.Id), x.DefaultBranch, "true", x.Path, 1, x.Name, task, &wait)
//		})
//		if err != nil {
//			glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
//		}
//	}
//	wait.Wait()
//}

//项目文件解读
func startProjectFileData(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectsFileChan)
	defer task.wait.Done()
	var wait sync.WaitGroup
	defer wait.Done()
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectsFileChan {
		wait.Add(1)
		//err := task.TaskPool.Submit(func() {
		//	GetProjectsFile(c, x.Id, x.Tag, "true", x.Path, x.ProjectsName, task, &wait)
		//	fmt.Println("项目内部获取：", x)
		//})
		//if err != nil {
		//	glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
		//}
		//if x.Type == "tree" {
		//err := task.TaskPool.Submit(func() {
		GetProjectsList(c, x.EnvID, x.Id, x.Tag, "true", x.Path, x.ProjectsName, 1, task)
		//})
		//if err != nil {
		//	glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
		//}
		//} else {
		//	task.ProjectsFileListChan <- x
		//}
		wait.Done()

	}
	wait.Wait()
}

func startEsAddBatch(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectsFileListChan)
	defer task.wait.Done()
	var wait sync.WaitGroup
	defer wait.Done()
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectsFileListChan {
		wait.Add(1)
		//err := task.TaskPool.Submit(func() {
		raw := GetFileRaw(c, x.Id, x.Path, x.Tag)
		if raw != "" {
			x.Content = raw
			mongo.AddOne(c, viper.GetString("git.gitlab_code_name"), x)
		}
		//})
		//if err != nil {
		//	glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
		//}
		wait.Done()
	}
	wait.Wait()
}
