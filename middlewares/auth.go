package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	libToken "github.com/zainulbr/simple-loan-engine/libs/token"
)

const BEARER_SCHEMA = "Bearer"
const contextClaimKey = "ctx.mw.auth.claim"

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.Contains(authHeader, BEARER_SCHEMA) {
			tokenString = authHeader[len(BEARER_SCHEMA):]
		}

		if tokenString == "" {
			tokenString, _ = c.Cookie("token")
		}

		if tokenString == "" {
			tokenString = c.Query("token")
		}
		tokenString = strings.TrimSpace(tokenString)
		// TBD passsing encrypted via setting
		token, err := libToken.NewService().ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			log.Println(err)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			c.Set(contextClaimKey, claims)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Token invalid")
		}
	}
}

func GetClaim(c *gin.Context) jwt.MapClaims {
	out, ok := c.Get(contextClaimKey)
	if !ok {
		return make(jwt.MapClaims)
	}

	return out.(jwt.MapClaims)
}
