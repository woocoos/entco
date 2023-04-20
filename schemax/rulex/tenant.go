package rulex

import (
	"context"
	"entgo.io/ent/entql"
	"entgo.io/ent/privacy"
	"fmt"
	"github.com/woocoos/entco/pkg/identity"
)

type (
	// Filter is the interface same as privacy.Filter
	Filter interface {
		Where(entql.P)
	}
	// TenantsFilter matches tenants filter
	TenantsFilter interface {
		WhereTenantID(entql.IntP)
	}
	FilterFunc func(context.Context, Filter) error
)

// FilterTenantRule is a privacy rule that filters the query to return only
// Example:
//
//	 func (World) Policy() ent.Policy {
//		    return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
//			    return rulex.FilterTenantRule(ctx, f)
//		    })
//	 }
func FilterTenantRule(ctx context.Context, f Filter) error {
	tid, err := identity.TenantIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get tenant id from context: %s %w", err, privacy.Deny)
	}
	tf, ok := f.(TenantsFilter)
	if !ok {
		return fmt.Errorf("unexpected filter type %T %w", f, privacy.Deny)
	}
	tf.WhereTenantID(entql.IntEQ(tid))
	return privacy.Skip
}
