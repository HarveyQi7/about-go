package auth

import (
	"about-go/app/default/models"
	"about-go/config"
	"about-go/pkg/args"
	"about-go/pkg/errs"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const JWT = "JWT"

type UserJwtClaims struct {
	UserId             primitive.ObjectID   `json:"userId"`
	Username           string               `json:"username"`
	PhonePrefix        string               `json:"phonePrefix"`
	PhoneNumber        string               `json:"phoneNumber"`
	Authorities        models.UserAuthority `json:"authorities"`
	jwt.StandardClaims `json:"standardClaims"`
}

type JwtAuthChecker func(claims UserJwtClaims, args args.ReqData) bool

func ParseJwt(tokenStr string) (*jwt.Token, *UserJwtClaims, error) {
	userJwtClaims := &UserJwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, userJwtClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Get().JwtSecretKey), nil
	})
	return token, userJwtClaims, err
}

func JwtAuthorizer(funcs ...JwtAuthChecker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkers := funcs
		if ctx.GetHeader(AUTHORIZATION_TYPE) != JWT {
			ctx.JSON(errs.ERR2003.Ret("authorization type error"))
			ctx.Abort()
			return
		}
		tokenStr := ctx.GetHeader(AUTHORIZATION)
		if tokenStr == "" {
			ctx.JSON(errs.ERR2003.Ret("no token"))
			ctx.Abort()
			return
		}
		jwtTokenPrefix := "Bearer "
		if !strings.HasPrefix(tokenStr, jwtTokenPrefix) {
			ctx.JSON(errs.ERR2005.Ret("token string should start with ' " + jwtTokenPrefix + "'"))
			ctx.Abort()
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, jwtTokenPrefix)
		token, userJwtClaims, err := ParseJwt(tokenStr)
		if err != nil {
			ctx.JSON(errs.ERR2005.Ret(err.Error()))
			ctx.Abort()
			return
		}
		if !token.Valid {
			ctx.JSON(errs.ERR2003.Ret("invalid token"))
			ctx.Abort()
			return
		}
		if len(checkers) > 0 {
			if !checkers[0](*userJwtClaims, args.ParseGinCtx(ctx)) {
				ctx.JSON(errs.ERR2006.Ret("no permission"))
				ctx.Abort()
				return
			}
		}
		ctx.Set(AUTHORIZATION_METADATA, *userJwtClaims)
		ctx.Next()
	}
}
