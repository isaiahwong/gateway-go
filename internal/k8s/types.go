package k8s

import "google.golang.org/grpc"

// Ports represent k8s ports
type Port struct {
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
	Path           string            `json:"path"`
	DNSPath        string            `validate:"required"`
	Ports          []Port            `validate:"required"`
	Authentication APIAuthentication `json:"authentication"`
	ServiceName    string            `validate:"required"`
	Namespace      string            `validate:"required"`
	APIversion     string            `json:"apiVersion"`
	ResourceType   string            `validate:"required"`
	GRPCClientConn *grpc.ClientConn
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
	Ports []Port `json:"Ports" validate:"required"`
}

// K8SObject represents k8s object
type K8SObject struct {
	Metadata Metadata `json:"Metadata" validate:"required"`
	Spec     Spec     `json:"Spec" validate:"required"`
}

// Kind represents k8s kind
type Kind struct {
	Kind string `json:"kind"`
}

// AdmissionRequest represents k8s AdmissionRegistration request
type AdmissionRequest struct {
	Operation string    `json:"operation" validate:"required"`
	Kind      Kind      `json:"kind" validate:"required"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Object    K8SObject `json:"object"`
	OldObject K8SObject `json:"oldObject"`
}

// AdmissionRegistration represents k8s admissionregistration.k8s.io/v1
type AdmissionRegistration struct {
	Request AdmissionRequest `json:"request" validate:"required"`
}
