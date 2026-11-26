package broadcastsocket

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebSocketUpgraderChecksConfiguredOrigins(t *testing.T) {
	upgrader := newWebSocketUpgrader([]string{"http://localhost:3000"})

	tests := []struct {
		name        string
		origin      string
		wantAllowed bool
	}{
		{name: "configured origin", origin: "http://localhost:3000", wantAllowed: true},
		{name: "unconfigured origin", origin: "http://localhost:3001", wantAllowed: false},
		{name: "missing origin", wantAllowed: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8000/ws", nil)
			if tt.origin != "" {
				request.Header.Set("Origin", tt.origin)
			}

			if got := upgrader.CheckOrigin(request); got != tt.wantAllowed {
				t.Errorf("CheckOrigin() = %v, want %v", got, tt.wantAllowed)
			}
		})
	}
}
