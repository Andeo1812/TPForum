package pkg

import (
	"context"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

//go:generate easyjson  -disallow_unknown_fields -omit_empty wrapper.go

//easyjson:json
type ErrResponse struct {
	ErrMassage string `json:"message,omitempty"`
}

// DefaultHandlerHTTPError is error handler that detects the type of error and gives an error response.
func DefaultHandlerHTTPError(ctx context.Context, w http.ResponseWriter, err error) {
	CreateLog(ctx, err)

	errCause := errors.Cause(err)

	code, exist := GetErrorCodeHTTP(errCause)
	if !exist {
		errCause = errors.Wrap(errCause, "Undefined error")
	}

	errResp := ErrResponse{
		ErrMassage: errCause.Error(),
	}

	Response(ctx, w, code, errResp)
}

// NoBody is function designed to give a response that
// has no request body, only the code and/or headers matter
func NoBody(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func getEasyJSON(someStruct interface{}) ([]byte, error) {
	someStructUpdate, ok := someStruct.(easyjson.Marshaler)
	if !ok {
		return []byte{}, ErrGetEasyJSON
	}

	out, err := easyjson.Marshal(someStructUpdate)
	if !ok {
		return []byte{}, ErrJSONUnexpectedEnd
	}

	return out, err
}

// Response is a function for giving any response with a JSON body
func Response(ctx context.Context, w http.ResponseWriter, statusCode int, someStruct interface{}) {
	out, err := getEasyJSON(someStruct)
	if err != nil {
		DefaultHandlerHTTPError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)

	w.WriteHeader(statusCode)

	_, err = w.Write(out)
	if err != nil {
		DefaultHandlerHTTPError(ctx, w, err)
		return
	}
}
