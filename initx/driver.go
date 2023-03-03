package initx

import (
	"ariga.io/entcache"
	"entgo.io/ent/dialect"
	"github.com/tsingsun/woocoo/pkg/conf"
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
