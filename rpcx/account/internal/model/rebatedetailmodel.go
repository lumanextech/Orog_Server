package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ RebateDetailModel = (*customRebateDetailModel)(nil)

type (
	// RebateDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRebateDetailModel.
	RebateDetailModel interface {
		rebateDetailModel
		InsertMany(ctx context.Context, data []RebateDetail) (sql.Result, error)
		FindAddressRebateAmount(ctx context.Context, address string) (float64, error)
	}

	customRebateDetailModel struct {
		*defaultRebateDetailModel
	}
)

// NewRebateDetailModel returns a model for the database table.
func NewRebateDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RebateDetailModel {
	return &customRebateDetailModel{
		defaultRebateDetailModel: newRebateDetailModel(conn, c, opts...),
	}
}

func (m *defaultRebateDetailModel) InsertMany(ctx context.Context, data []RebateDetail) (sql.Result, error) {
	if len(data) == 0 {
		return nil, errors.New("插入数据不能为空")
	}

	var keys []string
	for _, d := range data {
		keys = append(keys, fmt.Sprintf("%s%v", cachePublicRebateDetailIdPrefix, d.Address))
	}

	query := fmt.Sprintf("INSERT INTO %s (rebate_amount, address) VALUES ", m.table)

	var args []interface{}
	var placeholders []string
	for i, d := range data {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, d.RebateAmount, d.Address)
	}

	query += strings.Join(placeholders, ", ")

	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query, args...)
	}, keys...)

	return ret, err
}

func (m *defaultRebateDetailModel) FindAddressRebateAmount(ctx context.Context, address string) (float64, error) {
	var totalAmount float64
	query := fmt.Sprintf("SELECT SUM(rebate_amount) FROM %s WHERE address = $1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &totalAmount, query, address)
	if err != nil {
		return 0, err
	}
	return totalAmount, nil
}
