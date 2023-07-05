package genx

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

type DecimalExtension struct {
	entc.DefaultExtension
}

func (DecimalExtension) Hooks() []gen.Hook {
	return []gen.Hook{
		DecimalHook(),
	}
}

func (DecimalExtension) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(gen.NewTemplate("runtime").
			Funcs(entgql.TemplateFuncs).
			ParseFS(_templates, "template/runtime.tmpl")),
		gen.MustParse(gen.NewTemplate("meta").
			Funcs(entgql.TemplateFuncs).
			ParseFS(_templates, "template/meta.tmpl")),
		gen.MustParse(gen.NewTemplate("create").
			Funcs(entgql.TemplateFuncs).
			ParseFS(_templates, "template/create.tmpl")),
		gen.MustParse(gen.NewTemplate("update").
			Funcs(entgql.TemplateFuncs).
			ParseFS(_templates, "template/update.tmpl")),
	}
}

func DecimalHook() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, nodes := range g.Nodes {
				for _, f := range nodes.Fields {
					if f.Type.RType != nil && f.Type.RType.String() == "decimal.Decimal" {
						f.Type.Type = field.TypeFloat64
					}
				}
			}
			return next.Generate(g)
		})
	}
}
