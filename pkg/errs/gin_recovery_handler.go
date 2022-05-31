package errs

import "github.com/gin-gonic/gin"

func DefaultRecoveryHandler() gin.RecoveryFunc {
	return func(ctx *gin.Context, err interface{}) {
		ctx.AbortWithStatusJSON(ERR1000.Ret(err.(error).Error()))
	}
}
