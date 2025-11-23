package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("|||||| REQUEST LOG  ||||||")
		log.Println("METHOD: ", r.Method)
		log.Println("URL: ", r.URL)
		log.Println("REMOTE_ADDR: ", r.RemoteAddr)
		log.Println("USER_AGENT: ", r.UserAgent())
		log.Println("PARAMS: ", r.URL.Query(), "")
		log.Println("HTTP VERSION: ", r.Proto)

		start := time.Now()
		next.ServeHTTP(w, r)
		dur := time.Since(start)

		log.Printf("Request took %v", dur)
		fmt.Println("|||||| END OF REQUEST LOG  ||||||")
	}
}
