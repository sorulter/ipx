package main

import (
	"encoding/binary"
	"net"
)

func ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}
func long2ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}

func logshash(id uint64) string {
	var tab = map[int]string{0: "a", 1: "b", 2: "c", 3: "d", 4: "e"}
	return tab[int(id%5)]
}
