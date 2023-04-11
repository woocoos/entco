package integration

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"

	"entgo.io/contrib/entgql"
	"github.com/woocoos/entco/genx/integration/ent"
)

func (r *queryResolver) Node(ctx context.Context, id string) (ent.Noder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Nodes(ctx context.Context, ids []string) ([]ent.Noder, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, orderBy *ent.UserOrder, where *ent.UserWhereInput) (*ent.UserConnection, error) {
	gctx := ctx.Value(gin.ContextKey).(*gin.Context)
	sp, err := ent.NewSimplePagination(gctx.Query("p"), gctx.Query("c"))
	if err != nil {
		return nil, err
	}
	return r.client.User.Query().SimplePaginate(ctx, sp, after, first, before, last)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
