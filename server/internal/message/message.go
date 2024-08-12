package message

import (
	"context"
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

type GetPrivateMessagesReq struct {
	SenderID    int `json:"sender_id" bson:"sender_id"`
	RecipientID int `json:"recipient_id" bson:"recipient_id"`
}

type PrivateMessageRes struct {
	ID          int       `json:"id" bson:"id"`
	Type        string    `json:"type" bson:"type"`
	SenderID    int       `json:"sender_id" bson:"sender_id"`
	RecipientID int       `json:"recipient_id" bson:"recipient_id"`
	Content     string    `json:"content" bson:"content"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type GroupMessageRes struct {
	ID        int       `json:"id" bson:"id"`
	Type      string    `json:"type" bson:"type"`
	SenderID  int       `json:"sender_id" bson:"sender_id"`
	GroupID   int       `json:"group_id" bson:"group_id"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type MessageWrapper struct {
	PrivateMsg *PrivateMessageRes `json:",omitempty,inline"`
	GroupMsg   *GroupMessageRes   `json:",omitempty,inline"`
}

type GetGroupMessagesReq struct {
	GroupID int `json:"group_id" bson:"group_id"`
}

type Repository interface {
	InsertMessage(c context.Context, msg *Message) (*Message, error)
	GetPrivateMessages(c context.Context, SenderID, RecipientID int) (*[]PrivateMessageRes, error)
	GetGroupMessages(c context.Context, groupID int) (*[]GroupMessageRes, error)
}

type Service interface {
	SendMessage(c context.Context, req *SendMessageReq) (*Message, error)
	GetPrivateMessages(c context.Context, req *GetPrivateMessagesReq) (*[]PrivateMessageRes, error)
	GetGroupMessages(c context.Context, req *GetGroupMessagesReq) (*[]GroupMessageRes, error)
}
