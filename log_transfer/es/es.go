package es

// 将日志数据写入Elasticsearch
import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type EsClient struct {
	client      *elastic.Client
	index       string
	logDataChan chan interface{}
}

var (
	esClient *EsClient
)

func Init(addr, index string, goroutineNum, maxSize int) (err error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		fmt.Printf(err.Error())
		// Handle error
		panic(err)
	}
	fmt.Printf("%#v\n", client)
	esClient = &EsClient{
		client:      client,
		index:       index,
		logDataChan: make(chan interface{}, maxSize),
	}
	fmt.Println("connect to es success")
	// 从通道中取出数据,写入到kafka中去
	for i := 0; i < goroutineNum; i++ {
		go sendToES()
	}
	return
}
func sendToES() {
	for m1 := range esClient.logDataChan {

		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(m1).
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

// 通过一个首字母大写的函数从包外接收msg,发送到chan中
func PutLogData(msg interface{}) {
	esClient.logDataChan <- msg
}
