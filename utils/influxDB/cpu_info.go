package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"log"
	"time"
)

// getCpuPercent 获取cpu占用率
func getCpuPercent() []float64 {
	percent, _ := cpu.Percent(time.Second, false) // 获取cpu使用率
	fmt.Printf("cpu percent:%v\n", percent[0])
	return percent
}

// getCpuPercent 写入CPU的Points数据到InfluxDB中
func writesCpuPoints(percent float64) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", // 精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{"cpu": "cpu_percent"}
	fields := map[string]interface{}{
		"cpu_percent": percent,
	}
	pt, err := client.NewPoint("cpu_percent", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert cpu percent success")
}
