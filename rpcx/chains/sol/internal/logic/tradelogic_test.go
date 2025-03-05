package logic

import (
	"context"
	"flag"
	"log"
	"testing"

	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/chains/common/trade"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/config"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestRaydiumAmm_trade(t *testing.T) {
	type args struct {
		TradeRequest *sol.TradeRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test buy market 0.001 wsol",
			args: args{
				TradeRequest: &sol.TradeRequest{
					TradeRequest: &trade.TradeRequest{
						MarketId:         1,
						AccountId:        0,
						TradeType:        common.TradeBuy,
						TradeBaseAmount:  0.001,
						TradeQuoteAmount: 0,
						TradeMode:        common.TradeModeFaster,
						TradeFee:         0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test buy market 0.001 wsol, new trade mode",
			args: args{
				TradeRequest: &sol.TradeRequest{
					TradeRequest: &trade.TradeRequest{
						MarketId:         2,
						AccountId:        0,
						TradeType:        common.TradeBuy,
						TradeBaseAmount:  0.001,
						TradeQuoteAmount: 0,
						TradeMode:        common.TradeModeFaster,
						TradeFee:         0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test sell market 0.001 wsol, new trade mode",
			args: args{
				TradeRequest: &sol.TradeRequest{
					TradeRequest: &trade.TradeRequest{
						MarketId:         3,
						AccountId:        0,
						TradeType:        common.TradeSell,
						TradeBaseAmount:  0.001,
						TradeQuoteAmount: 20,
						TradeMode:        common.TradeModeFaster,
						TradeFee:         0,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := processTrade(tt.args.TradeRequest); (err != nil) != tt.wantErr {
				t.Errorf("processTrade() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func processTrade(request *sol.TradeRequest) error {

	var configFile = flag.String("f", "./rpcx/chains/sol/etc/sol.yaml", "the config file")

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	ctxB := context.Background()
	logic := NewTradeLogic(ctxB, ctx)

	resp, err := logic.Trade(request)
	if err != nil {
		return err
	}

	if resp != nil {
		log.Println("Trade success: ", resp)
	}

	return nil
}
