package helloapp

import (
	"context"
	"entgo.io/ent/privacy"
	"errors"
	"github.com/woocoos/entco/integration/helloapp/ent"
	"github.com/woocoos/entco/pkg/identity"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/woocoos/entco/integration/helloapp/ent/runtime"
)

func Test_CreateWorld(t *testing.T) {
	ctx := context.Background()
	client := open(ctx)
	defer client.Close()

	if err := client.World.Create().Exec(ctx); !errors.Is(err, privacy.Deny) {
		t.Fatal("expect tenant creation to fail, but got:", err)
	}

	tctx := identity.WithTenantID(ctx, 1)
	if err := client.World.Create().SetName("woocoo").SetTenantID(1).Exec(tctx); err != nil {
		t.Fatal("expect tenant creation to succeed, but got:", err)
	}

	if _, err := client.World.Query().Count(ctx); !errors.Is(err, privacy.Deny) {
		t.Fatal("expect tenant query to fail, but got:", err)
	}

	if _, err := client.World.Query().Count(tctx); err != nil {
		t.Fatal("expect tenant query to succeed, but got:", err)
	}
}

func open(ctx context.Context) *ent.Client {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", ent.Debug())
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
