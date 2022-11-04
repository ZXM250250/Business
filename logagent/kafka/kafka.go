package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// kafka相关操作

var (
	producer sarama.SyncProducer
	msgChan  chan *sarama.ProducerMessage
)

//初始化全局的kafka Client

func Init(address []string, chanSize int64) (err error) {
	logrus.Infof("初始化Kafka%s", address)
	//生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // ACK
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 分区
	config.Producer.Return.Successes = true                   // 确认
	//2.连接kafka
	producer, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka:producer closed, err:", err)
		return
	}
	logrus.Infof("连接KafKa成功!")
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg()
	return
}

func sendMsg() {
	for {
		select {
		case msg := <-msgChan:

			pid, offset, err := producer.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed, err:", err)
				logrus.Infof("send msg to kafka success. pid:%v offset:%v", pid, offset)
				return
			}

		}

	}

}

// 定义一个函数向外暴露msgChan
func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
