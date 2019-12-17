package enum

type operation string

// K8S operation type
const (
	Create operation = "CREATE"
	Delete           = "DELETE"
)
