package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	// 连接etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Put(ctx, "collect_log_192.168.184.158_conf", `[{"path":"/Users/zhongsiru/GoProjects/src/logs-collection-system/logs/web.log","topic":"web"},{"path":"/Users/zhongsiru/GoProjects/src/logs-collection-system/logs/net.log","topic":"net"},{"path":"/Users/zhongsiru/GoProjects/src/logs-collection-system/logs/video.log","topic":"videos"},{"path":"/Users/zhongsiru/GoProjects/src/logs-collection-system/logs/photo.log","topic":"photo"},{"path":"/Users/zhongsiru/GoProjects/src/logs-collection-system/logs/tea.log","topic":"tea"}]`)
	cli.Put(ctx, "collect_log_192.168.51.158_conf", `[{"path":"/Users/zhongsiru/GoProjects/src/log-collection-system/log_agent/logs/web.log","topic":"web"}]`)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "collect_log_192.168.58.158_conf")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
