package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	cacheOrog1OrderIdPrefix = "cache:orog1:order:id:"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel
		FindOrders(ctx context.Context, pageSize, page, userId int64, status, chainName string) ([]Order, int64, error)
		FindRebateOrders(ctx context.Context) ([]Order, error)
		UpdateRebateOrders(ctx context.Context, orderId []int64) error
		FindOrderById(ctx context.Context, id int64) (*Order, error)
		FindOrderTradeNumber(ctx context.Context, address string) (int64, int64, error)
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c, opts...),
	}
}

func (c *customOrderModel) FindOrderById(ctx context.Context, id int64) (*Order, error) {
	orog1OrderIdKey := fmt.Sprintf("%s%v", cacheOrog1OrderIdPrefix, id)
	var resp Order
	err := c.QueryRowCtx(ctx, &resp, orog1OrderIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", orderRows, c.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customOrderModel) FindOrders(ctx context.Context, pageSize, page, userId int64, status, chainName string) ([]Order, int64, error) {
	var orders []Order
	var totalCount int64

	// Construct the base query
	query := fmt.Sprintf("SELECT %s FROM %s WHERE 1=1", orderRows, c.table)

	// Add filters to the query
	if userId != 0 {
		query += " AND user_id = ?"
	}
	if status != "" {
		query += " AND status = ?"
	}
	if chainName != "" {
		query += " AND chain_name = ?"
	}

	// Add pagination
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	// Execute the query
	err := c.QueryRowsNoCacheCtx(ctx, &orders, query, userId, status, chainName)
	if err != nil {
		return nil, 0, err
	}

	// Get the total count of orders
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE 1=1", c.table)
	if userId != 0 {
		countQuery += " AND user_id = ?"
	}
	if status != "" {
		countQuery += " AND status = ?"
	}
	if chainName != "" {
		countQuery += " AND chain_name = ?"
	}
	err = c.QueryRowNoCacheCtx(ctx, &totalCount, countQuery, userId, status, chainName)
	if err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func (c *customOrderModel) FindRebateOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	status := int64(1)        // 已成交
	side := int64(0)          // buy
	paymentStatus := int64(2) // 已上链
	rebateStatus := int64(1)  // ！= 1 未返佣
	// Construct the base query
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE status = $1 AND side = $2 AND payment_status = $3 AND COALESCE(rebate_status, 0) != $4`, orderRows, c.table)
	err := c.QueryRowsNoCacheCtx(ctx, &orders, query, status, side, paymentStatus, rebateStatus)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (c *customOrderModel) UpdateRebateOrders(ctx context.Context, orderIds []int64) error {
	if len(orderIds) == 0 {
		return errors.New("订单 ID 不能为空")
	}

	// 构造 SQL 语句（避免 SQL 注入）
	query := fmt.Sprintf("UPDATE %s SET rebate_status = $1 WHERE id = ANY($2)", c.table)

	// 执行更新
	_, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query, 1, pq.Array(orderIds)) // rebate_status 设为 1
	})

	if err != nil {
		logx.Errorf("批量更新订单返佣状态失败: %v", err)
		return err
	}

	logx.Infof("成功更新 %d 条订单的返佣状态", len(orderIds))
	return nil
}

func (c *customOrderModel) FindOrderTradeNumber(ctx context.Context, address string) (int64, int64, error) {
	var buyTotal, sellTotal int64
	status := int64(1)        // 已成交
	paymentStatus := int64(2) // 已上链

	// 查询买单（side = 0）
	buyQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE status = %d AND side = %d AND payment_status = %d AND account_address = '%s'`,
		c.table, status, 0, paymentStatus, address)
	err := c.QueryRowNoCacheCtx(ctx, &buyTotal, buyQuery)
	if err != nil {
		return 0, 0, err
	}

	// 查询卖单（side = 1）
	sellQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE status = %d AND side = %d AND payment_status = %d AND account_address = '%s'`,
		c.table, status, 1, paymentStatus, address)
	err = c.QueryRowNoCacheCtx(ctx, &sellTotal, sellQuery)
	if err != nil {
		return 0, 0, err
	}

	return buyTotal, sellTotal, nil
}
