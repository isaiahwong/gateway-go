package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/ptypes/any"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
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
	Param   string `json:"param"`
	Message string `json:"message"`
	Value   string `json:"value"`
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

// HTTPError replies to the request with the error.
// Overrides runtime.error HTTPError
func HTTPError(ctx context.Context, _ *gwruntime.ServeMux, marshaler gwruntime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`
	code := gwruntime.HTTPStatusFromCode(grpc.Code(err))
	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(code)

	eb := errorBody{
		Error: grpc.ErrorDesc(err),
	}
	md, ok := gwruntime.ServerMetadataFromContext(ctx)
	if ok {
		details := md.TrailerMD.Get("errors-bin")[0]
		// Maps json values to error body
		json.Unmarshal([]byte(details), &eb)
	}
	jErr := json.NewEncoder(w).Encode(eb)

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}

// OtherErrorHandler handles the following error used by the gateway: StatusMethodNotAllowed StatusNotFound and StatusBadRequest
// Overrides runtime.error OtherErrorHandler
func OtherErrorHandler(w http.ResponseWriter, _ *http.Request, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorBody{
		Error: msg,
	})
}
