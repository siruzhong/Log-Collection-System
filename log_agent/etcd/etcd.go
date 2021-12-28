package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"logs-collection-system/conf"
	"logs-collection-system/tail"
	"time"
)

// client 全局etcd连接
var client *clientv3.Client

// Init 初始化全局etcd连接
func Init(address []string) (err error) {
	// 连接etcd
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	//defer client.Close()
	return
}

// GetConf 从etcd中拉取配置项
func GetConf(key string) (logConfigList []conf.LogConfig, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err:%v\n", err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by key:%s\n", key)
		return
	}
	err = json.Unmarshal(resp.Kvs[0].Value, &logConfigList)
	if err != nil {
		logrus.Errorf("json unmarshal failed,err:%v", err)
		return
	}
	return
}

// WatchConf 循环监听etcd中日志配置项变化
func WatchConf(key string) {
	for {
		watchChan := client.Watch(context.Background(), key) // <-chan WatchResponse
		for wresp := range watchChan {
			logrus.Info("get new conf from etcd!")
			for _, ev := range wresp.Events {
				fmt.Printf("T ype: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				var newConf []conf.LogConfig // 新配置
				// 如果是删除操作,直接发送空配置
				if ev.Type == clientv3.EventTypeDelete {
					logrus.Warning("FBI warning:etcd delete the key!!!")
					tail.SendNewConf(newConf)
					continue
				}
				err := json.Unmarshal(ev.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshal new conf err:%v", err)
					continue
				}
				// 告知tail配置变更
				tail.SendNewConf(newConf)
			}
		}
	}
}
