package service

import (
	"fmt"
	"github.com/micro/go-micro/client/selector"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	std "github.com/PKUJohnson/solar/std"
	tureg "github.com/PKUJohnson/solar/std/service/registry"
	"github.com/PKUJohnson/solar/std/service/selector/cache"
	stdnet "github.com/PKUJohnson/solar/std/toolkit/net"
	gclient "github.com/micro/go-micro/client/grpc"
	//"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/server"
	gserver "github.com/micro/go-micro/server/grpc"
	"github.com/micro/go-os/trace"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/go-plugins/registry/zookeeper"
	"github.com/micro/go-plugins/transport/grpc"
	"github.com/micro/go-plugins/transport/http"
	"github.com/micro/go-plugins/transport/tcp"
	"github.com/micro/go-plugins/transport/utp"
)

const (
	grpcBetaListenAddr  = ":43210"
	grpcBetaSuffix      = "-grpc"
	grpcServer          = "grpc"
)

func NewGRPCBetaService(cfg std.ConfigService) micro.Service {
	cfg.RPCServer = grpcServer
	cfg.SvcAddr = grpcBetaListenAddr
	return NewService(cfg)
}

func NewService(svcCfg std.ConfigService, opts ...micro.Option) micro.Service {
	var reg registry.Registry

	if len(svcCfg.EtcdAddrs) > 0 {
		reg = tureg.NewRegistry(registry.Addrs(svcCfg.EtcdAddrs...))
	} else if len(svcCfg.Consul) > 0 {
		addrs := make([]string, 0, len(svcCfg.Consul))
		for _, c := range svcCfg.Consul {
			addrs = append(addrs, c.Addr)
		}
		reg = consul.NewRegistry(registry.Addrs(addrs...))
	} else if len(svcCfg.ZookeeperAddrs) > 0 {
		reg = zookeeper.NewRegistry(registry.Addrs(svcCfg.ZookeeperAddrs...))
	} else {
		panic("Failed to create registry: no config found")
	}

	initMetrics(svcCfg.Prometheus)
	initCircuitBreaker(svcCfg.Hystrix, svcCfg.AdvertiseSubnets)
	initCache(svcCfg.RpcCacheRedis)
	initRedisClient(svcCfg.Redis)
	return newService(svcCfg, configIns, reg, opts...)
}

func newService(svcCfg std.ConfigService,
	microCfg microConfig, reg registry.Registry, opts ...micro.Option) micro.Service {
	// server options
	var serverOpts []server.Option
	if svcCfg.SvcName != "" {
		name := svcCfg.SvcName
		if svcCfg.RPCServer == grpcServer && !strings.HasSuffix(svcCfg.SvcName, grpcBetaSuffix) {
			name = name + grpcBetaSuffix
		}
		serverOpts = append(serverOpts, server.Name(name))
	}

	advIp, advPort := getServiceIPPort(svcCfg)
	if advIp != "" && advPort != "" {
		addrStr := advIp + ":" + advPort
		serverOpts = append(serverOpts, server.Address(addrStr))
	} else {
		// as a fallback, mostly for local dev
		serverOpts = append(serverOpts, server.Address(svcCfg.SvcAddr))
	}

	serverOpts = append(serverOpts, server.WrapHandler(LogWrapper))
	serverOpts = append(serverOpts, server.WrapHandler(CacheWrapper))
	serverOpts = append(serverOpts, server.WrapHandler(MetricsWrapper))

	// client options
	var clientOpts []client.Option
	sel, reg := getSelector(reg)
	clientOpts = append(clientOpts, client.Selector(sel))
	clientOpts = append(clientOpts, client.RequestTimeout(60*time.Second))
	clientOpts = append(clientOpts, client.PoolSize(microCfg.ClientPoolSize))
	clientOpts = append(clientOpts, client.Wrap(newCircuitBreakerWrapper()))
	//clientOpts = append(clientOpts, client.Wrap(newMetricsClientWrapper()))

	// zipkin
	if svcCfg.Zipkin.Enabled {
		co, so, err := getZipkinOptions(svcCfg, advIp, advPort)
		if err != nil {
			std.LogErrorc("zipkin", err, "invalid zipkin config")
		} else {
			clientOpts = append(clientOpts, co)
			serverOpts = append(serverOpts, so)
			std.LogInfoc("zipkin", "success to init zipkin")
		}
	}

	// transport option for both server and client
	if microCfg.Transport != "" {
		if microCfg.Transport == "grpc" {
			grpcTrans := grpc.NewTransport()
			serverOpts = append(serverOpts, server.Transport(grpcTrans))
			clientOpts = append(clientOpts, client.Transport(grpcTrans))
			log.Println("Micro transport: grpc")
		} else if microCfg.Transport == "tcp" {
			tcpTrans := tcp.NewTransport()
			serverOpts = append(serverOpts, server.Transport(tcpTrans))
			clientOpts = append(clientOpts, client.Transport(tcpTrans))
			log.Println("Micro transport: tcp")
		} else if microCfg.Transport == "utp" {
			utpTrans := utp.NewTransport()
			serverOpts = append(serverOpts, server.Transport(utpTrans))
			clientOpts = append(clientOpts, client.Transport(utpTrans))
		} else if microCfg.Transport == "http" {
			httpTrans := http.NewTransport()
			serverOpts = append(serverOpts, server.Transport(httpTrans))
			clientOpts = append(clientOpts, client.Transport(httpTrans))
		}
	}

	serverOpts = append(serverOpts, server.Registry(reg))
	clientOpts = append(clientOpts, client.Registry(reg))

	// register ttl
	serverOpts = append(serverOpts, server.RegisterTTL(time.Duration(microCfg.RegisterTTL)*time.Second))

	// finally set options
	options := []micro.Option{
		micro.Registry(reg),
		micro.RegisterInterval(time.Duration(microCfg.RegisterInterval) * time.Second),
		micro.AfterStop(func() error {
			time.Sleep(time.Second)
			return nil
		}),
	}

	if os.Getenv("RPC_CLIENT") == "grpc" {
		options = append(options, micro.Client(gclient.NewClient(clientOpts...)))
	} else {
		options = append(options, micro.Client(client.NewClient(clientOpts...)))
	}

	if svcCfg.RPCServer == "grpc" {
		options = append(options, micro.Server(gserver.NewServer(serverOpts...)))
	} else {
		options = append(options, micro.Server(server.NewServer(serverOpts...)))
	}

	options = append(options, opts...)
	return micro.NewService(options...)
}

func getSelector(reg registry.Registry) (selector.Selector, registry.Registry) {
	std.LogInfoLn("service: default selector is chosen")
	return cache.NewSelector(selector.Registry(reg), cache.TTL(time.Second*30)), reg
}

func getZipkinOptions(
	svcCfg std.ConfigService, advIp, advPort string) (client.Option, server.Option, error) {
	port, err := strconv.Atoi(advPort)
	fmt.Println(port)
	if err != nil {
		return nil, nil, err
	}

	if advIp == "" {
		if advIp, err = stdnet.LocalIPAddr(); err != nil {
			std.LogWarnc("net", err, "fail to get my IP")
		}
	}

	service := &registry.Service{
		Name: svcCfg.SvcName,
		Nodes: []*registry.Node{
			{
				Address: fmt.Sprintf("%s:%d", advIp, port),
			},
		},
	}
	t := NewTrace(
		trace.Collectors(svcCfg.Zipkin.BrokerAddrs...),
	)
	co := client.Wrap(trace.ClientWrapper(t, service))
	so := server.WrapHandler(trace.HandlerWrapper(t, service))
	return co, so, nil
}

func getServiceIPPort(svcCfg std.ConfigService) (string, string) {
	if len(svcCfg.AdvertiseSubnets) > 0 {
		// find the container's IP within the specified subnet
		advertiseIP, err := stdnet.LocalIPAddrWithin(svcCfg.AdvertiseSubnets)
		if err != nil {
			std.LogErrorLn("failed to retrieve the advertise IP within subnets [", svcCfg.AdvertiseSubnets, "]: ", err)
			return "", ""
		}
		parts := strings.Split(svcCfg.SvcAddr, ":")
		if len(parts) < 2 {
			std.LogErrorLn("failed to split service address: ", svcCfg.SvcAddr)
			return "", ""
		}
		return advertiseIP, parts[1]
	}

	addr := svcCfg.SvcAddr
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
