package order

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/thebackendgrip/ecommerce-app/internal/common/observability"
	pb "github.com/thebackendgrip/ecommerce-app/internal/grpc/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
)

func NewCmdAPI() *cobra.Command {
	return &cobra.Command{
		Use:   "order-api",
		Short: "Order Service API",
		RunE: func(cmd *cobra.Command, args []string) error {
			listener, err := net.Listen("tcp", ":50002")
			if err != nil {
				log.Fatal(err)
			}

			dbFile := "orderdb.sqlite"
			if _, err := os.Create(dbFile); err != nil {
				log.Fatal("could not create database file: ", err)
			}

			db, err := sql.Open("sqlite3", dbFile)
			if err != nil {
				log.Fatal("could not open db: ", err)
			}

			ctx := context.Background()

			dbConn, err := db.Conn(ctx)
			if err != nil {
				log.Fatal("could not open new db connection: ", err)
			}

			if _, err := dbConn.ExecContext(ctx, `
			CREATE TABLE orders (
				id uuid PRIMARY KEY NOT NULL,
				user_id uuid NOT NULL,
				order_items JSON NOT NULL
			);`); err != nil {
				log.Fatal("could not setup database: ", err)
			}

			catalogClientConn, err := grpc.NewClient(
				"localhost:50003",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithStatsHandler(otelgrpc.NewClientHandler(
					otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
				)),
			)
			if err != nil {
				log.Fatalf("could not create catalog connection: %v", err)
			}
			defer catalogClientConn.Close()

			catalogClient := pb.NewCatalogServiceClient(catalogClientConn)

			logger, err := observability.NewLogger()
			if err != nil {
				log.Fatal("could not create logger")
			}

			tp, err := observability.InitTracing("order-api")
			if err != nil {
				log.Fatal("could not init trace provider: %w", err)
			}

			s := grpc.NewServer()
			pb.RegisterOrderServiceServer(s, &OrderServer{
				repo:          NewSqlRepository(dbConn),
				catalogClient: catalogClient,
				logger:        logger,
				traceProvider: tp,
			})

			go observability.InitMetrics(2114)

			log.Print("starting grpc server at: ", listener.Addr())
			if err := s.Serve(listener); err != nil {
				log.Fatal("failed to serve: ", err)
			}

			return nil
		},
	}
}
