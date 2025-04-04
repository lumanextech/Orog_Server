// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dbx

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
)

func newMarketKline15m(db *gorm.DB, opts ...gen.DOOption) marketKline15m {
	_marketKline15m := marketKline15m{}

	_marketKline15m.marketKline15mDo.UseDB(db, opts...)
	_marketKline15m.marketKline15mDo.UseModel(&model.MarketKline15m{})

	tableName := _marketKline15m.marketKline15mDo.TableName()
	_marketKline15m.ALL = field.NewAsterisk(tableName)
	_marketKline15m.MarketAddress = field.NewString(tableName, "market_address")
	_marketKline15m.O = field.NewFloat64(tableName, "o")
	_marketKline15m.H = field.NewFloat64(tableName, "h")
	_marketKline15m.L = field.NewFloat64(tableName, "l")
	_marketKline15m.C = field.NewFloat64(tableName, "c")
	_marketKline15m.V = field.NewFloat64(tableName, "v")
	_marketKline15m.Timestamp = field.NewTime(tableName, "timestamp")
	_marketKline15m.UpdatedAt = field.NewTime(tableName, "updated_at")
	_marketKline15m.CreatedAt = field.NewTime(tableName, "created_at")

	_marketKline15m.fillFieldMap()

	return _marketKline15m
}

type marketKline15m struct {
	marketKline15mDo marketKline15mDo

	ALL           field.Asterisk
	MarketAddress field.String
	O             field.Float64
	H             field.Float64
	L             field.Float64
	C             field.Float64
	V             field.Float64
	Timestamp     field.Time
	UpdatedAt     field.Time
	CreatedAt     field.Time

	fieldMap map[string]field.Expr
}

func (m marketKline15m) Table(newTableName string) *marketKline15m {
	m.marketKline15mDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m marketKline15m) As(alias string) *marketKline15m {
	m.marketKline15mDo.DO = *(m.marketKline15mDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *marketKline15m) updateTableName(table string) *marketKline15m {
	m.ALL = field.NewAsterisk(table)
	m.MarketAddress = field.NewString(table, "market_address")
	m.O = field.NewFloat64(table, "o")
	m.H = field.NewFloat64(table, "h")
	m.L = field.NewFloat64(table, "l")
	m.C = field.NewFloat64(table, "c")
	m.V = field.NewFloat64(table, "v")
	m.Timestamp = field.NewTime(table, "timestamp")
	m.UpdatedAt = field.NewTime(table, "updated_at")
	m.CreatedAt = field.NewTime(table, "created_at")

	m.fillFieldMap()

	return m
}

func (m *marketKline15m) WithContext(ctx context.Context) *marketKline15mDo {
	return m.marketKline15mDo.WithContext(ctx)
}

func (m marketKline15m) TableName() string { return m.marketKline15mDo.TableName() }

func (m marketKline15m) Alias() string { return m.marketKline15mDo.Alias() }

func (m marketKline15m) Columns(cols ...field.Expr) gen.Columns {
	return m.marketKline15mDo.Columns(cols...)
}

func (m *marketKline15m) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *marketKline15m) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 9)
	m.fieldMap["market_address"] = m.MarketAddress
	m.fieldMap["o"] = m.O
	m.fieldMap["h"] = m.H
	m.fieldMap["l"] = m.L
	m.fieldMap["c"] = m.C
	m.fieldMap["v"] = m.V
	m.fieldMap["timestamp"] = m.Timestamp
	m.fieldMap["updated_at"] = m.UpdatedAt
	m.fieldMap["created_at"] = m.CreatedAt
}

func (m marketKline15m) clone(db *gorm.DB) marketKline15m {
	m.marketKline15mDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m marketKline15m) replaceDB(db *gorm.DB) marketKline15m {
	m.marketKline15mDo.ReplaceDB(db)
	return m
}

type marketKline15mDo struct{ gen.DO }

func (m marketKline15mDo) Debug() *marketKline15mDo {
	return m.withDO(m.DO.Debug())
}

func (m marketKline15mDo) WithContext(ctx context.Context) *marketKline15mDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m marketKline15mDo) ReadDB() *marketKline15mDo {
	return m.Clauses(dbresolver.Read)
}

func (m marketKline15mDo) WriteDB() *marketKline15mDo {
	return m.Clauses(dbresolver.Write)
}

func (m marketKline15mDo) Session(config *gorm.Session) *marketKline15mDo {
	return m.withDO(m.DO.Session(config))
}

func (m marketKline15mDo) Clauses(conds ...clause.Expression) *marketKline15mDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m marketKline15mDo) Returning(value interface{}, columns ...string) *marketKline15mDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m marketKline15mDo) Not(conds ...gen.Condition) *marketKline15mDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m marketKline15mDo) Or(conds ...gen.Condition) *marketKline15mDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m marketKline15mDo) Select(conds ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m marketKline15mDo) Where(conds ...gen.Condition) *marketKline15mDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m marketKline15mDo) Order(conds ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m marketKline15mDo) Distinct(cols ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m marketKline15mDo) Omit(cols ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m marketKline15mDo) Join(table schema.Tabler, on ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m marketKline15mDo) LeftJoin(table schema.Tabler, on ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m marketKline15mDo) RightJoin(table schema.Tabler, on ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m marketKline15mDo) Group(cols ...field.Expr) *marketKline15mDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m marketKline15mDo) Having(conds ...gen.Condition) *marketKline15mDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m marketKline15mDo) Limit(limit int) *marketKline15mDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m marketKline15mDo) Offset(offset int) *marketKline15mDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m marketKline15mDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *marketKline15mDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m marketKline15mDo) Unscoped() *marketKline15mDo {
	return m.withDO(m.DO.Unscoped())
}

func (m marketKline15mDo) Create(values ...*model.MarketKline15m) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m marketKline15mDo) CreateInBatches(values []*model.MarketKline15m, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m marketKline15mDo) Save(values ...*model.MarketKline15m) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m marketKline15mDo) First() (*model.MarketKline15m, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.MarketKline15m), nil
	}
}

func (m marketKline15mDo) Take() (*model.MarketKline15m, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.MarketKline15m), nil
	}
}

func (m marketKline15mDo) Last() (*model.MarketKline15m, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.MarketKline15m), nil
	}
}

func (m marketKline15mDo) Find() ([]*model.MarketKline15m, error) {
	result, err := m.DO.Find()
	return result.([]*model.MarketKline15m), err
}

func (m marketKline15mDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.MarketKline15m, err error) {
	buf := make([]*model.MarketKline15m, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m marketKline15mDo) FindInBatches(result *[]*model.MarketKline15m, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m marketKline15mDo) Attrs(attrs ...field.AssignExpr) *marketKline15mDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m marketKline15mDo) Assign(attrs ...field.AssignExpr) *marketKline15mDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m marketKline15mDo) Joins(fields ...field.RelationField) *marketKline15mDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m marketKline15mDo) Preload(fields ...field.RelationField) *marketKline15mDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m marketKline15mDo) FirstOrInit() (*model.MarketKline15m, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.MarketKline15m), nil
	}
}

func (m marketKline15mDo) FirstOrCreate() (*model.MarketKline15m, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.MarketKline15m), nil
	}
}

func (m marketKline15mDo) FindByPage(offset int, limit int) (result []*model.MarketKline15m, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m marketKline15mDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m marketKline15mDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m marketKline15mDo) Delete(models ...*model.MarketKline15m) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *marketKline15mDo) withDO(do gen.Dao) *marketKline15mDo {
	m.DO = *do.(*gen.DO)
	return m
}
