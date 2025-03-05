package logic

import (
	"context"
	xerror2 "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"
	"gorm.io/gen/field"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarketActivityListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarketActivityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketActivityListLogic {
	return &MarketActivityListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarketActivityListLogic) MarketActivityList(in *sol.GetMarketListRequest) (*sol.MarketActivityListResponse, error) {
	page := int(in.Page)
	size := int(in.Size)
	offset := (page - 1) * size
	orderBy := in.OrderBy
	direction := in.Direction
	marketAddress := in.MarketAddress
	if page < 0 {
		return &sol.MarketActivityListResponse{}, xerror2.NewErrCode(xerror2.ErrPageInvalid)
	}

	if size < 0 {
		return &sol.MarketActivityListResponse{}, xerror2.NewErrCode(xerror2.ErrSizeInvalid)
	}

	if size > 1000 {
		return &sol.MarketActivityListResponse{}, xerror2.NewErrCode(xerror2.ErrMaxSizeInvalid)
	}

	marketDB := dbx.Use(l.svcCtx.PgDB).MarketTx
	orderExpr, exist := marketDB.GetFieldByName(orderBy)
	if !exist {
		orderBy = marketDB.CreatedTimestamp.ColumnName().String()
	}
	var directionExpr field.Expr
	if direction == "desc" {
		directionExpr = orderExpr.Desc()
	} else if direction == "asc" {
		directionExpr = orderExpr.Asc()
	}
	count, err := marketDB.WithContext(l.ctx).ReadDB().
		Where(marketDB.MarketAddress.Eq(marketAddress)).
		Count()
	if err != nil {
		return nil, err // 处理错误
	}

	// 获取分页数据
	marketTxInfo, err := marketDB.WithContext(l.ctx).ReadDB().
		Where(marketDB.MarketAddress.Eq(marketAddress)).
		Order(directionExpr).
		Limit(size).Offset(offset).Find()
	if err != nil {
		return &sol.MarketActivityListResponse{}, xerror2.NewErrCodeMsg(xerror2.ErrDBQueryError, err.Error())
	}

	marketListResp := make([]*sol.Activity, 0)
	for _, marketTx := range marketTxInfo {
		marketListResp = append(marketListResp, &sol.Activity{
			Maker:        marketTx.MakerAddress,
			BaseAmount:   marketTx.BaseAmount,
			QuoteAmount:  marketTx.QuoteAmount,
			AmountUsd:    marketTx.Volume,
			Timestamp:    marketTx.CreatedTimestamp.Unix(),
			Type:         int64(marketTx.TxType),
			TxHash:       marketTx.TxHash,
			PriceUsd:     marketTx.QuotePrice,
			TokenAddress: marketTx.MarketAddress,
			QuoteAddress: marketTx.QuoteAddress,

			MakerTags:            nil,
			MakerTwitterName:     "",
			MakerTwitterUsername: "",
			MakerName:            "",
			MakerAvatar:          "",
			MakerEns:             "",
			MakerTokenTags:       nil,
			TotalTrade:           0,
			IsFollowing:          0,
			IsOpenOrClose:        0,
			BuyCostUsd:           0,
			Balance:              "",
			Cost:                 0,
			HistoryBoughtAmount:  0,
			HistorySoldIncome:    0,
			HistorySoldAmount:    0,
			UnrealizedProfit:     0,
			RealizedProfit:       0,
		})
	}

	return &sol.MarketActivityListResponse{
		List:  marketListResp,
		Total: count,
	}, nil
}
