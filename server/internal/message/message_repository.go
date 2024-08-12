package message

import (
	"context"
	"server/db"
)

type repository struct {
	db db.DBTX
}

func NewRepository(db db.DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) InsertMessage(c context.Context, msg *Message) (*Message, error) {
	var newmsg Message
	var query string
	var err error

	if msg.Type == "group" {
		query = "INSERT INTO messages (type, sender_id, group_id, recipient_id, content) VALUES ($1, $2, $3, NULL, $4) RETURNING id, type, sender_id, group_id, content, created_at"
		err = r.db.QueryRowContext(c, query, msg.Type, msg.SenderID, msg.GroupID, msg.Content).Scan(&newmsg.ID, &newmsg.Type, &newmsg.SenderID, &newmsg.GroupID, &newmsg.Content, &newmsg.CreatedAt)
	} else if msg.Type == "private" {
		query = "INSERT INTO messages (type, sender_id, group_id, recipient_id, content) VALUES ($1, $2, NULL, $3, $4) RETURNING id, type, sender_id, recipient_id, content, created_at"
		err = r.db.QueryRowContext(c, query, msg.Type, msg.SenderID, msg.RecipientID, msg.Content).Scan(&newmsg.ID, &newmsg.Type, &newmsg.SenderID, &newmsg.RecipientID, &newmsg.Content, &newmsg.CreatedAt)
	}

	if err != nil {
		return nil, err
	}

	return &newmsg, nil
}

func (r *repository) GetPrivateMessages(c context.Context, SenderID, RecipientID int) (*[]PrivateMessageRes, error) {
	msgs := []PrivateMessageRes{}

	query := `
		SELECT id, type, sender_id, recipient_id, content, created_at
		FROM messages
		WHERE (recipient_id = $2 AND sender_id = $1) OR (recipient_id = $1 AND sender_id = $2) AND type = 'private'
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(c, query, SenderID, RecipientID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var msg PrivateMessageRes

		if err := rows.Scan(&msg.ID, &msg.Type, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return &msgs, nil
}

func (r *repository) GetGroupMessages(c context.Context, groupID int) (*[]GroupMessageRes, error) {
	msgs := []GroupMessageRes{}

	query := `
		SELECT id, type, sender_id, group_id, content, created_at
		FROM messages
		WHERE group_id = $1 AND type = 'group'
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(c, query, groupID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var msg GroupMessageRes

		if err := rows.Scan(&msg.ID, &msg.Type, &msg.SenderID, &msg.GroupID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return &msgs, nil
}
