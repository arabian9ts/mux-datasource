package plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const muxHost = "api.mux.com"

func (d *Datasource) queryMetric(ctx context.Context, q MetricQuery, timeRange backend.TimeRange) (*data.Frame, error) {
	return d.timeseries(ctx, q, timeRange)
}

func (d *Datasource) newRequest(ctx context.Context, path string, params map[string][]string) (*http.Request, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   muxHost,
		Path:   path,
	}

	q := u.Query()
	for k, v := range params {
		switch len(v) {
		case 0:
			continue
		case 1:
			q.Add(k, v[0])
		default:
			kk := fmt.Sprintf("%s[]", k)
			for _, vv := range v {
				q.Add(kk, vv)
			}
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(d.tokenID, d.tokenSecret)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
