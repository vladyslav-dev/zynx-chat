package group

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

func (r *repository) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	var lastInsertId int
	query := "INSERT INTO groups(name) VALUES ($1) returning id"
	err := r.db.QueryRowContext(ctx, query, group.Name).Scan(&lastInsertId)
	if err != nil {
		return &Group{}, err
	}

	return group, nil
}

func (r *repository) GetGroupById(ctx context.Context, groupID int) (*Group, error) {
	var g Group
	query := "SELECT id, name, created_at FROM groups WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(&g.ID, &g.Name, &g.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *repository) GetAllGroups(ctx context.Context) (*[]Group, error) {
	groups := []Group{}

	query := "SELECT * FROM groups"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var g Group

		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
			return nil, err
		}

		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &groups, nil
}

func (r *repository) GroupExists(ctx context.Context, groupID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM groups WHERE id = $1)"
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(&exists)
	return exists, err
}

func (r *repository) JoinGroup(ctx context.Context, g *JoinGroupReq) (*GroupMember, error) {
	var gm GroupMember
	query := "INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) RETURNING group_id, user_id, joined_at"
	err := r.db.QueryRowContext(ctx, query, g.GroupId, g.UserId).Scan(&gm.GroupId, &gm.UserId, &gm.JoinedAt)
	if err != nil {
		return nil, err
	}
	return &gm, nil
}

// func (r *repository) AddUserToGroup(ctx context.Context, user *user.User) (*GroupMember, error) {

// }

// func (r *repository) GetGroupsByUserId(ctx context.Context, userId *string) (*[]Group, error) {

// }
