package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

func test1() {
	test2()
}

func test2() {
	test3()
}

func test3() {
	logrus.WithField("stacktrace", string(debug.Stack())).Info("asfas")

	//	debug.PrintStack()
}

func main() {
	//	test1()
	fmt.Println(time.Now().In(time.Now().Location()).Format("20060102"))
}
