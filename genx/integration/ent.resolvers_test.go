package integration

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"github.com/tsingsun/woocoo/pkg/gds"
	"github.com/woocoos/entco/genx/integration/ent"
	"net/http/httptest"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type TestSuite struct {
	suite.Suite
	client        *ent.Client
	queryResolver queryResolver
}

func (s *TestSuite) SetupSuite() {
	dr, err := sql.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	s.Require().NoError(err)
	s.client = ent.NewClient(ent.Driver(dr), ent.Debug())
	s.queryResolver = queryResolver{&Resolver{s.client}}
	s.NoError(s.client.Schema.Create(context.Background()))
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestUsers_SimplePagination() {
	builder := make([]*ent.UserCreate, 0)
	for i := 0; i < 20; i++ {
		builder = append(builder, s.client.User.Create().SetName("user"+strconv.Itoa(i)))
	}
	s.NoError(s.client.User.CreateBulk(builder...).Exec(context.Background()))

	s.Run("pagination after", func() {
		ctx, _ := gin.CreateTestContext(nil)
		users, err := s.queryResolver.Users(ctx, &ent.Cursor{ID: 2}, gds.Ptr(2), nil, nil, nil, nil)
		s.NoError(err)
		s.Len(users.Edges, 2)
		s.Equal(3, users.Edges[0].Node.ID)
		s.Equal(4, users.Edges[1].Node.ID)

	})

	s.Run("simple after", func() {
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = httptest.NewRequest("GET", "/?p=3&c=1", nil)
		users, err := s.queryResolver.Users(ctx, &ent.Cursor{ID: 2}, gds.Ptr(2), nil, nil, nil, nil)
		s.NoError(err)
		s.Len(users.Edges, 2)
		s.Equal(5, users.Edges[0].Node.ID)
		s.Equal(6, users.Edges[1].Node.ID)

	})

	s.Run("pagination befor", func() {
		ctx, _ := gin.CreateTestContext(nil)
		users, err := s.queryResolver.Users(ctx, nil, nil, &ent.Cursor{ID: 5}, gds.Ptr(2), nil, nil)
		s.NoError(err)
		s.Len(users.Edges, 2)
		s.Equal(3, users.Edges[0].Node.ID)
		s.Equal(4, users.Edges[1].Node.ID)

	})

	s.Run("simple before", func() {
		ctx, _ := gin.CreateTestContext(nil)
		ctx.Request = httptest.NewRequest("GET", "/?p=1&c=3", nil)
		users, err := s.queryResolver.Users(ctx, nil, nil, &ent.Cursor{ID: 5}, gds.Ptr(2), nil, nil)
		s.NoError(err)
		s.Len(users.Edges, 2)
		s.Equal(1, users.Edges[0].Node.ID)
		s.Equal(2, users.Edges[1].Node.ID)
	})

}
