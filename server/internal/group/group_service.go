package group

import (
	"context"
	"errors"
	"server/internal/user"
	"strconv"
	"time"
)

type service struct {
	GroupRepo Repository
	UserRepo  user.Repository
	timeout   time.Duration
}

func NewService(groupRepo Repository, userRepo user.Repository) Service {
	return &service{
		GroupRepo: groupRepo,
		UserRepo:  userRepo,
		timeout:   time.Duration(2) * time.Second,
	}
}

func (s *service) CreateGroup(c context.Context, req *CreateGroupReq) (*GroupRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	g := &Group{
		Name: req.Name,
	}

	r, err := s.GroupRepo.CreateGroup(ctx, g)
	if err != nil {
		return nil, err
	}

	res := &GroupRes{
		ID:   strconv.Itoa(int(r.ID)),
		Name: r.Name,
	}

	return res, nil
}

func (s *service) GetAllGroups(c context.Context) (*[]Group, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.GroupRepo.GetAllGroups(ctx)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *service) JoinGroup(c context.Context, req *JoinGroupReq) (*GroupMember, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	groupExists, err := s.GroupRepo.GroupExists(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	if !groupExists {
		return nil, errors.New("group not found")
	}

	userExists, err := s.UserRepo.UserExists(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	jg := &JoinGroupReq{
		UserId:  req.UserId,
		GroupId: req.GroupId,
	}

	r, err := s.GroupRepo.JoinGroup(ctx, jg)
	if err != nil {
		return nil, err
	}

	return r, nil
}
