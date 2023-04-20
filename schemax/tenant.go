package schemax

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Tenant helps to generate a tenant_id field.
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
//	 // if you use policy, you can add this snippet:
//	 func (World) Policy() ent.Policy {
//		    return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
//			    return rulex.FilterTenantRule(ctx, f)
//		    })
//	 }
type Tenant struct {
	mixin.Schema
}

func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tenant_id").Immutable().SchemaType(IntID{}.SchemaType()),
	}
}
