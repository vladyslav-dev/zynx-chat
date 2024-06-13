package user

import (
	"context"
	"database/sql"
	"log"
)

type DBTX interface {
	// ExecContext(ctx context.Context, query string, args ...interface{} (sql.Result, error))
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)

	log.Println("Error retrieving user:", err)

	if err != nil {
		return &User{}, err
	}

	log.Println("User retrieved:", u)

	return &u, nil
}

func (r *repository) GetAllUsers(ctx context.Context) (*[]User, error) {
	var us []User

	query := "SELECT id, username, email, password FROM users"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User

		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
			return nil, err
		}

		us = append(us, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &us, nil
}
