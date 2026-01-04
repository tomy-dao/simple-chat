package common

import (
	"local/client"
	"local/infra/repo"
)

type Params struct {
	Repo   *repo.Repository
	Client *client.Client
}