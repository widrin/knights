package ip

import (
	"bytes"
	"net"
	"tzgit.kaixinxiyou.com/utils/common/log"
)

func GetAllIP() []net.IP {
	iFaces, err := net.Interfaces()
	ret := make([]net.IP, 0)
	if err != nil {
		return ret
	}

	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iFace.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iFace.Addrs()
		if err != nil {
			log.Error("%v", err)
			continue
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			ret = append(ret, ip)
		}
	}
	return ret
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}

//判断本机ip 是否是有C类ip
func GetHaveIpC() bool {
	allIp := GetAllIP()
	ip1 := net.ParseIP("192.168.0.0").To4()
	ip2 := net.ParseIP("192.168.255.255").To4()
	for _, v := range allIp {
		if bytes.Compare(v, ip1) >= 0 && bytes.Compare(v, ip2) <= 0 {
			return true
		}
	}
	return false
}

//判断本机ip 是否是有B类ip
func GetHaveIpB() bool {
	allIp := GetAllIP()
	ip1 := net.ParseIP("172.16.0.0").To4()
	ip2 := net.ParseIP("172.31.255.255").To4()
	for _, v := range allIp {
		if bytes.Compare(v, ip1) >= 0 && bytes.Compare(v, ip2) <= 0 {
			return true
		}
	}
	return false
}

//判断本机ip 是否是有A类ip
func GetHaveIpA() bool {
	allIp := GetAllIP()
	ip1 := net.ParseIP("10.0.0.0").To4()
	ip2 := net.ParseIP("10.255.255.255").To4()
	for _, v := range allIp {
		if bytes.Compare(v, ip1) >= 0 && bytes.Compare(v, ip2) <= 0 {
			return true
		}
	}
	return false
}
