package nexmo

import "net"

var masks = []string{
	"174.37.245.32/29",
	"174.36.197.192/28",
	"173.193.199.16/28",
	"119.81.44.0/28",
}

var subnets []*net.IPNet

func init() {
	subnets = make([]*net.IPNet, len(masks))
	for i, mask := range masks {
		_, net, _ := net.ParseCIDR(mask)
		subnets[i] = net
	}
}

// IsTrustedIP returns true if the provided IP address came from
// a trusted Nexmo server.
func IsTrustedIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)

	for _, net := range subnets {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}
