package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/animeshdas2000/privc/config"
	"github.com/gin-gonic/gin"
)

const testAPIKey = "test-api-key-0123456789abcdef0123456789"

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return newRouter(&config.Config{
		EncryptionKey: testKey,
		APIKeys:       []string{"other-key-0123456789abcdef0123456789ab", testAPIKey},
	})
}

type apiResponse struct {
	Success      bool              `json:"success"`
	ErrorMessage string            `json:"error_message"`
	Data         map[string]string `json:"data"`
}

func do(t *testing.T, r http.Handler, path, auth string, body any) (int, apiResponse) {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var resp apiResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return w.Code, resp
}

func TestHealthIsPublic(t *testing.T) {
	w := httptest.NewRecorder()
	testRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200", w.Code)
	}
}

func TestEndpointsRequireAPIKey(t *testing.T) {
	r := testRouter()
	body := map[string]any{"data": map[string]string{"field1": "4111"}}
	for _, path := range []string{"/tokenize", "/detokenize"} {
		for name, auth := range map[string]string{
			"missing":      "",
			"wrong key":    "Bearer not-the-right-key-0123456789abcdef",
			"wrong scheme": "Basic " + testAPIKey,
			"empty bearer": "Bearer ",
		} {
			if code, _ := do(t, r, path, auth, body); code != http.StatusUnauthorized {
				t.Errorf("%s %s: got %d, want 401", path, name, code)
			}
		}
	}
}

func TestTokenizeDetokenizeOverHTTP(t *testing.T) {
	r := testRouter()
	auth := "Bearer " + testAPIKey
	in := map[string]string{"field1": "4111111111111111", "field2": "jane@example.com"}

	code, tok := do(t, r, "/tokenize", auth, map[string]any{"data": in})
	if code != http.StatusCreated || !tok.Success {
		t.Fatalf("tokenize: %d %+v", code, tok)
	}
	for k, v := range in {
		if tok.Data[k] == "" || tok.Data[k] == v {
			t.Fatalf("tokenize: field %s not tokenized: %q", k, tok.Data[k])
		}
	}

	code, detok := do(t, r, "/detokenize", auth, map[string]any{"data": tok.Data})
	if code != http.StatusOK || !detok.Success {
		t.Fatalf("detokenize: %d %+v", code, detok)
	}
	for k, v := range in {
		if detok.Data[k] != v {
			t.Fatalf("detokenize %s: got %q, want %q", k, detok.Data[k], v)
		}
	}
}

func TestDetokenizeBadTokenIs400(t *testing.T) {
	code, resp := do(t, testRouter(), "/detokenize", "Bearer "+testAPIKey,
		map[string]any{"data": map[string]string{"field1": "RF6qeMtPh8CuCcieDe2ZWw=="}})
	if code != http.StatusBadRequest || resp.Success {
		t.Fatalf("got %d %+v, want 400", code, resp)
	}
}
