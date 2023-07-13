package gentest

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/woocoos/entco/genx"
	"github.com/woocoos/entco/integration/gentest/ent"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSkipMigration(t *testing.T) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&_fk=1")
	assert.NoError(t, err)

	err = client.Schema.Create(context.Background(), genx.SkipTablesDiffHook("ref_table"))
	assert.NoError(t, err)

	_, err = client.RefSchema.Query().All(context.Background())
	assert.Error(t, err)
	assert.ErrorContains(t, err, "no such table: ref_table")
}
