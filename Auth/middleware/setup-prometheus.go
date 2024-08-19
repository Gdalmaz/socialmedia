package middleware



import (
    "bytes"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/gofiber/fiber/v2"
    "net/http"
)

var (
    // Define your metrics
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "handler"},
    )
)

func InitMetrics(app *fiber.App) {
    prometheus.MustRegister(httpRequestsTotal)

    app.Use(func(c *fiber.Ctx) error {
        httpRequestsTotal.WithLabelValues(c.Method(), c.Path()).Inc()
        return c.Next()
    })

    app.Get("/metrics", func(c *fiber.Ctx) error {
        var buf bytes.Buffer
        handler := promhttp.Handler()
        writer := &responseWriter{buf: &buf}
        handler.ServeHTTP(writer, &http.Request{})
        c.Set("Content-Type", "text/plain; charset=utf-8")
        return c.Send(buf.Bytes())
    })
}

type responseWriter struct {
    buf *bytes.Buffer
}

func (rw *responseWriter) Header() http.Header {
    return make(http.Header)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
    return rw.buf.Write(p)
}

func (rw *responseWriter) WriteHeader(statusCode int) {}
