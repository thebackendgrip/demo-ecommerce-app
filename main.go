package main

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/thebackendgrip/ecommerce-app/internal/catalog"
	"github.com/thebackendgrip/ecommerce-app/internal/order"
	"github.com/thebackendgrip/ecommerce-app/internal/user"
)

func main() {
	cmd := &cobra.Command{
		Use:   "ecommerce [command]",
		Short: "Ecommerce App",
	}

	cmd.AddCommand(user.NewCmdAPI())
	cmd.AddCommand(catalog.NewCmdAPI())
	cmd.AddCommand(order.NewCmdAPI())

	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
