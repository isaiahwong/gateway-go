package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/isaiahwong/gateway-go/internal/common"
	"github.com/isaiahwong/gateway-go/internal/common/log"
	"github.com/isaiahwong/gateway-go/internal/k8s"
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
	AppEnv           string
	Production       bool
	Host             string
	Address          string
	WebhookAddress   string
	WebhookSecretKey string
	DisableK8S       bool

	WebhookKeyDir  string
	WebhookCertDir string

	EnableStackdriver bool

	DBUri      string
	DBName     string
	DBUser     string
	DBPassword string
	DBTimeout  time.Duration
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

	config = &EnvConfig{
		AppEnv:            mapEnvWithDefaults("APP_ENV", "development"),
		Production:        mapEnvWithDefaults("APP_ENV", "development") == "production",
		DisableK8S:        mapEnvWithDefaults("DISABLE_K8S_CLIENT", "true") == "true",
		Address:           mapEnvWithDefaults("ADDRESS", ":5000"),
		WebhookAddress:    mapEnvWithDefaults("WEBHOOK_ADDRESS", ":8443"),
		WebhookKeyDir:     mapEnvWithDefaults("WEBHOOK_KEY_DIR", ""),
		WebhookCertDir:    mapEnvWithDefaults("WEBHOOK_CERT_DIR", ""),
		EnableStackdriver: mapEnvWithDefaults("ENABLE_STACKDRIVER", "true") == "true",
	}
}

var logger *logrus.Logger

func init() {
	loadEnv()
	logger = log.NewLogger()
}

// Execute the entry point for gateway
func main() {
	var k *k8s.Client
	var err error

	if !config.DisableK8S {
		k, err = k8s.NewClient()
		if err != nil {
			logger.Fatalf("K8SClient Error: %v", err)
		}
	}

	gs, err := server.NewGatewayServer(
		server.WithAddress(config.Address),
		server.WithLogger(logger),
		server.WithK8SClient(k),
		server.WithAppEnv(config.Production),
	)
	if err != nil {
		logger.Fatalf("New Gateway error: %v", err)
	}

	ws, err := server.NewWebhook(
		server.WithAddress(config.WebhookAddress),
		server.WithTLSCredentials(config.WebhookCertDir, config.WebhookKeyDir),
		server.WithAppEnv(config.Production),
	)
	if err != nil {
		logger.Fatalf("New Webhook error: %v", err)
	}

	// Registers gateway as an observer
	ws.Notifier.Register(gs)

	gc := common.SignalContext(context.Background())
	wc := common.SignalContext(context.Background())

	// Start gateway Server
	gs.Run(gc)
	// Start Webhook server
	ws.Run(wc)

	select {
	case <-gc.Done():
	case <-wc.Done():
		gs.Gracefully(gc)
		ws.Gracefully(wc)
	}
}
