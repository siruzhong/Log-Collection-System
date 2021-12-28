package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

func main() {
	// 创建消费者
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer,err=%v\n", err)
		return
	}
	// 获取指定topic下的所有分区
	partitionList, err := consumer.Partitions("net")
	if err != nil {
		fmt.Printf("fail to get list of partition,err=%v\n", err)
		return
	}
	var waitGroup sync.WaitGroup
	// 遍历所有分区
	for partition := range partitionList {
		// 为每个分区创建一个对应的分区消费者
		partitionConsumer, err := consumer.ConsumePartition("web", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("fail to start consumer for partition,err=%v\n", err)
			return
		}
		defer partitionConsumer.AsyncClose()
		waitGroup.Add(1)
		// 异步从每个分区消费消息
		go func(sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(partitionConsumer)
	}
	waitGroup.Wait()
}
