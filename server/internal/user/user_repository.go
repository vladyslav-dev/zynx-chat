package user

import (
	"context"
	"log"
	"server/db"
)

type repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt)

	log.Println("Error retrieving user:", err)

	if err != nil {
		return &User{}, err
	}

	log.Println("User retrieved:", u)

	return &u, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	u := User{}

	query := "SELECT id, email, username, password FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

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

func (r *repository) UserExists(ctx context.Context, userID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)"
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	return exists, err
}

func (r *repository) CreateSession(c context.Context, sessionReq SessionReq) (*SessionRes, error) {
	var res SessionRes
	query := `
		INSERT INTO sessions (user_id, refresh_token, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(c, query, sessionReq.UserID, sessionReq.RefreshToken, sessionReq.UserAgent, sessionReq.IPAddress).Scan(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *repository) GetSession(c context.Context, token RefreshToken) (*Session, error) {
	var session Session
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, created_at
		FROM sessions
		WHERE refresh_token = $1
	`
	err := r.db.QueryRowContext(c, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *repository) UpdateSession(c context.Context, sessionReq SessionReq, newToken RefreshToken) (*Session, error) {
	var session Session
	query := `
		UPDATE sessions
		SET refresh_token = $1
		WHERE user_id = $2 AND refresh_token = $3 AND user_agent = $4 AND ip_address = $5
		RETURNING id, user_id, refresh_token, user_agent, ip_address, created_at
	`
	err := r.db.QueryRowContext(
		c,
		query,
		newToken,
		sessionReq.UserID,
		sessionReq.RefreshToken,
		sessionReq.UserAgent,
		sessionReq.IPAddress,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *repository) DeleteSession(c context.Context, token RefreshToken) error {
	query := `
		DELETE FROM sessions
		WHERE refresh_token = $1
	`
	_, err := r.db.ExecContext(c, query, token)
	if err != nil {
		return err
	}

	return nil
}
