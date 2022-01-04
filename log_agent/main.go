package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"logs-collection-system/conf"
	"logs-collection-system/etcd"
	"logs-collection-system/kafka"
	"logs-collection-system/tail"
)

func main() {
	// 1. 使用go-ini包读取config.ini获取初始化配置(etcd、kafka)
	var config = new(conf.Config)
	err := ini.MapTo(config, "./conf/config.ini") // 将配置文件转换为Config结构体
	if err != nil {
		logrus.Infof("load init config failed,err=%v", err)
		return
	}
	logrus.Infof("load init config success")

	// 2. 连接kafka、初始化msgChan、启动一个后台goroutine负责读取msgChan的数据写入kafka
	err = kafka.Init([]string{config.KafkaConfig.Address}, config.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka err=%v", err)
		return
	}
	logrus.Info("init kafka success")

	// 3. 初始化etcd
	err = etcd.Init([]string{config.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd err=%v", err)
		return
	}
	logrus.Info("init etcd success")

	// 4. 从etcd中获取日志配置项(etcd中存放json日志配置项)
	logConfigList, err := etcd.GetConf(config.EtcdConfig.Key)
	if err != nil {
		logrus.Errorf("get conf from etde err=%v", err)
		return
	}

	// 5. 启动后台goroutine监听etcd中日志配置项是否变化
	go etcd.WatchConf(config.EtcdConfig.Key)

	// 6. 初始化全局tailTaskManager,为每个配置项创建一个日志收集任务,读取其中的内容写入msgChan,并循环监听获取新配置
	err = tail.Init(logConfigList)
	if err != nil {
		logrus.Errorf("init tail err=%v", err)
		return
	}
	logrus.Info("init tail success")

	select {}
}
