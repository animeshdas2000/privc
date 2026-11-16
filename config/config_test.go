package config

import (
	"strings"
	"testing"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

const goodAPIKey = "0123456789abcdef0123456789abcdef"

func TestFromEnvValid(t *testing.T) {
	cfg, err := fromEnv(env(map[string]string{
		"ENCRYPTION_KEY": strings.Repeat("k", 32),
		"API_KEYS":       " " + goodAPIKey + " , " + strings.Repeat("b", 40) + ",",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.EncryptionKey) != 32 || len(cfg.APIKeys) != 2 || cfg.APIKeys[0] != goodAPIKey {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestFromEnvRejects(t *testing.T) {
	cases := map[string]map[string]string{
		"no encryption key":  {"API_KEYS": goodAPIKey},
		"bad key length":     {"ENCRYPTION_KEY": "short", "API_KEYS": goodAPIKey},
		"no api keys":        {"ENCRYPTION_KEY": strings.Repeat("k", 32)},
		"blank api keys":     {"ENCRYPTION_KEY": strings.Repeat("k", 32), "API_KEYS": " , "},
		"api key too short":  {"ENCRYPTION_KEY": strings.Repeat("k", 32), "API_KEYS": "abc"},
		"one short among ok": {"ENCRYPTION_KEY": strings.Repeat("k", 32), "API_KEYS": goodAPIKey + ",abc"},
	}
	for name, m := range cases {
		if _, err := fromEnv(env(m)); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}
