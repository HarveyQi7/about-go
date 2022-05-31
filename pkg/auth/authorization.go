package auth

import "github.com/gin-gonic/gin"

const (
	AUTHORIZATION      = "Authorization"
	AUTHORIZATION_TYPE = "Authorization-Type"
	ACCESS_TOKEN       = "Access-Token"
)

const AUTHORIZATION_METADATA = "AUTHORIZATION_METADATA"

func GetMetaData[T interface{}](ctx *gin.Context) (T, bool) {
	var t T
	value, exist := ctx.Get(AUTHORIZATION_METADATA)
	if exist {
		t, exist = value.(T)
	}
	return t, exist
}
