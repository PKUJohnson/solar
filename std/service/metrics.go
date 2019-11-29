package service

import (
	//"github.com/micro/go-micro"

	"sync"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/metrics"
	// "github.com/micro/go-plugins/metrics/prometheus"
	"math/rand"
	"strconv"

	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/service/metrics/prometheus"
	"golang.org/x/net/context"
)

var (
	MetricsMgrIns *MetricsManager
)

type MetricsManager struct {
	enabled     bool
	id          string
	metrics     metrics.Metrics
	counters    map[string]metrics.Counter
	counterLock sync.Mutex
	gauges      map[string]metrics.Gauge
	gaugeLock   sync.Mutex
	histograms  map[string]metrics.Histogram
	histoLock   sync.Mutex
	concurrency map[string]int
	concurLock  sync.RWMutex
}

func (mm *MetricsManager) GetCounter(id string) metrics.Counter {
	mm.counterLock.Lock()
	defer mm.counterLock.Unlock()

	if c, ok := mm.counters[id]; ok {
		return c
	}
	c := mm.metrics.Counter(id)
	mm.counters[id] = c
	return c
}

func (mm *MetricsManager) GetGauge(id string) metrics.Gauge {
	mm.gaugeLock.Lock()
	defer mm.gaugeLock.Unlock()

	if g, ok := mm.gauges[id]; ok {
		return g
	}
	g := mm.metrics.Gauge(id)
	mm.gauges[id] = g
	return g
}

func (mm *MetricsManager) GetHistogram(id string) metrics.Histogram {
	mm.histoLock.Lock()
	defer mm.histoLock.Unlock()

	if h, ok := mm.histograms[id]; ok {
		return h
	}
	h := mm.metrics.Histogram(id)
	mm.histograms[id] = h
	return h
}

func (mm *MetricsManager) UpdateConcurrency(id string, change int) int {
	mm.concurLock.Lock()
	defer mm.concurLock.Unlock()
	result := mm.concurrency[id] + change
	mm.concurrency[id] = result
	return result
}

func (mm *MetricsManager) GetConcurrency(id string) int {
	mm.concurLock.RLock()
	defer mm.concurLock.RUnlock()
	return mm.concurrency[id]
}

func initMetrics(conf std.ConfigPrometheus) {
	if !conf.Enabled {
		MetricsMgrIns = &MetricsManager{}
	} else {
		MetricsMgrIns = &MetricsManager{
			id:          strconv.FormatInt(rand.Int63()&0x001fffffffffffff, 16),
			enabled:     conf.Enabled,
			metrics:     newMetrics(conf),
			counters:    make(map[string]metrics.Counter),
			gauges:      make(map[string]metrics.Gauge),
			histograms:  make(map[string]metrics.Histogram),
			concurrency: make(map[string]int),
		}
		std.LogInfoLn("prometheus initialized for " + MetricsMgrIns.id)
	}
}

func newMetrics(conf std.ConfigPrometheus) metrics.Metrics {
	options := make([]metrics.Option, 0)
	options = append(options, metrics.WithFields(metrics.Fields{
		"service": "",
		// "method":  "",
		// "source":  "",
		"status": "",
		// "id":     "",
	}))
	if conf.Namespace != "" {
		options = append(options, metrics.Namespace(conf.Namespace))
	}
	if conf.BatchInterval > 0 {
		options = append(options, metrics.BatchInterval(time.Duration(conf.BatchInterval)*time.Millisecond))
	}
	if len(conf.Collectors) > 0 {
		addrs := make([]string, 0, len(conf.Collectors))
		for _, col := range conf.Collectors {
			addrs = append(addrs, col.Addr)
		}
		options = append(options, metrics.Collectors(addrs...))
	}
	return prometheus.NewMetrics(options...)
}

func MetricsWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		mm := MetricsMgrIns
		if !mm.enabled {
			return fn(ctx, req, rsp)
		}

		begin := time.Now()
		labels := map[string]string{
			"service": req.Method(),
		}

		err := fn(ctx, req, rsp)
		if err != nil {
			labels["status"] = "fail"
		} else {
			labels["status"] = "success"
		}

		mm.GetCounter("rpc_requests_total").WithFields(labels).Incr(1)

		duration := time.Since(begin)
		mm.GetGauge("rpc_response_time_nonseconds").WithFields(labels).Set(duration.Nanoseconds())
		mm.GetHistogram("rpc_response_time_nonseconds_his").WithFields(labels).Record(duration.Nanoseconds())

		return err
	}
}

// MetricsClientWrapper represents the client prometheus wrapper
type MetricsClientWrapper struct {
	client.Client
}

func (c *MetricsClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if !MetricsMgrIns.enabled {
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// mm := MetricsMgrIns
	labels := map[string]string{
		"service": req.Service(),
		// "method":  req.Method(),
		// "source":  "client",
		// "status":  "",
		// "id": mm.id,
	}

	// begin := time.Now()
	err := c.Client.Call(ctx, req, rsp, opts...)
	if err != nil {
		labels["status"] = "fail"
	} else {
		labels["status"] = "success"
	}
	// mm.GetCounter("rpc_requests_total").WithFields(labels).Incr(1)

	// duration := time.Since(begin)
	// mm.GetHistogram("rpc_response_time_nanoseconds").WithFields(labels).Record(int64(duration))
	return err
}

// newMetricsClientWrapper returns a hystrix client Wrapper.
func newMetricsClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &MetricsClientWrapper{c}
	}
}
