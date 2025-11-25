package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dev-zapi/docker-simple-panel/config"
)

// responseWriter wraps http.ResponseWriter to capture the status code and response body
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	body         *bytes.Buffer
	captureBody  bool
	bytesWritten int
}

func newResponseWriter(w http.ResponseWriter, captureBody bool) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           &bytes.Buffer{},
		captureBody:    captureBody,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.bytesWritten += len(b)
	if rw.captureBody {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b)
}

// Logging creates a logging middleware with configurable log levels
func Logging(configManager *config.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logLevel := configManager.GetLogLevel()

			// Skip logging for OPTIONS requests at INFO level and below
			if r.Method == "OPTIONS" && logLevel < config.LogLevelDebug {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			// Read request body for DEBUG level
			var requestBody []byte
			if logLevel == config.LogLevelDebug && r.Body != nil {
				requestBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			// Create response writer wrapper
			captureBody := logLevel == config.LogLevelDebug
			rw := newResponseWriter(w, captureBody)

			// Process request
			next.ServeHTTP(rw, r)

			// Calculate duration
			duration := time.Since(start)

			// Log based on level
			switch logLevel {
			case config.LogLevelError:
				// Only log errors (5xx responses)
				if rw.statusCode >= 500 {
					logError(r, rw, duration)
				}
			case config.LogLevelWarn:
				// Log warnings: errors and slow requests (>1s) and client errors (4xx)
				if rw.statusCode >= 400 || duration > time.Second {
					logWarn(r, rw, duration)
				}
			case config.LogLevelInfo:
				// Log basic request info
				logInfo(r, rw, duration)
			case config.LogLevelDebug:
				// Log detailed request/response info
				logDebug(r, rw, duration, requestBody)
			}
		})
	}
}

// logError logs error level messages (5xx errors)
func logError(r *http.Request, rw *responseWriter, duration time.Duration) {
	log.Printf("[ERROR] %s %s - Status: %d - Duration: %v - IP: %s",
		r.Method, r.URL.Path, rw.statusCode, duration, getClientIP(r))
}

// logWarn logs warning level messages (4xx errors, slow requests)
func logWarn(r *http.Request, rw *responseWriter, duration time.Duration) {
	msg := "[WARN] %s %s - Status: %d - Duration: %v - IP: %s"
	if duration > time.Second {
		msg = "[WARN] SLOW REQUEST: %s %s - Status: %d - Duration: %v - IP: %s"
	}
	log.Printf(msg, r.Method, r.URL.Path, rw.statusCode, duration, getClientIP(r))
}

// logInfo logs info level messages (basic request info)
func logInfo(r *http.Request, rw *responseWriter, duration time.Duration) {
	log.Printf("[INFO] %s %s - Status: %d - Duration: %v - Size: %d bytes",
		r.Method, r.URL.Path, rw.statusCode, duration, rw.bytesWritten)
}

// logDebug logs debug level messages (detailed request/response info)
func logDebug(r *http.Request, rw *responseWriter, duration time.Duration, requestBody []byte) {
	// Log request
	log.Printf("[DEBUG] --> %s %s", r.Method, r.URL.String())
	log.Printf("[DEBUG]     IP: %s", getClientIP(r))
	log.Printf("[DEBUG]     User-Agent: %s", r.Header.Get("User-Agent"))

	// Log selected headers (excluding sensitive ones)
	for key, values := range r.Header {
		lowerKey := strings.ToLower(key)
		// Skip sensitive headers
		if lowerKey == "authorization" || lowerKey == "cookie" {
			log.Printf("[DEBUG]     Header %s: [REDACTED]", key)
			continue
		}
		log.Printf("[DEBUG]     Header %s: %s", key, strings.Join(values, ", "))
	}

	// Log request body (truncated if too long)
	if len(requestBody) > 0 {
		bodyStr := string(requestBody)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "... (truncated)"
		}
		log.Printf("[DEBUG]     Body: %s", bodyStr)
	}

	// Log response
	log.Printf("[DEBUG] <-- %d %s - Duration: %v - Size: %d bytes",
		rw.statusCode, http.StatusText(rw.statusCode), duration, rw.bytesWritten)

	// Log response body (truncated if too long)
	if rw.body.Len() > 0 {
		bodyStr := rw.body.String()
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "... (truncated)"
		}
		log.Printf("[DEBUG]     Response: %s", bodyStr)
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
