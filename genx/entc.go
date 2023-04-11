package genx

import (
	"embed"
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	//go:embed template/*
	_templates embed.FS
)

// GlobalID is a global id template for Noder Query. Use with ChangeRelayNodeType().
//
// if you use GlobalID, must use GID as scalar type.
// and use ChangeRelayNodeType() in entgql.WithSchemaHook()
func GlobalID() entc.Option {
	return func(g *gen.Config) error {
		g.Templates = append(g.Templates, gen.MustParse(gen.NewTemplate("gql_globalid").
			Funcs(entgql.TemplateFuncs).
			ParseFS(_templates, "template/globalid.tmpl")))
		return nil
	}
}

func SimplePagination() entc.Option {
	return func(g *gen.Config) error {
		g.Templates = append(g.Templates, gen.MustParse(gen.NewTemplate("gql_pagination_simple").
			Funcs(entgql.TemplateFuncs).
			ParseFS(_templates, "template/gql_pagination_simple.tmpl")))
		return nil
	}
}

// ChangeRelayNodeType is a schema hook for change relay node type to GID. Use with GlobalID().
//
// add it to entgql.WithSchemaHook()
func ChangeRelayNodeType() entgql.SchemaHook {
	idType := ast.NonNullNamedType("GID", nil)
	return func(graph *gen.Graph, schema *ast.Schema) error {
		for _, field := range schema.Types["Query"].Fields {
			if field.Name == "node" {
				field.Arguments[0].Type = idType
			}
			if field.Name == "nodes" {
				field.Arguments[0].Type = ast.NonNullListType(idType, nil)
			}
		}
		return nil
	}
}

func ReplaceGqlMutationInput() entgql.ExtensionOption {
	rt := gen.MustParse(gen.NewTemplate("gql_mutation_input").
		Funcs(entgql.TemplateFuncs).
		ParseFS(_templates, "template/gql_mutation_input.tmpl")).SkipIf(skipMutationTemplate)
	return entgql.WithTemplates([]*gen.Template{
		entgql.CollectionTemplate,
		entgql.EnumTemplate,
		entgql.NodeTemplate,
		entgql.PaginationTemplate,
		entgql.TransactionTemplate,
		entgql.EdgeTemplate,
		entgql.WhereTemplate,
		rt,
	}...)
}

func skipMutationTemplate(g *gen.Graph) bool {
	for _, n := range g.Nodes {
		ant, err := annotation(n.Annotations)
		if err != nil {
			continue
		}
		for _, i := range ant.MutationInputs {
			if (i.IsCreate && !ant.Skip.Is(entgql.SkipMutationCreateInput)) ||
				(!i.IsCreate && !ant.Skip.Is(entgql.SkipMutationUpdateInput)) {
				return false
			}
		}
	}
	return true
}

// annotation extracts the entgql.Annotation or returns its empty value.
func annotation(ants gen.Annotations) (*entgql.Annotation, error) {
	ant := &entgql.Annotation{}
	if ants != nil && ants[ant.Name()] != nil {
		if err := ant.Decode(ants[ant.Name()]); err != nil {
			return nil, err
		}
	}
	return ant, nil
}
