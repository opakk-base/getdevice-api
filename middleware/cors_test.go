package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	corsHandler := CORS(dummyHandler)

	t.Run("GET request from external origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/getdevice", nil)
		req.Header.Set("Origin", "https://fifo-rsisa.serah-studio.my.id")

		rec := httptest.NewRecorder()
		corsHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		origin := rec.Header().Get("Access-Control-Allow-Origin")
		if origin != "*" {
			t.Errorf("expected Access-Control-Allow-Origin '*', got '%s'", origin)
		}

		if rec.Body.String() != "OK" {
			t.Errorf("expected body 'OK', got '%s'", rec.Body.String())
		}
	})

	t.Run("OPTIONS preflight request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/getdevice", nil)
		req.Header.Set("Origin", "https://fifo-rsisa.serah-studio.my.id")
		req.Header.Set("Access-Control-Request-Method", "GET")

		rec := httptest.NewRecorder()
		corsHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, rec.Code)
		}

		origin := rec.Header().Get("Access-Control-Allow-Origin")
		if origin != "*" {
			t.Errorf("expected Access-Control-Allow-Origin '*', got '%s'", origin)
		}

		methods := rec.Header().Get("Access-Control-Allow-Methods")
		if methods == "" {
			t.Errorf("expected Access-Control-Allow-Methods header to be set")
		}

		headers := rec.Header().Get("Access-Control-Allow-Headers")
		if headers == "" {
			t.Errorf("expected Access-Control-Allow-Headers header to be set")
		}
	})

	t.Run("OPTIONS preflight with Private Network Access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/getdevice", nil)
		req.Header.Set("Origin", "https://fifo-rsisa.serah-studio.my.id")
		req.Header.Set("Access-Control-Request-Private-Network", "true")

		rec := httptest.NewRecorder()
		corsHandler.ServeHTTP(rec, req)

		pna := rec.Header().Get("Access-Control-Allow-Private-Network")
		if pna != "true" {
			t.Errorf("expected Access-Control-Allow-Private-Network 'true', got '%s'", pna)
		}
	})
}
