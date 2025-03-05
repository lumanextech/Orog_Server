package okx_os

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type Client struct {
	conf Config
	api  *resty.Client
	logx.Logger
}

func NewClient(conf Config) *Client {
	api := resty.New()

	api.SetBaseURL(conf.Host)
	client := &Client{
		conf:   conf,
		api:    api,
		Logger: logx.WithContext(context.Background()),
	}

	api.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		var body []byte
		var err error
		path := r.URL
		if r.Method == "GET" {
			if r.QueryParam.Encode() != "" {
				path = path + "?" + r.QueryParam.Encode()
			}
			body = []byte{}
		} else {
			body, err = json.Marshal(r.Body)
			if err != nil {
				return err
			}
		}

		timestamp, sign := client.sign(r.Method, path, string(body))
		r.Header.Add("OK-ACCESS-KEY", conf.ApiKey)
		r.Header.Add("OK-ACCESS-PASSPHRASE", conf.Passphrase)
		r.Header.Add("OK-ACCESS-SIGN", sign)
		r.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		r.Header.Add("OK-ACCESS-PROJECT", conf.Name)
		r.Header.Add("Content-Type", "application/json")
		return nil
	})

	return client
}

func (c *Client) sign(method, path, body string) (string, string) {
	c.Debugf("method: %s, path: %s, body: %s", method, path, body)

	format := "2006-01-02T15:04:05.999Z07:00"
	t := time.Now().UTC().Format(format)
	ts := fmt.Sprint(t)
	s := ts + method + path + body
	p := []byte(s)
	h := hmac.New(sha256.New, []byte(c.conf.ApiSecret))
	h.Write(p)
	return ts, base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c *Client) Get(path string, params map[string]string) (*resty.Response, error) {
	return c.api.R().SetQueryParams(params).Get(path)
}

func (c *Client) Post(path string, body map[string]interface{}) (*resty.Response, error) {
	return c.api.R().SetBody(body).Post(path)
}
