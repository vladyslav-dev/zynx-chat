package user

import (
	"context"
	"server/db"

	"github.com/lib/pq"
)

type repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, phone, password) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Phone, user.Password).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	u := User{}

	query := "SELECT id, phone, username, password, created_at FROM users WHERE phone = $1"

	err := r.db.QueryRowContext(ctx, query, phone).Scan(&u.ID, &u.Phone, &u.Username, &u.Password, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repository) GetUserByID(ctx context.Context, id int) (*User, error) {
	u := User{}

	query := "SELECT id, phone, username, password, created_at FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Phone, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repository) GetUsersByIDs(ctx context.Context, usersIDs []int) (*[]BaseUserResponse, error) {
	us := []BaseUserResponse{}

	query := `
		SELECT id, username, phone 
		FROM users 
		WHERE id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(usersIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u BaseUserResponse
		if err := rows.Scan(&u.ID, &u.Username, &u.Phone); err != nil {
			return nil, err
		}
		us = append(us, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &us, nil
}

func (r *repository) GetUsersByGroupID(ctx context.Context, groupID int) (*[]BaseUserResponse, error) {
	users := []BaseUserResponse{}

	query := `
		SELECT u.id, u.username, u.phone
		FROM users u
		INNER JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u BaseUserResponse

		if err := rows.Scan(&u.ID, &u.Username, &u.Phone); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *repository) GetAllUsers(ctx context.Context) (*[]BaseUserResponse, error) {
	us := []BaseUserResponse{}

	query := "SELECT id, username, phone FROM users"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u BaseUserResponse

		if err := rows.Scan(&u.ID, &u.Username, &u.Phone); err != nil {
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
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(c, query, sessionReq.UserID, sessionReq.RefreshToken, sessionReq.UserAgent, sessionReq.IPAddress).Scan(&res.ID, &res.CreatedAt)
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
