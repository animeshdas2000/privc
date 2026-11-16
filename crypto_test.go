package main

import (
	"encoding/base64"
	"strings"
	"testing"
)

const (
	testKey = "0123456789abcdef0123456789abcdef" // 32 bytes -> AES-256
	testIV  = "fedcba9876543210"                 // 16 bytes
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	cases := []string{
		"",
		"4111",
		"411111111111111",       // 15 bytes
		"4111111111111111",      // 16 bytes: used to panic on decrypt
		"41111111111111111",     // 17 bytes
		strings.Repeat("9", 32), // two full blocks
		"José Ñandú 💳",
	}
	for _, in := range cases {
		token, err := AESEncrypt(testKey, testIV, in)
		if err != nil {
			t.Fatalf("encrypt %q: %v", in, err)
		}
		raw, _ := base64.StdEncoding.DecodeString(token)
		if len(raw) == 0 || len(raw)%16 != 0 {
			t.Fatalf("encrypt %q: ciphertext length %d not a positive multiple of 16", in, len(raw))
		}
		out, err := AESDecrypt(testKey, testIV, token)
		if err != nil {
			t.Fatalf("decrypt %q: %v", in, err)
		}
		if string(out) != in {
			t.Fatalf("round trip: got %q, want %q", out, in)
		}
	}
}

func TestDecryptRejectsBadTokens(t *testing.T) {
	good, _ := AESEncrypt(testKey, testIV, "4111111111111111")
	otherKey := "ffffffffffffffffffffffffffffffff"

	bad := map[string]func() (string, string){
		"not base64":   func() (string, string) { return testKey, "%%%" },
		"empty":        func() (string, string) { return testKey, "" },
		"wrong length": func() (string, string) { return testKey, base64.StdEncoding.EncodeToString([]byte("short")) },
		"wrong key":    func() (string, string) { return otherKey, good },
		"no padding":   func() (string, string) { return testKey, legacyUnpaddedToken(t) },
	}
	for name, mk := range bad {
		key, token := mk()
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("%s: panicked: %v", name, r)
				}
			}()
			if _, err := AESDecrypt(key, testIV, token); err == nil && name != "wrong key" {
				t.Fatalf("%s: expected error", name)
			}
		}()
	}
}

// legacyUnpaddedToken builds a token the way the old code did for 16-byte
// inputs (no padding block), which is what used to crash detokenize.
func legacyUnpaddedToken(t *testing.T) string {
	t.Helper()
	full, err := AESEncrypt(testKey, testIV, "4111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := base64.StdEncoding.DecodeString(full)
	return base64.StdEncoding.EncodeToString(raw[:16])
}
