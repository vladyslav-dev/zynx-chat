package user

import (
	"context"
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

func (s *service) Register(c context.Context, req *CreateUserReq) (*UserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &UserRes{
		ID:       int(r.ID),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) Login(c context.Context, userInfo *UserInfo) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}

	err = util.CheckPassword(userInfo.Password, u.Password)
	if err != nil {
		return &LoginUserRes{}, err
	}

	tokens, err := GenerateTokens(JWTUser{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
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

	return &LoginUserRes{
		Username:     u.Username,
		Email:        u.Email,
		ID:           int(u.ID),
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

func (s *service) ValidateSession(c context.Context, token RefreshToken) (*Tokens, error) {
	_, err := ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	session, err := s.Repository.GetSession(c, token)
	if err != nil {
		return nil, err
	}

	user, err := s.Repository.GetUserByID(c, string(session.UserID))
	if err != nil {
		return nil, err
	}

	tokens, err := GenerateTokens(JWTUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return nil, err
	}

	sessionReq := SessionReq{
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		IPAddress:    session.IPAddress,
	}

	updatedSession, err := s.Repository.UpdateSession(c, sessionReq, RefreshToken(tokens.RefreshToken))
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: RefreshToken(updatedSession.RefreshToken),
	}, nil
}

func (s *service) GetAllUsers(c context.Context) (*[]User, error) {
	r, err := s.Repository.GetAllUsers(c)
	if err != nil {
		return r, err
	}

	return r, nil
}
