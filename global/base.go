package global

import (
	"cloud-search/conf"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v6"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Rdb *redis.Client
var Mdb *mongo.Client
var EdbClient *elastic.Client

func GetDB(c *gin.Context, dbNmae ...string) *gorm.DB {
	key := conf.ConfInfo.DBConf.DbName
	if len(dbNmae) > 0 {
		key = dbNmae[0]
	}
	if key == "" {
		key = "grom_cxt"
	}
	return DB.Set(key, c).WithContext(c)
}
