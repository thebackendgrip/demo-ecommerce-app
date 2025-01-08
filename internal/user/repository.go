package user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

var _ Repository = (*SqlRepository)(nil)

type User struct {
	ID    string
	Email string
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
func (s *SqlRepository) createUser(ctx context.Context, email string) (err error) {
	userID := uuid.NewString()
	if _, err := s.conn.ExecContext(ctx, `
		INSERT INTO users (id, email) values (?, ?)
	`, userID, email); err != nil {
		return err
	}

	return nil
}
