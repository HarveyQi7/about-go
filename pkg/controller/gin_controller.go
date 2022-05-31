package controller

import "github.com/gin-gonic/gin"

type GinController interface {
	Register(router *gin.RouterGroup) *gin.RouterGroup
}

func RegisterForGin(router *gin.RouterGroup, controllers ...GinController) *gin.RouterGroup {
	for _, c := range controllers {
		c.Register(router)
	}
	return router
}
