package authorization

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	"github.com/tsingsun/woocoo/pkg/authz"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/pkg/security"
	"github.com/woocoos/casbin-ent-adapter"
	"github.com/woocoos/casbin-ent-adapter/ent"
	"github.com/woocoos/entco/pkg/identity"
)

// SetAuthorization 设置授权器
func SetAuthorization(cnf *conf.Configuration, driver dialect.Driver) (authorizer *authz.Authorization, err error) {
	casbinClient := ent.NewClient(ent.Driver(driver))
	adp, err := entadapter.NewAdapterWithClient(casbinClient)
	if err != nil {
		return
	}
	err = casbinClient.Schema.Create(context.Background())
	if err != nil {
		return
	}
	authz.SetAdapter(adp)
	authorizer, err = authz.NewAuthorization(cnf, authz.WithRequestParseFunc(RBACWithDomainRequestParserFunc))
	if err != nil {
		return
	}
	authz.SetDefaultAuthorization(authorizer)
	return
}

// RBACWithDomainRequestParserFunc 以RBAC with domain模型生成casbin请求
//
// ctx: 一般就是gin.Context
func RBACWithDomainRequestParserFunc(ctx context.Context, id security.Identity, item *security.PermissionItem) []any {
	gctx := ctx.Value(gin.ContextKey).(*gin.Context)
	domain := gctx.GetHeader(identity.TenantHeaderKey)
	p := item.AppCode + ":" + item.Action
	return []any{id.Name(), domain, p, item.Operator}
}
