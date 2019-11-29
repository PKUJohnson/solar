package std

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	rate "github.com/anacrolix/rate/redis"
	"go.etcd.io/etcd/clientv3"
	"github.com/jinzhu/configor"
	"gopkg.in/yaml.v2"
)

const (
	ConfEnvDevelopment = "development"
	ConfEnvTest        = "ivktest"
	ConfEnvProduction  = "ivkprod"
	ConfEnvStage       = "ivkstage"
)

// ConfigService sets the service's configurations
type ConfigService struct {
	SvcName          string   `yaml:"svc_name"`
	SvcAddr          string   `yaml:"svc_addr"`
	AdvertiseSubnets []string `yaml:"advertise_subnets"`
	EtcdAddrs        []string `yaml:"etcd_addrs"`
	Consul           []struct {
		Addr string `yaml:"addr"`
	} `yaml:"consul"`
	ZookeeperAddrs []string            `yaml:"zookeeper_addrs"`
	RpcCacheRedis  ConfigRpcCacheRedis `yaml:"rpc_cache_redis"`
	Redis          ConfigRedis         `yaml:"redis"`
	Prometheus     ConfigPrometheus    `yaml:"prometheus"`
	Hystrix        ConfigHystrix       `yaml:"hystrix"`
	Zipkin         ConfigZipkin        `yaml:"zipkin"`
	RateLimiter    ConfigRateLimiter   `yaml:"rate_limiter"`
	RPCServer      string              `yaml:"rpc_server"`
}

// ConfigRateLimiter defines the rate limiter
type ConfigRateLimiter struct {
	// Enabled controls the rate limiter's switch
	Enabled bool `yaml:"enabled"`

	// Redis indicates the rate limiter's redis.
	Redis rate.ConfigRedis `yaml:"redis"`

	// HTTP message when limit is reached.
	Message string `yaml:"message"`

	// Content-Type for Message
	MessageContentType string `yaml:"message_content_type"`

	// HTTP status code when limit is reached.
	StatusCode int `yaml:"status_code"`

	// Maximum number of requests to limit per duration.
	Max int64 `yaml:"max"`

	// Duration of rate-limiter.
	TTL time.Duration `yaml:"ttl"`

	// List of places to look up IP address.
	// Default is "RemoteAddr", "X-Forwarded-For", "X-Real-IP".
	// You can rearrange the order as you like.
	IPLookups []string `yaml:"ip_lookups"`

	// List of HTTP Methods to limit (GET, POST, PUT, etc.).
	// Empty means limit all methods.
	Methods []string `yaml:"methods"`

	// List of HTTP headers to limit.
	// Empty means skip headers checking.
	Headers []string `yaml:"headers"`

	APISettings []APIRateLimit `yaml:"api_settings"`
}

// APIRateLimit defines the specified API's rate limit.
type APIRateLimit struct {
	Path     string `yaml:"path"`
	Method   string `yaml:"method"`
	Max      int64  `yaml:"max"`
	Duration int64  `yaml:"duration"`
}

// ConfigZipkin configures the rpc call tracing system
type ConfigZipkin struct {
	Enabled     bool     `yaml:"enabled"`
	BrokerAddrs []string `yaml:"broker_addrs"`
}

// ConfigLog sets the logger level and destination
type ConfigLog struct {
	Level            int      `yaml:"level"`
	Path             string   `yaml:"path"`
	FileName         string   `yaml:"file_name"`
	RotationDuration string   `yaml:"rotation_duration"`
	RotationCount    uint      `yaml:"rotation_count"`
	OutputDest       string   `yaml:"output_dest"`
	CIDRs            []string `yaml:"cidrs"`
	SentryDSN        string   `yaml:"sentry_dsn"`
}

// ConfigNSQ sets the NSQ address
type ConfigNSQ struct {
	LookupAddress []string `yaml:"lookup_address_list"`
	NsqdAddress   string   `yaml:"addr"`
}

// ConfigMysql sets the MySQL
type ConfigMysql struct {
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	DBName         string `yaml:"db_name"`
	MaxIdle        int    `yaml:"max_idle"`
	MaxConn        int    `yaml:"max_conn"`
	LogType        string `yaml:"log_type"`
	NotPrintSql    bool   `yaml:"not_print_sql"`
	NotCreateTable bool   `yaml:"not_create_table"`
}

type ConfigUpload struct {
	Path    string `yaml:"path"`
	Url     string `yaml:"url"`
	CdnUrl  string `yaml:"cdn_url"`
	FakeCDN bool   `yaml:"fake_cdn"`
}

type ConfigPicService struct {
	ThumbUrl      string `yaml:"thumb_url"`
}

type ConfigPubnoService struct {
	BaseUrl      string `yaml:"base_url"`
}

type ConfigFormulaService struct {
	BaseUrl      string `yaml:"base_url"`
}

type ConfigWebPageShare struct {
	BaseUrl      string `yaml:"base_url"`
}

type ConfigAuthCode struct {
	Url           string `yaml:"url"`
	Expiration    int64 `yaml:"expiration"`
	SendInterval  int64 `yaml:"send_interval"`
}

// ConfigRedis sets the Redis
type ConfigRedis struct {
	Enabled     bool   `yaml:"enabled"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Auth        string `yaml:"auth"`
	IdleTimeout int    `yaml:"idle_timeout"`
}

// ConfigRpcCacheRedis sets the RPC cache backend by Redis
type ConfigRpcCacheRedis struct {
	Enabled     bool              `yaml:"enabled"`
	Addrs       map[string]string `yaml:"addrs"`
	Password    string            `yaml:"password"`
	IdleTimeout int               `yaml:"idle_timeout"`
}

// ConfigHystrix sets RPC client error handler and circuit breaker
type ConfigHystrix struct {
	Enabled              bool        `yaml:"enabled"`
	DashboardPort        int         `yaml:"dashboard_port"`
	ServiceName          string      `yaml:"service_name"`
	ConsulAddrs          []string    `yaml:"consul_addrs"`
	Redis                ConfigRedis `yaml:"redis"`
	CacheExpiration      int         `yaml:"cache_expiration"`
	CacheRefreshInterval int         `yaml:"cache_refresh_interval"`

	DefaultTimeout               int `yaml:"default_timeout"`
	DefaultMaxConcurrent         int `yaml:"default_max_concurrent"`
	DefaultErrorPercentThreshold int `yaml:"default_error_percent_threshold"`
	DefaultSleepWindow           int `yaml:"default_sleep_window"`
	DefaultVolumeThreshold       int `yaml:"default_volume_threshold"`
}

// ConfigPrometheus sets server metrics collectors
type ConfigPrometheus struct {
	Enabled       bool   `yaml:"enabled"`
	Namespace     string `yaml:"namespace"`
	BatchInterval int    `yaml:"batch_interval"`
	Collectors    []struct {
		Addr string `yaml:"addr"`
	} `yaml:"collectors"`
}

// ConfigCassandra sets the Cassandra
type ConfigCassandra struct {
	Hosts    []string `yaml:"hosts"`
	Keyspace string   `yaml:"keyspace"`
	Port     string   `yaml:"port"`
}

// config nsq consumer
type ConfigNsqConsumer struct {
	TopicName    string   `yaml:"topic_name"`
	ChannelName  string   `yaml:"channel_name"`
	IsLookupAddr bool     `yaml:"is_lookup_addr"`
	AddressList  []string `yaml:"address_list"`
}

// config wechat pa account
type ConfigWechatAccount struct {
	AppId       string   `yaml:"app_id"`
	AppSecret   string   `yaml:"app_secret"`
	Token       string   `yaml:"token"`
	AesKey      string   `yaml:"aes_key"`
}

// config wechat mp account
type ConfigMiniProgramAccount struct {
	AppId       string   `yaml:"mp_id"`
	AppSecret   string   `yaml:"mp_secret"`
}

// config wechat test token account
type ConfigTestToken struct {
	AppId       string   `yaml:"app_id"`
}

// ConfEnv fetches the current runtime environment
func ConfEnv() string {
	return configor.ENV()
}

// LoadConf loads configuration from specified file
func LoadConf(dest interface{}, path string) {
	if os.Getenv("CONFIG_ETCDS") != "" {
		loadFromEtcd(dest, path)
		return
	}

	err := configor.Load(dest, path)
	if err != nil {
		msg := "Failed to load config file !!!"
		LogError(LogFields{"error": err}, msg)
		panic(msg)
	}
}

func loadFromEtcd(dest interface{}, path string) {
	etcdList := os.Getenv("CONFIG_ETCDS")
	if etcdList == "" {
		LogPanic(LogFields{
			TagCategory: "config",
			"error":     "etcd list is empty",
		}, "fail to load configuration")
	}

	path = getConfigurationWithENV(path, ConfEnv())
	LogInfo(LogFields{
		TagCategory: "config",
	}, "start to load configuration "+path)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdList, ","),
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		LogPanic(LogFields{
			TagCategory: "config",
			"error":     err,
		}, "fail to load configuration")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	rsp, err := cli.Get(ctx, path)
	cancel()
	if err != nil {
		LogPanic(LogFields{
			TagCategory: "config",
			"error":     err,
		}, "fail to load configuration")
	}

	err = load(dest, path, rsp.Kvs[0].Value)
	if err != nil {
		LogPanic(LogFields{
			TagCategory: "config",
			"error":     err,
		}, "fail to load configuration")
	}
	//watchEtcdKey(strings.Split(etcdList, ","), dest, path)
}

func watchEtcdKey(etcdlist []string, dest interface{}, keyName string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdlist,
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		LogPanic(LogFields{
			TagCategory: "config",
			"error":     err,
		}, "fail to watch "+keyName)
	}

	watcher := cli.Watch(context.Background(), keyName)
	go func() {
		for result := range watcher {
			if len(result.Events) != 0 {
				err := load(dest, keyName, result.Events[0].Kv.Value)
				if err != nil {
					LogErrorc("config", err, "update config failed")
				} else {
					LogInfoc("config", "update config success")
				}
			}
		}
	}()
}

func load(config interface{}, file string, data []byte) (err error) {
	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		return yaml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".toml"):
		return toml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".json"):
		return json.Unmarshal(data, config)
	default:
		if toml.Unmarshal(data, config) != nil {
			if json.Unmarshal(data, config) != nil {
				if yaml.Unmarshal(data, config) != nil {
					return errors.New("failed to decode config")
				}
			}
		}
		return nil
	}
}

func getConfigurationWithENV(file, env string) string {
	var envFile string
	var extname = path.Ext(file)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}
	return filepath.Base(envFile)
}
