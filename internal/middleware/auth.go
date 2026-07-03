package middleware

import (
	"net/http"
	"pendekin_go/pkg/errs"
	"pendekin_go/pkg/jwt"
	"pendekin_go/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	JWT *jwt.JWT
}

func NewAuthMiddleware(jwt *jwt.JWT) *AuthMiddleware {
	return &AuthMiddleware{JWT: jwt}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ResponseNOK(c, http.StatusUnauthorized, "authorization header is required", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.ResponseNOK(c, http.StatusUnauthorized, "invalid authorization format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			response.ResponseNOK(c, http.StatusUnauthorized, "token is required", nil)
			c.Abort()
			return
		}

		claims, err := m.JWT.GetClaims(tokenString)
		if err != nil {
			switch {
			case errs.IsTokenExpired(err):
				response.ResponseNOK(c, http.StatusUnauthorized, "token expired", nil)
			case errs.IsTokenInvalid(err):
				response.ResponseNOK(c, http.StatusUnauthorized, "invalid token", nil)
			default:
				response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
			}
			c.Abort()
			return
		}

		if claims == nil {
			response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)

		c.Next()

	}
}
