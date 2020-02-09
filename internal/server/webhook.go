package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaiahwong/gateway-go/internal/observer"
)

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
	Name     string
	Server   *http.Server
	Notifier *AdmissionNotifier
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

// NewWebhook returns a new webhook server
func NewWebhook(port string) (*WebhookServer, error) {
	// Initialize a new Notifier
	an := &AdmissionNotifier{
		observers: map[observer.Observer]struct{}{},
	}
	r, err := newRouter(webhookMiddleware(an))
	if err != nil {
		return nil, err
	}
	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}
	ws := &WebhookServer{
		Name:     "Admission Webhook Server",
		Server:   s,
		Notifier: an,
	}
	return ws, nil
}
