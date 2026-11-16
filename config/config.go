// Package config loads and validates service configuration once at startup.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// minAPIKeyLen is the shortest API key accepted. Keys should be random,
// e.g. `openssl rand -hex 32`.
const minAPIKeyLen = 32

// dotenvPaths are optional .env locations. Real environment variables always
// take precedence over values from these files. /etc/secrets is where Render
// mounts Secret Files.
var dotenvPaths = []string{".env", "/etc/secrets/.env"}

type Config struct {
	// EncryptionKey is the raw AES key (16, 24 or 32 bytes).
	EncryptionKey []byte
	// APIKeys are the accepted bearer tokens. Several may be set
	// (comma-separated) so keys can be rotated without downtime.
	APIKeys []string
}

// Load reads configuration from the environment (plus optional .env files)
// and returns an error if anything required is missing or invalid.
func Load() (*Config, error) {
	for _, p := range dotenvPaths {
		if err := godotenv.Load(p); err != nil && !errors.Is(err, os.ErrNotExist) {
			// Do not wrap err: dotenv parse errors quote file contents,
			// which would print secrets into the logs.
			return nil, fmt.Errorf("could not parse %s (check quoting and syntax)", p)
		}
	}
	return fromEnv(os.Getenv)
}

func fromEnv(getenv func(string) string) (*Config, error) {
	key := getenv("ENCRYPTION_KEY")
	switch len(key) {
	case 16, 24, 32:
	case 0:
		return nil, errors.New("ENCRYPTION_KEY is not set")
	default:
		return nil, fmt.Errorf("ENCRYPTION_KEY must be 16, 24 or 32 bytes, got %d", len(key))
	}

	var apiKeys []string
	for _, k := range strings.Split(getenv("API_KEYS"), ",") {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if len(k) < minAPIKeyLen {
			return nil, fmt.Errorf("each API_KEYS entry must be at least %d characters", minAPIKeyLen)
		}
		apiKeys = append(apiKeys, k)
	}
	if len(apiKeys) == 0 {
		return nil, errors.New("API_KEYS is not set")
	}

	return &Config{EncryptionKey: []byte(key), APIKeys: apiKeys}, nil
}
