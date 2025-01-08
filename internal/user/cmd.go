package user

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/thebackendgrip/ecommerce-app/internal/common/observability"
	v1 "github.com/thebackendgrip/ecommerce-app/internal/grpc/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func NewCmdAPI() *cobra.Command {
	return &cobra.Command{
		Use:   "user-api",
		Short: "User Service API",
		RunE: func(cmd *cobra.Command, args []string) error {

			listener, err := net.Listen("tcp", ":50001")
			if err != nil {
				log.Fatal(err)
			}

			dbFile := "userdb.sqlite"
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
			CREATE TABLE users (
				id uuid PRIMARY KEY NOT NULL,
				email varchar(320) NOT NULL
			);`); err != nil {
				log.Fatal("could not setup database: ", err)
			}

			logger, err := observability.NewLogger()
			if err != nil {
				log.Fatal("could not create logger")
			}

			s := grpc.NewServer(
				grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
				grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
			)
			v1.RegisterUserServiceServer(s, UserServer{
				repo:   NewSqlRepository(dbConn),
				logger: logger,
			})

			go observability.InitMetrics(2113)

			log.Print("starting grpc server at: ", listener.Addr())
			if err := s.Serve(listener); err != nil {
				log.Fatal("failed to serve: ", err)
			}

			return nil
		},
	}
}
