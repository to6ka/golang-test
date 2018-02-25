package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Service struct{}

func (s Service) Handler() *httptransport.Server {
	return httptransport.NewServer(
		s.MakeEndpoint(),
		s.DecodeRequest,
		s.EncodeResponse,
	)
}

type Request struct {
	A, B, C int
}

type Response struct {
	Result int
}

func (Service) MakeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, in interface{}) (response interface{}, err error) {
		req, ok := in.(Request)
		if !ok {
			return nil, fmt.Errorf("unknown endpoint input type: %#v", in)
		}

		return Response{
			Result: req.A*req.B + req.C,
		}, nil
	}
}

func (Service) DecodeRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	buf := bytes.NewBuffer([]byte{})
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return nil, fmt.Errorf("can't read req body: %s", err)
	}

	var req Request
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		return nil, fmt.Errorf("can't unmarshal req body in json: %s", err)
	}

	return req, nil
}

func (Service) EncodeResponse(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("can't marshal resp body in json: %s", err)
	}

	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("can't write resp body in output: %s", err)
	}

	return nil
}
