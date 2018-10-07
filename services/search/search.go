package search

import (
	"context"

	searchpb "github.com/fvdveen/mu2-proto/go/proto/search"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/transport"
)

// Service searches youtube and returns a Video
type Service interface {
	Search(context.Context, string) (*Video, error)
}

// Video represents a youtube video
type Video struct {
	ID           string
	Name         string
	URL          string
	ThumbnailURL string
}

type service struct {
	s    searchpb.SearchService
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

	s := searchpb.NewSearchService(loc, srv.Client())

	return &service{
		s:    s,
		opts: opts,
	}
}

func (s *service) Search(ctx context.Context, n string) (*Video, error) {
	res, err := s.s.Search(ctx, &searchpb.SearchRequest{
		Name: n,
	}, s.opts...)
	if err != nil {
		return nil, err
	}

	v := &Video{
		ID:           res.Video.Id,
		Name:         res.Video.Name,
		URL:          res.Video.Url,
		ThumbnailURL: res.Video.Thumbnail,
	}

	return v, nil
}
