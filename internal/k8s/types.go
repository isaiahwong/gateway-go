package k8s

import (
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Ports represent k8s ports
type Ports struct {
	Name       string      `json:"name"`
	Protocol   string      `json:"protocol"`
	Port       int         `json:"port"`
	TargetPort interface{} `json:"targetPort"`
}

// APIAuthentication defines a set of rules
// to which gateway should adhere
// Exclude routes from authentication if authentication is enabled
type APIAuthentication struct {
	Required interface{} `json:"required"`
	Exclude  []string    `json:"exclude"`
}

// APIService encapsulates k8s object into a simpler format
// Path path to route to service '/api/v1/payment
type APIService struct {
	Path           string `json:"path"`
	DNSPath        string
	Ports          []Ports
	Authentication APIAuthentication `json:"authentication"`
	ServiceName    string
	Namespace      string
	APIversion     string `json:"apiVersion"`
	ResourceType   string
	GRPCClientConn *gwruntime.ServeMux
}

// Labels represent k8s labels
type Labels struct {
	ResourceType string `json:"resourceType" validate:"required"`
}

// Annotations represents k8s annotations
type Annotations struct {
	Config string `json:"config"`
}

// Metadata represents k8s metadata
type Metadata struct {
	Namespace   string      `json:"namespace" validate:"required"`
	Labels      Labels      `json:"Labels" validate:"required"`
	Name        string      `json:"name" validate:"required"`
	Annotations Annotations `json:"Annotations"`
}

// Spec represents k8s spec
type Spec struct {
	Ports []Ports `json:"Ports" validate:"required"`
}

// Object represents k8s object
type Object struct {
	Metadata Metadata `json:"Metadata" validate:"required"`
	Spec     Spec     `json:"Spec" validate:"required"`
}

// Kind represents k8s kind
type Kind struct {
	Kind string `json:"kind"`
}

// AdmissionRequest represents k8s AdmissionRegistration request
type AdmissionRequest struct {
	Operation string `json:"operation" validate:"required"`
	Kind      Kind   `json:"kind" validate:"required"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Object    Object `json:"object" validate:"required"`
}

// AdmissionRegistration represents k8s admissionregistration.k8s.io/v1
type AdmissionRegistration struct {
	Request AdmissionRequest `json:"request" validate:"required"`
}
