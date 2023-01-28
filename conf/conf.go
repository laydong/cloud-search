package conf

import (
	"github.com/spf13/viper"
)

var ConfInfo *Config

func InitDoAfter() error {
	viper.SetConfigFile("./conf/app.toml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.WatchConfig()

	ConfInfo = &Config{
		AppConf: AppConf{
			Name:       viper.GetString("app.name"),
			Mode:       viper.GetString("app.mode"),
			Level:      viper.GetString("app.level"),
			HttpListen: viper.GetString("app.http_listen"),
			Url:        viper.GetString("app.url"),
			Pprof:      viper.GetBool("app.pprof"),
			Params:     viper.GetBool("app.params"),
			Logger:     viper.GetString("app.logger"),
			Version:    viper.GetString("app.version"),
		},
		MGConf: MGConf{
			Dsn:             viper.GetString("mongodb.dsn"),
			ConnTimeOut:     viper.GetInt("mongodb.conn_time_out"),
			ConnMaxPoolSize: viper.GetInt("mongodb.conn_max_pool_size"),
		},
		DBConf: DBConf{
			Dsn:             viper.GetString("mysql.dsn"),
			DbName:          viper.GetString("mysql.db_name"),
			MaxIdleConn:     viper.GetInt("mysql.max_idle_conn"),
			MaxOpenConn:     viper.GetInt("mysql.max_open_conn"),
			ConnMaxLifetime: viper.GetInt("mysql.conn_max_lifetime"),
		},
		RDConf: RDConf{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
	}
	return nil
}

type Config struct {
	AppConf AppConf
	DBConf  DBConf
	MGConf  MGConf
	RDConf  RDConf
}

type AppConf struct {
	Name       string `json:"name"`
	Mode       string `json:"mode"`
	Level      string `json:"level"`
	HttpListen string `json:"http_listen"`
	Url        string `json:"url"`
	Pprof      bool   `json:"pprof"`
	Params     bool   `json:"params"`
	Logger     string `json:"logger"`
	Version    string `json:"version"`
}

type MGConf struct {
	ConnTimeOut     int    `json:"conn_time_out"`
	ConnMaxPoolSize int    `json:"conn_max_pool_size"`
	Dsn             string `json:"dsn"`
}

type DBConf struct {
	MaxIdleConn     int    `json:"max_idle_conn"`
	MaxOpenConn     int    `json:"max_open_conn"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
	Dsn             string `json:"dsn"`
	DbName          string `json:"db_name"`
}

type RDConf struct {
	DB       int    `json:"db"`       // redis的哪个数据库
	Addr     string `json:"addr"`     // 服务器地址:端口
	Password string `json:"password"` // 密码
}
