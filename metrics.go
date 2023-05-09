package nitriding

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	reqPath    = "http_req_path"
	reqMethod  = "http_req_method"
	respStatus = "http_resp_status"
	respErr    = "http_resp_error"

	notAvailable = "n/a"
	namespace    = "nitriding"
)

// metrics contains our Prometheus metrics.
type metrics struct {
	proxiedReqs *prometheus.CounterVec
}

// newMetrics initializes our Prometheus metrics.
func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		proxiedReqs: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "proxy_responses",
				Help:      "HTTP responses of the enclave application backend",
			},
			[]string{reqPath, reqMethod, respStatus, respErr},
		),
	}
	reg.MustRegister(m.proxiedReqs)

	return m
}

// checkRevProxyResp captures Prometheus metrics for HTTP responses from our
// enclave application backend.
func (m *metrics) checkRevProxyResp(resp *http.Response) error {
	m.proxiedReqs.With(prometheus.Labels{
		reqPath:    resp.Request.URL.Path,
		reqMethod:  resp.Request.Method,
		respStatus: fmt.Sprint(resp.StatusCode),
		respErr:    notAvailable,
	}).Inc()

	return nil
}

// checkRevProxyErr captures Prometheus metrics for errors that occurred when
// we tried to talk to the enclave application backend.
func (m *metrics) checkRevProxyErr(_ http.ResponseWriter, r *http.Request, err error) {
	m.proxiedReqs.With(prometheus.Labels{
		reqPath:    r.URL.Path,
		reqMethod:  r.Method,
		respStatus: notAvailable,
		respErr:    err.Error(),
	}).Inc()
}

// middleware implements a chi middleware that records each request as part of
// our Prometheus metrics.
func (m *metrics) middleware(h http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		h.ServeHTTP(ww, r)
		m.proxiedReqs.With(prometheus.Labels{
			reqPath:    r.URL.Path,
			reqMethod:  r.Method,
			respStatus: fmt.Sprint(ww.Status()),
			respErr:    notAvailable,
		}).Inc()
	}
	return http.HandlerFunc(f)
}
