package recovery

import (
	"net/http"

	l "github.com/duartqx/ttgowebdd/application/middleware/logger"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				rl := l.NewRequestLoggerBuilder().
					SetMethod(r.Method).
					SetStatus(http.StatusInternalServerError).
					SetPath(r.URL.Path)

				rl.LogErr(err)

				w.WriteHeader(rl.GetStatus())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
