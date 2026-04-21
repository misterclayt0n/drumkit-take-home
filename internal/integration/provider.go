package integration

import (
	"context"

	"drumkit-take-home/internal/load"
)

type Provider interface {
	ListLoads(ctx context.Context, params load.ListParams) (load.ListResponse, error)
	CreateLoad(ctx context.Context, input load.Load) (load.CreateResponse, error)
}
