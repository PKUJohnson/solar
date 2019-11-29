package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/PKUJohnson/solar/std/ec"
)

// Config wraps the sarama kafka configuration
func Config(c *sarama.Config) ec.Option {
	return func(o *ec.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, ec.KafkaConfig, c)
	}
}
