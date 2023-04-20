//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/woocoos/entco/genx"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithWhereInputs(true),
		entgql.WithConfigPath("./gqlgen.yml"),
		entgql.WithSchemaPath("./ent.graphql"),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
		genx.SimplePagination(),
	}
	err = entc.Generate("./ent/schema", &gen.Config{},
		opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
