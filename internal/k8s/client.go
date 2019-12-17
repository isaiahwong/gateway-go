package k8s

import (
	"encoding/json"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientset *kubernetes.Clientset
	coreAPI   *CoreAPI
}

func (c *Client) GetServices(ns string) ([][]byte, error) {
	svcs := [][]byte{}
	cs := c.clientset
	list, err := cs.CoreV1().Services(ns).List(v1.ListOptions{
		TimeoutSeconds: &[]int64{10}[0],
	})

	if err != nil {
		return nil, err
	}

	for _, s := range list.Items {
		b, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		svcs = append(svcs, b)
	}
	return svcs, nil
}

func (c *Client) CoreAPI() CoreAPIInterface {
	return c.coreAPI
}

// NewClient creates a new Client wrapper for gateway
func NewClient() (*Client, error) {
	var c Client
	var err error
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	c.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	c.coreAPI, err = newCoreApi()
	return &c, nil
}
