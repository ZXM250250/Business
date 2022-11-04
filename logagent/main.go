package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"io"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
	"os"
)

func main() {
	InitLog()
	var configObj = new(Config)
	// 0. 读配置文件 `go-ini`
	// 0. 读配置文件 `go-ini`
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config failed,err:%v", err)
		return
	}
	fmt.Printf("%#v\n", configObj)
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("init kafka success!")
	// 初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})

	if err != nil {
		logrus.Errorf("init etcd failed, err:%v", err)
		return
	}
	logrus.Info("init etcd success!")
	// 从etcd中拉取要收集日志的配置项
	allConf, err := etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err:%v", err)
		return
	}
	fmt.Println(allConf)
	// 派一个小弟去监控etcd中 configObj.EtcdConfig.CollectKey 对应值的变化
	go etcd.WatchConf(configObj.EtcdConfig.CollectKey)
	// 2. 根据配置中的日志路径初始化tail
	err = tailfile.Init(allConf) // 把从etcd中获取的配置项传到Init
	if err != nil {
		logrus.Errorf("init tailfile failed, err:%v", err)
		return
	}
	logrus.Info("init tailfile success!")
	run()
}

func InitLog() {
	// 设置日志格式为json格式
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	file, err := os.OpenFile("E:\\Projects\\Golang\\business\\logs\\s1.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		logrus.SetOutput(fileAndStdoutWriter)
	} else {
		logrus.Info("failed to logrus. to file.")
	}
	//设置最低loglevel
	logrus.SetLevel(logrus.InfoLevel)

}
func run() {
	select {}
}

// 整个logaent的配置结构体
type Config struct {
	conf.KafkaConfig   `ini:"kafka"`
	conf.CollectConfig `ini:"collect"`
	conf.EtcdConfig    `ini:"etcd"`
}
