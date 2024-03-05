package xrp

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	url string
	rpc *http.Client
}

func (t *Client) Init(url string) (err error) {
	// init client http
	t.url = url
	t.rpc = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return nil
}
