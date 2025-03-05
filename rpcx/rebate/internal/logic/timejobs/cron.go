package timejob

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/simance-ai/smdx/rpcx/rebate/internal/logic"
	"github.com/simance-ai/smdx/rpcx/rebate/internal/svc"
	"log"
	"os"
	"sync"
)

var once sync.Once

func Run(ctx *svc.ServiceContext) (cronInstance *cron.Cron, err error) {

	once.Do(func() {
		logg := log.New(os.Stdout, "cron: ", log.LstdFlags)
		logger := cron.VerbosePrintfLogger(logg)
		cronInstance = cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.SkipIfStillRunning(logger),
				cron.Recover(logger)))

		err = registerJob(cronInstance, ctx)
		if err != nil {
			log.Fatalln(err)
		}

		cronInstance.Start()

	})

	return
}

func registerJob(cronInstance *cron.Cron, ctx *svc.ServiceContext) (err error) {

	rebateLogic := logic.NewRebateLogic(context.Background(), ctx)
	// 返佣 "0 0 0 * * *"
	_, err = cronInstance.AddFunc(ctx.Config.Corn.Rebate, rebateLogic.Rebate) //每天执行一次
	if err != nil {
		log.Printf("❌ 添加定时任务失败: %v", err)
		return err
	}

	return
}
