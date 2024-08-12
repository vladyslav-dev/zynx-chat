package message

import (
	"context"
	"encoding/json"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) SendMessage(c context.Context, req *SendMessageReq) (*Message, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	msg := &Message{
		Type:        req.Type,
		SenderID:    req.SenderID,
		GroupID:     &req.GroupID,
		RecipientID: &req.RecipientID,
		Content:     req.Content,
	}

	res, err := s.Repository.InsertMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *service) GetPrivateMessages(c context.Context, req *GetPrivateMessagesReq) (*[]PrivateMessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	msgs, err := s.Repository.GetPrivateMessages(ctx, req.SenderID, req.RecipientID)
	if err != nil {
		return nil, err
	}

	return msgs, nil

}

func (s *service) GetGroupMessages(c context.Context, req *GetGroupMessagesReq) (*[]GroupMessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	msgs, err := s.Repository.GetGroupMessages(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func MarshalMessageJSON(m *MessageWrapper) ([]byte, error) {
	if m.PrivateMsg != nil {
		return json.Marshal(m.PrivateMsg)
	}
	if m.GroupMsg != nil {
		return json.Marshal(m.GroupMsg)
	}
	return nil, nil
}
