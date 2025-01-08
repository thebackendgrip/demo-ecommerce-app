package user

import (
	"context"
	"fmt"

	"github.com/thebackendgrip/ecommerce-app/internal/common/observability"
	v1 "github.com/thebackendgrip/ecommerce-app/internal/grpc/v1"
	"go.uber.org/zap"
)

type Repository interface {
	createUser(ctx context.Context, email string) error
}

type UserServer struct {
	v1.UnimplementedUserServiceServer
	repo   Repository
	logger *zap.Logger
}

func NewUserServer(repo Repository) UserServer {
	return UserServer{
		repo: repo,
	}
}

func (s UserServer) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.User, error) {
	observability.OpsProcessed.Inc()

	if err := s.repo.createUser(ctx, in.Email); err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	s.logger.Info("created user", zap.String("email", in.Email))
	return nil, nil
}
