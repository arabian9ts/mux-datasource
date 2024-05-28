package plugin

import (
	"fmt"
)

type MetricQuery struct {
	MetricID    string         `json:"metricId"`
	Measurement string         `json:"measurement"`
	Filters     []MetricFilter `json:"filters"`
}

type MetricFilter struct {
	Group    string   `json:"group"`
	Operator operator `json:"operator"`
	Value    string   `json:"value"`
}

type operator string

var operatorMap = map[operator]string{
	"=":  "",
	"!=": "!",
}

func (o operator) String() string {
	return operatorMap[o]
}

func (m MetricFilter) String() string {
	return fmt.Sprintf("%s%s:%s", m.Operator.String(), m.Group, m.Value)
}
