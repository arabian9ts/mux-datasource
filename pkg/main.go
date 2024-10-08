package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"

	"github.com/arabian9ts/mux-datasource/pkg/plugin"
)

func main() {
	if err := datasource.Manage("arabian9ts-mux-datasource", plugin.NewDatasource, datasource.ManageOpts{}); err != nil {
		os.Exit(1)
	}
}
