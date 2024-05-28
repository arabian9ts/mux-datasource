package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"

	"github.com/arabian9ts/mux-datasource/pkg/constants"
)

func (d *Datasource) resourceHandler() backend.CallResourceHandler {
	d.rcHandlerOnce.Do(func() {
		d.rcHandler = httpadapter.New(d.resources())
	})
	return d.rcHandler
}

func (d *Datasource) resources() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("GET /metrics", d.metrics)
	r.HandleFunc("GET /metrics/{metricId}/breakdowns", d.breakdowns)
	r.HandleFunc("GET /measurements", d.measurements)
	r.HandleFunc("GET /groups", d.groups)
	return r
}

func (d *Datasource) metrics(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(constants.MetricIDs)
}

func (d *Datasource) measurements(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(constants.Measurements)
}

func (d *Datasource) groups(w http.ResponseWriter, r *http.Request) {
	dims, err := d.listGroups(r.Context())
	if err != nil {
		log.DefaultLogger.Error("Failed to list groups", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(dims)
}

type listDimensionResponse struct {
	Data struct {
		Basic    []string `json:"basic"`
		Advanced []string `json:"advanced"`
	}
}

func (d *Datasource) listGroups(ctx context.Context) ([]string, error) {
	const path = "/data/v1/dimensions"
	listReq, err := d.newRequest(ctx, path, nil)
	listResp, err := d.client.Do(listReq)
	if err != nil {
		return nil, err
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		var msg bytes.Buffer
		_, _ = msg.ReadFrom(listResp.Body)
		log.DefaultLogger.Error("Request error", "status", listResp.Status, "body", msg.String())
		return nil, fmt.Errorf("request error: %s", msg.String())
	}

	var dims listDimensionResponse
	err = json.NewDecoder(listResp.Body).Decode(&dims)
	if err != nil {
		log.DefaultLogger.Error("Failed to decode dimensions", "error", err)
		return nil, err
	}

	resp := make([]string, 0, len(dims.Data.Basic)+len(dims.Data.Advanced))
	resp = append(resp, dims.Data.Basic...)
	resp = append(resp, dims.Data.Advanced...)
	return resp, nil
}

type (
	breakdownResponse struct {
		Data []breakdownData
	}

	breakdownData struct {
		Field          string        `json:"field"`
		Value          float64       `json:"value"`
		Views          int64         `json:"views"`
		TotalWatchTime time.Duration `json:"total_watch_time"`
	}
)

func (m *breakdownResponse) Fields() []string {
	v := make([]string, 0, len(m.Data))
	for _, data := range m.Data {
		v = append(v, data.Field)
	}
	return v
}

func (d *Datasource) breakdowns(w http.ResponseWriter, r *http.Request) {
	metricID := r.PathValue("metricId")
	groupBy := r.URL.Query().Get("groupBy")

	if metricID == "" || groupBy == "" {
		http.Error(w, "metricId or groupBy is blank", http.StatusBadRequest)
		return
	}

	path, err := url.JoinPath("/data/v1/metrics/", metricID, "/breakdown")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := map[string][]string{
		"order_direction": {"desc"},
		"limit":           {"50"},
		"group_by":        {groupBy},
		"timeframe[]":     {"1:hour"},
	}
	req, err := d.newRequest(r.Context(), path, params)
	resp, err := d.client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg bytes.Buffer
		_, _ = msg.ReadFrom(resp.Body)
		http.Error(w, fmt.Sprintf("request error: %s", msg.String()), resp.StatusCode)
		return
	}

	breakdowns := new(breakdownResponse)
	err = json.NewDecoder(resp.Body).Decode(breakdowns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(breakdowns.Fields())
}
