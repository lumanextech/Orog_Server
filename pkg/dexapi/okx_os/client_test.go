package okx_os

import "testing"

func TestClient(t *testing.T) {
	conf := Config{
		Name:       "smdx",
		Host:       "https://www.okx.com",
		ApiKey:     "e381056b-3c47-43c5-aa7f-1971547760aa",
		ApiSecret:  "D6058FC92B599F55D0883B5F89E9BF71",
		Passphrase: "555IloveZz@",
	}

	client := NewClient(conf)

	resp, err := client.Get("/api/v5/wallet/chain/supported-chains", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.String())
}

func TestPost(t *testing.T) {
	conf := Config{
		Name:       "smdx",
		Host:       "https://www.okx.com",
		ApiKey:     "e381056b-3c47-43c5-aa7f-1971547760aa",
		ApiSecret:  "D6058FC92B599F55D0883B5F89E9BF71",
		Passphrase: "555IloveZz@",
	}

	client := NewClient(conf)

	resp, err := client.Post("/api/v5/wallet/chain/supported-chains", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp.String())
}

func TestWalletCases(t *testing.T) {
	conf := Config{
		//Name:       "smdx",
		Host:       "https://www.okx.com",
		ApiKey:     "e381056b-3c47-43c5-aa7f-1971547760aa",
		ApiSecret:  "D6058FC92B599F55D0883B5F89E9BF71",
		Passphrase: "555IloveZz@",
	}
	cases := []struct {
		method string
		path   string
		params map[string]string
		resp   interface{}
	}{
		{method: "GET", path: "/api/v5/wallet/chain/supported-chains", params: nil, resp: nil},
		{method: "GET", path: "/api/v5/wallet/asset/all-token-balances-by-address",
			params: map[string]string{"chains": "501",
				"address": "1ZRxijoyTqok5zWU5ofo6FdGPwZrvyGPTFttAmeYTZi"}, resp: nil}, //获取用户资产
		{method: "GET", path: "/api/v5/wallet/token/token-detail",
			params: map[string]string{"chainIndex": "501", "tokenAddress": "GJAFwWjJ3vnTsrQVabjBVK2TYB1YtRCQXRDfDgUnpump"}, resp: nil}, //获取代币详情

	}

	for _, c := range cases {
		client := NewClient(conf)
		resp, err := client.Get(c.path, c.params)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("==============CASE ", c.path, "===============")

		t.Log(resp)
	}
}

func TestDexCases(t *testing.T) {
	conf := Config{
		Name:       "smdx",
		Host:       "https://www.okx.com",
		ApiKey:     "e381056b-3c47-43c5-aa7f-1971547760aa",
		ApiSecret:  "D6058FC92B599F55D0883B5F89E9BF71",
		Passphrase: "555IloveZz@",
	}
	cases := []struct {
		method string
		path   string
		params map[string]string
		resp   interface{}
	}{
		{method: "GET", path: "/api/v5/dex/aggregator/all-tokens", params: map[string]string{"chainId": "501"}, resp: nil},    //获取支持的代币列表
		{method: "GET", path: "/api/v5/dex/aggregator/get-liquidity", params: map[string]string{"chainId": "501"}, resp: nil}, //获取支持的流动池提供商列表
		{method: "POST", path: "/api/v5/dex/aggregator/approve-transaction",
			params: map[string]string{
				"chainId":              "501",
				"tokenContractAddress": "7NoYCzhP3UzjSHTgVoveRCvgLPUhXfyYqFJAnyM6g9e9",
				"approveAmount":        "1000000000000000000"}, resp: nil}, //获取授权合约地址信息(sol chain not need approveAmount)
		//获取兑换价格
		{method: "GET", path: "/api/v5/dex/aggregator/quote",
			params: map[string]string{
				"chainId":                         "501",
				"amount":                          "100000000",
				"fromTokenAddress":                "GJAFwWjJ3vnTsrQVabjBVK2TYB1YtRCQXRDfDgUnpump", //mint address
				"toTokenAddress":                  "So11111111111111111111111111111111111111112",  //wsol mint address
				"dexIds":                          "277",
				"priceImpactProtectionPercentage": "0.01",
				"feePercent":                      "0",
			},
			resp: nil},
	}

	for _, c := range cases {
		client := NewClient(conf)
		resp, err := client.Get(c.path, c.params)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("==============CASE ", c.path, "===============")

		t.Log(resp)
	}
}
