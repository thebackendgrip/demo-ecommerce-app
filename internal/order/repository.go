package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

var _ Repository = (*SqlRepository)(nil)

type Order struct {
	ID     string
	UserID string
	Items  []Item
}

type Item struct {
	Name string
	Qty  int64
}

type SqlRepository struct {
	conn *sql.Conn
}

func NewSqlRepository(conn *sql.Conn) *SqlRepository {
	return &SqlRepository{
		conn: conn,
	}
}

// createUser implements Repository.
func (s *SqlRepository) createOrder(ctx context.Context, o Order) (err error) {
	orderItemsJson, err := json.Marshal(o.Items)
	if err != nil {
		return fmt.Errorf("could not marshal order items")
	}

	if _, err := s.conn.ExecContext(ctx, `
		INSERT INTO orders (id, user_id, order_items) values (?, ?, ?)
	`, o.ID, o.UserID, orderItemsJson); err != nil {
		return err
	}

	return nil
}
