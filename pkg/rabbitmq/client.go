package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type ConsumeHandler interface {
	Consume(ctx context.Context, key, value string) error
}

type RabbitMQClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	conf    RabbitMQConf
	handler ConsumeHandler
}

// NewClient 创建一个新的 RabbitMQ 客户端
func NewClient(conf RabbitMQConf, handler ConsumeHandler) (*RabbitMQClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		conf.Username,	
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Vhost)

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明队列
	_, err = ch.QueueDeclare(
		conf.Queue,
		true,  // durable 持久化true：队列会持久化到磁盘，false：队列只在内存中，服务器重启后会消失
		false, // delete when unused 自动删除，true：当最后一个消费者断开连接后，队列会被自动删除 false：队列会一直存在，除非手动删除
		false, // exclusive 排他性，true：队列只能被一个消费者使用，false：队列可以被多个消费者使用
		false, // no-wait true：不等待服务器确认就返回，false：等待服务器确认后再返回
		nil,   // arguments 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		conf.Exchange,
		"fanout",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 绑定队列到交换机
	err = ch.QueueBind(
		conf.Queue,
		conf.RoutingKey,
		conf.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: ch,
		conf:    conf,
		handler: handler,
	}, nil
}

func (c *RabbitMQClient) Start() {
	for i := 0; i < c.conf.Consumers; i++ {
		go c.consume()
	}
}

func (c *RabbitMQClient) Stop() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *RabbitMQClient) consume() {
	msgs, err := c.channel.Consume(
		c.conf.Queue,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	for msg := range msgs {
		ctx := context.Background()
		err := c.handler.Consume(ctx, msg.RoutingKey, string(msg.Body))
		if err != nil {
			log.Printf("Failed to process message: %v", err)
			msg.Nack(false, true) // 重新入队
		} else {
			msg.Ack(false)
		}
	}
}

func (c *RabbitMQClient) Publish(ctx context.Context, routingKey string, body interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.channel.PublishWithContext(
		ctx,
		c.conf.Exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
			Timestamp:   time.Now(),
		},
	)
}
