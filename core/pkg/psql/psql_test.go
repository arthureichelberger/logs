package psql_test

import (
	"context"
	"testing"

	"github.com/arthureichelberger/logs/pkg/psql"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("it should return an error if the credentials are not right", func(t *testing.T) {
		db, err := psql.Connect(ctx, "wrong", "wrong", "localhost", "5432", "wrong")
		require.Error(t, err)
		require.Nil(t, db)
	})

	t.Run("it should return the queryable interface if credentials are right", func(t *testing.T) {
		db, err := psql.Connect(ctx, "logs", "logs", "localhost", "5432", "logs")
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Implements(t, (*psql.Queryable)(nil), db)
	})
}
