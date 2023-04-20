package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/woocoos/entco/integration/helloapp/ent/privacy"
	"github.com/woocoos/entco/schemax"
	"github.com/woocoos/entco/schemax/rulex"
)

type World struct {
	ent.Schema
}

func (World) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schemax.IntID{},
		schemax.Tenant{},
	}
}

func (World) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

func (World) Policy() ent.Policy {
	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		return rulex.FilterTenantRule(ctx, f)
	})
}
