package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/simance-ai/smdx/rpcx/order_consumer/internal/config"
	"github.com/simance-ai/smdx/rpcx/order_consumer/internal/model"
	"github.com/simance-ai/smdx/rpcx/order_consumer/internal/server"
	"github.com/simance-ai/smdx/rpcx/order_consumer/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order_consumer/order_consumer"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PriceData struct {
	Price     float64
	Timestamp time.Time
	Chain     string
}

var configFile = flag.String("f", "etc/orderconsumer.yaml", "the config file")

func fetchPrices(priceChan chan<- map[string]PriceData, wg *sync.WaitGroup) {
	defer wg.Done()
	tokenAddress := "7AcpfFuRJMRGhxvLQg4g7TCWCK22Fv83CLVrt9hEpump"
	chains := []string{"Ethereum", "sol"}

	rand.Seed(time.Now().UnixNano())

	for {
		priceBatch := make(map[string]PriceData)
		for _, chain := range chains {
			for i := 0; i < 10; i++ {
				priceBatch[fmt.Sprintf("%s_%d_%s", tokenAddress, i, chain)] = PriceData{
					Price:     rand.Float64(),
					Timestamp: time.Now(),
					Chain:     chain,
				}
			}
		}
		priceChan <- priceBatch
	}
}

type orderModel struct {
	svcCtx *svc.ServiceContext
}

func (m *orderModel) loadOrders(ctx context.Context, chainId string, orderType, side, paymentStatus int64) []model.Order {
	orders, err := m.svcCtx.OrderModel.QueryOrders(ctx, chainId, orderType, side, paymentStatus)
	if err != nil {
		fmt.Printf("Error querying orders: %v\n", err)
		return nil
	}
	return orders
}

func (m *orderModel) processOrders(priceChan <-chan map[string]PriceData, solBuyOrderChan, solSellOrderChan chan<- model.Order, wg *sync.WaitGroup) {
	defer wg.Done()

	for prices := range priceChan {
		for _, data := range prices {
			orders := m.loadOrders(context.Background(), data.Chain, 0, 0, 0)
			go func(data PriceData) {
				for _, order := range orders {
					if order.ChainId == "sol" && data.Chain == "sol" {
						if order.Type == 0 && data.Price < order.Amount {
							solBuyOrderChan <- order
						} else if order.Type == 1 && data.Price > order.Amount {
							solSellOrderChan <- order
						}
					}
				}
			}(data)
		}
	}
}

func processBuyOnChain(solBuyOrderChan <-chan model.Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range solBuyOrderChan {
		fmt.Printf("Processing Buy On-Chain for Order ID: %s, Token: %s, Chain: %s\n", order.Id, order.MarketAddress, order.ChainId)
	}
}

func processSellOnChain(solSellOrderChan <-chan model.Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range solSellOrderChan {
		fmt.Printf("Processing Sell On-Chain for Order ID: %s, Token: %s, Chain: %s\n", order.Id, order.MarketAddress, order.ChainId)
	}
}

func (m *orderModel) startPriceFetchingAndOrderMatching() {
	priceChan := make(chan map[string]PriceData, 100000)
	solBuyOrderChan := make(chan model.Order, 100000)
	solSellOrderChan := make(chan model.Order, 100000)
	var wg sync.WaitGroup

	go fetchPrices(priceChan, &wg)

	go m.processOrders(priceChan, solBuyOrderChan, solSellOrderChan, &wg)

	go processBuyOnChain(solBuyOrderChan, &wg)

	go processSellOnChain(solSellOrderChan, &wg)

}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	model := &orderModel{svcCtx: ctx}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order_consumer.RegisterOrderConsumerServer(grpcServer, server.NewOrderConsumerServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	model.startPriceFetchingAndOrderMatching()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
