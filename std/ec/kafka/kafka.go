package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/ec"
)

// Collector wraps a kafka collector.
type Collector struct {
	Client sarama.AsyncProducer
	Topic  string
}

func defaultConfig() *sarama.Config {
	return sarama.NewConfig()
}

func useConfig(cfg *ec.Config) ec.Option {
	return func(o *ec.Options) {
		o.Topic = cfg.Topic
		o.Addrs = cfg.Addrs

		c := defaultConfig()
		switch cfg.Compression {
		case "snappy":
			c.Producer.Compression = sarama.CompressionSnappy
		case "lz4":
			c.Producer.Compression = sarama.CompressionLZ4
		case "gzip":
			c.Producer.Compression = sarama.CompressionGZIP
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, ec.KafkaConfig, c)

		o.Context = ctx
	}
}

// NewCollectorWithConfig contructs a kafka event collector with `Config`.
func NewCollectorWithConfig(cfg *ec.Config) (ec.EventCollector, error) {
	return NewCollector(useConfig(cfg))
}

// NewCollector creates a kafka collector.
func NewCollector(opts ...ec.Option) (*Collector, error) {
	var options ec.Options
	for _, o := range opts {
		o(&options)
	}

	config := defaultConfig()
	if options.Context != nil {
		if c, ok := options.Context.Value(ec.KafkaConfig).(*sarama.Config); ok {
			config = c
		}
	}

	if options.Topic == "" {
		panic("kafka topic must not be empty")
	}

	asyncp, err := sarama.NewAsyncProducer(options.Addrs, config)
	if err != nil {
		std.LogErrorc("kafka", err, "fail to create data collector with kafka")
		return nil, err
	}

	collector := &Collector{
		Client: asyncp,
		Topic:  options.Topic,
	}

	go run(asyncp)

	return collector, nil
}

func run(asyncp sarama.AsyncProducer) {
	for {
		select {
		case err := <-asyncp.Errors():
			std.LogErrorc("kafka", err, "fail to record event")
		}
	}
}

// CollectByTopic send events to kafka by the specified topic.
func (c *Collector) CollectByTopic(topic string, e *ec.Event) error {
	if e == nil || topic == "" {
		return nil
	}

	select {
	case c.Client.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(e.Body)}:
		std.LogDebugc("kafka", "message sent")
	}

	return nil
}

// Collect send events to kafka.
func (c *Collector) Collect(e *ec.Event) error {
	if e == nil {
		return nil
	}

	select {
	case c.Client.Input() <- &sarama.ProducerMessage{Topic: c.Topic, Key: nil, Value: sarama.StringEncoder(e.Body)}:
		std.LogDebugc("kafka", "message sent")
	}

	return nil
}

// Close the kafka collector.
func (c *Collector) Close() error {
	return c.Client.Close()
}
