package tools

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

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
