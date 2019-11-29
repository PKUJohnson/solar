package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/micro/go-micro/metadata"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/protobuf/proto"
	"github.com/afex/hystrix-go/hystrix"
	std "github.com/PKUJohnson/solar/std"
	consts "github.com/PKUJohnson/solar/std/common"
	"github.com/PKUJohnson/solar/std/toolkit"
	stdnet "github.com/PKUJohnson/solar/std/toolkit/net"
	stdreflect "github.com/PKUJohnson/solar/std/toolkit/reflect"
	"golang.org/x/net/context"
	"gopkg.in/redis.v5"
)

var (
	defaultHystrixSetting hystrix.CommandConfig

	// CircuitBreakerMgrIns holds the global circuit breaker settings.
	CircuitBreakerMgrIns *CircuitBreakerManager
)

// FallbackFunc defines the API fallback, receives the response as parameter,
// and returns a hystrix's fallback.
type FallbackFunc func(req client.Request, rsp interface{}, err error) error

// CacheKeyFunc generates the cache key
type CacheKeyFunc func(req client.Request) string

// CircuitBreakerManager controls the hystrix status and settings.
type CircuitBreakerManager struct {
	enabled          bool
	redis            *redis.Client
	fallbackFuncs    map[string]FallbackFunc
	cachekeyFuncs    map[string]CacheKeyFunc
	cacheRefresher   chan keyValuePair
	cacheLastUpdated map[string]time.Time
}

type keyValuePair struct {
	key string
	val proto.Message
}

// CircuitSetting defines the hystrix settings for the specified command.
type CircuitSetting struct {
	Method                 interface{}
	Timeout                int
	MaxConcurrentRequests  int
	RequestVolumeThreshold int
	SleepWindow            int
	ErrorPercentThreshold  int
	Fallback               FallbackFunc
	CacheKey               CacheKeyFunc
}

type logger struct {
}

func (l logger) Printf(format string, items ...interface{}) {
	std.LogInfoc("hystrix", fmt.Sprintf(format, items...))
}

func init() {
	defaultHystrixSetting = hystrix.CommandConfig{
		Timeout:                3000,
		MaxConcurrentRequests:  1000,
		RequestVolumeThreshold: 20,
		SleepWindow:            5000,
		ErrorPercentThreshold:  80,
	}
	//hystrix.SetLogger(logger{})
}

func initCircuitBreaker(conf std.ConfigHystrix, advertiseSubnets []string) {
	CircuitBreakerMgrIns = &CircuitBreakerManager{
		enabled: conf.Enabled,
	}
	if conf.Enabled {
		std.LogInfoc("hystrix", "start to initialize settings")

		initDefaultSettings(&conf)

		serveMetrics(&conf, advertiseSubnets)

		redis := toolkit.InitRedis(conf.Redis)
		if redis == nil {
			std.LogErrorc("hystrix", nil, "fail to initialize redis")
			return
		}

		CircuitBreakerMgrIns.redis = redis
		CircuitBreakerMgrIns.cacheLastUpdated = make(map[string]time.Time)
		CircuitBreakerMgrIns.cacheRefresher = make(chan keyValuePair, 1000)
		go refreshCache(&conf)
	}
}

// refreshCache saves service snapshot to cache when they work fine
func refreshCache(conf *std.ConfigHystrix) {
	defer func() {
		if err := recover(); err != nil {
			std.LogErrorc("hystrix", fmt.Errorf("%v", err), "fail to save snapshot, stop the snapshot now")
		}
	}()

	cacheExpiration := time.Hour * 6
	cacheRefreshInterval := time.Second * 3
	if conf.CacheExpiration != 0 {
		cacheExpiration = time.Duration(conf.CacheExpiration) * time.Second
	}

	if conf.CacheRefreshInterval != 0 {
		cacheExpiration = time.Duration(conf.CacheRefreshInterval) * time.Second
	}

	for kv := range CircuitBreakerMgrIns.cacheRefresher {
		now := time.Now()
		if now.Sub(CircuitBreakerMgrIns.cacheLastUpdated[kv.key]) < cacheRefreshInterval {
			continue
		}
		encoded, err := proto.Marshal(kv.val)
		if err != nil {
			std.LogErrorc("hystrix", err, "fail to encode response")
			continue
		}
		CircuitBreakerMgrIns.redis.Set(kv.key, encoded, cacheExpiration)
		CircuitBreakerMgrIns.cacheLastUpdated[kv.key] = now
	}
}

func tryRefreshCache(key string, protoMsg proto.Message) {
	select {
	case CircuitBreakerMgrIns.cacheRefresher <- keyValuePair{key, protoMsg}:
	default:
		std.LogWarnc("hystrix", nil, "tryRefreshCache channel is full, discarding: "+key)
	}
}

// serveMetrics provides hystrix metrics stream for watching
func serveMetrics(conf *std.ConfigHystrix, advertiseSubnets []string) {
	if conf.ConsulAddrs == nil {
		std.LogErrorc("hystrix", errors.New("empty registry"), "fail to register")
		return
	}

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.DashboardPort), hystrixStreamHandler)

	// ip, err := stdnet.LocalIPAddr()
	// exclude the docker's bridge & docker_gwbridge & ingress ips
	ip, err := stdnet.LocalIPAddrWithin(advertiseSubnets)
	if err != nil {
		std.LogErrorc("hystrix", err, "fail to get IP address")
		return
	}

	r := consul.NewRegistry(registry.Addrs(conf.ConsulAddrs...))
	service := registry.Service{
		Name:    conf.ServiceName,
		Version: "v1",
		Nodes: []*registry.Node{
			{
				Id:      os.Getenv("HOSTNAME"),
				Address: fmt.Sprintf("%s:%d", ip, conf.DashboardPort),
			},
		},
	}
	err = r.Register(&service, registry.RegisterTTL(30*time.Second))
	if err != nil {
		std.LogErrorc("hystrix", err, "fail to register hystrix")
		return
	}

	go func(reg registry.Registry, svc *registry.Service) {
		t := time.NewTicker(20 * time.Second)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				reg.Register(svc, registry.RegisterTTL(30*time.Second))
			}
		}
	}(r, &service)
	std.LogInfoc("hystrix", "success to register hystrix")
}

func initDefaultSettings(conf *std.ConfigHystrix) {
	if conf.DefaultErrorPercentThreshold > 0 {
		hystrix.DefaultErrorPercentThreshold = conf.DefaultErrorPercentThreshold
	}

	if conf.DefaultMaxConcurrent > 0 {
		hystrix.DefaultMaxConcurrent = conf.DefaultMaxConcurrent
	}

	if conf.DefaultTimeout > 0 {
		hystrix.DefaultTimeout = conf.DefaultTimeout
	}

	if conf.DefaultSleepWindow > 0 {
		hystrix.DefaultSleepWindow = conf.DefaultSleepWindow
	}

	if conf.DefaultVolumeThreshold > 0 {
		hystrix.DefaultVolumeThreshold = conf.DefaultVolumeThreshold
	}
}

// CacheValue fetches the cache for the latest successful request.
func CacheValue(method interface{}, req client.Request, rsp interface{}) error {
	cacheKeyFunc, ok := CircuitBreakerMgrIns.cachekeyFuncs[stdreflect.PbMethodName(method)]
	if ok && CircuitBreakerMgrIns.redis != nil {
		cacheKey := cacheKeyFunc(req)
		encoded, err := CircuitBreakerMgrIns.redis.Get(cacheKey).Bytes()
		if err != nil {
			return err
		}

		protoMsg, ok := rsp.(proto.Message)
		if !ok {
			return errors.New("hystrix cache miss")
		}

		err = proto.Unmarshal(encoded, protoMsg)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("hystrix cache not open")
}

// ConfigCircuitBreaker configures the circuit breaker's settings for multiple commands
func ConfigCircuitBreaker(settings []CircuitSetting) {
	if CircuitBreakerMgrIns == nil {
		std.LogErrorc("hystrix", nil, "forget to call initCircuitBreaker?")
		return
	}

	fallbackFuncs := make(map[string]FallbackFunc)
	cachekeyFuncs := make(map[string]CacheKeyFunc)
	for _, s := range settings {
		if name, hs := s.toHystrixSetting(); name == "" {
			std.LogErrorc("hystrix", errors.New("setting is wrong"), "fail to init hystrix settings")
		} else {
			if s.Fallback != nil {
				fallbackFuncs[name] = s.Fallback
			}
			if s.CacheKey != nil {
				cachekeyFuncs[name] = s.CacheKey
			}
			hystrix.ConfigureCommand(name, *hs)
			std.LogInfoc("hystrix", "success to register command: "+name)
		}
	}
	CircuitBreakerMgrIns.fallbackFuncs = fallbackFuncs
	CircuitBreakerMgrIns.cachekeyFuncs = cachekeyFuncs
}

func (c *CircuitSetting) toHystrixSetting() (string, *hystrix.CommandConfig) {
	if c == nil || c.Method == nil {
		return "", nil
	}
	name := stdreflect.PbMethodName(c.Method)
	if name == "" {
		return "", nil
	}
	ret := defaultHystrixSetting
	if c.Timeout > 0 {
		ret.Timeout = c.Timeout
	}
	if c.MaxConcurrentRequests > 0 {
		ret.MaxConcurrentRequests = c.MaxConcurrentRequests
	}
	if c.RequestVolumeThreshold > 0 {
		ret.RequestVolumeThreshold = c.RequestVolumeThreshold
	}
	if c.SleepWindow > 0 {
		ret.SleepWindow = c.SleepWindow
	}
	if c.ErrorPercentThreshold > 0 {
		ret.ErrorPercentThreshold = c.ErrorPercentThreshold
	}
	return name, &ret
}

type cbWrapper struct {
	client.Client
}

func (c *cbWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// close hystrix when disabled or API is admin
	hystrixStatus := consts.Md_HYSTRIX_OPEN
	md, ok1 := metadata.FromContext(ctx)
	if ok1 {
		hystrixStatus2, ok2 := md[consts.Md_HYSTRIX_STATUS]
		if ok2 {
			hystrixStatus = hystrixStatus2
		}
	}

	if !CircuitBreakerMgrIns.enabled || hystrixStatus == consts.Md_HYSTRIX_CLOSED {
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	var lastErr error
	callFunc := func() error {
		fmt.Println("enter call function 1: ", reflect.TypeOf(c.Client))
		err := c.Client.Call(ctx, req, rsp, opts...)
		if err != nil {
			fmt.Println("enter call function error: ", err)
		}
		lastErr = err
		if err != nil {
			std.LogWarnc("hystrix", err, "original error")
		}
		return nil
	}

	fallbackFunc := func(err error) error {
		if err == hystrix.ErrTimeout || err == hystrix.ErrCircuitOpen || err == hystrix.ErrMaxConcurrency {
			std.LogErrorc("hystrix", nil, fmt.Sprintf("fail to call %s.%s: %s", req.Service(), req.Method(), err.Error()))
			lastErr = std.ErrServerTooBusy
			return std.ErrServerTooBusy
		}
		lastErr = err
		return err
	}

	id := req.Method()
	cb, ok := CircuitBreakerMgrIns.fallbackFuncs[id]
	if !ok {
		hystrix.Do(id, callFunc, fallbackFunc)
		return lastErr
	}

	hystrix.Do(id, func() error {
		cacheKeyFunc, ok := CircuitBreakerMgrIns.cachekeyFuncs[id]
		if ok && CircuitBreakerMgrIns.redis != nil {
			cacheKey := cacheKeyFunc(req)
			err := c.Client.Call(ctx, req, rsp, opts...)
			if err == nil {
				// start to cache when backend service works fine
				protoMsg, ok := rsp.(proto.Message)
				if !ok {
					return nil
				}
				tryRefreshCache(cacheKey, protoMsg)
			}
			lastErr = err
			return nil
		}

		err := c.Client.Call(ctx, req, rsp, opts...)
		lastErr = err
		if err != nil {
			std.LogWarnc("hystrix", err, "original error")
		}
		return nil
	}, func(err error) error {
		if err == hystrix.ErrTimeout || err == hystrix.ErrCircuitOpen || err == hystrix.ErrMaxConcurrency {
			std.LogErrorc("hystrix", nil, fmt.Sprintf("fail to call %s.%s: %s", req.Service(), req.Method(), err.Error()))
			lastErr = std.ErrServerTooBusy
			return std.ErrServerTooBusy
		}

		lastErr = cb(req, rsp, err)
		return lastErr
	})
	return lastErr
}

// NewClientWrapper returns a hystrix client Wrapper.
func newCircuitBreakerWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &cbWrapper{c}
	}
}
