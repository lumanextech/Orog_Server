package model_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/simance-ai/smdx/rpcx/account/internal/config"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestDefaultUserTokenFollowModel_Insert(t *testing.T) {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// 使用 filepath.Join 拼接路径
	configFile := filepath.Join(dir, "../../", "etc", "account.yaml")

	var c config.Config
	conf.MustLoad(configFile, &c)
	ctx := svc.NewServiceContext(c)

	result, err := ctx.UserTokenFollowModel.Insert(context.Background(), &model.UserTokenFollow{
		Address: "EJZ3LVmtuBVnfUbshd6YHHaqKtMjm7EWC3YgMZhmMoMW",
		Chain: sql.NullString{
			String: "sol",
			Valid:  true,
		},
		TokenAddress: sql.NullString{
			String: "test",
			Valid:  true,
		},
		Status: sql.NullString{
			String: "test",
			Valid:  true,
		},
		FollowedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UnfollowedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Id: 99999,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestDefaultUserTokenFollowModel_FindAllByAddressAndStatus(t *testing.T) {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// 使用 filepath.Join 拼接路径
	configFile := filepath.Join(dir, "../../", "etc", "account.yaml")

	var c config.Config
	conf.MustLoad(configFile, &c)
	ctx := svc.NewServiceContext(c)

	result, err := ctx.UserTokenFollowModel.FindAllByAddressAndStatus(context.Background(),
		"EJZ3LVmtuBVnfUbshd6YHHaqKtMjm7EWC3YgMZhmMoMW", "1")
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range result {
		t.Log(i)
	}
}
