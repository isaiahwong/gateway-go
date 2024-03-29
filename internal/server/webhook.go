package server

import (
	"bytes"
	"context"
	"gitlab.com/eco_system/gateway/internal/k8s"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"gitlab.com/eco_system/gateway/internal/common/log"
	"gitlab.com/eco_system/gateway/internal/observer"
)

type RedisChannel string

var (
	RGateway RedisChannel = "gatewayadmission"
)

var defaultWebhookOptions = options{
	logger: log.NewLogger(),
	addr:   ":8443",
}

// AdmissionNotifier keep track of observers and notifies
// observers when webhook receives from AdmissionController
type AdmissionNotifier struct {
	observers   map[observer.Observer]struct{}
	redisClient *redis.Client
}

// Register subscribes observers for an admission event
func (an *AdmissionNotifier) Register(l observer.Observer) {
	an.observers[l] = struct{}{}
}

// Deregister removes observers from the notifier lists
func (an *AdmissionNotifier) Deregister(l observer.Observer) {
	delete(an.observers, l)
}

// Notify broadcasts to all observers
func (an *AdmissionNotifier) Notify(e observer.Event) {
	if an.redisClient == nil {
		for o := range an.observers {
			o.OnNotify(e)
		}
		return
	}
	an.redisClient.Publish(string(RGateway), string(e.Data))
}

// WebhookServer encapsulates Webhookserver and Notifier
type WebhookServer struct {
	Name       string
	Server     *http.Server
	Notifier   *AdmissionNotifier
	production bool
	logger     *logrus.Logger
	certFile   string
	keyFile    string
}

func webhookMiddleware(an *AdmissionNotifier, k *k8s.Client) func(*gin.Engine) {
	return func(r *gin.Engine) {
		var ar *k8s.AdmissionRegistration

		res := gin.H{
			"success":    true,
			"apiVersion": "admission.k8s.io/v1",
			"kind":       "AdmissionReview",
			"response": map[string]interface{}{
				"allowed": true,
			},
		}
		r.Use(gin.Recovery())
		r.POST("/admission", func(c *gin.Context) {
			if c.Request.Body == nil {
				c.JSON(200, gin.H{
					"success": true,
				})
				return
			}

			// Read the content
			bodyBytes, err := ioutil.ReadAll(c.Request.Body)

			if err != nil {
				// TODO Log
				res["response"].(map[string]interface{})["allowed"] = false
				c.AbortWithStatusJSON(200, res)
				return
			}

			ar, err = k.CoreAPI().Admission().Unmarshal(bodyBytes)
			if err != nil {
				// TODO Log
				res["response"].(map[string]interface{})["allowed"] = false
				c.AbortWithStatusJSON(200, res)
				return
			}
			// Restore the io.ReadCloser to its original state
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			// Notify all subscribes
			an.Notify(observer.Event{Data: bodyBytes})
			res["response"].(map[string]interface{})["uid"] = ar.Request.UID
			c.JSON(200, res)
		})
	}
}

func (ws *WebhookServer) Gracefully(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := ws.Server.Shutdown(ctx); err != nil {
		ws.logger.Errorf("Gracefully: Error shutting down %v server: %v\n", ws.Name, err)
	}
	ws.logger.Infof("%v closed\n", ws.Name)
}

// Run executes WebhookServer
func (ws *WebhookServer) Run(ctx context.Context) error {
	_, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		cancelTLS := false
		if ws.certFile == "" || ws.keyFile == "" {
			cancelTLS = true
		}

		ws.logger.Infof("Running %v on [%v] - Production: %v\n", ws.Name, ws.Server.Addr, ws.production)
		if cancelTLS {
			ws.logger.Warnln("Running Webhook without TLS")
			if err := ws.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				ws.logger.Fatalf("Webhook server: %s\n", err)
			}
		} else {
			// Start webhook server
			if err := ws.Server.ListenAndServeTLS(ws.certFile, ws.keyFile); err != nil && err != http.ErrServerClosed {
				ws.logger.Fatalf("Webhook server: %s\n", err)
			}
		}
	}()

	return nil
}

// NewWebhook returns a new webhook server
func NewWebhook(opt ...Option) (*WebhookServer, error) {
	opts := defaultWebhookOptions
	for _, o := range opt {
		o(&opts)
	}
	// Initialize a new Notifier
	an := &AdmissionNotifier{
		observers:   map[observer.Observer]struct{}{},
		redisClient: opts.redisClient,
	}
	r, err := newRouter(webhookMiddleware(an, opts.k8sClient))
	if err != nil {
		return nil, err
	}
	s := &http.Server{
		Addr:    opts.addr,
		Handler: r,
	}
	ws := &WebhookServer{
		Name:       "Request Webhook Server",
		Server:     s,
		Notifier:   an,
		production: opts.production,
		logger:     opts.logger,
		certFile:   opts.certFile,
		keyFile:    opts.keyFile,
	}
	return ws, nil
}
