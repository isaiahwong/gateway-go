package internal

type MethodDesc struct {
	MethodName string `json:"methodName"`
	Path       string `json:"path"`
}

type ServiceDesc struct {
	Package      string       `json:"package"`
	ServiceName  string       `json:"serviceName"`
	OriginalName string       `json:"originalName"`
	Methods      []MethodDesc `json:"methods"`
}
