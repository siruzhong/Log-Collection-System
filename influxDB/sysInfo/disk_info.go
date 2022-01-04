package sysInfo

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/disk"
	"influxDB/dao"
	"log"
	"time"
)

// DiskInfo 磁盘相关数据结构体
type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat // 分区使用统计情况(key为挂载点目录，value为UsageStat结构体)
}

// GetDiskInfo 获取磁盘相关信息
func GetDiskInfo() *DiskInfo {
	var diskInfo = &DiskInfo{PartitionUsageStat: make(map[string]*disk.UsageStat, 16)}
	parts, _ := disk.Partitions(true) // 获取所有磁盘分区
	for _, part := range parts {
		usageStat, _ := disk.Usage(part.Mountpoint) // 获取该分区挂载点的使用统计
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
	}
	return diskInfo
}

// WritesDiskPoints 写入disk的Points数据到InfluxDB中
func WritesDiskPoints(diskInfo *DiskInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", // 精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	for path, usageStat := range diskInfo.PartitionUsageStat {
		tags := map[string]string{"mountPoint": path}
		fields := map[string]interface{}{
			"path":              usageStat.Path,
			"fstype":            usageStat.Fstype,
			"total":             int64(usageStat.Total),
			"free":              int64(usageStat.Free),
			"used":              int64(usageStat.Used),
			"usedPercent":       usageStat.UsedPercent,
			"inodesTotal":       int64(usageStat.InodesTotal),
			"inodesUsed":        int64(usageStat.InodesUsed),
			"inodesFree":        int64(usageStat.InodesFree),
			"inodesUsedPercent": usageStat.InodesUsedPercent,
		}
		pt, err := client.NewPoint("disk_info", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	err = dao.DB.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert disk info success")
}
