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
	"strconv"
	"strings"
	"sync"
	"time"
)

type TaskEsPool struct {
	projectNumber     *atomic.Int64
	projectFileNumber *atomic.Int64
	projectPathNumber *atomic.Int64
	//projectFileAllRowsNumber *atomic.Int64
	start                time.Time
	wait                 *sync.WaitGroup
	TaskPool             *ants.Pool
	ProjectChan          chan gitlab.Projects
	ProjectsTag          chan gitlab.ProjectsTag
	ProjectsFileListChan chan gitlab.ProjectsFileList
	ProjectsFileChan     chan gitlab.ProjectsFileList
}

var (
	//ProjectChan          = make(chan gitlab.Projects, 1000)
	ProjectsTag          = make(chan gitlab.ProjectsTag, 1000)
	ProjectsFileListChan = make(chan gitlab.ProjectsFileList, 10000)
	ProjectsFileChan     = make(chan gitlab.ProjectsFileList, 10000)
)

// 批量更新项目数据
func UpProjects(c *gin.Context) (err error) {
	var task = &TaskEsPool{
		projectPathNumber:    atomic.NewInt64(0),
		projectFileNumber:    atomic.NewInt64(0),
		projectNumber:        atomic.NewInt64(0),
		start:                time.Now(),
		wait:                 &sync.WaitGroup{},
		ProjectChan:          make(chan gitlab.Projects, 1000),
		ProjectsTag:          make(chan gitlab.ProjectsTag, 1000),
		ProjectsFileListChan: make(chan gitlab.ProjectsFileList, 10000),
		ProjectsFileChan:     make(chan gitlab.ProjectsFileList, 10000),
	}
	task.TaskPool, err = ants.NewPool(runtime.NumCPU() * 100)
	if err != nil {
		return err
	}
	defer task.TaskPool.Release()
	startProject(c, task)
	task.wait.Add(3)
	//go startProjectTag(c, task)          //拆分项目
	//go startProjectFileListData(c, task) //查询目录 获取文件列表
	////go startProjectFileData(c, task)     //获取文件详情
	//go startEsAddBatch(c, task) //项目文件入库
	err = task.TaskPool.Submit(func() {
		startProjectTag(c, task)
	})
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	err = task.TaskPool.Submit(func() {
		startProjectFileListData(c, task)
	})
	//err = task.TaskPool.Submit(func() {
	//	startProjectFileData(c, task)
	//})
	err = task.TaskPool.Submit(func() {
		startEsAddBatch(c, task)
	})
	//go startEsAddBatch(c, task)          //项目文件入库
	task.wait.Wait()
	stop := time.Since(task.start)
	fmt.Printf("task over, 耗时：%v", stop)

	return nil
}

//ProjectCodeUp 更新指定项目
func ProjectCodeUp(c *gin.Context, code string) (err error) {

	project, _ := gitlab.QueryByName(c, ProjectReplace(c, code))
	if project.Id > 0 {
		//var task = &TaskEsPool{
		//	projectNumber:        atomic.NewInt64(0),
		//	start:                time.Now(),
		//	wait:                 &sync.WaitGroup{},
		//	ProjectChan:          make(chan gitlab.Projects, 1000),
		//	ProjectsTag:          make(chan gitlab.ProjectsTag, 1000),
		//	ProjectsFileListChan: make(chan gitlab.ProjectsFileList, 10000),
		//	ProjectsFileChan:     make(chan gitlab.ProjectsFileList, 10000),
		//}
		task, err := ants.NewPool(runtime.NumCPU() * 100)
		if err != nil {
			return err
		}
		//defer task.Release()
		//删除项目数据
		mongo.CodeInitCodeName(c)
		//err = mongo.DelCodeAll(c, project.Name)
		//if err != nil {
		//	glogs.ErrorF(c, err.Error())
		//}
		project.EnvID = 3
		project.Tag = "dev"
		fmt.Println(project)
		//go ProjectList(c, gitlab.ProjectsTag{
		//	Id:    project.Id,
		//	Code:  project.Name,
		//	EnvID: project.EnvID,
		//	Tag:   project.Tag,
		//})

		defer task.Release()

		//runTimes := runtime.NumCPU() * 100

		// Use the common pool.
		//go func() {
		//	fmt.Println(2222)
		//	var wg sync.WaitGroup
		//	for x := range ProjectsFileChan {
		//		fmt.Println(2222)
		//		wg.Add(1)
		//		ProjectTree(c, x.Id, x.Tag, x.Path, x.ProjectsName, 1, x.EnvID)
		//		wg.Done()
		//	}
		//	wg.Wait()
		//}()
		//_ = task.Submit(syncCalculateSum)
		//for i := 0; i < runTimes; i++ {
		//	wg.Add(1)

		//}

		//for x := range ProjectsTag {
		//	ProjectList(c, x)
		//}
		syncCalculate := func() {
			var wgt sync.WaitGroup
			for x := range ProjectsFileListChan {
				wgt.Add(1)
				mongo.AddOne(c, viper.GetString("git.gitlab_code_name"), x)
				wgt.Done()
			}
			wgt.Wait()
		}
		_ = task.Submit(syncCalculate)
		//var wgt sync.WaitGroup
		//syncCalculate := func() {
		//	for x := range ProjectsFileListChan {
		//		mongo.AddOne(c, viper.GetString("git.gitlab_code_name"), x)
		//	}
		//	wgt.Done()
		//}
		//for i := 0; i < runTimes; i++ {
		//	wgt.Add(1)
		//	_ = ants.Submit(syncCalculate)
		//}

		//for x := range ProjectsFileListChan {
		//	mongo.AddOne(c, viper.GetString("git.gitlab_code_name"), x)
		//}
		//var wg sync.WaitGroup
		//wg.Add(3)
		//err = task.Submit(func() {
		//	startProjectFileListData(c, task)
		//	wg.Done()
		//})
		//err = task.Submit(func() {
		//	startProjectFileData(c, task)
		//	wg.Done()
		//})
		//err = task.Submit(func() {
		//	startEsAddBatch(c, task)
		//	wg.Done()
		//})
		//wg.Wait()
		//ProjectTag(c, project, task)
		////go startProjectTag(c, task)          //拆分项目
		//go startProjectFileListData(c, task) //查询目录 获取文件列表
		////go startProjectFileData(c, task)     //获取文件详情
		//go startEsAddBatch(c, task) //项目文件入库

		//task.wait.Wait()
		//stop := time.Since(task.start)
		//fmt.Printf("task over, 耗时：%v", stop)
		fmt.Println(11111)

	} else {
		err = errors.New("更新失败-项目不存在")
	}
	return
}

//startProject 获取项目信息
//func startProject(c *gin.Context, task *TaskEsPool) {
func startProject(c *gin.Context, task *TaskEsPool) {
	data, _ := GetProjectsAll(c, 1, 50, []gitlab.Projects{})
	var projects []interface{}
	//获取项目信息
	if len(data) > 0 {
		for _, v := range data {
			projects = append(projects, v)
			task.ProjectChan <- v
			fmt.Println("项目列表：" + v.Name + "====" + v.DefaultBranch)
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
		err := task.TaskPool.Submit(func() {
			ProjectTag(c, x)
		})
		if err != nil {
			glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
		}
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
		err := task.TaskPool.Submit(func() {
			fmt.Println("处理项目" + x.Code + "---" + x.Tag)
			ProjectList(c, x, task)
			wait.Done()
		})
		if err != nil {
			glogs.ErrorF(c, "项目执行错误: "+strconv.Itoa(x.Id), err.Error())
		}
	}
	wait.Wait()
}

//项目文件解读
//func startProjectFileData(c *gin.Context, task *TaskEsPool) {
//	defer close(ProjectsFileChan)
//	//defer task.wait.Done()
//	var wait sync.WaitGroup
//	//defer wait.Done()
//	// 从ch中接收值并赋值给变量x
//	for x := range ProjectsFileChan {
//		wait.Add(1)
//		err := task.TaskPool.Submit(func() {
//			ProjectTree(c, x.Id, x.Tag, x.Path, x.ProjectsName, 1, x.EnvID, task)
//			fmt.Println("项目文件夹处理：", x)
//			wait.Done()
//		})
//		if err != nil {
//			glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
//		}
//
//	}
//	wait.Wait()
//}

func startEsAddBatch(c *gin.Context, task *TaskEsPool) {
	defer close(task.ProjectsFileListChan)
	defer task.wait.Done()
	var wait sync.WaitGroup
	// 从ch中接收值并赋值给变量x
	for x := range task.ProjectsFileListChan {
		wait.Add(1)
		err := task.TaskPool.Submit(func() {
			fmt.Println("入库数据" + x.Name + "-" + x.Tag)
			mongo.AddOne(c, viper.GetString("git.gitlab_code_name"), x)
			wait.Done()
		})
		if err != nil {
			glogs.ErrorF(c, "项目执行错误: "+x.Name, err.Error())
		}
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
