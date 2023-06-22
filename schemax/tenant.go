package schemax

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"errors"
	"fmt"
	casbinerr "github.com/casbin/casbin/v2/errors"
	"github.com/tsingsun/woocoo/pkg/authz"
	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/woocoos/entco/pkg/authorization"
	"github.com/woocoos/entco/pkg/identity"
	"strconv"
	"strings"
)

var (
	FieldTenantID = "tenant_id"
)

// TenantMixin helps to generate a tenant_id field and inject resource query.
//
//	 type World struct {
//		    ent.Schema
//	 }
//
//	 func (World) Mixin() []ent.Mixin {
//		    return []ent.Mixin{
//		    	schemax.Tenant{},
//		    }
//	 }
type TenantMixin[T Query, Q Mutator] struct {
	mixin.Schema
	QueryFunc func(ent.Query) (T, error)
}

func NewTenantMixin[T Query, Q Mutator](qf func(ent.Query) (T, error)) TenantMixin[T, Q] {
	return TenantMixin[T, Q]{
		QueryFunc: qf,
	}
}

func (TenantMixin[T, Q]) Fields() []ent.Field {
	return []ent.Field{
		field.Int(FieldTenantID).Immutable().SchemaType(IntID{}.SchemaType()),
	}
}

type tenantKey struct{}

// SkipTenantKey returns a new context that skips the soft-delete interceptor/mutators.
func SkipTenantKey(parent context.Context) context.Context {
	return context.WithValue(parent, tenantKey{}, true)
}

// Interceptors of the SoftDeleteMixin.
func (d TenantMixin[T, Q]) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			// Skip soft-delete, means include soft-deleted entities.
			if skip, _ := ctx.Value(tenantKey{}).(bool); skip {
				return nil
			}

			df, err := d.QueryFunc(q)
			if err != nil {
				return err
			}
			return d.QueryRulesP(ctx, df)
		}),
	}
}

type tenant[Q Mutator] interface {
	Query
	Client() Q
	SetTenantID(int)
}

// Hooks of the SoftDeleteMixin.
func (d TenantMixin[T, Q]) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if skip, _ := ctx.Value(tenantKey{}).(bool); skip {
					return next.Mutate(ctx, m)
				}

				tid, err := identity.TenantIDFromContext(ctx)
				if err != nil {
					return nil, fmt.Errorf("get tenant id from context: %w", err)
				}

				mx, ok := m.(tenant[Q])
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				switch m.Op() {
				case ent.OpCreate:
					mx.SetTenantID(tid)
				default:
					d.P(mx, tid)
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}

// P adds a storage-level predicate to the queries and mutations.
func (d TenantMixin[T, Q]) P(w Query, tid int) {
	w.WhereP(
		sql.FieldEQ(FieldTenantID, tid),
	)
}

func (d TenantMixin[T, Q]) QueryRulesP(ctx context.Context, w Query) error {
	typ := ent.QueryFromContext(ctx).Type
	uid, err := identity.UserIDFromContext(ctx)
	if err != nil {
		return err
	}
	tid, err := identity.TenantIDFromContext(ctx)
	if err != nil {
		return err
	}

	if authz.DefaultAuthorization == nil {
		d.P(w, tid)
		return nil
	}
	tidstr := strconv.Itoa(tid)
	uidstr := strconv.Itoa(uid)
	prefix := authorization.FormatArnPrefix(conf.Global().AppName(), tidstr, typ)
	flts, err := authorization.GetAllowedObjectConditions(uidstr, "read", prefix, tidstr)
	if err != nil && !errors.Is(err, casbinerr.ErrEmptyCondition) {
		return err
	}

	w.WhereP(func(selector *sql.Selector) {
		rules := GetTenantRules(flts, tidstr, selector)
		if len(rules) > 0 {
			selector.Where(sql.Or(rules...))
		}
	})
	return nil
}

// GetTenantRules returns the tenant resource conditions for the current user.
// if field rule is not has value after "/", it will be ignore, and like * effect.
func GetTenantRules(filers []string, tid string, selector *sql.Selector) []*sql.Predicate {
	v := make([]*sql.Predicate, 0, len(filers))
	for _, flt := range filers {
		if flt == "" {
			v = append(v, sql.EQ(selector.C(FieldTenantID), tid))
			continue
		}
		fs := strings.Split(flt, ":")
		fv := make([]*sql.Predicate, 0, len(fs))
		for _, f := range fs {
			kvs := strings.Split(f, "/")
			if len(kvs) != 2 {
				continue
			}
			fv = append(fv, sql.EQ(selector.C(kvs[0]), kvs[1]))
		}
		if len(fv) == 1 {
			v = append(v, fv...)
		} else {
			v = append(v, sql.And(fv...))
		}
	}
	return v
}
