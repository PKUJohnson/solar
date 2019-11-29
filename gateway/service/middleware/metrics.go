package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/micro/go-os/metrics"
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/service/metrics/prometheus"
)

// APIMetrics collects the API metrics
func APIMetrics(conf std.ConfigPrometheus) func(echo.HandlerFunc) echo.HandlerFunc {
	if !conf.Enabled {
		return func(handle echo.HandlerFunc) echo.HandlerFunc {
			return handle
		}
	}

	std.LogDebugLn("start to open api metrics")
	metric := newMetrics(conf)
	apiCounter := metric.Counter("http_requests_total")
	apiLatency := metric.Histogram("http_requests_latency_in_seconds")
	return func(handle echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			start := time.Now()
			err = handle(c)
			end := time.Now()

			if err == echo.ErrNotFound {
				return
			}

			// pass through the non-API requests
			p := c.Path()
			if p == "" || p == "/" || !strings.HasPrefix(p, "/apiv1/") {
				return
			}

			// send metrics to prometheus
			apiCounter.WithFields(metrics.Fields{
				"path": p,
				"code": strconv.Itoa(c.Response().Status),
			}).Incr(1)

			paths := strings.Split(p, "/")
			if len(paths) > 2 {
				apiLatency.WithFields(metrics.Fields{
					"path": paths[2],
				}).Record(end.Sub(start).Nanoseconds())
			}
			return
		}
	}
}

func newMetrics(conf std.ConfigPrometheus) metrics.Metrics {
	options := make([]metrics.Option, 0)
	options = append(options, metrics.WithFields(metrics.Fields{
		"path": "",
		"code": "",
	}))
	if conf.Namespace != "" {
		options = append(options, metrics.Namespace(conf.Namespace))
	}
	if conf.BatchInterval > 0 {
		options = append(options, metrics.BatchInterval(time.Duration(conf.BatchInterval) * time.Millisecond))
	}
	if len(conf.Collectors) > 0 {
		addrs := make([]string, 0, len(conf.Collectors))
		for _, col := range conf.Collectors {
			std.LogInfoLn("prometheus pushgateway " + col.Addr + " added")
			addrs = append(addrs, col.Addr)
		}
		options = append(options, metrics.Collectors(addrs...))
	}
	return prometheus.NewMetrics(options...)
}
