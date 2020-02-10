package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isaiahwong/gateway-go/internal/k8s"
	"github.com/isaiahwong/gateway-go/internal/util/log"
)

// notFoundMW aka not found middleware
func notFoundMW(c *gin.Context) {
	c.JSON(404, gin.H{
		"success": false,
		"error":   "not_found",
	})
}

// ReverseProxy routes traffic to the intended target
func ReverseProxy(target string) gin.HandlerFunc {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return func(c *gin.Context) {
		// If empty, the Request.Write method uses
		// the value of URL.Host. Host may contain an international
		// domain name
		c.Request.Host = target
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// WebhookRequests intercepts webhook requests and repackage it to fit grpc requirements
func WebhookRequests(c *gin.Context) {
	if strings.Contains(c.Request.URL.Path, "webhook") {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		p := &Payload{
			Body: buf,
		}
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			c.Next()
			return
		}
		rc := ioutil.NopCloser(bytes.NewBuffer(b))
		c.Request.Body = rc
	}
	c.Next()
}

func requestLogger(l log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		if l == nil {
			fmt.Println(c.Request.URL.Path, readBody(rdr1)) // Print request body
		} else {
			l.Infof(c.Request.URL.Path)
		}
		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

func authMW(services *map[string]*k8s.APIService) gin.HandlerFunc {
	// Retrieves service
	return func(c *gin.Context) {
		c.Next()
	}
}
