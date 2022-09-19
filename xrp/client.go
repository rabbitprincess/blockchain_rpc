package xrp

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
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

func (t *Client) cmdSendAndRecv(req CmdReq, res CmdRes) (err error) {
	// marshal
	btReq, err := req.Marshal()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(btReq)

	// httpReq
	httpReq, err := http.NewRequest("POST", t.url, buf)
	if err != nil {
		return err
	}
	httpRes, err := t.rpc_client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()
	btRes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	// unmarshal
	err = res.Unmarshal(btRes)
	if err != nil {
		return err
	}
	// 후처리 - error
	pt_err := res.Error()
	if pt_err.ErrCode != Ok {
		return pt_err
	}

	return nil
}
