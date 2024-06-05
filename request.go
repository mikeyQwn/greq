package greq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request[T any] struct {
	headers map[string]string

	body io.Reader
	err  error
}

func NewRequest[T any]() *Request[T] {
	return &Request[T]{
		headers: map[string]string{},

		body: nil,
		err:  nil,
	}
}

func (r *Request[T]) WithJson(body any) *Request[T] {
	encoded, err := json.Marshal(body)
	r.err = fmt.Errorf("%w: WithJson could not marshal request's body", err)
	r.body = bytes.NewReader(encoded)
	return r
}

func (r *Request[T]) WithHeader(key, value string) *Request[T] {
	r.headers[key] = value
	return r
}

func (r *Request[T]) WithHeaders(headers map[string]string) *Request[T] {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

func (r *Request[T]) Copy() *Request[T] {
	return NewRequest[T]().WithHeaders(r.headers)
}

func (r *Request[T]) Get(url string) (*Response[T], error) {
	return r.GetContext(context.Background(), url)
}

func (r *Request[T]) MustGet(url string) *Response[T] {
	resp, err := r.GetContext(context.Background(), url)
	if err != nil {
		return &Response[T]{
			err: err,
		}
	}
	return resp
}

func (r *Request[T]) GetContext(ctx context.Context, url string) (*Response[T], error) {
	return r.doReqContext(ctx, url, http.MethodGet)
}

func (r *Request[T]) Post(url string) (*Response[T], error) {
	return r.PostContext(context.Background(), url)
}

func (r *Request[T]) MustPost(url string) *Response[T] {
	resp, err := r.PostContext(context.Background(), url)
	if err != nil {
		return &Response[T]{
			err: err,
		}
	}
	return resp
}

func (r *Request[T]) PostContext(ctx context.Context, url string) (*Response[T], error) {
	return r.doReqContext(ctx, url, http.MethodPost)
}

func (r *Request[T]) Err() error {
	return r.err
}

func (r *Request[T]) doReqContext(ctx context.Context, url string, method string) (*Response[T], error) {
	if r.err != nil {
		return nil, r.err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, r.body)
	if err != nil {
		return nil, fmt.Errorf("%w: error creating a new request", err)
	}
	for k, v := range r.headers {
		req.Header.Add(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: error sending a request", err)
	}
	return NewResponse[T](resp), nil
}
