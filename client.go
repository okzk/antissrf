package antissrf

import (
	"errors"
	"net"
	"net/http"
	"syscall"
	"time"
)

var errorNotTCPAddress = errors.New("not a tcp address")
var errorIsNotGlobalUnicastAddress = errors.New("not a global unicast address")
var errorPrivateAddress = errors.New("private address")
var errorUniqueLocalAddress = errors.New("unique local address")
var errorBlackListedAddress = errors.New("blacklisted by user")

var privateV4 = []*net.IPNet{
	MustParseCIDR("10.0.0.0/8"),
	MustParseCIDR("172.16.0.0/12"),
	MustParseCIDR("192.168.0.0/16"),
}

var localUnicast = MustParseCIDR("fc00::/7")

type netList []*net.IPNet

func MustParseCIDR(cidr string) *net.IPNet {
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return n
}

func (list netList) filter(network, address string, _ syscall.RawConn) error {
	if network != "tcp4" && network != "tcp6" {
		return errorNotTCPAddress
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return errorNotTCPAddress
	}

	addr := net.ParseIP(host)
	if addr == nil {
		return errorNotTCPAddress
	} else if !addr.IsGlobalUnicast() {
		return errorIsNotGlobalUnicastAddress
	} else if addrV4 := addr.To4(); addrV4 != nil {
		for _, n := range privateV4 {
			if n.Contains(addrV4) {
				return errorPrivateAddress
			}
		}
	} else {
		if localUnicast.Contains(addr) {
			return errorUniqueLocalAddress
		}
	}
	for _, n := range list {
		if n.Contains(addr) {
			return errorBlackListedAddress
		}
	}
	return nil
}

func Transport(ngNetwork ...*net.IPNet) *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	d := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
		Control:   netList(ngNetwork).filter,
	}
	t.DialContext = d.DialContext
	return t
}

func Client(ngNetwork ...*net.IPNet) *http.Client {
	return &http.Client{
		Transport: Transport(ngNetwork...),
	}
}
