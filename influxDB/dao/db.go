package dao

import (
	client "github.com/influxdata/influxdb1-client/v2"
)

// DB 全局InfluxDB连接对象
var DB client.Client

// init 初始化InfluxDB连接
func init() {
	DB, _ = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
}
