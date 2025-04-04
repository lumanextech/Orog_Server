package rabbitmq

import "github.com/zeromicro/go-zero/core/service"

type RabbitMQConf struct {
	service.ServiceConf
	Host       string
	Port       int
	Username   string
	Password   string
	Vhost      string
	Queue      string
	Exchange   string
	RoutingKey string
	Consumers  int `json:",default=8"`
	Processors int `json:",default=8"`
}
