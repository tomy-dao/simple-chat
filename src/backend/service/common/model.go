package common

import (
	"local/client"
	"local/repository"
)

type Params struct {
	Repo   repository.RepositoryInterface
	Client *client.Client
}