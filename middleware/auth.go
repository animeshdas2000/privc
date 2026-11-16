// Package middleware holds HTTP middleware for the vault API.
package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/animeshdas2000/privc/utils"
	"github.com/gin-gonic/gin"
)

// RequireAPIKey rejects requests that do not carry one of apiKeys as
// "Authorization: Bearer <key>".
func RequireAPIKey(apiKeys []string) gin.HandlerFunc {
	// Compare fixed-length digests so the comparison time does not reveal
	// the key length or how many leading bytes matched.
	digests := make([][32]byte, len(apiKeys))
	for i, k := range apiKeys {
		digests[i] = sha256.Sum256([]byte(k))
	}

	return func(c *gin.Context) {
		presented, ok := bearerToken(c.GetHeader("Authorization"))
		if ok {
			got := sha256.Sum256([]byte(presented))
			match := 0
			for i := range digests {
				match |= subtle.ConstantTimeCompare(got[:], digests[i][:])
			}
			if match == 1 {
				c.Next()
				return
			}
		}
		c.Header("WWW-Authenticate", `Bearer realm="privc"`)
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
			Success:      false,
			ErrorMessage: "missing or invalid API key",
		})
	}
}

func bearerToken(header string) (string, bool) {
	scheme, token, found := strings.Cut(header, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	return token, token != ""
}
