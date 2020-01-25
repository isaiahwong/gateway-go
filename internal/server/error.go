package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/ptypes/any"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
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

type errorBody struct {
	Error string `protobuf:"bytes,100,name=error" json:"error"`
	// This is to make the error more compatible with users that expect errors to be Status objects:
	// https://github.com/grpc/grpc/blob/master/src/proto/grpc/status/status.proto
	// It should be the exact same message as the Error field.
	Code    int32      `protobuf:"varint,1,name=code" json:"code"`
	Message string     `protobuf:"bytes,2,name=message" json:"message"`
	Details []*any.Any `protobuf:"bytes,3,rep,name=details" json:"details,omitempty"`
}

// HTTPStatusFromCode converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	}

	grpclog.Infof("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}

func CustomHTTPError(ctx context.Context, _ *gwruntime.ServeMux, marshaler gwruntime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(gwruntime.HTTPStatusFromCode(grpc.Code(err)))
	jErr := json.NewEncoder(w).Encode(errorBody{
		Error: grpc.ErrorDesc(err),
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
