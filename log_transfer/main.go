package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"log_transfer/conf"
	"log_transfer/es"
	"log_transfer/kafka"
)

func main() {
	// 1.使用go-ini包读取config.ini获取初始化配置
	var config = new(conf.Config)
	err := ini.MapTo(config, "./conf/config.ini") // 将配置文件转换为Config结构体
	if err != nil {
		logrus.Infof("load config failed,err=%v", err)
		return
	}
	logrus.Infof("load config success")

	// 2.初始化kafka连接,从kafka中取出数据
	err = kafka.Init([]string{config.KafkaConfig.Address}, config.KafkaConfig.Topic)
	if err != nil {
		logrus.Errorf("init kafka err=%v", err)
		return
	}
	logrus.Info("init kafka success")

	// 3.初始化es
	err = es.Init(config.ESConfig.Address, config.ESConfig.MaxChanSize, config.ESConfig.Index)
	if err != nil {
		logrus.Errorf("init es err=%v", err)
		return
	}
	logrus.Info("init es success")

	select {}
}
