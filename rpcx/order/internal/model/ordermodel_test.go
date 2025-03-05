package model_test

import (
	"context"
	"fmt"
	"github.com/simance-ai/smdx/rpcx/order/internal/config"
	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"os"
	"path/filepath"
	"testing"
)

func TestCustomOrderModel_FindRebateOrders(t *testing.T) {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// 使用 filepath.Join 拼接路径
	configFile := filepath.Join(dir, "../../", "etc", "order.yaml")

	var c config.Config
	conf.MustLoad(configFile, &c)
	ctx := svc.NewServiceContext(c)

	orders, err := ctx.OrderModel.FindRebateOrders(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(orders)
}
