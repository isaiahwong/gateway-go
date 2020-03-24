package services

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type clientOption struct {
	address string
	timeout time.Duration
}

var defaultOptions clientOption = clientOption{
	address: ":50051",
	timeout: 1 * time.Second,
}

// Option defines the APIClient options
type Option func(*clientOption)

// WithAddress returns an Option which sets the APIClient remote address
func WithAddress(a string) Option {
	return func(o *clientOption) {
		o.address = a
	}
}

// WithTimeout returns an option which sets the connection timeout on connection initially
func WithTimeout(t int) Option {
	return func(o *clientOption) {
		o.timeout = time.Duration(t) * time.Second
	}
}

// CreateClient returns a Client Connection
func CreateClient(opt ...Option) (*grpc.ClientConn, error) {
	var opts = &defaultOptions
	for _, o := range opt {
		o(opts)
	}

	conn, err := grpc.Dial(
		opts.address,
		grpc.WithInsecure(),
		grpc.WithTimeout(opts.timeout),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("CreateClient: %v", err)
	}
	return conn, nil
}
