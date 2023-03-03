package initx

import (
	"ariga.io/entcache"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/XSAM/otelsql"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/pkg/store/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

// BuildEntCacheDriver 构建ent缓存驱动
func BuildEntCacheDriver(cnf *conf.AppConfiguration, preDriver dialect.Driver) dialect.Driver {
	var cacheOpts []entcache.Option
	switch cnf.String("entcache.level") {
	case "context":
		// 使用Context缓存,但是不使用缓存的ttl
		cacheOpts = append(cacheOpts, entcache.ContextLevel())
	case "db":
		// 使用db缓存,如不设置TTL.
		if entCacheTTL := cnf.Duration("entcache.ttl"); entCacheTTL > 0 {
			cacheOpts = append(cacheOpts, entcache.TTL(entCacheTTL))
		}
	default:
		return preDriver // no cache
	}
	return entcache.NewDriver(preDriver, cacheOpts...)
}

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
