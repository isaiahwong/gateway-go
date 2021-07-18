package server

import (
	"github.com/go-redis/redis/v7"
	"gitlab.com/eco_system/gateway/internal/k8s"
	"github.com/sirupsen/logrus"
)

type options struct {
	logger           *logrus.Logger
	k8sClient        *k8s.Client
	redisClient      *redis.Client
	production       bool
	certFile         string
	keyFile          string
	addr             string
	accountsDisabled bool
	accountsAddr     string
	accountsTimeout  int
}

// Option sets options for Server.
type Option func(*options)

// WithAddress returns an Option which sets the address the server will be listening to.
func WithAddress(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

// WithPubSub returns an Option which sets the pubsub that the server will utilise
func WithPubSub(r *redis.Client) Option {
	return func(o *options) {
		o.redisClient = r
	}
}

// WithLogger returns an Option sets logger for gateway
func WithLogger(l *logrus.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// WithK8SClient returns an Option sets k8s client for GatewayServer.
// Though there isn't a generic type interface :(
func WithK8SClient(k *k8s.Client) Option {
	return func(o *options) {
		o.k8sClient = k
	}
}

// WithAppEnv returns an Option sets Gateway server running mode
func WithAppEnv(e bool) Option {
	return func(o *options) {
		o.production = e
	}
}

// WithAccountsDisabled returns an Option determines if accounts authentication is enabled
func WithAccountsDisabled(disable bool) Option {
	return func(o *options) {
		o.accountsDisabled = disable
	}
}

// WithAccountsAddr returns an Option which sets the accounts address
func WithAccountsAddr(addr string) Option {
	return func(o *options) {
		o.accountsAddr = addr
	}
}

// WithAccountsTimeout returns an Option sets the connection towards
// the accounts before deadline
func WithAccountsTimeout(t int) Option {
	return func(o *options) {
		o.accountsTimeout = t
	}
}

// WithTLSCredentials returns an Option that sets TLS file directory
func WithTLSCredentials(certFile, keyFile string) Option {
	return func(o *options) {
		o.certFile = certFile
		o.keyFile = keyFile
	}
}
