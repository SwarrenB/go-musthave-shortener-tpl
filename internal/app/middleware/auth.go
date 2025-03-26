package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const cookieName = "userID"

func AuthMiddleware(secretKey string, logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			logger.Debug("JWT cookie not found, generating new one")
			issueNewToken(c, secretKey, logger)
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			logger.Debug("Invalid JWT, issuing new one", zap.Error(err))
			issueNewToken(c, secretKey, logger)
			return
		}

		c.Set(cookieName, claims.UserID)
	}
}

func issueNewToken(c *gin.Context, secretKey string, logger zap.Logger) {
	userID := uuid.New().String()

	tokenString, err := generateJWT(userID, secretKey)
	if err != nil {
		logger.Error("Failed to generate JWT", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    tokenString,
		HttpOnly: true,
	})

	c.Set(cookieName, userID)
}

func generateJWT(userID, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	})

	return token.SignedString([]byte(secret))
}
