package server

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaiahwong/gateway-go/internal/common/log"
	"github.com/isaiahwong/gateway-go/internal/observer"
)

var defaultWebhookOptions = options{
	logger: log.NewLogger(),
	addr:   ":8443",
}

// AdmissionNotifier keep track of observers and notifies
// observers when webhook receives from AdmissionController
type AdmissionNotifier struct {
	observers map[observer.Observer]struct{}
}

// Register subscribes observers for an admission event
func (en *AdmissionNotifier) Register(l observer.Observer) {
	en.observers[l] = struct{}{}
}

// Deregister removes observers from the notifier lists
func (en *AdmissionNotifier) Deregister(l observer.Observer) {
	delete(en.observers, l)
}

// Notify broadcasts to all observers
func (en *AdmissionNotifier) Notify(e observer.Event) {
	for o := range en.observers {
		o.OnNotify(e)
	}
}

// WebhookServer encapsulates Webhookserver and Notifier
type WebhookServer struct {
	Name       string
	Server     *http.Server
	Notifier   *AdmissionNotifier
	production bool
	logger     log.Logger
	certFile   string
	keyFile    string
}

func webhookMiddleware(an *AdmissionNotifier) func(*gin.Engine) {
	return func(r *gin.Engine) {
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
				c.AbortWithStatusJSON(400, gin.H{"success": false})
				return
			}

			// Restore the io.ReadCloser to its original state
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			// Notify all subscribes
			an.Notify(observer.Event{Data: bodyBytes})
			c.JSON(200, gin.H{
				"success": true,
			})
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
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
		cancel()
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
		observers: map[observer.Observer]struct{}{},
	}
	r, err := newRouter(webhookMiddleware(an))
	if err != nil {
		return nil, err
	}
	s := &http.Server{
		Addr:    opts.addr,
		Handler: r,
	}
	ws := &WebhookServer{
		Name:       "Admission Webhook Server",
		Server:     s,
		Notifier:   an,
		production: opts.production,
		logger:     opts.logger,
		certFile:   opts.certFile,
		keyFile:    opts.keyFile,
	}
	return ws, nil
}
