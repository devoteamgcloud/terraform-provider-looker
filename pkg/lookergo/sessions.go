package lookergo

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

const sessionBasePath = "4.0/session"

type SessionsResource interface {
	Get(ctx context.Context) (*Session, *Response, error)
	SetWorkspaceId(ctx context.Context, workspaceId string) (*Session, *Response, error)
	GetCurrentUser(ctx context.Context) (*User, *Response, error)
	GetLoginUserToken(ctx context.Context, userId string) (*oauth2.Token, *Response, error)
}

type SessionsResourceOp struct {
	client *Client
}

var _ SessionsResource = &SessionsResourceOp{}

type Session struct {
	WorkspaceId string `json:"workspace_id"`
	SudoUserId  *int   `json:"sudo_user_id,string,omitempty"`
}

// Get -
func (s *SessionsResourceOp) Get(ctx context.Context) (*Session, *Response, error) {
	return doGet(ctx, s.client, sessionBasePath, new(Session))
}

// SetWorkspaceId -
func (s *SessionsResourceOp) SetWorkspaceId(ctx context.Context, workspaceId string) (session *Session, resp *Response, err error) {
	updateReq := Session{WorkspaceId: workspaceId}
	req, err := s.client.NewRequest(ctx, http.MethodPatch, sessionBasePath, updateReq)
	if err != nil {
		return nil, nil, err
	}

	resp, err = s.client.Do(ctx, req, &session)
	if err != nil {
		return nil, resp, err
	}

	return
}

// GetCurrentUser -
func (s *SessionsResourceOp) GetCurrentUser(ctx context.Context) (*User, *Response, error) {
	return doGet(ctx, s.client, "4.0/user", new(User))
}

// GetLoginUserToken -
func (s *SessionsResourceOp) GetLoginUserToken(ctx context.Context, userId string) (*oauth2.Token, *Response, error) {
	path := fmt.Sprintf("4.0/login/%s", userId)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	token := new(oauth2.Token)
	resp, err := s.client.Do(ctx, req, &token)
	if err != nil {
		return nil, resp, err
	}

	return token, nil, nil
}
