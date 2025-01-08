package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func NewSqlRepository(conn *sql.Conn) *SqlRepository {
	return &SqlRepository{
		conn: conn,
	}
}

type operation int

const (
	OP_ADD operation = iota
	OP_REMOVE
)

type Item struct {
	Name string
	Qty  int64
}

type SqlRepository struct {
	conn *sql.Conn
}

var _ Repository = (*SqlRepository)(nil)

// updateInventory implements Repository.
func (r *SqlRepository) UpdateInventory(ctx context.Context, items []Item, op operation) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not create transaction: %w", err)
	}

	query := `SELECT qty FROM inventory WHERE name = ?`
	insertStmt := `REPLACE INTO inventory (name, qty) VALUES (?, ?)`
	for _, i := range items {
		var qty int
		if err := tx.QueryRowContext(ctx, query, i.Name).Scan(&qty); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				if err := tx.Rollback(); err != nil {
					// log error message
				}
				return fmt.Errorf("could not update item (%s): %w", i.Name, err)
			}
			if errors.Is(err, sql.ErrNoRows) && op == OP_REMOVE {
				// nothing to do, records not found to subtract from
				continue
			}
		}

		var updatedQty int
		if op == OP_ADD {
			updatedQty = qty + int(i.Qty)
		} else {
			updatedQty = qty - int(i.Qty)
		}

		if updatedQty < 0 {
			if err := tx.Rollback(); err != nil {
				// log error message
			}
			return fmt.Errorf("item (%s) qty cannot be negative", i.Name)
		}

		if _, err := tx.ExecContext(ctx, insertStmt, i.Name, updatedQty); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("could not update inventory item")
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not update inventory")
	}

	return nil
}
