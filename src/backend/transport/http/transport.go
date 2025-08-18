package httpTransoprt

import (
	"local/endpoint"
	"local/model"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MakeHttpTransport(initParams *model.InitParams, endpoints *endpoint.Endpoints) http.Handler {
	r := chi.NewRouter()

	SetupMiddleware(r)

	handleRouter(r, endpoints)

	return r
}
