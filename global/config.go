package global

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port int
}
type Configs struct {
	ServerCfg ServerConfig
}

var (
	Cfg *Configs
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
	Cfg = &Configs{}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(Cfg)
	if err != nil {
		panic(err)
	}
	logrus.Info("配置文件:", Cfg)
}
