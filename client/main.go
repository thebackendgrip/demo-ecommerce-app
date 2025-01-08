package main

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/thebackendgrip/ecommerce-app/internal/grpc/v1"
)

func main() {
	userClientConn, err := grpc.NewClient(
		"localhost:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not create grpc client conn: %v", err)
	}
	defer userClientConn.Close()

	ctx := context.Background()
	userClient := pb.NewUserServiceClient(userClientConn)
	userResp, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Email: gofakeit.Email(),
	})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}

	// add item to inventory
	catalogClientConn, err := grpc.NewClient(
		"localhost:50003",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not create catalog connection: %v", err)
	}
	defer catalogClientConn.Close()

	catalogClient := pb.NewCatalogServiceClient(catalogClientConn)
	_, err = catalogClient.UpdateInventory(ctx, &pb.UpdateInventoryRequest{
		Items: []*pb.Item{
			{
				Name: "nike-shoe-air-jordan-4",
				Qty:  3,
			},
		},
		Op: pb.Op_Add,
	})
	if err != nil {
		log.Fatalf("could not update inventory: %v", err)
	}

	// Create order
	orderClientConn, err := grpc.NewClient(
		"localhost:50002",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not create grpc client conn: %v", err)
	}
	defer orderClientConn.Close()

	orderClient := pb.NewOrderServiceClient(orderClientConn)
	_, err = orderClient.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId: userResp.Id,
		Items: []*pb.Item{
			{
				Name: "nike-shoe-air-jordan-5",
				Qty:  1,
			},
		},
	})
	if err != nil {
		log.Fatalf("could not create order: %v", err)
	}
}
