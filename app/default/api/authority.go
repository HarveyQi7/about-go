package api

import (
	"about-go/app/default/models"
	"about-go/config"
	"about-go/pkg/args"
	"about-go/pkg/auth"
	"about-go/pkg/errs"
	"about-go/pkg/mongodb"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthorityController struct{}

func (c *AuthorityController) Register(router *gin.RouterGroup) *gin.RouterGroup {

	router.POST(
		"/v1/authority",
		auth.JwtAuthorizer(func(claims auth.UserJwtClaims, args args.ReqData) bool {
			return (claims.Authorities.Role == "ADMIN" && claims.Authorities.Level >= 3) ||
				claims.Authorities.AddAuth
		}),
		func(ctx *gin.Context) {
			var body struct {
				Role    string `json:"role" binding:"required,oneof=ADMIN STAFF"`
				Level   int64  `json:"level" binding:"min=0"`
				UserMgm int64  `json:"userMgm" binding:"min=0"`
				AddAuth bool   `json:"addAuth"`
				UpdAuth bool   `json:"updAuth"`
				DelAuth bool   `json:"delAuth"`
			}
			if err := ctx.ShouldBindJSON(&body); err != nil {
				ctx.JSON(errs.ERR2000.Ret(err.Error()))
				return
			}
			conn, disconn := mongodb.Default()
			defer disconn()
			db := conn.Database(config.Get().Datasources.Get().Database)
			authorityCol := db.Collection(models.Authority{}.ColName())
			authority := models.Authority{
				Role:      body.Role,
				Level:     body.Level,
				UserMgm:   body.UserMgm,
				AddAuth:   body.AddAuth,
				UpdAuth:   body.UpdAuth,
				DelAuth:   body.DelAuth,
				CreatedAt: time.Now(),
			}
			ctx.JSON(http.StatusOK, mongodb.RetW(authorityCol.InsertOne(context.TODO(), authority)))
		},
	)

	router.GET(
		"/v1/authority",
		auth.JwtAuthorizer(),
		func(ctx *gin.Context) {
			conn, disconn := mongodb.Default()
			defer disconn()
			db := conn.Database(config.Get().Datasources.Get().Database)
			authorityCol := db.Collection(models.Authority{}.ColName())
			filter := bson.M{}
			result := mongodb.RetA[[]bson.M](authorityCol.Aggregate(context.TODO(), filter))
			ctx.JSON(http.StatusOK, result)
		},
	)

	return router
}
