// This package is used to Debug
package main

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	entcachegen "github.com/woocoos/entcache/gen"
	"github.com/woocoos/entco/genx"
	"log"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		genx.WithGqlWithTemplates(),
		entgql.WithWhereInputs(true),
		entgql.WithConfigPath("gentest/gqlgen.yml"),
		entgql.WithSchemaPath("gentest/ent.graphql"),
		entgql.WithSchemaHook(genx.ChangeRelayNodeType()),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex, genx.DecimalExtension{}),
		genx.GlobalID(),
		genx.SimplePagination(),
		entcachegen.QueryCache(),
	}
	err = entc.Generate("./gentest/ent/schema", &gen.Config{
		Package: "github.com/woocoos/entco/integration/gentest/ent",
		Target:  "./gentest/ent",
	},
		opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
