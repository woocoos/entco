package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"github.com/woocoos/entco/schemax/fieldx"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// annotation
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField("users"),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(45).Comment("姓名"),
		field.Time("created_at").Immutable().Default(time.Now).Immutable().
			Annotations(entgql.OrderField("createdAt"), entgql.Skip(entgql.SkipMutationCreateInput),
				entproto.Field(3)),
		fieldx.Decimal("money").Precision(10, 6).Optional().
			Range(decimal.NewFromInt(1), decimal.NewFromInt(100000)).
			Comment("money").Nillable().Default("2"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
