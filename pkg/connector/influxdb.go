// +build influxdb

package connector

import (
	"fmt"

	"github.com/facette/facette/pkg/catalog"
	"github.com/facette/facette/pkg/config"
	"github.com/facette/facette/pkg/plot"
	influxdb "github.com/facette/facette/thirdparty/github.com/influxdb/influxdb/client"
)

// InfluxDBConnector represents the main structure of the InfluxDB connector.
type InfluxDBConnector struct {
	name     string
	url      string
	username string
	password string
	database string
	client   *influxdb.Client
}

func init() {
	Connectors["influxdb"] = func(name string, settings map[string]interface{}) (Connector, error) {
		var (
			pattern string
			err     error
		)

		connector := &InfluxDBConnector{
			name:     name,
			url:      "http://localhost:8086",
			username: "root",
			password: "root",
		}

		if connector.url, err = config.GetString(settings, "url", false); err != nil {
			return nil, err
		}

		if connector.username, err = config.GetString(settings, "username", false); err != nil {
			return nil, err
		}

		if connector.password, err = config.GetString(settings, "password", false); err != nil {
			return nil, err
		}

		if connector.database, err = config.GetString(settings, "database", true); err != nil {
			return nil, err
		}

		connector.client, err = influxdb.NewClient(&influxdb.Config{
			URL:      connector.url,
			Username: connector.username,
			Password: connector.password,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to create client: %s", err)
		}

		return connector, nil
	}
}

// GetName returns the name of the current connector.
func (connector *InfluxDBConnector) GetName() string {
	return connector.name
}

// GetPlots retrieves time series data from provider based on a query and a time interval.
func (connector *InfluxDBConnector) GetPlots(query *plot.Query) ([]plot.Series, error) {
	return nil, nil
}

// Refresh triggers a full connector data update.
func (connector *InfluxDBConnector) Refresh(originName string, outputChan chan<- *catalog.Record) error {
	return nil
}
