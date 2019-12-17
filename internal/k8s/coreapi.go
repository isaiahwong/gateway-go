package k8s

type CoreAPIInterface interface {
	AdmissionGetter
	APIServicesGetter
}

type CoreAPI struct{}

func (c *CoreAPI) Admission() AdmissionInterface {
	return newAdmission()
}

func (c *CoreAPI) APIServices() APIServicesInterface {
	return newAPIServices()
}

func newCoreApi() (*CoreAPI, error) {
	return &CoreAPI{}, nil
}
