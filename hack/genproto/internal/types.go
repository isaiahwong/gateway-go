package internal

type MethodDesc struct {
	MethodName string `json:"methodName"`
	Path       string `json:"path"`
}

type ServiceDesc struct {
	ServiceName     string       `json:"serviceName"`
	Methods         []MethodDesc `json:"methods"`
	PackageSVC      string       `json:"packageSvc"`
	Package         string       `json:"package"`
	CurrentPackage  string       `json:"currentPackage"`
	OriginalPackage string       `json:"originalPackage"`
	Path            string       `json:"path"`
	FilePath        string       `json:"filePath"`
}

type ProtoMap struct {
	Protos     []string `json:"protos"`
	Includes   []string `json:"includes"`
	Googleapis bool     `json:"googleapis"`
}
