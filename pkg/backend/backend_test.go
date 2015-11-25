package backend

import (
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	date                          time.Time
	mysqlCfg, pgsqlCfg, sqliteCfg map[string]interface{}
)

func init() {
	date = time.Now().Round(time.Second)

	mysqlCfg = map[string]interface{}{
		"driver":   "mysql",
		"dbname":   "facette",
		"user":     "facette",
		"password": "facette",
	}

	if v := os.Getenv("TEST_MYSQL_DBNAME"); v != "" {
		mysqlCfg["dbname"] = v
	}
	if v := os.Getenv("TEST_MYSQL_ADDRESS"); v != "" {
		mysqlCfg["address"] = v
	}
	if v := os.Getenv("TEST_MYSQL_USER"); v != "" {
		mysqlCfg["user"] = v
	}
	if v := os.Getenv("TEST_MYSQL_PASSWORD"); v != "" {
		mysqlCfg["password"] = v
	}

	pgsqlCfg = map[string]interface{}{
		"driver":   "postgres",
		"dbname":   "facette",
		"user":     "facette",
		"password": "facette",
	}

	if v := os.Getenv("TEST_PGSQL_DBNAME"); v != "" {
		pgsqlCfg["dbname"] = v
	}
	if v := os.Getenv("TEST_PGSQL_HOST"); v != "" {
		pgsqlCfg["host"] = v
	}
	if v := os.Getenv("TEST_PGSQL_PORT"); v != "" {
		pgsqlCfg["port"] = v
	}
	if v := os.Getenv("TEST_PGSQL_USER"); v != "" {
		pgsqlCfg["user"] = v
	}
	if v := os.Getenv("TEST_PGSQL_PASSWORD"); v != "" {
		pgsqlCfg["password"] = v
	}

	sqliteCfg = map[string]interface{}{
		"driver": "sqlite3",
		"path":   ":memory:",
	}

	if v := os.Getenv("TEST_SQLITE_PATH"); v != "" {
		sqliteCfg["path"] = v
	}
}

func Test_Sqlite3_Graph(t *testing.T) {
	execTestGraph(sqliteCfg, t)
}

func Test_Sqlite3_SourceGroup(t *testing.T) {
	execTestSourceGroup(sqliteCfg, t)
}

func Test_Sqlite3_MetricGroup(t *testing.T) {
	execTestMetricGroup(sqliteCfg, t)
}

func Test_Sqlite3_Scale(t *testing.T) {
	execTestScale(sqliteCfg, t)
}

func Test_Sqlite3_Unit(t *testing.T) {
	execTestUnit(sqliteCfg, t)
}

func Test_Postgres_Graph(t *testing.T) {
	execTestGraph(pgsqlCfg, t)
}

func Test_Postgres_SourceGroup(t *testing.T) {
	execTestSourceGroup(pgsqlCfg, t)
}

func Test_Postgres_MetricGroup(t *testing.T) {
	execTestMetricGroup(pgsqlCfg, t)
}

func Test_Postgres_Scale(t *testing.T) {
	execTestScale(pgsqlCfg, t)
}

func Test_Postgres_Unit(t *testing.T) {
	execTestUnit(pgsqlCfg, t)
}

func Test_MySQL_Graph(t *testing.T) {
	execTestGraph(mysqlCfg, t)
}

func Test_MySQL_SourceGroup(t *testing.T) {
	execTestSourceGroup(mysqlCfg, t)
}

func Test_MySQL_MetricGroup(t *testing.T) {
	execTestMetricGroup(mysqlCfg, t)
}

func Test_MySQL_Scale(t *testing.T) {
	execTestScale(mysqlCfg, t)
}

func Test_MySQL_Unit(t *testing.T) {
	execTestUnit(mysqlCfg, t)
}

func execTestGraph(cfg map[string]interface{}, t *testing.T) {
	b, err := NewBackend(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, []interface{}{
		Graph{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000000",
				Name:        "graph1",
				Description: "A great graph description",
				Created:     date,
				Modified:    date,
			},
			Series: []SeriesGroup{
				SeriesGroup{
					Series: []Series{
						Series{
							Name:   "series1",
							Origin: "origin1",
							Source: "source1",
							Metric: "metric1",
						},
					},
				},
			},
			Options: map[string]interface{}{
				"title": "A great graph title",
			},
			Template: false,
		},
		Graph{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000001",
				Name:        "graph2",
				Description: "Another great graph description",
				Created:     date,
				Modified:    date,
			},
			Series: []SeriesGroup{
				SeriesGroup{
					Series: []Series{
						Series{
							Name:   "series1",
							Origin: "origin1",
							Source: "source1",
							Metric: "metric1",
						},
					},
				},
			},
			Options: map[string]interface{}{
				"title": "Another great graph title",
			},
			Template: false,
		},
	}, Graph{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "graph1",
			Description: "A great graph description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Series: []SeriesGroup{
			SeriesGroup{
				Series: []Series{
					Series{
						Name:   "series1",
						Origin: "origin1",
						Source: "source1",
						Metric: "metric1",
					},
					Series{
						Name:   "series2",
						Origin: "origin1",
						Source: "source1",
						Metric: "metric2",
					},
				},
			},
		},
		Options: map[string]interface{}{
			"title": "A great graph title (updated)",
		},
		Template: false,
	}, t)
}

func execTestSourceGroup(cfg map[string]interface{}, t *testing.T) {
	b, err := NewBackend(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, []interface{}{
		SourceGroup{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000000",
				Name:        "sourcegroup1",
				Description: "A great sourcegroup description",
				Created:     date,
				Modified:    date,
			},
			Entries: []GroupEntry{
				GroupEntry{
					Pattern: "glob:host*.example.net",
					Origin:  "origin1",
				},
			},
		},
		SourceGroup{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000001",
				Name:        "sourcegroup2",
				Description: "Another great sourcegroup description",
				Created:     date,
				Modified:    date,
			},
			Entries: []GroupEntry{
				GroupEntry{
					Pattern: "host3.example.net",
					Origin:  "origin1",
				},
			},
		},
	}, SourceGroup{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "sourcegroup1",
			Description: "A great sourcegroup description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Entries: []GroupEntry{
			GroupEntry{
				Pattern: "glob:host*.example.net",
				Origin:  "",
			},
		},
	}, t)
}

func execTestMetricGroup(cfg map[string]interface{}, t *testing.T) {
	b, err := NewBackend(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, []interface{}{
		MetricGroup{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000000",
				Name:        "metricgroup1",
				Description: "A great metricgroup description",
				Created:     date,
				Modified:    date,
			},
			Entries: []GroupEntry{
				GroupEntry{
					Pattern: "glob:metric1.*",
					Origin:  "origin1",
				},
			},
		},
		MetricGroup{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000001",
				Name:        "metricgroup2",
				Description: "Another great metricgroup description",
				Created:     date,
				Modified:    date,
			},
			Entries: []GroupEntry{
				GroupEntry{
					Pattern: "metric2",
					Origin:  "origin1",
				},
			},
		},
	}, MetricGroup{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "metricgroup1",
			Description: "A great metricgroup description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Entries: []GroupEntry{
			GroupEntry{
				Pattern: "glob:metric1.*",
				Origin:  "",
			},
		},
	}, t)
}

func execTestScale(cfg map[string]interface{}, t *testing.T) {
	b, err := NewBackend(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, []interface{}{
		Scale{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000000",
				Name:        "scale1",
				Description: "A great scale description",
				Created:     date,
				Modified:    date,
			},
			Value: 0.123,
		},
		Scale{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000001",
				Name:        "scale2",
				Description: "Another great scale description",
				Created:     date,
				Modified:    date,
			},
			Value: 0.456,
		},
	}, Scale{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "scale1",
			Description: "A great scale description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Value: 0.1234,
	}, t)
}

func execTestUnit(cfg map[string]interface{}, t *testing.T) {
	b, err := NewBackend(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, []interface{}{
		Unit{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000000",
				Name:        "unit1",
				Description: "A great unit description",
				Created:     date,
				Modified:    date,
			},
			Label: "a",
		},
		Unit{
			Item: Item{
				ID:          "00000000-0000-0000-0000-000000000001",
				Name:        "unit2",
				Description: "Another great unit description",
				Created:     date,
				Modified:    date,
			},
			Label: "b",
		},
	}, Unit{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "unit1",
			Description: "A great unit description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Label: "c",
	}, t)
}

func execTest(b *Backend, items []interface{}, update interface{}, t *testing.T) {
	// Check items insertion
	for _, item := range items {
		rv := reflect.ValueOf(item)

		if err := b.Insert(item); err != nil {
			t.Fatal(err)
		}

		result := reflect.New(reflect.TypeOf(item)).Interface()
		if err := b.Get(rv.FieldByName("ID").String(), result); err != nil {
			t.Fatal(err)
		}

		result = reflect.Indirect(reflect.ValueOf(result)).Interface()
		if !deepEqual(item, result) {
			t.Logf("\nExpected %#v\nbut got  %#v", item, result)
			t.Fail()
		}
	}

	// Check items list
	sv := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(items[0])), 0, 0)

	s := reflect.New(sv.Type())
	if err := b.List(s.Interface(), ""); err != nil {
		t.Fatal(err)
	}

	if reflect.Indirect(s).Len() != len(items) {
		t.Logf("\nExpected %d\nbut got  %d", len(items), reflect.Indirect(s).Len())
		t.Fail()
	}

	for i, item := range items {
		r := reflect.Indirect(s).Index(i).Interface()
		if !deepEqual(item, r) {
			t.Logf("\nExpected %#v\nbut got  %#v", item, r)
			t.Fail()
		}
	}

	// Check item update
	if err := b.Update(update); err != nil {
		t.Fatal(err)
	}

	result := reflect.New(reflect.TypeOf(update)).Interface()
	if err := b.Get("00000000-0000-0000-0000-000000000000", result); err != nil {
		t.Fatal(err)
	}

	result = reflect.Indirect(reflect.ValueOf(result)).Interface()
	if !deepEqual(update, result) {
		t.Logf("\nExpected %#v\nbut got  %#v", update, result)
		t.Fail()
	}

	for _, item := range items {
		// Check item deletion
		if err := b.Delete(item); err != nil {
			t.Fatal(err)
		}
	}

	// Check for empty items list
	s = reflect.New(sv.Type())
	if err := b.List(s.Interface(), ""); err != nil {
		t.Fatal(err)
	}

	if reflect.Indirect(s).Len() != 0 {
		t.Logf("\nExpected %d\nbut got  %d", 0, reflect.Indirect(s).Len())
		t.Fail()
	}
}

func deepEqual(a, b interface{}) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() == reflect.Struct {
		if va.NumField() != vb.NumField() {
			return false
		}

		for i, n := 0, va.NumField(); i < n; i++ {
			ia := va.Field(i).Interface()
			ib := vb.Field(i).Interface()

			if ta, ok := ia.(time.Time); ok {
				if tb, ok := ib.(time.Time); !ok {
					return false
				} else if !ta.Equal(tb) {
					return false
				}
			}

			return deepEqual(ia, ib)
		}
	}

	return reflect.DeepEqual(a, b)
}
