package service

import (
	"context"
	"errors"
	"local/client"
	"local/model"
	"local/repository"
	"strings"
)

type Service interface {
	GetUsers(ctx context.Context) ([]*model.User, error)
	Uppercase(ctx context.Context, s string) (string, error)
	Count(string) int
}

type Params struct {
	Repo   repository.RepositoryInterface
	Client *client.Client
}

type service struct {
	repo   repository.RepositoryInterface
	client *client.Client
}

func (svc *service) GetUsers(ctx context.Context) ([]*model.User, error) {
	users := svc.repo.User().QueryMany(ctx, &model.User{})
	return users, nil
}

func (svc *service) Uppercase(ctx context.Context, s string) (string, error) {
	if s == "" {
		return "", errors.New("empty string")
	}

	return strings.ToUpper(s), nil
}

func (svc *service) Count(s string) int {
	return len(s)
}

func NewService(params *Params) Service {
	return &service{
		repo:   params.Repo,
		client: params.Client,
	}
}
