package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

type Adapter func(http.Handler) http.Handler

// Compatible with http.HandlerFunc
func Auth(re *redis.Client) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(os.Getenv("COOKIE_NAME"))
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			v := c.Value

			// Check if user is authenticated
			_, err = re.Get(v).Result()
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func Log(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
			h.ServeHTTP(w, r)
		})
	}
}
