package server

import (
	"codesearch/global/glogs"
	"codesearch/model/gitlab"
	"codesearch/model/mongo"
	"codesearch/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"runtime"
	"strings"
	"sync"
	"time"
)

type TaskEsPool struct {
	projectNumber     *atomic.Int64
	projectFileNumber *atomic.Int64
	projectPathNumber *atomic.Int64
	start             time.Time
	wait              *sync.WaitGroup
	TaskPool          *ants.Pool
	ProjectChan       chan gitlab.Projects
	ProjectsTag       chan gitlab.ProjectsTag
	ProjectsPathChan  chan gitlab.ProjectsFileList
	ProjectsFileChan  chan gitlab.ProjectsFileList
}

// 批量更新项目数据
func UpProjects(c *gin.Context) (err error) {
	var task = &TaskEsPool{
		projectPathNumber: atomic.NewInt64(0),
		projectFileNumber: atomic.NewInt64(0),
		projectNumber:     atomic.NewInt64(0),
		start:             time.Now(),
		wait:              &sync.WaitGroup{},
		ProjectChan:       make(chan gitlab.Projects, 1000),
		ProjectsTag:       make(chan gitlab.ProjectsTag, 10000),
		ProjectsPathChan:  make(chan gitlab.ProjectsFileList, 10000),
		ProjectsFileChan:  make(chan gitlab.ProjectsFileList, 10000),
	}
	task.TaskPool, err = ants.NewPool(runtime.NumCPU() * 1000)
	if err != nil {
		return err
	}
	defer task.TaskPool.Release()
	startProject(c, task)
	task.TaskPool.Submit(func() {
		task.wait.Add(1)
		startProjectTag(c, task)
	})
	task.TaskPool.Submit(func() {
		task.wait.Add(1)
		startProjectFileListData(c, task)
	})
	task.TaskPool.Submit(func() {
		task.wait.Add(1)
		startProjectFileData(c, task)
	})
	task.TaskPool.Submit(func() {
		task.wait.Add(1)
		startEsAddBatch(c, task)
	})
	task.wait.Wait()
	stop := time.Since(task.start)
	fmt.Printf("项目数 ：%v", task.projectNumber.String())
	fmt.Printf("task over, 耗时：%v", stop)

	return nil
}

//ProjectCodeUp 更新指定项目
func ProjectCodeUp(c *gin.Context, envID int, code, tag string) (err error) {
	project, _ := gitlab.QueryByName(c, ProjectReplace(c, code))
	if project.Id > 0 {
		var task = &TaskEsPool{
			projectPathNumber: atomic.NewInt64(0),
			projectFileNumber: atomic.NewInt64(0),
			projectNumber:     atomic.NewInt64(0),
			start:             time.Now(),
			wait:              &sync.WaitGroup{},
			ProjectsTag:       make(chan gitlab.ProjectsTag, 10000),
			ProjectsFileChan:  make(chan gitlab.ProjectsFileList, 10000),
			ProjectsPathChan:  make(chan gitlab.ProjectsFileList, 10000),
		}
		task.TaskPool, err = ants.NewPool(runtime.NumCPU() * 100)
		if err != nil {
			return err
		}
		defer task.TaskPool.Release()
		//删除项目数据
		err = mongo.DelCodeAll(c, project.Name, tag)
		if err != nil {
			glogs.ErrorF(c, err.Error())
		}
		task.ProjectsTag <- gitlab.ProjectsTag{
			Id:    project.Id,
			Tag:   tag,
			EnvID: envID,
			Code:  code,
		}
		task.TaskPool.Submit(func() {
			task.wait.Add(1)
			startProjectFileListData(c, task)
		})

		task.TaskPool.Submit(func() {
			task.wait.Add(1)
			startProjectFileData(c, task)
		})
		task.TaskPool.Submit(func() {
			task.wait.Add(1)
			startEsAddBatch(c, task)
		})
		task.wait.Wait()
		stop := time.Since(task.start)
		fmt.Printf("task over, 耗时：%v", stop)
	} else {
		err = errors.New("更新失败-项目不存在")
	}
	return
}

//startProject 获取项目信息
func startProject(c *gin.Context, task *TaskEsPool) {
	data, _ := GetProjectsAll(c, 1, 50, []gitlab.Projects{})
	var projects []interface{}
	//获取项目信息
	if len(data) > 0 {
		for _, v := range data {
			task.ProjectChan <- v
			projects = append(projects, v)
		}
	}
	if len(projects) > 0 {
		//删除原数据
		mongo.CodeInitCodeName(c)
		//写入数据库
		mongo.AddALL(c, viper.GetString("git.gitlab_project_name"), projects)
	}
}

func startProjectTag(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectChan)
	defer task.wait.Done()
	var wait sync.WaitGroup
	defer wait.Done()
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectChan {
		wait.Add(1)
		ProjectTag(c, x, task)
		wait.Done()
	}
	wait.Wait()
}

//项目文件列表
func startProjectFileListData(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectsTag)
	defer task.wait.Done()
	var wait sync.WaitGroup
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectsTag {
		wait.Add(1)
		fmt.Println("项目消费" + x.Code + "-" + x.Tag)
		ProjectList(c, x, task)
		wait.Done()
	}
	wait.Wait()
}

//项目文件解读
func startProjectFileData(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectsFileChan)
	defer task.wait.Done()
	var wait sync.WaitGroup
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectsFileChan {
		wait.Add(1)
		ProjectTree(c, x.Id, x.Tag, x.Path, x.ProjectName, 1, x.EnvID, task)
		wait.Done()
	}
	wait.Wait()
}

func startEsAddBatch(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectsPathChan)
	defer task.wait.Done()
	var wait sync.WaitGroup
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectsPathChan {
		wait.Add(1)
		mongo.AddOne(c, viper.GetString("git.gitlab_code_name"), x)
		wait.Done()
	}
	wait.Wait()
}

func GetProjectsAll(c *gin.Context, page, perPage int, req []gitlab.Projects) (data []gitlab.Projects, err error) {
	resp, _ := gitlab.GetProjectsList(c, page, perPage)
	if len(resp) > 0 {
		for _, v := range resp {
			v.Tag = v.DefaultBranch
			req = append(req, v)
		}
		fmt.Println(page)
		data, _ = GetProjectsAll(c, page+1, perPage, req)
	} else {
		data = req
	}
	return
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
