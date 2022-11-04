package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log_transfer/es"
)

func Init(add []string, topic string) (err error) {

	consumer, err := sarama.NewConsumer(add, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	// 拿到指定topic下面的所有分区列表
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList {
		var pc sarama.PartitionConsumer
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		fmt.Println("start to consume...")
		go func(partitionConsumer sarama.PartitionConsumer) {
			fmt.Println("in sarama.PartitionConsumer")
			for msg := range pc.Messages() {
				//logDataChan<-msg // 为了将同步流程异步化,所以将取出的日志数据先放到channel中
				fmt.Println(msg.Topic, string(msg.Value))
				var m1 map[string]interface{}
				err = json.Unmarshal(msg.Value, &m1)
				if err != nil {
					fmt.Printf("unmarshal msg failed, err:%v\n", err)
					continue
				}
				es.PutLogData(m1)
			}
		}(pc)

	}

	return
}
