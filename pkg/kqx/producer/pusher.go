package producer

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/simance-ai/smdx/pkg/kqx"
	internal2 "github.com/simance-ai/smdx/pkg/kqx/internal"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	saslk "github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"

	"context"
)

type (
	PushOption func(options *pushOptions)

	Pusher struct {
		topic    string
		producer kafkaWriter
		executor *executors.ChunkExecutor
	}

	kafkaWriter interface {
		Close() error
		WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	}

	pushOptions struct {
		// kafka.Writer options
		allowAutoTopicCreation bool
		balancer               kafka.Balancer

		// executors.ChunkExecutor options
		chunkSize     int
		flushInterval time.Duration

		// syncPush is used to enable sync push
		syncPush bool
	}
)

// NewPusher returns a Pusher with the given Kafka addresses and topic.
func NewPusher(c kqx.KqConf, topic string, opts ...PushOption) *Pusher {

	//producer := &kafka.Writer{
	//	Addr:        kafka.TCP(addrs...),
	//	Topic:       topic,
	//	Balancer:    &kafka.LeastBytes{},
	//	Compression: kafka.Snappy,
	//}
	kafka.NewWriter(kafka.WriterConfig{
		Brokers:          c.Brokers,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		CompressionCodec: kafka.Snappy.Codec(),
	})

	writeConfig := kafka.WriterConfig{
		Brokers: c.Brokers,
		Topic:   c.Topic,
	}

	var err error
	var mechanism saslk.Mechanism
	switch c.Protocol {
	case "plain":
		mechanism = plain.Mechanism{
			Username: c.Username,
			Password: c.Password,
		}
	case "scram-sha-256":
		mechanism, err = scram.Mechanism(scram.SHA256, c.Username, c.Password)
	case "scram-sha-512":
		mechanism, err = scram.Mechanism(scram.SHA512, c.Username, c.Password)
	default:
		log.Fatalf("unknown sasl protocol: %s", c.Protocol)
	}
	if err != nil {
		log.Fatalf("failed to create sasl mechanism: %v", err)
	}

	if len(c.Username) > 0 && len(c.Password) > 0 {
		writeConfig.Dialer = &kafka.Dialer{
			SASLMechanism: mechanism,
		}
	}
	if len(c.CaFile) > 0 {
		caCert, err := os.ReadFile(c.CaFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			log.Fatal(err)
		}

		writeConfig.Dialer.TLS = &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}
	}
	producer := kafka.NewWriter(writeConfig)

	var options pushOptions
	for _, opt := range opts {
		opt(&options)
	}

	// apply kafka.Writer options
	producer.AllowAutoTopicCreation = options.allowAutoTopicCreation
	if options.balancer != nil {
		producer.Balancer = options.balancer
	}

	pusher := &Pusher{
		producer: producer,
		topic:    topic,
	}

	// if syncPush is true, return the pusher directly
	if options.syncPush {
		return pusher
	}

	// apply ChunkExecutor options
	var chunkOpts []executors.ChunkOption
	if options.chunkSize > 0 {
		chunkOpts = append(chunkOpts, executors.WithChunkBytes(options.chunkSize))
	}
	if options.flushInterval > 0 {
		chunkOpts = append(chunkOpts, executors.WithFlushInterval(options.flushInterval))
	}

	pusher.executor = executors.NewChunkExecutor(func(tasks []interface{}) {
		chunk := make([]kafka.Message, len(tasks))
		for i := range tasks {
			chunk[i] = tasks[i].(kafka.Message)
		}
		if err := pusher.producer.WriteMessages(context.Background(), chunk...); err != nil {
			logx.Error(err)
		}
	}, chunkOpts...)

	return pusher
}

// Close closes the Pusher and releases any resources used by it.
func (p *Pusher) Close() error {
	if p.executor != nil {
		p.executor.Flush()
	}

	return p.producer.Close()
}

// Name returns the name of the Kafka topic that the Pusher is sending messages to.
func (p *Pusher) Name() string {
	return p.topic
}

// KPush sends a message to the Kafka topic.
func (p *Pusher) KPush(ctx context.Context, k, v string) error {
	msg := kafka.Message{
		Key:   []byte(k), // current timestamp
		Value: []byte(v),
	}
	if p.executor != nil {
		return p.executor.Add(msg, len(v))
	} else {
		return p.producer.WriteMessages(ctx, msg)
	}
}

// Push sends a message to the Kafka topic.
func (p *Pusher) Push(ctx context.Context, v string) error {
	return p.PushWithKey(ctx, strconv.FormatInt(time.Now().UnixNano(), 10), v)
}

// PushWithKey sends a message with the given key to the Kafka topic.
func (p *Pusher) PushWithKey(ctx context.Context, key, v string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(v),
	}

	// wrap message into message carrier
	mc := internal2.NewMessageCarrier(internal2.NewMessage(&msg))
	// inject trace context into message
	otel.GetTextMapPropagator().Inject(ctx, mc)

	if p.executor != nil {
		return p.executor.Add(msg, len(v))
	} else {
		return p.producer.WriteMessages(ctx, msg)
	}
}

// WithAllowAutoTopicCreation allows the Pusher to create the given topic if it does not exist.
func WithAllowAutoTopicCreation() PushOption {
	return func(options *pushOptions) {
		options.allowAutoTopicCreation = true
	}
}

// WithBalancer customizes the Pusher with the given balancer.
func WithBalancer(balancer kafka.Balancer) PushOption {
	return func(options *pushOptions) {
		options.balancer = balancer
	}
}

// WithChunkSize customizes the Pusher with the given chunk size.
func WithChunkSize(chunkSize int) PushOption {
	return func(options *pushOptions) {
		options.chunkSize = chunkSize
	}
}

// WithFlushInterval customizes the Pusher with the given flush interval.
func WithFlushInterval(interval time.Duration) PushOption {
	return func(options *pushOptions) {
		options.flushInterval = interval
	}
}

// WithSyncPush enables the Pusher to push messages synchronously.
func WithSyncPush() PushOption {
	return func(options *pushOptions) {
		options.syncPush = true
	}
}
