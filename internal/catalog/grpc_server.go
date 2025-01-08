package catalog

import (
	"context"

	"github.com/thebackendgrip/ecommerce-app/internal/common/observability"
	v1 "github.com/thebackendgrip/ecommerce-app/internal/grpc/v1"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository interface {
	UpdateInventory(ctx context.Context, items []Item, op operation) error
}

type CatalogServer struct {
	v1.UnimplementedCatalogServiceServer
	repo   Repository
	logger *zap.Logger

	traceProvider trace.TracerProvider
}

func (s CatalogServer) UpdateInventory(ctx context.Context, in *v1.UpdateInventoryRequest) (
	*v1.UpdateInventoryResponse, error,
) {
	observability.OpsProcessed.Inc()

	ctx, span := s.traceProvider.Tracer("catalog-api").Start(ctx, "UpdateInventory")
	defer span.End()

	var items []Item
	for _, i := range in.Items {
		items = append(items, Item{Name: i.Name, Qty: i.Qty})
	}

	var op operation
	if in.Op == v1.Op_Add {
		op = OP_ADD
	} else {
		op = OP_REMOVE
	}

	if err := s.repo.UpdateInventory(ctx, items, op); err != nil {
		s.logger.Error("could not update inventory", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.UpdateInventoryResponse{}, nil
}
