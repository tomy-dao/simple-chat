package common

import (
	"local/client"
	"local/infra/repo"
)

type Params struct {
	Repo   repo.RepositoryInterface
	Client *client.Client
}