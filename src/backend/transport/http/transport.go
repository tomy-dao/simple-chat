package httpTransoprt

import (
	"local/endpoint"
	"local/model"

	"github.com/gin-gonic/gin"
)

func MakeHttpTransport(initParams *model.InitParams, endpoints *endpoint.Endpoints) *gin.Engine {
	r := gin.Default()

	SetupMiddleware(r)

	handleRouter(r, endpoints)

	return r
}
