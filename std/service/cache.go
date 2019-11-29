package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/server"
	std "github.com/PKUJohnson/solar/std"
	stdc "github.com/PKUJohnson/solar/std/common"
	ivkreflect "github.com/PKUJohnson/solar/std/toolkit/reflect"
	"golang.org/x/net/context"
	"gopkg.in/go-redis/cache.v5"
	"gopkg.in/redis.v5"
	"github.com/PKUJohnson/solar/std/toolkit"
)

var (
	CacheMgrIns *CacheManager
	RedisClient *redis.Client
)

type CacheSetting struct {
	Method     interface{}
	Expiration time.Duration
	AnonOnly   bool
}

type CacheManager struct {
	enabled bool
	codec   *cache.Codec
	config  map[string]CacheSetting
}

type cachedObj struct {
	Data   []byte
	IvkErr *std.Err
}

func initRedisClient(conf std.ConfigRedis) {
	if conf.Enabled {
		RedisClient = toolkit.InitRedis(conf)
	} else {
		RedisClient = nil
	}
}

func initCache(conf std.ConfigRpcCacheRedis) {

	if !conf.Enabled {
		CacheMgrIns = &CacheManager{}
	} else {
		ring := redis.NewRing(&redis.RingOptions{
			Addrs:       conf.Addrs,
			Password:    conf.Password,
			IdleTimeout: time.Second * time.Duration(conf.IdleTimeout),
		})
		marshal := func(v interface{}) ([]byte, error) {
			if _, ok := v.(*cachedObj); !ok {
				return nil, std.ErrRpcCacheMarshal
			} else if b, err := json.Marshal(v); err != nil {
				return nil, std.ErrRpcCacheMarshal
			} else {
				return b, nil
			}
		}
		unmarshal := func(b []byte, v interface{}) error {
			if _, ok := v.(*cachedObj); !ok {
				return std.ErrRpcCacheUnmarshal
			} else if err := json.Unmarshal(b, v); err != nil {
				return std.ErrRpcCacheUnmarshal
			} else {
				return nil
			}
		}
		codec := &cache.Codec{
			Redis:     ring,
			Marshal:   marshal,
			Unmarshal: unmarshal,
		}
		CacheMgrIns = &CacheManager{
			enabled: true,
			codec:   codec,
		}
	}
}

func ConfigRpcCache(settings []CacheSetting) {
	if CacheMgrIns == nil {
		std.LogErrorLn("Call initCache first to initialize")
		return
	}

	conf := make(map[string]CacheSetting)
	for _, setting := range settings {
		if n := ivkreflect.PbMethodName(setting.Method); n != "" {
			conf[n] = setting
		}
	}

	CacheMgrIns.config = conf
}

func CacheWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		if !CacheMgrIns.enabled {
			return fn(ctx, req, rsp)
		}
		settings, ok := CacheMgrIns.config[req.Method()]
		if !ok || settings.Expiration <= 0 {
			return fn(ctx, req, rsp)
		}
		if settings.AnonOnly {
			userID := stdc.PbGetUser(ctx)
			if userID > 0 {
				return fn(ctx, req, rsp)
			}
		}
		if settings.Expiration <= 0 {
			return fn(ctx, req, rsp)
		}
		if key, err := getCacheKey(ctx, req.Method(), req); err != nil {
			// NOTE: we probably should NOT go ahead and call the logic
			// because that could cause unbearable traffic
			std.LogErrorc("redis", err, "fail to fetch rpc content from cache")
			return err
		} else {
			v, err := CacheMgrIns.codec.Do(&cache.Item{
				Key:        key,
				Object:     new(cachedObj), // destination
				Expiration: settings.Expiration,
				Func: func() (interface{}, error) {
					retChan := make(chan error, 1)
					go func() {
						defer func() {
							if e := recover(); e != nil {
								std.LogRecover(e)
							}
						}()
						retChan <- fn(ctx, req, rsp)
					}()

					timeout := settings.Expiration / 2
					select {
					case err := <-retChan:
						var cobj cachedObj
						if err != nil {
							cobj.IvkErr = std.ErrFromGoErr(err)
							if std.IsSysErr(cobj.IvkErr, std.ErrInternalFromString) {
								return nil, err
							}
						} else {
							cobj.Data, err = marshalResp(rsp)
							if err != nil {
								return nil, err
							}
						}
						return &cobj, nil
					case <-time.After(timeout):
						std.LogErrorc("rpc", nil, fmt.Sprintf("fail to call rpc %s: timeout", req.Method()))
						return nil, std.ErrRpcCacheTimeout
					}
				},
			})
			if err != nil {
				if std.IsSysErr(err, std.ErrRpcCacheTimeout) {
					return err
				} else {
					std.LogErrorc("rpc", err, "fail to call rpc")
					return std.ErrRpcCache
				}
			} else {
				if cobj, ok := v.(*cachedObj); !ok {
					std.LogErrorc("rpc", nil, "rpc cache: invalid return type")
					return std.ErrRpcCache
				} else {
					if cobj.IvkErr != nil {
						return cobj.IvkErr
					} else {
						if err := unmarshalResp(cobj.Data, rsp); err != nil {
							std.LogErrorc("rpc", err, "fail to unmarshal response")
							return err
						}
					}
				}
				return nil
			}
		}
	}
}

func marshalResp(v interface{}) ([]byte, error) {
	if msg, ok := v.(proto.Message); !ok {
		return nil, std.ErrRpcCacheMarshal
	} else if b, err := proto.Marshal(msg); err != nil {
		return nil, std.ErrRpcCacheMarshal
	} else {
		return b, nil
	}
}

func unmarshalResp(b []byte, v interface{}) error {
	if msg, ok := v.(proto.Message); !ok {
		return std.ErrRpcCacheUnmarshal
	} else if err := proto.Unmarshal(b, msg); err != nil {
		return std.ErrRpcCacheUnmarshal
	} else {
		return nil
	}
}

func getCacheKey(ctx context.Context, methodName string, req interface{}) (string, error) {
	if msg, ok := req.(proto.Message); !ok {
		return "", std.ErrRpcCacheMarshal
	} else if b, err := proto.Marshal(msg); err != nil {
		return "", std.ErrRpcCacheMarshal
	} else {
		platform := stdc.PbGetPlatform(ctx)
		//return "rpcc:" + methodName + ":" + base64.RawURLEncoding.EncodeToString(b), nil
		return "rpcc:" + methodName + ":" + base64.RawURLEncoding.EncodeToString(b) + ":" + platform, nil
	}
}

