package encode

import (
	"context"
	"fmt"
	"io"

	encodepb "github.com/fvdveen/mu2-proto/go/proto/encode"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/transport"
)

// Service downloads the video at the url and encodes it into opus
type Service interface {
	Encode(context.Context, string) (OpusReader, error)
}

// OpusReader returns an opus frame and an error
type OpusReader interface {
	OpusFrame() ([]byte, error)
}

type service struct {
	s    encodepb.EncodeService
	opts []client.CallOption
}

// NewService creates a new service
func NewService(loc string, cc *api.Config, opts ...client.CallOption) Service {
	srv := grpc.NewService(
		micro.Registry(consul.NewRegistry(consul.Config(cc))),
		micro.Transport(
			transport.NewTransport(
				transport.Secure(true),
			),
		),
	)

	s := encodepb.NewEncodeService(loc, srv.Client())

	return &service{
		s:    s,
		opts: opts,
	}
}

func (s *service) Encode(ctx context.Context, url string) (OpusReader, error) {
	str, err := s.s.Encode(ctx, &encodepb.EncodeRequest{
		Url: url,
	}, s.opts...)
	if err != nil {
		return nil, fmt.Errorf("call encode: %v", err)
	}

	return &opusReader{str}, nil
}

type opusReader struct {
	s encodepb.EncodeService_EncodeService
}

func (r *opusReader) OpusFrame() ([]byte, error) {
	res, err := r.s.Recv()
	if err == io.EOF {
		return nil, io.EOF
	} else if err != nil {
		return nil, fmt.Errorf("recieve message: %v", err)
	}

	return res.Opus, nil
}
