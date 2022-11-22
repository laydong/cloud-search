package mongo

import (
	"codesearch/global/glogs"
	"codesearch/global/gstore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type Project struct {
	//ID          string `bson:"_id"`
	ProjectName string `bson:"project_name"` //项目信息
	Path        string `bson:"path"`         //路径
	Version     string `bson:"version"`      //版本
	Content     string `bson:"content"`      // 文件内容
	CreatedAt   string `bson:"created_at"`   //创建时间
	UpdateTime  string `bson:"update_time"`  //更新时间
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
}

type RespProjectsList struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	Path         string     `json:"path"`
	Mode         string     `json:"mode"`
	Content      string     `json:"content"`
	Tag          string     `json:"tag"`
	ProjectsName string     `json:"projects_name"`
	LineData     []LineData `json:"line_data"`
}

type LineData struct {
	Line    int    `json:"line"`
	Content string `json:"content"`
}

// AddOne 单条插入
func AddOne(c *gin.Context, dbname string, data interface{}) (err error, resp interface{}) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(dbname)
	resp, err = collection.InsertOne(c, data)
	glogs.Info("入库记录", err)
	return
}

//AddALL 批量插入
func AddALL(c *gin.Context, name string, data []interface{}) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(name)
	_, err := collection.InsertMany(c, data)
	if err != nil {
		return
	}
}

//FindOne 单条查询
func FindOne(c *gin.Context) (resp Project, err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_project_name"))
	err = collection.FindOne(c, bson.M{"project_name": "user-web"}).Decode(&resp)
	return
}

//Find 多条查询
func Find(c *gin.Context, page, limit int) (data interface{}, err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_code_name"))
	count, err := collection.CountDocuments(c, bson.M{"version": "1.0.0"})
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(limit * (page - 1)))
	findOptions.SetSort(map[string]int{"updated_at": -1})
	find, err := collection.Find(c, bson.M{"version": "1.0.0"}, findOptions)
	if err != nil {
		return
	}
	var res []Project
	err = find.All(c, &res)
	if err != nil {
		return nil, err
	}
	return CutPageData(count, page, limit, len(res), res), err
}

//UpOne 单条更新
func UpOne(c *gin.Context) (err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_code_name"))
	_, err = collection.UpdateOne(c, bson.M{"project_name": "user-web"}, bson.M{"$set": bson.M{"version": "1.0.1"}})
	return
}

//UpAll 批量更新
func UpAll(c *gin.Context) (err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_code_name"))
	_, err = collection.UpdateMany(c, bson.M{"version": "1.0.0"}, bson.M{"$set": bson.M{"version": "1.0.1"}})
	return
}

//DelOne 删除单条
func DelOne(c *gin.Context) (err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_code_name"))
	_, err = collection.DeleteOne(c, bson.M{"project_name": "user-web"})
	return
}

//DelAll 删除多条
func DelAll(c *gin.Context) (err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_project_name"))
	_, err = collection.DeleteMany(c, bson.M{"version": "1.0.1"})
	return
}

//DelCodeAll 删除指定项目
func DelCodeAll(c *gin.Context, project string) (err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_project_name"))
	_, err = collection.DeleteMany(c, bson.M{"projectsname": project})
	return
}

//CodeInitCodeName 处理数据初始化
func CodeInitCodeName(c *gin.Context) (err error) {
	_, err = gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_code_name")).DeleteMany(c, bson.D{})
	_, err = gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_project_name")).DeleteMany(c, bson.D{})
	return
}

//CodeFind 多条查询
func CodeFind(c *gin.Context, envID uint, str string, page, limit int) (data interface{}, err error) {
	collection := gstore.Mdb.Database(viper.GetString("git.gitlab_depod_name")).Collection(viper.GetString("git.gitlab_code_name"))
	filter := bson.M{"envid": envID, "content": bson.M{"$regex": primitive.Regex{Pattern: str}}}
	count, err := collection.CountDocuments(c, filter)
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(limit * (page - 1)))
	findOptions.SetSort(map[string]int{"id": -1})
	find, err := collection.Find(c, filter, findOptions)
	if err != nil {
		return
	}
	var res []RespProjectsList
	err = find.All(c, &res)
	//处理关键字行数处理
	if len(res) > 0 {
		for k, v := range res {
			split := strings.Split(v.Content, "\n")
			if len(split) > 0 {
				for k1, v1 := range split {
					if v1 != "" && strings.Contains(strings.ToLower(v1), strings.ToLower(str)) {
						v.LineData = append(v.LineData, LineData{
							Line:    k1 + 1,
							Content: v1,
						})
					}
				}
			}
			res[k] = v
		}
	}
	return CutPageData(count, page, limit, len(res), res), err
}

func GetCodeComment(c *gin.Context, str string) (resp interface{}) {
	collection := gstore.Mdb.Database("app_gitlab").Collection("code_list")
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
