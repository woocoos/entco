package oteldriver

import (
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/XSAM/otelsql"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/pkg/store/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// TryRegisterOTEL 尝试注册otel,如果配置中有otel配置,则注册.
func TryRegisterOTEL(cfg *conf.AppConfiguration, storekey string) (*conf.Configuration, error) {
	storeCfg := cfg.Sub(storekey)
	driverName := storeCfg.String("driverName")
	if cfg.IsSet("otel") {
		var err error
		// Register the otelsql wrapper for the provided postgres driver.
		driverName, err = otelsql.Register("mysql",
			otelsql.WithAttributes(semconv.DBSystemMySQL),
			otelsql.WithAttributes(semconv.DBNameKey.String(storekey)),
			otelsql.WithSpanOptions(otelsql.SpanOptions{
				DisableErrSkip:  true,
				OmitRows:        true,
				OmitConnPrepare: true,
			}),
		)
		if err != nil {
			return nil, err
		}
		storeCfg.Parser().Set("driverName", driverName)
	}
	return storeCfg, nil
}

// BuildOTELDriver 构建具有otel的sql驱动,一般该驱动需要优先调用.
func BuildOTELDriver(cnf *conf.AppConfiguration, storekey string) dialect.Driver {
	storeCfg, err := TryRegisterOTEL(cnf, storekey)
	if err != nil {
		panic(err)
	}
	db := sqlx.NewSqlDB(storeCfg)
	return sql.OpenDB(storeCfg.String("driverName"), db)
}
