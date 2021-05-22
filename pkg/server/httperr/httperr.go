package httperr

import (
	"net/http"
)

// ErrResponse is an error response payload.
type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
}

// InternalError prepares internal server error.
func InternalError(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error",
		ErrorText:      err.Error(),
	}
}

// BadRequest prepares bad request error.
func BadRequest(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad request",
		ErrorText:      err.Error(),
	}
}

// NotFoundRequest prepares not found error.
func NotFoundRequest(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Resource not found",
		ErrorText:      err.Error(),
	}
}