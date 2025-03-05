package svc

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/simance-ai/smdx/rpcx/order_consumer/internal/config"
	"github.com/simance-ai/smdx/rpcx/order_consumer/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel model.OrderModel
}

// Deadline implements context.Context.
func (s *ServiceContext) Deadline() (deadline time.Time, ok bool) {
	panic("unimplemented")
}

// Done implements context.Context.
func (s *ServiceContext) Done() <-chan struct{} {
	panic("unimplemented")
}

// Err implements context.Context.
func (s *ServiceContext) Err() error {
	panic("unimplemented")
}

// Value implements context.Context.
func (s *ServiceContext) Value(key any) any {
	panic("unimplemented")
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
