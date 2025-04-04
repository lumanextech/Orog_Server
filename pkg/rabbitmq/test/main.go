package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/simance-ai/smdx/pkg/rabbitmq"
)

func main() {
	// 配置
	conf := rabbitmq.RabbitMQConf{
		Host:       "localhost",
		Port:       5672,
		Username:   "guest",
		Password:   "guest",
		Vhost:      "/",
		Queue:      "test_queue",
		Exchange:   "test_exchange",
		RoutingKey: "test_key",
		Consumers:  1,
		Processors: 1,
	}

	// 创建客户端
	client, err := rabbitmq.NewClient(conf, nil) // 不设置 handler，这样消息就不会被自动消费
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Stop()

	// 发送测试消息
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("测试消息 %d", i)
		err := client.Publish(context.Background(), "test_key", msg)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
		} else {
			log.Printf("发送消息成功: %s", msg)
		}
		time.Sleep(time.Second)
	}

	log.Println("消息已发送，请在 RabbitMQ 管理界面查看消息。按 Ctrl+C 退出程序。")
	select {} // 保持程序运行
}
