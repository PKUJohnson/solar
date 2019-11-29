package ec

import "context"

// Options configures the event collector.
type Options struct {
	// Addrs stores the collector's addresses, like kafka server `127.0.0.1:9092`.
	Addrs []string

	// Topic represents the collector's topic to receive event, like kafka server, optional.
	Topic string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option defines the option setting function.
type Option func(*Options)

// Addrs is the host:port pairs to use.
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Topic is the topic to use.
func Topic(topic string) Option {
	return func(o *Options) {
		o.Topic = topic
	}
}

// ConfigType represents the event collector's type.
type ConfigType int

const (
	// KafkaConfig is kafka configuration.
	KafkaConfig ConfigType = iota
)

// Config sets event collector.
type Config struct {
	Enabled  bool     `yaml:"enabled" json:"enabled"`
	Type     string   `yaml:"type" json:"type"`
	Addrs    []string `yaml:"addrs" json:"addrs"`
	Topic    string   `yaml:"topic" json:"topic"`
	Encoding string   `yaml:"encoding" json:"encoding"`

	Compression string `yaml:"compression" json:"compression"`
}
