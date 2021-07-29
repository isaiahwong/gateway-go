package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
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

type Error struct {
	Param   string      `json:"param"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
	Errors  [][]Error   `json:"errors"`
}

type errorBody struct {
	Error string `protobuf:"bytes,100,name=error" json:"error"`
	// This is to make the error more compatible with users that expect errors to be Status objects:
	// https://github.com/grpc/grpc/blob/master/src/proto/grpc/status/status.proto
	// It should be the exact same message as the Error field.
	// Message string     `protobuf:"bytes,2,name=message" json:"message"`
	Details []*any.Any `protobuf:"bytes,3,rep,name=details" json:"details,omitempty"`
	Errors  []Error    `json:"errors"`
}

// ProtoErrorWithLogger replies to the request with the error.
// Overrides runtime.error HTTPError
func ProtoErrorWithLogger(l *logrus.Logger, strict bool) func(context.Context, *runtime.ServeMux, runtime.Marshaler, http.ResponseWriter, *http.Request, error) {
	return func(ctx context.Context, m *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, protoErr error) {
		const fallback = `{"error": "An unexpected error occurred."}`
		eb := errorBody{}
		code := runtime.HTTPStatusFromCode(status.Code(protoErr))

		w.Header().Set("Content-type", marshaler.ContentType(""))
		eb.Error = status.Convert(protoErr).Message()
		md, ok := runtime.ServerMetadataFromContext(ctx)
		if ok {
			if details := md.TrailerMD.Get("errors-bin"); len(details) > 0 {
				e := []Error{}
				// Maps json values to error body
				err := json.Unmarshal([]byte(details[0]), &e)
				if err != nil {
					l.Errorf("ProtoErrorWithLogger: %v", err)
				} else {
					eb.Errors = e
				}
			}
		}
		jsonByte, err := json.Marshal(eb)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fallback))
			l.Errorf("ProtoErrorWithLogger: %v")
			return
		}

		// Log errors
		l.WithFields(logrus.Fields{
			"error":         string(jsonByte),
			"requestUrl":    req.URL,
			"requestMethod": req.Method,
			"remoteIp":      req.Header.Get("X-Forwarded-For"),
		}).Errorln(eb.Error)

		// Handle 501 error
		if code == 501 {
			code = 404
			eb.Error = "Not Found"
		} else if code >= 500 {
			code = 500
			eb.Error = "An unexpected error occurred"
			eb.Errors = nil
		}
		// Change message of sensitive messages
		// Usually gRPC marshal or type errors
		if strict {
			knownTypeErrors(&eb)
		}
		w.WriteHeader(code)

		err = json.NewEncoder(w).Encode(eb)

		if err != nil {
			w.Write([]byte(fallback))
		}
	}
}

// knownTypeErrors replaces sensitive errors with a said message
func knownTypeErrors(eb *errorBody) {
	known := []string{
		"json: cannot unmarshal",
		"proto: (line",
	}
	for _, k := range known {
		if strings.Contains(eb.Error, k) {
			eb.Error = "Malformed Request"
		}
	}
}
