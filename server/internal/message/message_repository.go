package message

import (
	"context"
	"server/db"
	"server/internal/group"
	"server/internal/user"
	"strconv"
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

func (r *repository) GetPrivateMessages(c context.Context, SenderID, RecipientID string) (*[]PrivateMessageRes, error) {
	var msgs []PrivateMessageRes

	sid, _ := strconv.Atoi(SenderID)
	rid, _ := strconv.Atoi(RecipientID)

	query := `
		SELECT 
		m.id, m.type, m.content, m.created_at,
		s.id, s.username, s.email,
		r.id, r.username, r.email
		FROM messages m
		JOIN users s ON m.sender_id = s.id
		LEFT JOIN users r ON m.recipient_id = r.id
		WHERE (m.recipient_id = $2 AND m.sender_id = $1) OR (m.recipient_id = $1 AND m.sender_id = $2) AND m.type = 'private'
		ORDER BY m.created_at
	`

	rows, err := r.db.QueryContext(c, query, sid, rid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var msg PrivateMessageRes
		var sender, recipient user.UserRes

		if err := rows.Scan(
			&msg.ID, &msg.Type, &msg.Content, &msg.CreatedAt,
			&sender.ID, &sender.Username, &sender.Email,
			&recipient.ID, &recipient.Username, &recipient.Email,
		); err != nil {
			return nil, err
		}

		msg.Sender = sender
		msg.Recipient = recipient

		msgs = append(msgs, msg)
	}

	return &msgs, nil
}

func (r *repository) GetGroupMessages(c context.Context, groupID string) (*[]GroupMessageRes, error) {
	var msgs []GroupMessageRes

	query := `
		SELECT 
		m.id, m.type, m.content, m.created_at,
		s.id, s.username, s.email,
		g.id, g.name
		FROM messages m
		JOIN users s ON sender_id = s.id
		LEFT JOIN groups g ON group_id = g.id
		WHERE m.group_id = $1 AND m.type = 'group'
		ORDER BY m.created_at
	`

	rows, err := r.db.QueryContext(c, query, groupID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var msg GroupMessageRes
		var sender user.UserRes
		var group group.GroupRes

		if err := rows.Scan(
			&msg.ID, &msg.Type, &msg.Content, &msg.CreatedAt,
			&sender.ID, &sender.Username, &sender.Email,
			&group.ID, &group.Name,
		); err != nil {
			return nil, err
		}

		msg.Sender = sender
		msg.Group = group
		msgs = append(msgs, msg)
	}

	return &msgs, nil

	// var msgs []GroupMessageRes
	// query := "SELECT * FROM messages WHERE group_id = $1 AND type = 'group'"
	// rows, err := r.db.QueryContext(c, query, groupID)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var msg GroupMessageRes
	// 	if err := rows.Scan(&msg.ID, &msg.Type, &msg.SenderID, &msg.GroupID, &msg.RecipientID, &msg.Content, &msg.CreatedAt); err != nil {
	// 		return nil, err
	// 	}

	// 	msgs = append(msgs, msg)
	// }

	// return &msgs, nil
}
