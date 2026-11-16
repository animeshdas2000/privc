package main

import (
	"encoding/base64"
	"strings"
	"testing"
)

var testKey = []byte("0123456789abcdef0123456789abcdef") // 32 bytes -> AES-256

func TestEncryptDecryptRoundTrip(t *testing.T) {
	cases := []string{
		"",
		"4111",
		"411111111111111",       // 15 bytes
		"4111111111111111",      // 16 bytes: crashed the old CBC code
		"41111111111111111",     // 17 bytes
		strings.Repeat("9", 32), // two full blocks
		"José Ñandú 💳",
	}
	for _, in := range cases {
		token, err := Encrypt(testKey, in)
		if err != nil {
			t.Fatalf("encrypt %q: %v", in, err)
		}
		out, err := Decrypt(testKey, token)
		if err != nil {
			t.Fatalf("decrypt %q: %v", in, err)
		}
		if out != in {
			t.Fatalf("round trip: got %q, want %q", out, in)
		}
	}
}

func TestEncryptIsNonDeterministic(t *testing.T) {
	a, _ := Encrypt(testKey, "4111111111111111")
	b, _ := Encrypt(testKey, "4111111111111111")
	if a == b {
		t.Fatal("same plaintext produced the same token twice")
	}
	// Tokens for values sharing a prefix must not share ciphertext beyond
	// the version byte (the old fixed-IV CBC leaked this).
	c, _ := Encrypt(testKey, "411111111111111122")
	ra, _ := base64.RawURLEncoding.DecodeString(a)
	rc, _ := base64.RawURLEncoding.DecodeString(c)
	if string(ra[1:17]) == string(rc[1:17]) {
		t.Fatal("tokens for values with a shared prefix share a nonce")
	}
}

func TestDecryptRejectsBadTokens(t *testing.T) {
	good, _ := Encrypt(testKey, "4111111111111111")
	raw, _ := base64.RawURLEncoding.DecodeString(good)

	flip := func(i int) string {
		b := append([]byte(nil), raw...)
		b[i] ^= 0x01
		return base64.RawURLEncoding.EncodeToString(b)
	}

	cases := map[string]struct {
		key   []byte
		token string
	}{
		"not base64":      {testKey, "%%%"},
		"empty":           {testKey, ""},
		"too short":       {testKey, base64.RawURLEncoding.EncodeToString(raw[:20])},
		"wrong key":       {[]byte("ffffffffffffffffffffffffffffffff"), good},
		"bad version":     {testKey, flip(0)},
		"tampered nonce":  {testKey, flip(5)},
		"tampered body":   {testKey, flip(len(raw) - 20)},
		"tampered tag":    {testKey, flip(len(raw) - 1)},
		"legacy CBC form": {testKey, "RF6qeMtPh8CuCcieDe2ZWw=="},
	}
	for name, tc := range cases {
		if _, err := Decrypt(tc.key, tc.token); err != errInvalidToken {
			t.Errorf("%s: got err %v, want errInvalidToken", name, err)
		}
	}
}
