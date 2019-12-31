package k8s

import (
	"encoding/json"
)

// AdmissionGetter has a method to return a AdmissionInterface.
type AdmissionGetter interface {
	Admission() AdmissionInterface
}

// AdmissionInterface has methods to work with K8S Admission Object
type AdmissionInterface interface {
	Unmarshal(d []byte) (*AdmissionRegistration, error)
	UnmarshalK8SObject(d []byte) (*K8SObject, error)
}

// admission implements Admission interface
type admission struct{}

func (a *admission) Unmarshal(d []byte) (*AdmissionRegistration, error) {
	if len(d) <= 0 {
		return nil, EmptyData()
	}
	var ar *AdmissionRegistration = &AdmissionRegistration{}
	err := json.Unmarshal(d, ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (a *admission) UnmarshalK8SObject(d []byte) (*K8SObject, error) {
	if len(d) <= 0 {
		return nil, EmptyData()
	}
	var o *K8SObject = &K8SObject{}
	err := json.Unmarshal(d, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func newAdmission() *admission {
	return &admission{}
}
