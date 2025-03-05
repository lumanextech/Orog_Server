package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTokenFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddTokenFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTokenFollowLogic {
	return &AddTokenFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddTokenFollowLogic) AddTokenFollow(in *account.AddTokenFollowRequest) (*account.AddTokenFollowResponse, error) {

	// 检查数据库中是否已存在相同的地址、链和代币地址
	existing, err := l.svcCtx.UserTokenFollowModel.FindOneByAddressAndTokenAddress(l.ctx, in.Address, in.TokenAddress)
	if err == nil && existing.Chain.String == in.Chain && existing.TokenAddress.String == in.TokenAddress {
		// 如果数据库中已存在相同的地址、链和代币地址，检查其状态
		if existing.Status.String != "1" {
			existing.Status = sql.NullString{String: "1", Valid: true}
			err = l.svcCtx.UserTokenFollowModel.Update(l.ctx, existing)
			if err != nil {
				return nil, err
			}
		}
		return &account.AddTokenFollowResponse{}, nil
	}

	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	// 初始化雪花节点
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}

	// 生成雪花 ID
	snowflakeID := node.Generate().Int64()
	// 如果数据库中不存在相同的地址、链和代币地址，则从请求中提取数据并插入新代币关注条目
	userTokenFollow := &model.UserTokenFollow{
		Id:           snowflakeID, //雪花id
		Address:      in.Address,
		Chain:        sql.NullString{String: in.Chain, Valid: in.Chain != ""},
		TokenAddress: sql.NullString{String: in.TokenAddress, Valid: in.TokenAddress != ""},
		Status:       sql.NullString{Valid: true, String: "1"},
		FollowedAt:   sql.NullTime{Time: time.Now(), Valid: true},
	}

	// 插入新的代币关注条目
	_, err = l.svcCtx.UserTokenFollowModel.InsertByAddressAndStatus(l.ctx, userTokenFollow)
	if err != nil {
		return nil, err
	}

	// 返回成功响应
	return &account.AddTokenFollowResponse{
		Message: "success",
	}, nil
}
