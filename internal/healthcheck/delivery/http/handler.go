package http

import (
	"github.com/go-rest-api/internal/healthcheck/service"
	"github.com/go-rest-api/pkg/response"
	"github.com/opentracing/opentracing-go"
	nethttp "net/http"
)

// operation constants
const (
	// Operation name
	APIHealthCheckOperation            = "handler.healthcheck.api"
	InfrastructureHealthCheckOperation = "handler.healthcheck.infrastructure"
)

// HealthCheckDelegate type
type HealthCheckHandler struct {
	service service.IHealthCheckService
}

// NewHealthCheckDelegate for initiate HealthCheckDelegate
func NewHealthCheckHandler(service service.IHealthCheckService) (*HealthCheckHandler, error) {
	delegate := &HealthCheckHandler{
		service: service,
	}

	return delegate, nil
}

// API is handler function for check the API healthiness
func (d *HealthCheckHandler) API(w nethttp.ResponseWriter, r *nethttp.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), APIHealthCheckOperation)
	defer span.Finish()

	d.service.API(ctx)

	response.WriteAPIOK(w)
}

// Infrastructure is handler function for check the Infrastructure healthiness
func (d *HealthCheckHandler) Infrastructure(w nethttp.ResponseWriter, r *nethttp.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), InfrastructureHealthCheckOperation)
	defer span.Finish()

	healthResult := d.service.Infrastructure(ctx)

	if !healthResult.IsOk {
		response.WriteAPIErrorWithData(w, response.APIErrNotFound, healthResult)
		return
	}
	response.WriteAPIOKWithData(w, healthResult)
}
