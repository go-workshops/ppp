package tracing

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func InstrumentHTTP(h http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(h, fmt.Sprintf("%s_endpoint", operation))
}

func HTTPMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		carrier := propagation.HeaderCarrier(w.Header())
		otel.GetTextMapPropagator().Inject(ctx, carrier)
		h.ServeHTTP(w, r)
	})
}

type HTTPTransport struct {
	Transport http.RoundTripper
}

func (t *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))
	return t.Transport.RoundTrip(req)
}
