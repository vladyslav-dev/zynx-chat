package group

import (
	"context"
	"time"
)

type Group struct {
	ID        int64     `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type GroupMember struct {
	GroupId  string    `json:"group_id" bson:"group_id"`
	UserId   string    `json:"user_id" bson:"user_id"`
	JoinedAt time.Time `json:"joined_at" bson:"joined_at"`
}

type CreateGroupReq struct {
	Name string `json:"name" bson:"name"`
}

type JoinGroupReq struct {
	UserId  int `json:"user_id" bson:"user_id`
	GroupId int `json:"group_id" bson:"group_id`
}

type GroupRes struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type Repository interface {
	CreateGroup(ctx context.Context, group *Group) (*Group, error)
	GetAllGroups(ctx context.Context) (*[]Group, error)
	GroupExists(ctx context.Context, groupID int) (bool, error)
	JoinGroup(ctx context.Context, req *JoinGroupReq) (*GroupMember, error)
	GetGroupById(ctx context.Context, groupID int) (*Group, error)
	// AddUserToGroup(ctx context.Context, user *user.User) (*GroupMember, error)
	// GetGroupsByUserId(ctx context.Context, userId *string) (*[]Group, error)
}

type Service interface {
	CreateGroup(ctx context.Context, req *CreateGroupReq) (*GroupRes, error)
	GetAllGroups(ctx context.Context) (*[]Group, error)
	JoinGroup(ctx context.Context, req *JoinGroupReq) (*GroupMember, error)
	GetGroupById(ctx context.Context, groupID int) (*Group, error)
	// AddUserToGroup(c context.Context, req *user.User) (*GroupMember, error)
	// GetGroupsByUserId(ctx context.Context, req *string) (*[]Group, error)
}
