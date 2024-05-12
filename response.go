package requests

import (
	"encoding/json"
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
		return nil, err
	}
	return &t, nil
}
