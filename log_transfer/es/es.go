package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type ESClient struct {
	client      *elastic.Client
	index       string
	logDataChan chan interface{}
}

// esClient 全局es连接对象
var esClient ESClient

// Init 初始化es
func Init(address string, maxChanSize int, index string) (err error) {
	client, err := elastic.NewClient(elastic.SetURL("http://" + address))
	if err != nil {
		fmt.Println("connect to ElasticSearch failed,err=", err)
	}
	fmt.Println("connect to ElasticSearch success!")
	// 初始化全局esClient
	esClient.client = client
	esClient.index = index
	esClient.logDataChan = make(chan interface{}, maxChanSize)
	// 从通道中取出数据,写入到es中
	go sendMsgToES()
	return
}

// sendMsgToES 写入消息到es中
func sendMsgToES() {
	for message := range esClient.logDataChan {
		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(message).
			Do(context.Background())
		if err != nil {
			fmt.Println("put data into index failed,err=", err)
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

// PutLogData 从包外接收msg,发送到msg中
func PutLogData(msg interface{}) {
	esClient.logDataChan <- msg
}
