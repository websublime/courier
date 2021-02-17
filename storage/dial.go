package storage

import (
	"net/url"
	"reflect"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/columns"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/websublime/courier/config"
	"github.com/websublime/courier/storage/namespace"
)

type Connection struct {
	*pop.Connection
}

func Dial(conf *config.EnvironmentConfig) (*Connection, error) {
	u, err := url.Parse(conf.CourierDatabaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "parsing db connection url")
	}

	db, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect: u.Scheme,
		URL:     conf.CourierDatabaseURL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "opening database connection")
	}
	if err := db.Open(); err != nil {
		return nil, errors.Wrap(err, "checking database connection")
	}

	namespace.SetNamespace(conf.CourierDatabaseNamespace)

	if logrus.StandardLogger().Level == logrus.DebugLevel {
		pop.Debug = true
	}
	pop.Debug = true

	return &Connection{db}, nil
}

func (c *Connection) Transaction(fn func(*Connection) error) error {
	if c.TX == nil {
		return c.Connection.Transaction(func(tx *pop.Connection) error {
			return fn(&Connection{tx})
		})
	}
	return fn(c)
}

func getExcludedColumns(model interface{}, includeColumns ...string) ([]string, error) {
	sm := &pop.Model{Value: model}
	st := reflect.TypeOf(model)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}

	// get all columns and remove included to get excluded set
	cols := columns.ForStructWithAlias(model, sm.TableName(), sm.As, "id")
	for _, f := range includeColumns {
		if _, ok := cols.Cols[f]; !ok {
			return nil, errors.Errorf("Invalid column name %s", f)
		}
		cols.Remove(f)
	}

	xcols := make([]string, len(cols.Cols))
	for n := range cols.Cols {
		xcols = append(xcols, n)
	}
	return xcols, nil
}
