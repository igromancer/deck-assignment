package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/igromancer/deck-assignment/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyHeader := c.Request.Header.Get("x-api-key")
		if apiKeyHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "x-api-key is required"})
			return
		}
		s := strings.Split(apiKeyHeader, ".")
		if len(s) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "x-api-key is formatted incorrectly"})
			return
		}
		publicId, secret := s[0], s[1]
		apiKeyRepo, err := data.NewApiKeyrepository()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		dbKey, err := apiKeyRepo.Get(c.Request.Context(), publicId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(dbKey.HashedSecret), []byte(secret))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid x-api-key"})
			return
		}
		c.Set("api_key_id", dbKey.ID)
		c.Next()
	}
}
