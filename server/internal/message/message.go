package message

import (
	"context"
	"server/internal/group"
	"server/internal/user"
	"time"
)

type Message struct {
	ID          int       `json:"id" bson:"id"`
	Type        string    `json:"type" bson:"type"`
	SenderID    int       `json:"sender_id" bson:"sender_id"`
	GroupID     *int      `json:"group_id" bson:"group_id"`
	RecipientID *int      `json:"recipient_id" bson:"recipient_id"`
	Content     string    `json:"content" bson:"content"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type SendMessageReq struct {
	Type        string `json:"type" bson:"type"`
	SenderID    int    `json:"sender_id" bson:"sender_id"`
	GroupID     int    `json:"group_id" bson:"group_id"`
	RecipientID int    `json:"recipient_id" bson:"recipient_id"`
	Content     string `json:"content" bson:"content"`
}

type InsertGroupMessageReq struct {
}

type InsertPrivateMessageReg struct{}

type GetPrivateMessagesReq struct {
	SenderID    string `json:"sender_id" bson:"sender_id"`
	RecipientID string `json:"recipient_id" bson:"recipient_id"`
}

type PrivateMessageRes struct {
	ID        int          `json:"id" bson:"id"`
	Type      string       `json:"type" bson:"type"`
	Sender    user.UserRes `json:"sender" bson:"sender"`
	Recipient user.UserRes `json:"recipient" bson:"recipient"`
	Content   string       `json:"content" bson:"content"`
	CreatedAt time.Time    `json:"created_at" bson:"created_at"`
}

type GroupMessageRes struct {
	ID        int            `json:"id" bson:"id"`
	Type      string         `json:"type" bson:"type"`
	Sender    user.UserRes   `json:"sender" bson:"sender"`
	Group     group.GroupRes `json:"group" bson:"group"`
	Content   string         `json:"content" bson:"content"`
	CreatedAt time.Time      `json:"created_at" bson:"created_at"`
}

type GetGroupMessagesReq struct {
	GroupID string `json:"group_id" bson:"group_id"`
}

type Repository interface {
	InsertMessage(c context.Context, msg *Message) (*Message, error)
	GetPrivateMessages(c context.Context, SenderID, RecipientID string) (*[]PrivateMessageRes, error)
	GetGroupMessages(c context.Context, groupID string) (*[]GroupMessageRes, error)
}

type Service interface {
	SendMessage(c context.Context, req *SendMessageReq) (*Message, error)
	GetPrivateMessages(c context.Context, req *GetPrivateMessagesReq) (*[]PrivateMessageRes, error)
	GetGroupMessages(c context.Context, req *GetGroupMessagesReq) (*[]GroupMessageRes, error)
}
