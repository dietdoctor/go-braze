package braze

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type PreferenceCenterEndpoint interface {
	CreateURL(ctx context.Context, req *PreferenceCenterCreateURLRequest) (*PreferenceCenterCreateURLResponse, error)
}

type PreferenceCenterCreateURLRequest struct {
	PreferenceCenterID string
	UserID             string
}

func (r *PreferenceCenterCreateURLRequest) validate() error {
	if r == nil {
		return errors.New("request must not be nil")
	}

	if r.PreferenceCenterID == "" {
		return errors.New("preferences center ID must not be empty")
	}

	if r.UserID == "" {
		return errors.New("user ID must not be empty")
	}

	return nil
}

type PreferenceCenterCreateURLResponse struct {
	URL string
}

type PreferenceCenterService struct {
	client *Client
}

func (s *PreferenceCenterService) CreateURL(ctx context.Context, r *PreferenceCenterCreateURLRequest) (*PreferenceCenterCreateURLResponse, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/preference_center/v1/%s/url/%s", r.PreferenceCenterID, r.UserID)

	req, err := s.client.http.newRequest(http.MethodPost, path, r)
	if err != nil {
		return nil, err
	}

	resp := struct {
		URL string `json:"preference_center_url"`
	}{}

	if err := s.client.http.do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &PreferenceCenterCreateURLResponse{
		URL: resp.URL,
	}, nil
}
