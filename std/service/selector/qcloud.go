package selector

import (
	"errors"
	"fmt"
	"strings"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/client/selector"
)

const (
	svcPrefix = "solar."
	svcPort   = 10087
)

type qcloud struct{}

func (r *qcloud) Init(opts ...selector.Option) error {
	return nil
}

func (r *qcloud) Options() selector.Options {
	return selector.Options{}
}

func (r *qcloud) Select(service string, opts ...selector.SelectOption) (selector.Next, error) {
	i := strings.Index(service, svcPrefix)
	if i < 0 {
		errors.New("solar selector: bad service name")
	}
	node := &registry.Node{
		Id:      service,
		Address: fmt.Sprintf("%s:%d", service[i+len(svcPrefix):], svcPort),
	}
	return func() (*registry.Node, error) {
		return node, nil
	}, nil
}

func (r *qcloud) Mark(service string, node *registry.Node, err error) {
	return
}

func (r *qcloud) Reset(service string) {
	return
}

func (r *qcloud) Close() error {
	return nil
}

func (r *qcloud) String() string {
	return "solar"
}

func NewSelector(opts ...selector.Option) selector.Selector {
	return &qcloud{}
}
