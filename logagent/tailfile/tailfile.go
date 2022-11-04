package tailfile

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
	"strings"
	"time"
)

import "github.com/hpcloud/tail"

//抽象起来的每一个日志任务
type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancelFunc,
	}

}

func (t *tailTask) init() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tObj, err = tail.TailFile(t.path, config)
	return

}

func (t *tailTask) run() {
	//读取日志发送到kafka
	logrus.Infof("collect for path:%s is running...", t.path)

	for {
		select {
		case <-t.ctx.Done():
			logrus.Infof("path:%s is stopping...", t.path)
			return
		case line, ok := <-t.tObj.Lines:
			//logrus.Infof("日志发生了变动%s", line.Text)
			if !ok {
				logrus.Warn("tail file close reopen, path:%s\n", t.path)
				time.Sleep(time.Second) // 读取出错等一秒
				continue
			}
			if len(strings.Trim(line.Text, "/r")) == 0 {
				logrus.Infof("出现空行")
				continue
			}
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			//把日志消息丢到通道当中去
			kafka.ToMsgChan(msg)
		}

	}
}
