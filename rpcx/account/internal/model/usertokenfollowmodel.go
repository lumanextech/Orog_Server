package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserTokenFollowModel = (*customUserTokenFollowModel)(nil)

var cachePublicUserTokenFollowAddressAndStatusPrefix = "cache:public:userTokenFollow:address:status"

type (
	// UserTokenFollowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTokenFollowModel.
	UserTokenFollowModel interface {
		userTokenFollowModel
		InsertByAddressAndStatus(ctx context.Context, data *UserTokenFollow) (sql.Result, error)
		FindAllByAddressAndStatus(ctx context.Context, address string, status string) ([]*UserTokenFollow, error)
	}

	customUserTokenFollowModel struct {
		*defaultUserTokenFollowModel
	}
)

// NewUserTokenFollowModel returns a model for the database table.
func NewUserTokenFollowModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserTokenFollowModel {
	return &customUserTokenFollowModel{
		defaultUserTokenFollowModel: newUserTokenFollowModel(conn, c, opts...),
	}
}

// InsertByAddressAndStatus
func (m *customUserTokenFollowModel) InsertByAddressAndStatus(ctx context.Context, data *UserTokenFollow) (sql.Result, error) {
	publicUserTokenFollowAddressKey := fmt.Sprintf("%s%v", cachePublicUserTokenFollowAddressPrefix, data.Address)
	publicUserTokenFollowIdKey := fmt.Sprintf("%s%v", cachePublicUserTokenFollowIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table, userTokenFollowRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Address, data.Chain, data.TokenAddress, data.Status, data.FollowedAt, data.UnfollowedAt, data.CreatedAt, data.UpdatedAt, data.Id)
	}, publicUserTokenFollowAddressKey, publicUserTokenFollowIdKey)
	return ret, err
}

func (m *customUserTokenFollowModel) FindAllByAddressAndStatus(ctx context.Context, address string, status string) ([]*UserTokenFollow, error) {
	query := fmt.Sprintf("select %s from %s where address = $1 and status = $2", userTokenFollowRows, m.table)
	var resp []*UserTokenFollow
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, address, status)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
