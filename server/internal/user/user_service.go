package user

import (
	"context"
	"errors"
	"fmt"
	"server/util"
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

func (s *service) Register(c context.Context, req *CreateUserReq) (*BaseUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &BaseUserResponse{
		ID:       int(r.ID),
		Username: r.Username,
		Phone:    r.Phone,
	}

	return res, nil
}

func (s *service) Login(c context.Context, userInfo *UserInfo) (*UserResponseWithTokens, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.GetUserByPhone(ctx, userInfo.Phone)
	if err != nil {
		return &UserResponseWithTokens{}, errors.New("invalid phone number")
	}

	err = util.CheckPassword(userInfo.Password, u.Password)
	if err != nil {
		return &UserResponseWithTokens{}, err
	}

	tokens, err := GenerateTokens(JWTUser{
		ID:       u.ID,
		Username: u.Username,
		Phone:    u.Phone,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.Repository.CreateSession(ctx, SessionReq{
		UserID:       u.ID,
		RefreshToken: string(tokens.RefreshToken),
		UserAgent:    userInfo.UserAgent,
		IPAddress:    userInfo.IPAddress,
	})

	if err != nil {
		return nil, err
	}

	return &UserResponseWithTokens{
		BaseUserResponse: BaseUserResponse{
			ID:       int(u.ID),
			Username: u.Username,
			Phone:    u.Phone,
		},
		AccessToken:  string(tokens.AccessToken),
		RefreshToken: string(tokens.RefreshToken),
	}, nil
}

func (s *service) Logout(c context.Context, token RefreshToken) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	err := s.Repository.DeleteSession(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) isSessionExist(c context.Context, token RefreshToken) bool {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	_, err := s.Repository.GetSession(ctx, token)
	return err == nil
}

func (s *service) ValidateSession(c context.Context, token RefreshToken) (*UserResponseWithTokens, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	baseUserResponse, _ := ValidateRefreshToken(token)

	session, _ := s.Repository.GetSession(ctx, token)

	if (baseUserResponse == nil) || (session == nil) {
		return nil, errors.New("Unauthorized")
	}

	user, err := s.Repository.GetUserByID(ctx, int(session.UserID))
	if err != nil {
		return nil, err
	}

	tokens, err := GenerateTokens(JWTUser{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
	})

	if err != nil {
		fmt.Println("Error generating tokens")
		return nil, err
	}

	sessionReq := SessionReq{
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		IPAddress:    session.IPAddress,
	}

	_, err = s.Repository.UpdateSession(ctx, sessionReq, RefreshToken(tokens.RefreshToken))
	if err != nil {
		fmt.Println("Error updating session")
		return nil, err
	}

	return &UserResponseWithTokens{
		BaseUserResponse: BaseUserResponse{
			ID:       int(user.ID),
			Username: user.Username,
			Phone:    user.Phone,
		},
		AccessToken:  string(tokens.AccessToken),
		RefreshToken: string(tokens.RefreshToken),
	}, nil
}

func (s *service) GetAllUsers(c context.Context) (*[]BaseUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r, err := s.Repository.GetAllUsers(ctx)
	if err != nil {
		return r, err
	}

	return r, nil
}

func (s *service) GetUsersByIDs(c context.Context, usersIDs []int) (*[]BaseUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	us, err := s.Repository.GetUsersByIDs(ctx, usersIDs)
	if err != nil {
		return nil, err
	}

	return us, err
}

func (s *service) GetUsersByGroupID(c context.Context, groupID int) (*[]BaseUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	us, err := s.Repository.GetUsersByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return us, err
}
