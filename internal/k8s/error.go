package k8s

type emptyData struct {
	s string
}

func (e *emptyData) Error() string {
	return e.s
}

func EmptyData() *emptyData {
	return &emptyData{"Data is empty"}
}
