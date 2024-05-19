package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arabian9ts/mux-datasource/pkg/config"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
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
	tokenID     string
	tokenSecret string
	client      *http.Client
}

func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	MetricName string
	Dimension  string
	//TimeRange
}

func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	frame := data.NewFrame("response")
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("values", nil, []int64{10, 20}),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	config, err := config.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)
	log.DefaultLogger.Info("CheckHealth", "datasource", req.PluginContext.DataSourceInstanceSettings.DecryptedSecureJSONData)

	res := &backend.CheckHealthResult{}
	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to load settings"
		return res, nil
	}

	if config.Secrets.TokenID == "" || config.Secrets.TokenSecret == "" {
		res.Status = backend.HealthStatusError
		res.Message = "TokenID or TokenSecret are missing"
		return res, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
