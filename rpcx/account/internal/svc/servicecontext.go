package svc

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/simance-ai/smdx/rpcx/account/internal/config"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	AccountModel         model.AccountModel
	UserTokenFollowModel model.UserTokenFollowModel
	RebateDetailModel    model.RebateDetailModel
	RebateRoleModel      model.RebateRoleModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	db, err := sql.Open("postgres", c.DBConf.DSN)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:               c,
		AccountModel:         model.NewAccountModel(sqlx.NewSqlConnFromDB(db), c.Cache),
		UserTokenFollowModel: model.NewUserTokenFollowModel(sqlx.NewSqlConnFromDB(db), c.Cache),
		RebateDetailModel:    model.NewRebateDetailModel(sqlx.NewSqlConnFromDB(db), c.Cache),
		RebateRoleModel:      model.NewRebateRoleModel(sqlx.NewSqlConnFromDB(db), c.Cache),
	}
}
