package service

import (
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-os/trace"
	"github.com/micro/go-plugins/trace/zipkin/thrift/gen-go/zipkincore"

	"strings"

	"github.com/apache/thrift/lib/go/thrift"
	std "github.com/PKUJohnson/solar/std"
	stdnet "github.com/PKUJohnson/solar/std/toolkit/net"
	sarama "gopkg.in/Shopify/sarama.v1"
)

type zipkinKey struct{}

type zipkin struct {
	opts  trace.Options
	spans chan *trace.Span
	exit  chan bool
}

var (
	TraceTopic = "zipkin"

	TraceHeader  = "X-B3-TraceId"
	SpanHeader   = "X-B3-SpanId"
	ParentHeader = "X-B3-ParentSpanId"
	SampleHeader = "X-B3-Sampled"
)

func random() int64 {
	return rand.Int63() & 0x001fffffffffffff
}

func newZipkin(opts ...trace.Option) trace.Trace {
	var opt trace.Options
	for _, o := range opts {
		o(&opt)
	}

	if opt.BatchSize == 0 {
		opt.BatchSize = trace.DefaultBatchSize
	}

	if opt.BatchInterval == time.Duration(0) {
		opt.BatchInterval = trace.DefaultBatchInterval
	}

	if len(opt.Topic) == 0 {
		opt.Topic = TraceTopic
	}

	z := &zipkin{
		exit:  make(chan bool),
		opts:  opt,
		spans: make(chan *trace.Span, 100),
	}

	go z.run()
	return z
}

func toInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func toEndpoint(s *registry.Service) *zipkincore.Endpoint {
	ep := zipkincore.NewEndpoint()
	if s == nil || len(s.Nodes) == 0 {
		ep.ServiceName = s.Name
		return ep
	}

	addrs, err := net.LookupIP(s.Nodes[0].Address)
	if err != nil {
		return ep
	}
	if len(addrs) == 0 {
		return ep
	}

	ep.Ipv4 = int32(stdnet.IP2Int(addrs[0]))
	//ep.Port = int16(s.Nodes[0].Port)
	if parts := strings.Split(s.Name, "."); len(parts) > 0 {
		ep.ServiceName = parts[len(parts)-1]
	} else {
		ep.ServiceName = s.Name
	}
	return ep
}

func toThrift(s *trace.Span) *zipkincore.Span {
	span := &zipkincore.Span{
		TraceID:   toInt64(s.TraceId),
		Name:      s.Name,
		ID:        toInt64(s.Id),
		Debug:     s.Debug,
		Timestamp: thrift.Int64Ptr(s.Timestamp.UnixNano() / 1e3),
		Duration:  thrift.Int64Ptr(s.Duration.Nanoseconds() / 1e3),
	}
	if parentID := toInt64(s.ParentId); parentID != 0 {
		span.ParentID = thrift.Int64Ptr(parentID)
	}

	for _, a := range s.Annotations {
		if len(a.Value) > 0 || a.Debug != nil {
			span.BinaryAnnotations = append(span.BinaryAnnotations, &zipkincore.BinaryAnnotation{
				Key:            a.Key,
				Value:          a.Value,
				AnnotationType: zipkincore.AnnotationType_BYTES,
				Host:           toEndpoint(a.Service),
			})
		} else {
			var val string
			switch a.Type {
			case trace.AnnClientRequest:
				val = zipkincore.CLIENT_SEND
			case trace.AnnClientResponse:
				val = zipkincore.CLIENT_RECV
			case trace.AnnServerRequest:
				val = zipkincore.SERVER_RECV
			case trace.AnnServerResponse:
				val = zipkincore.SERVER_SEND
			default:
				val = a.Key
			}

			if len(val) == 0 {
				continue
			}
			span.Annotations = append(span.Annotations, &zipkincore.Annotation{
				Timestamp: a.Timestamp.UnixNano() / 1e3,
				Value:     val,
				Host:      toEndpoint(a.Service),
			})
		}
	}

	return span
}

func (z *zipkin) pub(s *zipkincore.Span, pr sarama.SyncProducer) {
	t := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocolTransport(t)

	if err := s.Write(p); err != nil {
		std.LogErrorc("kafka", err, "fail to send to kafka")
		return
	}

	m := &sarama.ProducerMessage{
		Topic: z.opts.Topic,
		Key:   nil,
		Value: sarama.ByteEncoder(t.Buffer.Bytes()),
	}

	_, _, err := pr.SendMessage(m)
	if err != nil {
		std.LogErrorc("kafka", err, "fail to send to kafka")
	}
}

func (z *zipkin) run() {
	t := time.NewTicker(z.opts.BatchInterval)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	c, err := sarama.NewClient(z.opts.Collectors, config)
	if err != nil {
		std.LogErrorc("kafka", err, "fail to initialize kafka client")
		return
	}

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		std.LogErrorc("kafka", err, "fail to initialize kafka client")
		return
	}

	var buf []*trace.Span

	for {
		select {
		case s := <-z.spans:
			buf = append(buf, s)
			if len(buf) >= z.opts.BatchSize {
				go z.send(buf, p)
				buf = nil
			}
		case <-t.C:
			// flush
			if len(buf) > 0 {
				go z.send(buf, p)
				buf = nil
			}
		case <-z.exit:
			// exit
			t.Stop()
			p.Close()
			return
		}
	}
}

func (z *zipkin) send(b []*trace.Span, p sarama.SyncProducer) {
	for _, span := range b {
		z.pub(toThrift(span), p)
	}
}

func (z *zipkin) Close() error {
	select {
	case <-z.exit:
		return nil
	default:
		close(z.exit)
	}
	return nil
}

func (z *zipkin) Collect(s *trace.Span) error {
	select {
	case z.spans <- s:
	default:
		std.LogWarnLn("z.spans channel is full, discading data")
		//std.LogErrorLn("z.spans channel is full, discading data")
	}
	return nil
}

func (z *zipkin) NewSpan(s *trace.Span) *trace.Span {
	if s == nil {
		return &trace.Span{
			Id:        strconv.FormatInt(random(), 10),
			TraceId:   strconv.FormatInt(random(), 10),
			ParentId:  "0",
			Timestamp: time.Now(),
			Source:    z.opts.Service,
		}
	}

	if _, err := strconv.ParseInt(s.TraceId, 16, 64); err != nil {
		s.TraceId = strconv.FormatInt(random(), 10)
	}
	if _, err := strconv.ParseInt(s.ParentId, 16, 64); err != nil {
		s.ParentId = "0"
	}
	if _, err := strconv.ParseInt(s.Id, 16, 64); err != nil {
		s.Id = strconv.FormatInt(random(), 10)
	}

	if s.Timestamp.IsZero() {
		s.Timestamp = time.Now()
	}

	return &trace.Span{
		Id:        s.Id,
		TraceId:   s.TraceId,
		ParentId:  s.ParentId,
		Timestamp: s.Timestamp,
	}
}

func (z *zipkin) String() string {
	return "zipkin"
}

func NewTrace(opts ...trace.Option) trace.Trace {
	return newZipkin(opts...)
}
