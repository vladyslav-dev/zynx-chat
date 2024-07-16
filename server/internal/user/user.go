package user

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id" bson:"id"`
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
	ID       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginUserRes struct {
	accessToken string
	ID          string `json:"id" bson:"id"`
	Username    string `json:"username" bson:"username"`
	Email       string `json:"email" bson:"email"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) (*[]User, error)
	UserExists(ctx context.Context, userID int) (bool, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*UserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
	GetAllUsers(ctx context.Context) (*[]User, error)
}
