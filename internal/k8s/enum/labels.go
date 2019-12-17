package enum

type resourceType string

// resourceType is a type of label selector that acts as filter
const (
	LabelAPIService resourceType = "api-service"
)