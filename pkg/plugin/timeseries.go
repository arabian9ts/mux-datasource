package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type (
	timeSeriesResponse struct {
		Data []timeSeriesData
	}

	timeSeriesData struct {
		Time  time.Time
		Value float64
		Views int64
	}
)

func (m *timeSeriesResponse) Compact() *timeSeriesResponse {
	v := make([]timeSeriesData, 0, len(m.Data))
	for _, data := range m.Data {
		if data.Value != 0 {
			v = append(v, data)
		}
	}
	return &timeSeriesResponse{Data: v}
}

func (m *timeSeriesData) UnmarshalJSON(data []byte) error {
	var row []any
	err := json.Unmarshal(data, &row)
	if err != nil {
		return err
	}

	timestamp, ok := row[0].(string)
	if !ok {
		return errors.New("failed type cast timestamp")
	}

	m.Time, err = time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}

	value, ok := row[1].(float64)
	if !ok {
		return nil
	}
	m.Value = value

	views, ok := row[2].(int64)
	if !ok {
		return nil
	}
	m.Views = views

	return nil
}

func (d *Datasource) timeseries(ctx context.Context, q MetricQuery, timeRange backend.TimeRange) (*data.Frame, error) {
	path, err := url.JoinPath("/data/v1/metrics/", q.MetricID, "/timeseries")
	if err != nil {
		return nil, err
	}

	filters := make([]string, 0, len(q.Filters))
	for _, f := range q.Filters {
		filters = append(filters, f.String())
	}
	params := map[string][]string{
		"order_by":        {"group_by_value"},
		"order_direction": {"asc"},
		"measurement":     {q.Measurement},
		"filters":         filters,
		"timeframe": {
			strconv.FormatInt(timeRange.From.Unix(), 10),
			strconv.FormatInt(timeRange.To.Unix(), 10),
		},
	}
	req, err := d.newRequest(ctx, path, params)
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg bytes.Buffer
		_, _ = msg.ReadFrom(resp.Body)
		log.DefaultLogger.Error("Request error", "status", resp.Status, "body", msg.String())
		return nil, fmt.Errorf("request error: %s", msg.String())
	}

	ts := new(timeSeriesResponse)
	err = json.NewDecoder(resp.Body).Decode(ts)
	if err != nil {
		log.DefaultLogger.Error("Failed to decode timeseries", "error", err)
		return nil, err
	}

	ts = ts.Compact()
	times := make([]time.Time, 0, len(ts.Data))
	values := make([]float64, 0, len(ts.Data))
	for _, d := range ts.Data {
		times = append(times, d.Time)
		values = append(values, d.Value)
	}

	frame := data.NewFrame(q.MetricID)
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, times),
		data.NewField(q.FrameName(), nil, values),
	)

	return frame, nil
}
