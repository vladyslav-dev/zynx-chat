package user

import (
	"context"
	"time"
)

type User struct {
	ID        int       `json:"id" bson:"id"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"passowrd" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type CreateUserReq struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UserRes struct {
	ID       int    `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

type JWTUser = UserRes

type LoginUserReq struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginUserRes struct {
	ID           int    `json:"id" bson:"id"`
	Username     string `json:"username" bson:"username"`
	Email        string `json:"email" bson:"email"`
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type AccessToken string
type RefreshToken string

type Tokens struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

type UserInfo struct {
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	UserAgent string `json:"user_agent" bson:"user_agent"`
	IPAddress string `json:"ip_address" bson:"ip_address"`
}

type Session struct {
	ID           int       `json:"id" bson:"id"`
	UserID       int       `json:"user_id" bson:"user_id"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	UserAgent    string    `json:"user_agent" bson:"user_agent"`
	IPAddress    string    `json:"ip_address" bson:"ip_address"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

type SessionReq struct {
	UserID       int    `json:"user_id" bson:"user_id"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	UserAgent    string `json:"user_agent" bson:"user_agent"`
	IPAddress    string `json:"ip_address" bson:"ip_address"`
}

type SessionRes struct {
	ID        int       `json:"id" bson:"id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) (*[]User, error)
	UserExists(ctx context.Context, userID int) (bool, error)

	CreateSession(c context.Context, sessionReq SessionReq) (*SessionRes, error)
	UpdateSession(c context.Context, sessionReq SessionReq, newToken RefreshToken) (*Session, error)
	DeleteSession(c context.Context, token RefreshToken) error
	GetSession(c context.Context, token RefreshToken) (*Session, error)
}

type Service interface {
	Register(c context.Context, req *CreateUserReq) (*UserRes, error)
	Login(c context.Context, userInfo *UserInfo) (*LoginUserRes, error)
	Logout(c context.Context, token RefreshToken) error
	ValidateSession(c context.Context, token RefreshToken) (*Tokens, error)
	GetAllUsers(ctx context.Context) (*[]User, error)

	// RefreshToken(ctx context.Context) (*LoginUserReq, error)
}
