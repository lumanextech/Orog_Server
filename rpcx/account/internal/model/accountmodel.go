package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AccountModel = (*customAccountModel)(nil)

type (
	// AccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountModel.
	AccountModel interface {
		accountModel
		FindRebateAccount(ctx context.Context, address []string) ([]RebateAccount, error)
		UpdateRebateBalance(ctx context.Context, rebateRecords map[string]float64) error
		FindInviteAccount(ctx context.Context, address string) ([]string, error)
		FindOneWithoutCache(ctx context.Context, address string) (*Account, error)

		DelAccountAddressCache(address string) error
	}

	customAccountModel struct {
		*defaultAccountModel
	}

	RebateAccount struct {
		Address        string         `db:"address"`         // 用户地址
		RoleId         sql.NullInt64  `db:"role_id"`         // 角色
		InvitedAddress sql.NullString `db:"invited_address"` // 邀请人地址
	}
)

// NewAccountModel returns a model for the database table.
func NewAccountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AccountModel {
	return &customAccountModel{
		defaultAccountModel: newAccountModel(conn, c, opts...),
	}
}

func (m *customAccountModel) FindRebateAccount(ctx context.Context, addresses []string) ([]RebateAccount, error) {
	var accounts []RebateAccount
	// 将字符串数组转换为 PostgreSQL 兼容的格式
	formattedAddresses := "'" + strings.Join(addresses, "', '") + "'"

	query := fmt.Sprintf(`
	WITH RECURSIVE user_hierarchy AS (
		SELECT address, invited_address, role_id
		FROM account
		WHERE address = ANY(ARRAY[%s])
		UNION ALL
		SELECT u.address, u.invited_address, u.role_id
		FROM account u
		INNER JOIN user_hierarchy uh ON u.address = uh.invited_address
		WHERE uh.invited_address IS NOT NULL
		AND uh.address != uh.invited_address  
	)
	SELECT address, invited_address, role_id FROM user_hierarchy;
	`, formattedAddresses)

	err := m.QueryRowsNoCacheCtx(ctx, &accounts, query)

	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return accounts, nil
}

func (m *customAccountModel) UpdateRebateBalance(ctx context.Context, rebateRecords map[string]float64) error {
	if len(rebateRecords) == 0 {
		return nil // 如果没有数据，直接返回
	}

	// 构建 SQL 语句
	query := fmt.Sprintf("UPDATE %s SET bakance = CASE ", m.table)
	var args []interface{}
	i := 1

	// 生成 `CASE WHEN` 语句
	for address, amount := range rebateRecords {
		query += fmt.Sprintf("WHEN address = $%d THEN bakance + $%d ", i, i+1)
		args = append(args, address, amount)
		i += 2
	}

	// 添加 `WHERE` 条件，确保只更新涉及到的地址
	query += "END WHERE address IN ("
	placeholders := make([]string, 0, len(rebateRecords))
	for j := 1; j <= len(rebateRecords)*2; j += 2 {
		placeholders = append(placeholders, fmt.Sprintf("$%d", j))
	}
	query += strings.Join(placeholders, ", ") + ")"

	// 执行 SQL 更新
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query, args...)
	})
	return err
}

func (m *customAccountModel) FindInviteAccount(ctx context.Context, address string) ([]string, error) {
	var accountAddresses []string
	query := fmt.Sprintf("SELECT address FROM %s WHERE invited_address = '%s'", m.table, address)
	err := m.QueryRowsNoCacheCtx(ctx, &accountAddresses, query)
	if err != nil {
		return nil, err
	}
	return accountAddresses, nil
}

func (m *defaultAccountModel) FindOneWithoutCache(ctx context.Context, address string) (*Account, error) {
	var resp Account
	query := fmt.Sprintf("SELECT %s FROM %s WHERE address = '%s' LIMIT 1", accountRows, m.table, address)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query) // 直接查询数据库，不使用缓存

	if err == sqlc.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultAccountModel) DelAccountAddressCache(address string) error {
	publicAccountAddressKey := fmt.Sprintf("%s%v", cachePublicAccountAddressPrefix, address)
	err := m.DelCache(publicAccountAddressKey) // 传入 ctx
	return err
}
