package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type TomlConfig struct {
	AppName    string
	Sql        MySQLConfig
	StaticPath PathConfig
}

// MySQLConfig pool 日志
type MySQLConfig struct {
	Host        string
	Name        string
	Password    string
	Port        int
	TablePrefix string
	User        string
	TimeOut     string
}

// LogConfig 日志保存地址
type LogConfig struct {
	Path  string
	Level string
}

// PathConfig 相关地址信息，例如静态文件地址
type PathConfig struct {
	FilePath string
}

var c TomlConfig

func init() {
	fmt.Println(os.Getwd())
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.Unmarshal(&c)
	fmt.Println(viper.Get("mysql.username"))
}

func GetConfig() TomlConfig {
	fmt.Println(os.Getwd())
	return c
}
