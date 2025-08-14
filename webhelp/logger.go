package webhelp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s path=%s %s duration=%vms remote_addr=%s",
			r.Method,
			r.URL.Path,
			r.URL.Query(),
			time.Since(start).Milliseconds(),
			r.RemoteAddr,
		)
	})
}

type customLogger struct{}

func (l *customLogger) Write(p []byte) (n int, err error) {
	logEntryWithMediumPath := string(p[len(pwd)+1:])
	return fmt.Printf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), logEntryWithMediumPath)
}

var pwd string

func UseLogger() {
	pwd, _ = os.Getwd()
	log.SetFlags(log.Llongfile)
	log.SetOutput(new(customLogger))
}
