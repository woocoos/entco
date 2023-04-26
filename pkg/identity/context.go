package identity

import (
	"context"
	"errors"
	"github.com/tsingsun/woocoo/pkg/security"
	"strconv"
)

var (
	ErrInvalidUserID = errors.New("invalid user")
)

func UserIDFromContext(ctx context.Context) (int, error) {
	gp := security.GenericPrincipalFromContext(ctx)
	if gp == nil {
		return 0, ErrInvalidUserID
	}
	id, err := strconv.Atoi(gp.GenericIdentity.Name())
	if err != nil || id == 0 {
		return 0, ErrInvalidUserID
	}
	return id, nil
}
