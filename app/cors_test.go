package app

import (
	"maintenance-system-go/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSCredentials(t *testing.T) {
	origins := []string{"http://localhost:3000", "https://app.example.com"}
	router := gin.New()
	router.Use(corsMiddleware(config.CORSConfig{AllowedOrigins: origins}))
	router.POST("/login", func(c *gin.Context) { c.Status(http.StatusOK) })
	for _, origin := range origins {
		for _, method := range []string{http.MethodPost, http.MethodOptions} {
			req := httptest.NewRequest(method, "/login", nil)
			req.Header.Set("Origin", origin)
			if method == http.MethodOptions {
				req.Header.Set("Access-Control-Request-Method", "POST")
				req.Header.Set("Access-Control-Request-Headers", "content-type")
			}
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			want := http.StatusOK
			if method == http.MethodOptions {
				want = http.StatusNoContent
			}
			if res.Code != want || res.Header().Get("Access-Control-Allow-Origin") != origin || res.Header().Get("Access-Control-Allow-Credentials") != "true" {
				t.Fatalf("%s %s: status=%d headers=%v", method, origin, res.Code, res.Header())
			}
			if method == http.MethodOptions && !strings.Contains(res.Header().Get("Access-Control-Allow-Headers"), "Content-Type") {
				t.Fatal("preflight must allow JSON content type")
			}
		}
	}
	for _, origin := range []string{"https://untrusted.example", "http://localhost:3001"} {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.Header.Set("Origin", origin)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		if res.Code != http.StatusForbidden || res.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("unexpected access for %s: %d", origin, res.Code)
		}
	}
}
