package userhttp

import (
	"net/http"
	"strings"

	sharederrors "macabi-back/internal/shared/errors"
	userports "macabi-back/internal/user/application/ports"
	userdomain "macabi-back/internal/user/domain"

	"github.com/gin-gonic/gin"
)

const (
	AuthUserIDKey = "auth_user_id"
	AuthEmailKey  = "auth_email"
	AuthRoleKey   = "auth_role"
)

func AuthMiddleware(tokenPrv userports.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, sharederrors.NewErrorResponse("missing or invalid authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := tokenPrv.Validate(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, sharederrors.NewErrorResponse("invalid or expired token"))
			return
		}

		c.Set(AuthUserIDKey, claims.UserID)
		c.Set(AuthEmailKey, claims.Email)
		c.Set(AuthRoleKey, string(claims.Role))
		c.Next()
	}
}

func RequireRole(roles ...userdomain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr, exists := c.Get(AuthRoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, sharederrors.NewErrorResponse("unauthorized"))
			return
		}

		currentRole := userdomain.Role(roleStr.(string))
		for _, required := range roles {
			if currentRole == required {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, sharederrors.NewErrorResponse("insufficient permissions"))
	}
}
