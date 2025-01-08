package observability

import (
	"fmt"
	"log"
	"net"
	"net/http"

	cProm "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func InitMetrics(port int) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal("could not create prometheus metrics exporter: ", err)
	}

	metric := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)

	otel.SetMeterProvider(metric)

	metricsListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	handler := promhttp.Handler()
	http.Handle("/metrics", handler)

	log.Print("starting metrics http endpoint at: ", metricsListener.Addr())
	if err := http.Serve(metricsListener, handler); err != nil {
		log.Fatal("could not start metrics http endpoint: ", err)
	}
}

// number of requests processed
var OpsProcessed = promauto.NewCounter(cProm.CounterOpts{
	Name: "ops_requests_processed_total",
	Help: "Total number of requests processed",
})
