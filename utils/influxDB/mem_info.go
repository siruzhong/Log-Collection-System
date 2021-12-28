package main

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

// MemInfo 内存相关数据结构体
type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// getMemInfo 获取内存相关数据
func getMemInfo() *MemInfo {
	memData, _ := mem.VirtualMemory()
	memInfo := &MemInfo{Total: memData.Total, Available: memData.Available, Used: memData.Used, UsedPercent: memData.UsedPercent}
	return memInfo
}

// writesMemPoints 写入mem的Points数据到InfluxDB中
func writesMemPoints(memInfo *MemInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", // 精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{"mem": "mem_info"}
	fields := map[string]interface{}{
		"total":        int64(memInfo.Total),
		"available":    int64(memInfo.Available),
		"used":         int64(memInfo.Used),
		"used_percent": memInfo.UsedPercent,
	}
	pt, err := client.NewPoint("mem_info", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert mem info success")
}
