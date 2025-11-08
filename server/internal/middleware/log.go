package middleware

import(
	"log"
	"net/http"
	"time"
)

type LoggerMiddleware struct{}

func (logMiddleware *LoggerMiddleware) GetName() string {
	return "Logger"
}

func (logMiddleware *LoggerMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s in %v", logMiddleware.GetName(), r.Method, r.URL.Path, time.Since(start))
	})
}
