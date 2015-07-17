// +build influxdb

package connector

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/facette/facette/pkg/catalog"
	"github.com/facette/facette/pkg/config"
	"github.com/facette/facette/pkg/logger"
	"github.com/facette/facette/pkg/plot"
	influxdb "github.com/influxdb/influxdb/client"
	"github.com/influxdb/influxdb/influxql"
)

// InfluxDBConnector represents the main structure of the InfluxDB connector.
type InfluxDBConnector struct {
	name     string
	url      string
	username string
	password string
	database string
	mapping  influxDBMap
	client   *influxdb.Client
}

// influxDBMap represents the structure of InfluxDB series mapping.
type influxDBMap struct {
	source []string
	metric []string
	glue   string
	maps   map[string]map[string]influxDBMapEntry
}

type influxDBMapEntry struct {
	column string
	terms  map[string]string
}

func init() {
	Connectors["influxdb"] = func(name string, settings map[string]interface{}) (Connector, error) {
		var err error

		connector := &InfluxDBConnector{
			name: name,
			mapping: influxDBMap{
				glue: ".",
				maps: make(map[string]map[string]influxDBMapEntry),
			},
		}

		// Load provider configuration
		if connector.url, err = config.GetString(settings, "url", true); err != nil {
			return nil, err
		}

		if connector.username, err = config.GetString(settings, "username", true); err != nil {
			return nil, err
		}

		if connector.password, err = config.GetString(settings, "password", true); err != nil {
			return nil, err
		}

		if connector.database, err = config.GetString(settings, "database", true); err != nil {
			return nil, err
		}

		mapping, err := config.GetStringMap(settings, "mapping", true)
		if err != nil {
			return nil, err
		}

		if connector.mapping.source, err = config.GetStringSlice(mapping, "source", true); err != nil {
			return nil, err
		}

		if connector.mapping.metric, err = config.GetStringSlice(mapping, "metric", true); err != nil {
			return nil, err
		}

		glue, err := config.GetString(mapping, "glue", false)
		if err != nil {
			return nil, err
		} else if glue != "" {
			connector.mapping.glue = glue
		}

		// Create new client instance
		url, err := url.Parse(connector.url)
		if err != nil {
			return nil, fmt.Errorf("unable to parse URL: %s", err)
		}

		connector.client, err = influxdb.NewClient(influxdb.Config{
			URL:      *url,
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
func (connector *InfluxDBConnector) GetPlots(query *plot.Query) ([]*plot.Series, error) {
	var queries []string

	l := len(query.Series)
	if l == 0 {
		return nil, fmt.Errorf("influxdb[%s]: requested series list is empty", connector.name)
	}

	results := make([]*plot.Series, l)

	// Prepare query
	for _, s := range query.Series {
		var (
			series string
			parts  []string
		)

		if _, ok := connector.mapping.maps[s.Source]; !ok {
			return nil, fmt.Errorf("unknown series source `%s'", s.Source)
		} else if _, ok := connector.mapping.maps[s.Source][s.Metric]; !ok {
			return nil, fmt.Errorf("unknown series metric `%s' for source `%s'", s.Source, s.Metric)
		}

		mapping := connector.mapping.maps[s.Source][s.Metric]

		for term, value := range mapping.terms {
			if term == "" {
				series = value
			} else {
				parts = append(parts, fmt.Sprintf("%s = %s", influxql.QuoteIdent(term), influxql.QuoteString(value)))
			}
		}

		parts = append(parts, fmt.Sprintf("time > %ds and time < %ds order by asc", query.StartTime.Unix(),
			query.EndTime.Unix()))

		queries = append(queries, fmt.Sprintf("select %s, time from %s where %s", mapping.column, strconv.Quote(series),
			strings.Join(parts, " and ")))
	}

	q := influxdb.Query{
		Command:  strings.Join(queries, "; "),
		Database: connector.database,
	}

	// Execute query
	logger.Log(logger.LevelDebug, "connector", "influxdb[%s]: executing: %s", connector.name, q.Command)

	response, err := connector.client.Query(q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plots: %s", err)
	} else if response.Error() != nil {
		return nil, fmt.Errorf("failed to fetch plots: %s", response.Error())
	}

	// Parse results received from backend
	step := int(query.EndTime.Sub(query.StartTime) / time.Duration(query.Sample))

	for i, r := range response.Results {
		if r.Err != nil {
			continue
		}

		results[i] = &plot.Series{
			Step: step,
		}

		for _, s := range r.Series {
			for _, v := range s.Values {
				time, err := time.Parse(time.RFC3339Nano, v[0].(string))
				if err != nil {
					logger.Log(logger.LevelWarning, "connector", "influxdb[%s]: failed to parse time: %s",
						connector.name, v[0])
					continue
				}

				value, err := v[1].(json.Number).Float64()
				if err != nil {
					logger.Log(logger.LevelWarning, "connector", "influxdb[%s]: failed to parse value: %s",
						connector.name, v[1])
					continue
				}

				results[i].Plots = append(results[i].Plots, plot.Plot{
					Time:  time,
					Value: plot.Value(value),
				})
			}
		}
	}

	return results, nil
}

// Refresh triggers a full connector data update.
func (connector *InfluxDBConnector) Refresh(originName string, outputChan chan<- *catalog.Record) error {
	// Query backend for sample rows (used to detect numerical values)
	columnsMap := make(map[string][]string, 0)

	q := influxdb.Query{
		Command:  "select * from /.*/ limit 1",
		Database: connector.database,
	}

	response, err := connector.client.Query(q)
	if err != nil {
		return fmt.Errorf("failed to fetch sample rows: %s", err)
	} else if response.Error() != nil {
		return fmt.Errorf("failed to fetch sample rows: %s", response.Error())
	}

	if len(response.Results) != 1 {
		return fmt.Errorf("failed to retrieve sample rows: expected 1 result but got %d", len(response.Results))
	} else if response.Results[0].Err != nil {
		return fmt.Errorf("failed to retrieve sample rows: %s", response.Results[0].Err)
	}

	for _, s := range response.Results[0].Series {
		if len(s.Values) == 0 {
			continue
		}

		if _, ok := columnsMap[s.Name]; !ok {
			columnsMap[s.Name] = make([]string, 0)
		}

		for i, v := range s.Values[0] {
			if _, ok := v.(json.Number); !ok {
				continue
			}

			columnsMap[s.Name] = append(columnsMap[s.Name], s.Columns[i])
		}
	}

	// Query backend for series list
	q = influxdb.Query{
		Command:  "show series",
		Database: connector.database,
	}

	response, err = connector.client.Query(q)
	if err != nil {
		return fmt.Errorf("failed to fetch series: %s", err)
	} else if response.Error() != nil {
		return fmt.Errorf("failed to fetch series: %s", response.Error())
	}

	if len(response.Results) != 1 {
		return fmt.Errorf("failed to retrieve series: expected 1 result but got %d", len(response.Results))
	} else if response.Results[0].Err != nil {
		return fmt.Errorf("failed to retrieve series: %s", response.Results[0].Err)
	}

	// Parse results for sources and metrics
	for _, s := range response.Results[0].Series {
		for i := range s.Values {
			var parts []string

			terms := make(map[string]string)

			// Map source and metric names
			for _, item := range connector.mapping.source {
				term, part := mapKey(s, i, item)
				if part != "" {
					terms[term] = part
					parts = append(parts, part)
				}
			}
			sourceName := strings.Join(parts, connector.mapping.glue)

			parts = []string{}
			for _, item := range connector.mapping.metric {
				term, part := mapKey(s, i, item)
				if part != "" {
					terms[term] = part
					parts = append(parts, part)
				}
			}
			metricName := strings.Join(parts, connector.mapping.glue)

			terms[""] = s.Name

			// Initialize metric mapping terms if needed
			if _, ok := connector.mapping.maps[sourceName]; !ok {
				connector.mapping.maps[sourceName] = make(map[string]influxDBMapEntry)
			}

			for _, c := range columnsMap[s.Name] {
				connector.mapping.maps[sourceName][metricName+connector.mapping.glue+c] = influxDBMapEntry{
					column: c,
					terms:  terms,
				}

				// Send record to catalog
				outputChan <- &catalog.Record{
					Origin:    originName,
					Source:    sourceName,
					Metric:    metricName + connector.mapping.glue + c,
					Connector: connector,
				}
			}
		}
	}

	return nil
}

func mapKey(row influxql.Row, index int, item string) (string, string) {
	if item == "name" {
		return "", row.Name
	} else if strings.HasPrefix(item, "column:") {
		// Try to match row column
		name := strings.TrimPrefix(item, "column:")
		for i, value := range row.Columns {
			if name != value {
				continue
			}

			return name, fmt.Sprintf("%v", row.Values[index][i])
		}
	} else if strings.HasPrefix(item, "tag:") {
		// Try to match row tag
		name := strings.TrimPrefix(item, "tag:")
		return name, row.Tags[name]
	}

	// Nothing matched
	return "", ""
}
