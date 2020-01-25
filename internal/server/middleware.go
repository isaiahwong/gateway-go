package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httputil"
	"net/url"

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

func reverseProxyMW(target string) gin.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
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

	}
}
