package cmn

import (
	"encoding/binary"
	"net"
)

// https://gist.github.com/kotakanbe/d3059af990252ba89a82
// https://play.golang.org/p/fe-F2k6prlA

// Hosts get all IP address from CIDR
func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil

	default:
		return ips[1 : len(ips)-1], nil
	}
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IpList(start net.IP, finish net.IP) []net.IP {
	l := make([]net.IP, 0)
	si := ip2int(start)
	fi := ip2int(finish)

	// loop through addresses as uint32
	for i := si; i <= fi; i++ {
		// convert back to net.IP
		ip := int2ip(i)
		binary.BigEndian.PutUint32(ip, i)
		l = append(l, ip)
	}

	return l
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
