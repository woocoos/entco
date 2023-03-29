package identity

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/web"
	"github.com/tsingsun/woocoo/web/handler"
	"net/http"
	"strconv"
)

var tenantContextKey = "github.com_woocoos_entco_tenant_id"

type TenantOptions struct {
	Lookup     string
	RootDomain string
}

// TenantIDMiddleware returns a middleware to get tenant id from http request
func TenantIDMiddleware(cfg *conf.Configuration) gin.HandlerFunc {
	opts := TenantOptions{
		Lookup: "header:X-Tenant-ID",
	}
	if err := cfg.Unmarshal(&opts); err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		tid := ""
		switch opts.Lookup {
		case "host":
			host := c.Request.Host
			if len(opts.RootDomain) > 0 {
				tid = host[:len(host)-len(opts.RootDomain)-1]
			}
		default:
			extr, err := handler.CreateExtractors(opts.Lookup, "")
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid tenant id %v", err))
				return
			}
			for _, extractor := range extr {
				ts, err := extractor(c)
				if err == nil && len(ts) != 0 {
					tid = ts[0]
					break
				}
			}
		}
		v, err := strconv.Atoi(tid)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid tenant id %s:%v", tid, err))
			return
		}
		c.Set(tenantContextKey, v)
	}
}

// RegistryTenantIDMiddleware register a middleware to get tenant id from request header
func RegistryTenantIDMiddleware() web.Option {
	return web.RegisterMiddlewareByFunc("tenant", TenantIDMiddleware)
}

// TenantIDFromContext returns the tenant id from context.tenant id is int format
func TenantIDFromContext[T int | string](ctx context.Context) (val T) {
	ginCtx := ctx.Value(gin.ContextKey).(*gin.Context)
	tid := ginCtx.Value(tenantContextKey)
	if tid == nil {
		return
	}
	switch tid.(type) {
	case int, string:
		return tid.(T)
	}
	panic(fmt.Errorf("invalid tenant id type %T", tid))
}
