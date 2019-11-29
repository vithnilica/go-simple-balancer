package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	simplebalancer "github.com/vithnilica/go-simple-balancer"
)

type nullTransport struct {
}

func (tr *nullTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println(req.URL)
	return httptest.NewRecorder().Result(), nil
}

func main() {
	client := http.Client{
		Transport: &simplebalancer.BalancerTransport{
			Transport:       &nullTransport{},
			Expiration:      time.Second,
			CleanupInterval: time.Minute,
		},
	}

	for i := 0; i < 10; i++ {
		resp, err := client.Get("http://google.com")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}

}
