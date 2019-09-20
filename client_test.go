package antissrf

import (
	"fmt"
	"net"
	"net/url"
	"testing"
)

func TestClient(t *testing.T) {
	errTestCase := []struct {
		host     string
		expected error
	}{
		{"127.0.0.10", errorIsNotGlobalUnicastAddress},
		{"[::1]", errorIsNotGlobalUnicastAddress},
		{"229.10.10.10", errorIsNotGlobalUnicastAddress},
		{"10.11.12.13", errorPrivateAddress},
		{"172.20.10.10", errorPrivateAddress},
		{"192.168.10.10", errorPrivateAddress},
		{"169.254.169.254", errorIsNotGlobalUnicastAddress},
		{"localhost", errorIsNotGlobalUnicastAddress},
		{"192.0.2.100", errorBlackListedAddress},
	}

	client := Client(MustParseCIDR("192.0.2.0/24"))
	for _, c := range errTestCase {
		_, err := client.Get(fmt.Sprintf("http://%s/", c.host))
		uerr := err.(*url.Error)
		oerr := uerr.Err.(*net.OpError)
		if oerr.Err != c.expected {
			t.Errorf("unexpected err. expected: %v, actual: %v", c.expected, oerr.Err)
		}
	}

	res, err := client.Get("http://example.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		res.Body.Close()
	}
}
