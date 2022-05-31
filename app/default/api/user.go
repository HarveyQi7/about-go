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

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

func (c *UserController) Register(router *gin.RouterGroup) *gin.RouterGroup {

	router.POST("/v1/user/auth/token", func(ctx *gin.Context) {
		authType := ctx.GetHeader(auth.AUTHORIZATION_TYPE)
		var signInfo struct {
			PhonePrefix string `form:"phonePrefix" binding:"required"`
			PhoneNumber string `form:"phoneNumber" binding:"required"`
			Password    string `form:"password" binding:"required"`
		}
		if err := ctx.ShouldBind(&signInfo); err != nil {
			ctx.JSON(errs.ERR2000.Ret(err.Error()))
			return
		}
		conn, disconn := mongodb.Default()
		defer disconn()
		db := conn.Database(config.Get().Datasources.Get().Database)
		userCol := db.Collection(models.User{}.ColName())
		user, exist := mongodb.RetD[models.User](userCol.FindOne(context.TODO(), bson.M{
			"phonePrefix": signInfo.PhonePrefix,
			"phoneNumber": signInfo.PhoneNumber,
		}))
		if !exist {
			ctx.JSON(errs.ERR2003.Ret("wrong username or password"))
			return
		}
		if authType == "" {
			ctx.JSON(errs.ERR2001.Ret("no authorization type"))
		} else if authType == auth.JWT {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInfo.Password)); err != nil {
				ctx.JSON(errs.ERR2003.Ret("wrong username or password"))
				return
			}
			claims := &auth.UserJwtClaims{
				UserId:      user.Id,
				Username:    user.Username,
				PhonePrefix: user.PhonePrefix,
				PhoneNumber: user.PhoneNumber,
				Authorities: user.Authorities,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
					IssuedAt:  time.Now().Unix(),
					Issuer:    "about-go",
					Subject:   auth.ACCESS_TOKEN,
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			if tokenStr, err := token.SignedString([]byte(config.Get().JwtSecretKey)); err != nil {
				ctx.JSON(errs.ERR2004.Ret(err.Error()))
			} else {
				ctx.Header(auth.ACCESS_TOKEN, tokenStr)
				ctx.Status(http.StatusOK)
			}
		} else {
			ctx.JSON(errs.ERR2001.Ret("unknown authorization type"))
		}
	})

	router.GET(
		"/v1/user/auth/token/metadata",
		auth.JwtAuthorizer(),
		func(ctx *gin.Context) {
			metadate, ok := auth.GetMetaData[auth.UserJwtClaims](ctx)
			if ok {
				ctx.JSON(http.StatusOK, metadate)
			} else {
				ctx.JSON(http.StatusOK, nil)
			}
		},
	)

	router.GET(
		"/test/myself/:id",
		auth.JwtAuthorizer(func(claims auth.UserJwtClaims, args args.ReqData) bool {
			return claims.UserId.Hex() == args.Params["id"]
		}),
		func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "ok")
		},
	)

	router.POST(
		"/v1/user",
		auth.JwtAuthorizer(func(claims auth.UserJwtClaims, args args.ReqData) bool {
			return (claims.Authorities.Role == "ADMIN" && claims.Authorities.Level >= 3) ||
				claims.Authorities.UserMgm >= 3
		}),
		func(ctx *gin.Context) {
			var body struct {
				Username    string `json:"username" binding:"required"`
				PhonePrefix string `json:"phonePrefix" binding:"required"`
				PhoneNumber string `json:"phoneNumber" binding:"required"`
				Password    string `json:"password" binding:"required"`
				Authorities struct {
					AuthId  primitive.ObjectID `json:"authId"`
					Role    string             `json:"role" binding:"required,oneof=ADMIN STAFF"`
					Level   int64              `json:"level" binding:"min=0"`
					UserMgm int64              `json:"userMgm" binding:"min=0"`
					AddAuth bool               `json:"addAuth"`
					UpdAuth bool               `json:"updAuth"`
					DelAuth bool               `json:"delAuth"`
				} `json:"authorities" binding:"required"`
			}
			if err := ctx.ShouldBindJSON(&body); err != nil {
				ctx.JSON(errs.ERR2000.Ret(err.Error()))
				return
			}
			conn, disconn := mongodb.Default()
			defer disconn()
			db := conn.Database(config.Get().Datasources.Get().Database)
			if !body.Authorities.AuthId.IsZero() {
				authCol := db.Collection(models.Authority{}.ColName())
				count, err := authCol.CountDocuments(context.TODO(), bson.M{"_id": body.Authorities.AuthId})
				if count > 0 || err != nil {
					ctx.JSON(errs.ERR2002.Ret("authority not found"))
					return
				}
			}
			userCol := db.Collection(models.User{}.ColName())
			user := models.User{
				Username:    body.Username,
				PhonePrefix: body.PhonePrefix,
				PhoneNumber: body.PhoneNumber,
				Password:    body.Password,
				Authorities: models.UserAuthority(body.Authorities),
				CreatedAt:   time.Now(),
			}
			ctx.JSON(http.StatusOK, mongodb.RetW(userCol.InsertOne(context.TODO(), user)))
		},
	)

	router.DELETE(
		"/v1/user",
		auth.JwtAuthorizer(func(claims auth.UserJwtClaims, args args.ReqData) bool {
			return (claims.Authorities.Role == "ADMIN" && claims.Authorities.Level >= 3) ||
				claims.Authorities.UserMgm >= 3
		}),
		func(ctx *gin.Context) {
			var body struct {
				Ids []string `json:"ids" binding:"required,min=1"`
			}
			if err := ctx.ShouldBindJSON(&body); err != nil {
				ctx.JSON(errs.ERR2000.Ret(err.Error()))
				return
			}
			objIds := []primitive.ObjectID{}
			for _, id := range body.Ids {
				objId, _ := primitive.ObjectIDFromHex(id)
				objIds = append(objIds, objId)
			}
			conn, disconn := mongodb.Default()
			defer disconn()
			db := conn.Database(config.Get().Datasources.Get().Database)
			userCol := db.Collection(models.User{}.ColName())
			filter := bson.M{"_id": bson.M{"$in": objIds}}
			result := mongodb.RetW(userCol.DeleteMany(context.TODO(), filter))
			ctx.JSON(http.StatusOK, result)
		},
	)

	router.GET(
		"/v1/user",
		auth.JwtAuthorizer(func(claims auth.UserJwtClaims, args args.ReqData) bool {
			return claims.Authorities.UserMgm >= 1
		}),
		func(ctx *gin.Context) {
			conn, disconn := mongodb.Default()
			defer disconn()
			db := conn.Database(config.Get().Datasources.Get().Database)
			userCol := db.Collection(models.User{}.ColName())
			filter := bson.M{}
			result := mongodb.RetA[[]bson.M](userCol.Aggregate(context.TODO(), filter))
			ctx.JSON(http.StatusOK, result)
		},
	)

	router.GET(
		"/v1/user/:id",
		auth.JwtAuthorizer(func(claims auth.UserJwtClaims, args args.ReqData) bool {
			return claims.Authorities.UserMgm >= 1
		}),
		func(ctx *gin.Context) {
			id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
			conn, disconn := mongodb.Default()
			defer disconn()
			db := conn.Database(config.Get().Datasources.Get().Database)
			userCol := db.Collection(models.User{}.ColName())
			filter := bson.M{"_id": id}
			result, _ := mongodb.RetD[bson.M](userCol.FindOne(context.TODO(), filter))
			ctx.JSON(http.StatusOK, result)
		},
	)

	return router
}
