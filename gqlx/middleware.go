package gqlx

import (
	"ariga.io/entcache"
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/tsingsun/woocoo/contrib/gql"
	"github.com/tsingsun/woocoo/web"
	"github.com/tsingsun/woocoo/web/handler"
	"github.com/tsingsun/woocoo/web/handler/signer"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/woocoos/entco/pkg/pagination"
)

// SimplePagination is a middleware that parses the query string for the simple (similar limit,offset) pagination
// use it like:
//
//	gqlsrv.AroundResponses(gqlx.SimplePagination())
func SimplePagination() graphql.ResponseMiddleware {
	return func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		if op := graphql.GetOperationContext(ctx).Operation; op != nil && op.Operation != ast.Query {
			return next(ctx)
		}
		gctx, _ := gql.FromIncomingContext(ctx)
		if gctx != nil {
			sp, err := pagination.NewSimplePagination(gctx.Query("p"), gctx.Query("c"))
			if err != nil {
				return graphql.ErrorResponse(ctx, "pagination error:%v", err)
			}
			if sp != nil {
				ctx = pagination.WithSimplePagination(ctx, sp)
			}
		}
		return next(ctx)
	}
}

// ContextCache is a middleware for entcache
func ContextCache() graphql.ResponseMiddleware {
	return func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		if op := graphql.GetOperationContext(ctx).Operation; op != nil && op.Operation == ast.Query {
			ctx = entcache.NewContext(ctx)
		}
		return next(ctx)
	}
}

func RegisterTokenSignerMiddleware() web.Option {
	return web.WithMiddlewareNewFunc(signer.TokenSignerName, func() handler.Middleware {
		mw := signer.NewMiddleware(signer.TokenSignerName, handler.WithMiddlewareConfig(func(config any) {
			c := config.(*signer.Config)
			c.Skipper = func(c *gin.Context) bool {
				if c.IsWebsocket() {
					return true
				}
				return false
			}
		}))
		return mw
	})
}
