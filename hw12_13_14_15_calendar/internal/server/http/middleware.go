package internalhttp

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	handler http.Handler
	logg    *logrus.Logger
}

func NewLoggerMiddleware(handlerToWrap http.Handler, logg *logrus.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{handlerToWrap, logg}
}

type MyWriter struct {
	w          http.ResponseWriter
	StatusCode int
	DataLen    int
}

func (w *MyWriter) Header() http.Header {
	return w.w.Header()
}

func (w *MyWriter) Write(data []byte) (int, error) {
	w.DataLen += len(data)
	return w.w.Write(data)
}

func (w *MyWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	myWriter := &MyWriter{w: w, StatusCode: http.StatusOK}
	l.handler.ServeHTTP(myWriter, r)
	l.logg.Infof("%v %v %v %v %v %s %dbytes %vms %s", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, myWriter.StatusCode,
		http.StatusText(myWriter.StatusCode), myWriter.DataLen, time.Since(start).Microseconds(), r.UserAgent())
}

func loggingMiddleware(next http.Handler, logg *logrus.Logger) http.Handler {
	return NewLoggerMiddleware(next, logg)
}
