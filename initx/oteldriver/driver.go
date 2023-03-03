package oteldriver

import (
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/XSAM/otelsql"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/pkg/store/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// BuildOTELDriver 构建具有otel的sql驱动,一般该驱动需要优先调用.
func BuildOTELDriver(cnf *conf.AppConfiguration, storekey string) dialect.Driver {
	storeCfg := cnf.Sub(storekey)
	driverName := storeCfg.String("driverName")
	if cnf.IsSet("otel") {
		// Register the otelsql wrapper for the provided postgres driver.
		driverName, err := otelsql.Register("mysql",
			otelsql.WithAttributes(semconv.DBSystemMySQL),
			otelsql.WithAttributes(semconv.DBNameKey.String(storekey)),
			otelsql.WithSpanOptions(otelsql.SpanOptions{
				DisableErrSkip:  true,
				OmitRows:        true,
				OmitConnPrepare: true,
			}),
		)
		if err != nil {
			panic(err)
		}
		storeCfg.Parser().Set("driverName", driverName)
	}
	db := sqlx.NewSqlDB(storeCfg)
	return sql.OpenDB(driverName, db)
}
