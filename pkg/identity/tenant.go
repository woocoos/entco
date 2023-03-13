package identity

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/web"
	"net/http"
	"strconv"
)

var tenantContextKey = "github.com_woocoos_entco_tenant_id"

// RegistryTenantIDMiddleware register a middleware to get tenant id from request header
func RegistryTenantIDMiddleware() web.Option {
	return web.RegisterMiddlewareByFunc("tenant", func(cfg *conf.Configuration) gin.HandlerFunc {
		typ := cfg.String("type")
		return func(c *gin.Context) {
			tid := c.GetHeader("X-Tenant-ID")
			if tid != "" {
				if typ == "int" {
					if _, err := strconv.Atoi(tid); err != nil {
						c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid tenant id %s:%v", tid, err))
						return
					}
				} else {
					c.Set(tenantContextKey, tid)
				}
			}
		}
	})
}

// WithTenantID returns a new context with the tenant id.
func WithTenantID[T int | string](ctx context.Context, tid T) context.Context {
	return context.WithValue(ctx, tenantContextKey, tid)
}

// TenantIDFromContext returns the tenant id from context.tenant id is int format
func TenantIDFromContext[T int | string](ctx context.Context) (val T) {
	tid := ctx.Value(tenantContextKey)
	if tid == nil {
		return
	}
	return T(tid)
}
