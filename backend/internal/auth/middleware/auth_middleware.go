package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authjwt "github.com/abhinavkumar03/crm-lite/backend/internal/auth/jwt"
)

type AuthMiddleware struct {
	jwt *authjwt.Service
}

func New(jwt *authjwt.Service) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {

	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header missing",
			})
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header",
			})
			return
		}

		token := strings.TrimPrefix(header, prefix)

		claims, err := m.jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
