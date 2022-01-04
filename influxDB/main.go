package main

import (
	"influxDB/sysInfo"
	"time"
)

func main() {
	// 每隔一秒插入一条网卡IO速率数据
	for {
		percent := sysInfo.GetCpuPercent()       // 获取cpu使用率
		sysInfo.WritesCpuPoints(percent[0])      // 将cpu使用率写入到influxdb中
		memInfo := sysInfo.GetMemInfo()          // 获取内存相关数据
		sysInfo.WritesMemPoints(memInfo)         // 将内存数据写入到influxdb中
		diskInfo := sysInfo.GetDiskInfo()        // 获取磁盘相关数据
		sysInfo.WritesDiskPoints(diskInfo)       // 将磁盘数据写入到influxdb中
		netRateInfo := sysInfo.GetNetRateInfo()  // 获取网络相关数据
		sysInfo.WritesNetRatePoints(netRateInfo) // 将网卡数据写入到influxdb中
		time.Sleep(time.Second)
	}
}
