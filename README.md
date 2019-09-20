# antissrf

This package provides an anti SSRF(Server Side Request Forgery) http client.

This client will return an error if the remote address is:
- loopback address
- multicast address
- link local address
- private address
- unique local address
- additionally blacklisted by user
 
## Usage

`antissrf.Client()` just returns `*http.Client`

```golang
var client = antissrf.Client()

func main() {
    // OK
    res, err := client.Get("http://example.com")

    // NG
    res, err := client.Get("http://169.254.169.254/")
}
```

If you want to blacklist additional address spaces:

```golang
var client = antissrf.Client(
	antissrf.MustParseCIDR("192.0.2.0/24"),
	antissrf.MustParseCIDR("198.51.100.0/24"),
	antissrf.MustParseCIDR("203.0.113.0/24"))
```
