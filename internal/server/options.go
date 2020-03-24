package server

import (
	"github.com/isaiahwong/gateway-go/internal/common/log"
	"github.com/isaiahwong/gateway-go/internal/k8s"
)

type options struct {
	logger     log.Logger
	k8sClient  *k8s.Client
	production bool
	certFile   string
	keyFile    string
	addr       string
}

// Option sets options for Server.
type Option func(*options)

// WithAddress returns an Option which sets the address the server will be listening to.
func WithAddress(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

// WithLogger sets logger for gateway
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// WithK8SClient sets k8s client for GatewayServer.
// Though there isn't a generic type interface :(
func WithK8SClient(k *k8s.Client) Option {
	return func(o *options) {
		o.k8sClient = k
	}
}

// WithAppEnv sets Gateway server running mode
func WithAppEnv(e bool) Option {
	return func(o *options) {
		o.production = e
	}
}

// WithTLSCredentials returns an Option that sets TLS file directory
func WithTLSCredentials(certFile, keyFile string) Option {
	return func(o *options) {
		o.certFile = certFile
		o.keyFile = keyFile
	}
}
