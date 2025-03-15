package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zainulbr/simple-loan-engine/models/user"
)

const claimRoleKey = "loan.role"

// simple role validation
func RolePermission(role user.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaim(c)

		v, ok := claims[claimRoleKey]
		fmt.Println(v)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, ("Permission dined"))
			return
		}

		vv, ok := v.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, ("Permission dined"))
			return
		}
		fmt.Println(vv, role)

		if vv != string(role) {
			c.AbortWithStatusJSON(http.StatusForbidden, ("Permission dined"))
			return
		}
	}
}
