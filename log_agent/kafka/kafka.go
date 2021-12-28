package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer          // 全局kafka连接对象
	msgChan chan *sarama.ProducerMessage // 消息管道（这里数据类型选择指针类型，不采用string类型占用空间）
	// 读取到到日志信息写入该管道，然后client从中取写入kafka。实现读取文件和写入kafka是异步操作
)

// Init 初始化全局kafka连接对象、msgChan、启动一个后台goroutine负责将msgChan中读数据写入kafka
func Init(address []string, chanSize int) (err error) {
	// 1.初始化kafka配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 设置发送完数据需要leader和follower都确定
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 设置发送消息的分区为随机分区
	config.Producer.Return.Successes = true                   // 发送成功回复确认

	// 2.连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Errorf("kafka: producer err=%v", err)
		return
	}

	// 3.初始化消息管道
	msgChan = make(chan *sarama.ProducerMessage, chanSize)

	// 4.起一个后台goroutine从msgChan中读数据写入kafka
	go SendMsgFromChanToKafka()
	return
}

// SendMsgFromChanToKafka 读取msgChan中的消息写入kafka
func SendMsgFromChanToKafka() {
	for {
		select {
		case msg := <-msgChan:
			partition, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Errorf("send message err=%v", err)
				return
			}
			logrus.Infof("send msg to kafka success! patition:%d offset:%v", partition, offset)
		}
	}
}

// SendMesToChan 发送消息到管道中
func SendMesToChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
