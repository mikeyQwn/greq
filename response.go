package greq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response[T any] struct {
	core *http.Response

	body []byte
	err  error
}

func NewResponse[T any](core *http.Response) *Response[T] {
	body, err := io.ReadAll(core.Body)
	if err != nil {
		err = fmt.Errorf("%w: could not read the request's body", err)
	}
	return &Response[T]{
		core: core,
		body: body,
		err:  err,
	}
}

func (r *Response[T]) String() string {
	return string(r.body)
}

func (r *Response[T]) BaseType() (*T, error) {
	var t T
	if r.err != nil {
		return nil, r.err
	}
	err := json.Unmarshal(r.body, &t)
	if err != nil {
		return nil, fmt.Errorf("%w: could not convert response body to the base type", err)
	}
	return &t, nil
}

func (r *Response[T]) BodyRaw() []byte {
	return r.body
}

func (r *Response[T]) StatusCode() int {
	return r.core.StatusCode
}

func (r *Response[T]) Error() error {
	return r.err
}

func (r Response[T]) Header() http.Header {
	return r.core.Header
}
