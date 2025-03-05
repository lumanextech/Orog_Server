package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RebateRoleModel = (*customRebateRoleModel)(nil)

type (
	// RebateRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRebateRoleModel.
	RebateRoleModel interface {
		rebateRoleModel

		FindMany(ctx context.Context) ([]RebateRole, error)
	}

	customRebateRoleModel struct {
		*defaultRebateRoleModel
	}
)

// NewRebateRoleModel returns a model for the database table.
func NewRebateRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RebateRoleModel {
	return &customRebateRoleModel{
		defaultRebateRoleModel: newRebateRoleModel(conn, c, opts...),
	}
}

func (m *defaultRebateRoleModel) FindMany(ctx context.Context) ([]RebateRole, error) {
	var resp []RebateRole
	// 查询所有返佣角色
	query := fmt.Sprintf("SELECT %s FROM %s", rebateRoleRows, m.table)

	// 使用 QueryRowsNoCacheCtx 执行查询，不使用缓存
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return resp, nil
}
