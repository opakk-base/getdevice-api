package middleware

import "net/http"

// CORS wraps an http.Handler and adds CORS headers to all responses.
// It also handles OPTIONS preflight requests.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The API binds to localhost by default; keep CORS restricted to local app origins.
		origin := r.Header.Get("Origin")
		if origin == "http://localhost" || origin == "http://127.0.0.1" || origin == "wails://wails" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Support Private Network Access (PNA / CORS-RFC1918) preflight
		// so browsers allow requests from public/remote origins to localhost
		if r.Header.Get("Access-Control-Request-Private-Network") == "true" {
			w.Header().Set("Access-Control-Allow-Private-Network", "true")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
