package sol_scan

import (
	"context"

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

	api.SetHeader("token", conf.ApiKey)

	return client
}

func (c *Client) get(path string, params map[string]string) (*resty.Response, error) {
	return c.api.R().SetQueryParams(params).Get(path)
}

func (c *Client) post(path string, body map[string]interface{}) (*resty.Response, error) {
	return c.api.R().SetBody(body).Post(path)
}
