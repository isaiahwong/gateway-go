package k8s

import (
	"encoding/json"
	"fmt"

	"github.com/isaiahwong/gateway-go/internal/common/validator"
)

// APIServicesGetter has a method to return a APIServiceInterface.
type APIServicesGetter interface {
	APIServices() APIServicesInterface
}

// APIServicesInterface has methods to work with APIService
type APIServicesInterface interface {
	ObjectToAPI(*K8SObject) (*APIService, error)
}

type apiservices struct{}

// ObjectToAPI maps Object to APIService struct
func (a *apiservices) ObjectToAPI(o *K8SObject) (*APIService, error) {
	var s *APIService = &APIService{}
	validate := validator.Instance()
	err := validate.Struct(o)
	if err != nil {
		return nil, err
	}

	// Default values
	s.APIversion = "v1"
	s.Path = fmt.Sprintf("/%v/%v", s.APIversion, o.Metadata.Name)
	s.Authentication.Required = false

	// Parse annotations
	if o.Metadata.Annotations.Config != "" {
		err = json.Unmarshal([]byte(o.Metadata.Annotations.Config), s)
		// TODO: Handle JSON Error
		if err != nil {
			return nil, err
		}
	}

	dnsServiceName := fmt.Sprintf("%v.%v", o.Metadata.Name, o.Metadata.Namespace)
	s.DNSPath = dnsServiceName
	s.Ports = o.Spec.Ports
	s.ServiceName = o.Metadata.Name
	s.Namespace = o.Metadata.Namespace
	s.ResourceType = o.Metadata.Labels.ResourceType

	return s, nil
}

func newAPIServices() *apiservices {
	return &apiservices{}
}
