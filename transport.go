package simplebalancer

import (
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

type BalancerTransport struct {
	Transport       http.RoundTripper
	Expiration      time.Duration
	CleanupInterval time.Duration
	cache           *cache.Cache
}

func (tr *BalancerTransport) lookup(hostname string) []string {
	if tr.cache != nil {
		if addrs1, ok := tr.cache.Get(hostname); ok {
			addrs, _ := addrs1.([]string)
			return addrs
		}
	} else {
		tr.cache = cache.New(tr.Expiration, tr.CleanupInterval)
	}
	addrs, _ := net.LookupHost(hostname)
	tr.cache.SetDefault(hostname, addrs)
	return addrs
}

func (tr *BalancerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme == "http" {
		addrs := tr.lookup(req.URL.Hostname())
		if len(addrs) > 1 {
			req.URL.Host = addrs[rand.Int()%len(addrs)] + ":" + req.URL.Port()
		}
	}
	if tr.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return tr.Transport.RoundTrip(req)
}
