package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"log_transfer/es"
)

// Init 初始化全局kafka连接对象
func Init(address []string, topic string) (err error) {
	// 创建消费者
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		logrus.Errorf("fail to start consumer,err=%v\n", err)
		return
	}
	// 获取指定topic下的所有分区
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		logrus.Errorf("fail to get list of partition,err=%v\n", err)
		return
	}
	// 遍历所有分区
	for partition := range partitionList {
		// 为每个分区创建一个对应的分区消费者
		partitionConsumer, err1 := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err1 != nil {
			logrus.Errorf("fail to start consumer for partition,err=%v\n", err)
			return
		}
		// 异步从每个分区消费消息
		go func(sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				// 为了将从kafka中取出消息和将消息写入es中异步化，首先将取出对消息放入一个logMesChan通道中
				var message map[string]interface{}
				err = json.Unmarshal(msg.Value, &message)
				if err != nil {
					logrus.Errorf("json.Unmarshal failed,err=%v\n", err)
				}
				fmt.Println(message)
				logrus.Infof("Partition:%d Offset:%d Key:%v Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
				es.PutLogData(message)
			}
		}(partitionConsumer)
	}
	return
}
