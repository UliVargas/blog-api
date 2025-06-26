package middleware

import (
	"net/http"
	"strings"

	"github.com/UliVargas/blog-go/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			ctx.Abort()
			return
		}

		bearerToken := strings.SplitN(authHeader, " ", 2)

		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token invalido"})
			ctx.Abort()
			return
		}

		cfg := config.Load()
		if cfg.JWTSECRET == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar sesión"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSECRET), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"].(float64); ok {
				ctx.Set("user_id", uint(userID))
			}
		}

		ctx.Next()
	}
}
