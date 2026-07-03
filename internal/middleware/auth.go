package middleware

import (
	"net/http"
	"pendekin_go/pkg/errs"
	"pendekin_go/pkg/jwt"
	"pendekin_go/pkg/response"

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
		tokenString, err := c.Cookie("token")
		if err != nil {
			response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}

		if tokenString == "" {
			response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
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
