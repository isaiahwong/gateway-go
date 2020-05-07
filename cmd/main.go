package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	redisV7 "github.com/go-redis/redis/v7"
	"github.com/isaiahwong/gateway-go/internal/common"
	"github.com/isaiahwong/gateway-go/internal/common/log"
	"github.com/isaiahwong/gateway-go/internal/k8s"
	"github.com/isaiahwong/gateway-go/internal/redis"
	"github.com/isaiahwong/gateway-go/internal/server"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// EnvConfig Application wide env configurations
// AppEnv specifies if the app is in `development` or `production`
// Host specifies host address or dns
// Port specifies the port the server will run on
// EnableStackdriver specifies if google stackdriver will be enabled
// StripeSecret specifies Stripe api production key
// StripeSecretDev specifies Stripe api key for development
// StripeEndpointSecret specifies Stripe api key for webhook verification
// PaypalClientIDDev specifies Paypal api key for development
// PaypalSecretDev specifies Paypal api key secret for development
// PaypalClientID
// PaypalSecret         string
// PaypalURL specifies Paypal api URL for request
// DBUri
// DBUriDev
// DBUriTest
// DBName
// DBUser
// DBPassword
type EnvConfig struct {
	AppEnv            string
	Production        bool
	Host              string
	Address           string
	AccountsAddress   string
	WebhookAddress    string
	WebhookSecretKey  string
	DisableK8S        bool
	WebhookKeyDir     string
	WebhookCertDir    string
	AccountsTimeout   int
	EnableStackdriver bool
	DBUri             string
	DBName            string
	DBUser            string
	DBPassword        string
	DBTimeout         time.Duration
	RedisAddr         string
	RedisPassword     string
}

// AppConfig config from EnvConfig
var config *EnvConfig

func mapEnvWithDefaults(envKey string, defaults string) string {
	v := os.Getenv(envKey)
	if v == "" {
		return defaults
	}
	return v
}

// LoadEnv loads environment variables for Application
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not loaded", err)
	}

	at, err := strconv.Atoi(mapEnvWithDefaults("ACCOUNTS_TIMEOUT", "10"))
	if err != nil {
		at = 10
	}

	config = &EnvConfig{
		AppEnv:            mapEnvWithDefaults("APP_ENV", "development"),
		Production:        mapEnvWithDefaults("APP_ENV", "development") == "production",
		DisableK8S:        mapEnvWithDefaults("DISABLE_K8S_CLIENT", "true") == "true",
		Address:           mapEnvWithDefaults("ADDRESS", ":5000"),
		AccountsAddress:   mapEnvWithDefaults("ACCOUNTS_ADDRESS", ":50051"),
		WebhookAddress:    mapEnvWithDefaults("WEBHOOK_ADDRESS", ":8443"),
		WebhookKeyDir:     mapEnvWithDefaults("WEBHOOK_KEY_DIR", ""),
		WebhookCertDir:    mapEnvWithDefaults("WEBHOOK_CERT_DIR", ""),
		AccountsTimeout:   at,
		EnableStackdriver: mapEnvWithDefaults("ENABLE_STACKDRIVER", "true") == "true",
		RedisAddr:         mapEnvWithDefaults("REDIS_ADDR", ""),
		RedisPassword:     mapEnvWithDefaults("REDIS_ADDR", ""),
	}
}

var logger *logrus.Logger

func init() {
	loadEnv()
	logger = log.NewLogger()

	if config.Production {
		gin.SetMode(gin.ReleaseMode)
	}
}

// Execute the entry point for gateway
func main() {
	var k *k8s.Client
	var r *redisV7.Client
	var err error

	if !config.DisableK8S {
		k, err = k8s.NewClient()
		if err != nil {
			logger.Fatalf("K8SClient Error: %v", err)
		}
	}

	if config.RedisAddr != "" {
		// Initialize a new Redis Client
		r, err = redis.New(
			redis.WithAddress(config.RedisAddr),
			redis.WithPassword(config.RedisPassword),
			redis.WithDBTimeout(config.DBTimeout),
		)
		if err != nil {
			logger.Fatalf("Redis: %v", err)
		}
	}

	gs, err := server.NewGatewayServer(
		server.WithAddress(config.Address),
		server.WithAccountsAddr(config.AccountsAddress),
		server.WithAccountsTimeout(config.AccountsTimeout),
		server.WithLogger(logger),
		server.WithK8SClient(k),
		server.WithAppEnv(config.Production),
		server.WithPubSub(r),
	)
	if err != nil {
		logger.Fatalf("NewGatewayServer: %v", err)
	}

	ws, err := server.NewWebhook(
		server.WithAddress(config.WebhookAddress),
		server.WithTLSCredentials(config.WebhookCertDir, config.WebhookKeyDir),
		server.WithAppEnv(config.Production),
		server.WithPubSub(r),
	)
	if err != nil {
		logger.Fatalf("New Webhook error: %v", err)
	}

	// Registers gateway as an observer
	ws.Notifier.Register(gs)

	ctx := common.SignalContext(context.Background())

	// Start gateway Server
	gs.Run(ctx)
	// Start Webhook server
	ws.Run(ctx)

	// Terminate all servers
	select {
	case <-ctx.Done():
		gs.Gracefully(ctx)
		ws.Gracefully(ctx)
	}
}
