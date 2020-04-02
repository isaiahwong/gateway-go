package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/ptypes/any"
	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

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

type errors struct {
	Param   string      `json:"param"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

type errorBody struct {
	Error string `protobuf:"bytes,100,name=error" json:"error"`
	// This is to make the error more compatible with users that expect errors to be Status objects:
	// https://github.com/grpc/grpc/blob/master/src/proto/grpc/status/status.proto
	// It should be the exact same message as the Error field.
	Message string     `protobuf:"bytes,2,name=message" json:"message"`
	Details []*any.Any `protobuf:"bytes,3,rep,name=details" json:"details,omitempty"`
	Errors  []errors   `json:"errors"`
}

// HTTPErrorWithLogger replies to the request with the error.
// Overrides runtime.error HTTPError
func HTTPErrorWithLogger(l *logrus.Logger) func(context.Context, *runtime.ServeMux, runtime.Marshaler, http.ResponseWriter, *http.Request, error) {
	return func(ctx context.Context, m *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
		eb := errorBody{}

		code := runtime.HTTPStatusFromCode(grpc.Code(err))
		w.Header().Set("Content-type", marshaler.ContentType())
		eb.Error = grpc.ErrorDesc(err)

		md, ok := runtime.ServerMetadataFromContext(ctx)
		if ok {
			if details := md.TrailerMD.Get("errors-bin"); len(details) > 0 {
				e := []errors{}
				// Maps json values to error body
				err = json.Unmarshal([]byte(details[0]), &e)
				if err != nil {
					l.Errorf("HTTPErrorWithLogger: %v", err)
				} else {
					eb.Errors = e
				}
			}
		}
		processErrors(w, req, code, &eb, l)
	}
}

// OtherErrorWithLogger handles the following error used by the gateway: StatusMethodNotAllowed StatusNotFound and StatusBadRequest
// Overrides runtime.error OtherErrorHandler
func OtherErrorWithLogger(l *logrus.Logger) func(http.ResponseWriter, *http.Request, string, int) {
	return func(w http.ResponseWriter, req *http.Request, msg string, code int) {
		w.Header().Set("Content-Type", "application/json")
		eb := errorBody{}
		eb.Error = msg

		processErrors(w, req, code, &eb, l)
	}
}

func processErrors(w http.ResponseWriter, req *http.Request, code int, eb *errorBody, l *logrus.Logger) {
	const fallback = `{"error": "An unexpected error occurred."}`
	jsonByte, err := json.Marshal(eb)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fallback))
		l.Errorf("HTTPErrorWithLogger: %v")
		return
	}

	// Log errors
	l.WithFields(logrus.Fields{
		"error":         string(jsonByte),
		"requestUrl":    req.URL,
		"requestMethod": req.Method,
		"remoteIp":      req.Header.Get("X-Forwarded-For"),
	}).Errorln(eb.Error)

	if code > 500 {
		code = 500
		eb.Error = "An unexpected error occurred"
		eb.Errors = nil
	}
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(eb)

	if err != nil {
		w.Write([]byte(fallback))
	}
}
