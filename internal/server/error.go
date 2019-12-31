package server

type invalidParams struct {
	s string
}

func (e *invalidParams) Error() string {
	return e.s
}

func InvalidParams(msg string) *invalidParams {
	return &invalidParams{"Invalid Params: " + msg}
}

type notFound struct {
	s string
}

func (e *notFound) Error() string {
	return e.s
}

func NotFound(msg string) *notFound {
	return &notFound{"Not Found: " + msg}
}
