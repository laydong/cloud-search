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
	}
	return nil
}

type Config struct {
	AppConf AppConf
	MGConf  MGConf
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
