package middleware

import "net/http"

func APIKeyMiddle(token string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			apiKey := r.Header.Get("X-API-Key")

			if apiKey == "" || apiKey != token {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
