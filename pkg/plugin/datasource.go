package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"

	"github.com/arabian9ts/mux-datasource/pkg/config"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	cfg, err := config.LoadPluginSettings(settings)
	if err != nil {
		return nil, err
	}

	return &Datasource{
		tokenID:     cfg.Secrets.TokenID,
		tokenSecret: cfg.Secrets.TokenSecret,
		client:      http.DefaultClient,
	}, nil
}

type Datasource struct {
	tokenID       string
	tokenSecret   string
	client        *http.Client
	rcHandlerOnce sync.Once
	rcHandler     backend.CallResourceHandler
}

func (d *Datasource) Dispose() {}

func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return d.resourceHandler().CallResource(ctx, req, sender)
}

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *Datasource) query(ctx context.Context, _ backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse
	var qm MetricQuery

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	resp, err := d.queryMetric(ctx, qm, query.TimeRange)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("query metric: %v", err.Error()))
	}

	response.Frames = append(response.Frames, resp)

	return response
}

func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	cfg, err := config.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)
	if err != nil {
		return nil, err
	}

	res := &backend.CheckHealthResult{}
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to load settings"
		return res, nil
	}

	if cfg.Secrets.TokenID == "" || cfg.Secrets.TokenSecret == "" {
		res.Status = backend.HealthStatusError
		res.Message = "TokenID or TokenSecret are missing"
		return res, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
