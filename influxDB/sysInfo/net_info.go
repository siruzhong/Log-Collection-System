package sysInfo

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/net"
	"influxDB/dao"
	"log"
	"time"
)

var (
	// lastGetNetInfoTime 上次获取网卡数据的时间戳
	lastGetNetInfoTime int64

	// lastNetInfo 上次所获取的网卡数据
	lastNetInfo *NetInfo
)

// NetInfo 网卡相关数据结构体
type NetInfo struct {
	NetIOCountersStat map[string]*net.IOCountersStat // 网卡信息(key为网卡名称,value为网络IO速率)
}

// NetRateInfo 网络速率结构体
type NetRateInfo struct {
	NetRateInfo map[string]*NetIORate
}

// NetIORate 网络IO速率结构体
type NetIORate struct {
	BytesSentRate   float64 `json:"bytes_sent_rate"`   // 字节发送速率
	BytesRecvRate   float64 `json:"bytes_recv_rate"`   // 字节接收速率
	PacketsSentRate float64 `json:"packets_sent_rate"` // 包发送速率
	PacketsRecvRate float64 `json:"packets_recv_rate"` // 包接收速率
}

// GetNetRateInfo
func GetNetRateInfo() *NetRateInfo {
	var netInfo = &NetInfo{NetIOCountersStat: make(map[string]*net.IOCountersStat, 16)}
	var netRateInfo = &NetRateInfo{NetRateInfo: make(map[string]*NetIORate, 16)}
	currentTimeStamp := time.Now().Unix() // 获取当前时间戳
	ioCounters, _ := net.IOCounters(true)
	for _, netData := range ioCounters {
		netInfo.NetIOCountersStat[netData.Name] = &netData
		if lastGetNetInfoTime == 0 || lastNetInfo == nil {
			continue
		}
		// 获取时间差
		timeDifference := currentTimeStamp - lastGetNetInfoTime
		netIORate := &NetIORate{
			// 开始计算网络IO速率
			BytesSentRate:   (float64(netData.BytesSent) - float64(lastNetInfo.NetIOCountersStat[netData.Name].BytesSent)) / float64(timeDifference),
			BytesRecvRate:   (float64(netData.BytesRecv) - float64(lastNetInfo.NetIOCountersStat[netData.Name].BytesRecv)) / float64(timeDifference),
			PacketsSentRate: (float64(netData.PacketsSent) - float64(lastNetInfo.NetIOCountersStat[netData.Name].PacketsSent)) / float64(timeDifference),
			PacketsRecvRate: (float64(netData.PacketsRecv) - float64(lastNetInfo.NetIOCountersStat[netData.Name].PacketsRecv)) / float64(timeDifference),
		}
		netRateInfo.NetRateInfo[netData.Name] = netIORate
	}
	// 更新全局lastGetNetInfoTime和lastNetInfo变量
	lastGetNetInfoTime = currentTimeStamp
	lastNetInfo = netInfo
	return netRateInfo
}

// WritesNetRatePoints 写入netRate的Points数据到InfluxDB中
func WritesNetRatePoints(netRateInfo *NetRateInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", // 精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	for netCardName, netRate := range netRateInfo.NetRateInfo {
		tags := map[string]string{"netCard_name": netCardName}
		fields := map[string]interface{}{
			"bytes_sent_rate":   netRate.BytesSentRate,
			"bytes_recv_rate":   netRate.BytesRecvRate,
			"packets_sent_rate": netRate.PacketsSentRate,
			"packets_recv_rate": netRate.PacketsRecvRate,
		}
		pt, err := client.NewPoint("net_rate_info", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}
	err = dao.DB.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert net rate into InfluxDB success")
}
