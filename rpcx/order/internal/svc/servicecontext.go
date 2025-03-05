package svc

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/simance-ai/smdx/rpcx/order/internal/config"
	"github.com/simance-ai/smdx/rpcx/order/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel model.OrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sql.Open("postgres", c.DBConf.DSN)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewOrderModel(sqlx.NewSqlConnFromDB(db), c.Cache),
	}
}
