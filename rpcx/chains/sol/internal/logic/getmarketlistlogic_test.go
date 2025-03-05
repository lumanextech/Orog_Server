package logic

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/config"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/gen/field"
)

func Init() *svc.ServiceContext {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error: %v \n", err)
	}

	os.Setenv("KAFKA_CA_FILE", dir+"/Documents/GitHub/smdx/rpcx/chains/sol/etc/only-4096-ca-cert")

	configPath := dir + "/Documents/GitHub/smdx/rpcx/chains/sol/etc/sol.yaml"

	var configFile = flag.String("f", configPath, "the config file")

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	return ctx
}

func TestGetMarketListSql(t *testing.T) {
	svcCtx := Init()

	ctx := context.Background()

	marketDB := dbx.Use(svcCtx.PgDB).ReadDB().MarketRealTimeDatum

	orderBy := "volume_24h"
	orderOp, exist := marketDB.GetFieldByName(orderBy)
	if !exist {
		orderOp = marketDB.Liquidity
	}

	fieldVolume24h := marketDB.SellVolume24h.AddCol(marketDB.BuyVolume24h)

	if strings.EqualFold(orderBy, "volume_24h") {
		orderOp = fieldVolume24h.(field.OrderExpr)
	}

	page := 0
	size := 10000
	orderField := orderOp.Desc()

	result, count, err := marketDB.WithContext(ctx).
		Select(marketDB.ALL, fieldVolume24h).
		Order(orderField).
		FindByPage(int(page), int(size))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(count)
	for _, r := range result {
		t.Log(r.SellVolume24h + r.BuyVolume24h)
	}

}
