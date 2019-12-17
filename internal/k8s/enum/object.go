package enum

type objectType string

// K8S object types https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/
const (
	K8SServiceObject    objectType = "service"
	K8SDeploymentObject            = "deployment"
)
