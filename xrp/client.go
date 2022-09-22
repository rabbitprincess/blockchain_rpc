package xrp

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	url        string
	rpc_client *http.Client
}

func (t *Client) Init(_s_url string) (err error) {
	// init client http
	t.url = _s_url
	t.rpc_client = &http.Client{
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
