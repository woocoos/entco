//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	entcachegen "github.com/woocoos/entcache/gen"
	"github.com/woocoos/entco/genx"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		genx.WithGqlWithTemplates(),
		entgql.WithWhereInputs(true),
		entgql.WithConfigPath("./gqlgen.yml"),
		entgql.WithSchemaPath("./ent.graphql"),
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
	err = entc.Generate("./ent/schema", &gen.Config{},
		opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
