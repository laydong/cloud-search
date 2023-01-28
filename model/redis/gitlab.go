package redis

import (
	"cloud-search/global"
	"github.com/gin-gonic/gin"
	"time"
)

func SAdd(c *gin.Context, key, val string) int64 {
	return global.Rdb.SAdd(c, key, val).Val()
}

func Sismember(c *gin.Context, key, val string) bool {
	return global.Rdb.SIsMember(c, key, val).Val()
}

func Expire(c *gin.Context, key string, exp int64) bool {
	return global.Rdb.Expire(c, key, time.Duration(exp)*time.Second).Val()
}

func Del(c *gin.Context, key string) int64 {
	return global.Rdb.Del(c, key).Val()
}
