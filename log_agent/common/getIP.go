package common

import (
	"net"
	"strings"
)

// GetNativeIP 获取本机IP地址
func GetNativeIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "get native ip failed", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.IP.String(), ":")[0], err
}
