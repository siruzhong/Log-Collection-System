package main

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
)

// cli 全局InfluxDB连接对象
var cli client.Client

// connInflux 连接InfluxDB
func connInflux() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func main() {
	// 获取InfluxDB连接
	cli = connInflux()
	// 每隔一秒插入一条网卡IO速率数据
	for {
		percent := getCpuPercent()       // 获取cpu使用率
		writesCpuPoints(percent[0])      // 将cpu使用率写入到influxdb中
		memInfo := getMemInfo()          // 获取内存相关数据
		writesMemPoints(memInfo)         // 将内存数据写入到influxdb中
		diskInfo := getDiskInfo()        // 获取磁盘相关数据
		writesDiskPoints(diskInfo)       // 将磁盘数据写入到influxdb中
		netRateInfo := getNetRateInfo()  // 获取网络相关数据
		writesNetRatePoints(netRateInfo) // 将网卡数据写入到influxdb中
		time.Sleep(time.Second)
	}
}
