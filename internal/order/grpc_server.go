package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/thebackendgrip/ecommerce-app/internal/common/observability"
	v1 "github.com/thebackendgrip/ecommerce-app/internal/grpc/v1"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository interface {
	createOrder(ctx context.Context, o Order) error
}

const traceProviderName = "order-api"

type OrderServer struct {
	v1.UnimplementedOrderServiceServer
	repo          Repository
	catalogClient v1.CatalogServiceClient
	logger        *zap.Logger

	traceProvider *trace.TracerProvider
}

func (s *OrderServer) CreateOrder(ctx context.Context, in *v1.CreateOrderRequest) (*v1.Order, error) {
	observability.OpsProcessed.Inc()

	ctx, span := s.traceProvider.Tracer(traceProviderName).Start(ctx, "CreateOrder")
	defer span.End()

	var items []Item
	for _, i := range in.Items {
		items = append(items, Item{Name: i.Name, Qty: i.Qty})
	}

	o := Order{
		ID:     uuid.NewString(),
		UserID: in.UserId,
		Items:  items,
	}

	// update inventory
	if _, err := s.catalogClient.UpdateInventory(ctx, &v1.UpdateInventoryRequest{
		Items: in.Items,
		Op:    v1.Op_Remove,
	}); err != nil {
		s.logger.Error("could not remove items from inventory", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	// create order
	if err := s.repo.createOrder(ctx, o); err != nil {
		// try to revert changes made to inventory - add items
		if _, err := s.catalogClient.UpdateInventory(ctx, &v1.UpdateInventoryRequest{
			Items: in.Items,
			Op:    v1.Op_Add,
		}); err != nil {
			s.logger.Error("could not update inventory after order failure", zap.Any("items", o))
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("created order", zap.String("order-id", o.ID))
	return &v1.Order{}, nil
}
