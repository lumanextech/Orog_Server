package sol_scan

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var conf = Config{
	Name:   "smdx",
	Host:   "https://pro-api.solscan.io",
	ApiKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOjE3MzU1ODU1MDQ3NDMsImVtYWlsIjoiZWRjaGlyaXNAZ21haWwuY29tIiwiYWN0aW9uIjoidG9rZW4tYXBpIiwiYXBpVmVyc2lvbiI6InYyIiwiaWF0IjoxNzM1NTg1NTA0fQ.s4e5-R4VYnHP4dPV_EJQ1DRf5DhbTo4ntjLlFi8wvxQ",
}

func TestClient(t *testing.T) {

	client := NewClient(conf)

	resp, err := client.get("/v2.0/account/detail", map[string]string{
		"address": "7JFphvzpvcQjJV7Kv7FJnQX8h6GDKNrpUFMpKg7Zc9Tu",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.String())
}

func TestPost(t *testing.T) {

	client := NewClient(conf)

	resp, err := client.post("", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.String())
}

func TestTokenCases(t *testing.T) {

	cases := []struct {
		name   string
		method string
		path   string
		params map[string]string
		resp   interface{}
	}{
		// {method: "GET", path: "/v2.0/token/top", params: nil, resp: nil}, //获取代币排行
		//获取代币市场信息
		{name: "getTokenMarketInfo", method: "GET", path: "/v2.0/token/market/info", params: map[string]string{
			"address": "5zosizWVtXjGbwoXbjJyGD9CCgrxqPx7wEbL3uSqEh6w",
		}, resp: nil},

		//获取wsol的价格
		{name: "getWsolPrice", method: "GET", path: "/v2.0/token/price", params: map[string]string{
			"address": "So11111111111111111111111111111111111111112",
		}, resp: nil},
	}

	for _, c := range cases {
		client := NewClient(conf)
		resp, err := client.get(c.path, c.params)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("==============CASE ", c.name, "===============")

		switch c.name {
		case "getTokenMarketInfo":
			t.Log(resp.String())
		case "getWsolPrice":
			var v TokenPriceResp
			err = json.Unmarshal(resp.Body(), &v)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(v.Data[0].Price)
		}
	}
}

func TestAccount(t *testing.T) {

	client := NewClient(conf)

	//test case
	cases := []struct {
		method string
		path   string
		params map[string]string
		resp   interface{}
	}{
		{method: "GET", path: "/v2.0/account/balance_change", params: map[string]string{
			"address": "7JFphvzpvcQjJV7Kv7FJnQX8h6GDKNrpUFMpKg7Zc9Tu",
		}, resp: nil},
	}

	for _, c := range cases {
		resp, err := client.get(c.path, c.params)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("==============CASE ", c.path, "===============")

		t.Log(resp.String())
	}
}

func TestGetTokenMarketInfo(t *testing.T) {
	client := NewClient(conf)

	resp, err := client.GetTokenMarketInfo(context.Background(), "HCVjpYRgzMEqf2HuencjonT3puxPLV8ac6AaasWNUzXP")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
}
