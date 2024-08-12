package user

import (
	"context"
	"time"
)

type User struct {
	ID        int       `json:"id" bson:"id"`
	Username  string    `json:"username" bson:"username"`
	Phone     string    `json:"phone" bson:"phone"`
	Password  string    `json:"passowrd" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

/* User Request */

type CreateUserReq struct {
	Username string `json:"username" bson:"username"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"password"`
}

type LoginUserReq struct {
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"password"`
}

/* User Response */
type BaseUserResponse struct {
	ID       int    `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Phone    string `json:"phone" bson:"phone"`
}

type UserResponseWithAccess struct {
	BaseUserResponse
	AccessToken string `json:"access_token" bson:"access_token"`
}

type UserResponseWithRefresh struct {
	BaseUserResponse
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type UserResponseWithTokens struct {
	BaseUserResponse
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

/* JWT */

type JWTUser struct {
	ID       int    `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Phone    string `json:"phone" bson:"phone"`
}

type AccessToken string
type RefreshToken string

type Tokens struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

type UserInfo struct {
	Phone     string `json:"phone" bson:"phone"`
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
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	GetUsersByGroupID(ctx context.Context, groupID int) (*[]BaseUserResponse, error)
	GetAllUsers(ctx context.Context) (*[]BaseUserResponse, error)
	GetUsersByIDs(ctx context.Context, usersIDs []int) (*[]BaseUserResponse, error)
	UserExists(ctx context.Context, userID int) (bool, error)

	CreateSession(c context.Context, sessionReq SessionReq) (*SessionRes, error)
	UpdateSession(c context.Context, sessionReq SessionReq, newToken RefreshToken) (*Session, error)
	DeleteSession(c context.Context, token RefreshToken) error
	GetSession(c context.Context, token RefreshToken) (*Session, error)
}

type Service interface {
	Register(c context.Context, req *CreateUserReq) (*BaseUserResponse, error)
	Login(c context.Context, userInfo *UserInfo) (*UserResponseWithTokens, error)
	Logout(c context.Context, token RefreshToken) error
	ValidateSession(c context.Context, token RefreshToken) (*UserResponseWithTokens, error)
	GetAllUsers(ctx context.Context) (*[]BaseUserResponse, error)
	GetUsersByIDs(ctx context.Context, usersIDs []int) (*[]BaseUserResponse, error)
	GetUsersByGroupID(ctx context.Context, groupID int) (*[]BaseUserResponse, error)
	isSessionExist(ctx context.Context, token RefreshToken) bool
}
