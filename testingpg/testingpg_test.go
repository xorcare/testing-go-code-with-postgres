// Copyright (c) 2023 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testingpg_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/xorcare/testing-go-code-with-postgres/testingpg"
)

func TestNewPostgres(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Parallel()

	t.Run("Successfully connect by URL and get version", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres := testingpg.New(t)

		ctx := context.Background()
		dbPool, err := pgxpool.New(ctx, postgres.URL())
		require.NoError(t, err)

		// Act
		var version string
		err = dbPool.QueryRow(ctx, "SELECT version();").Scan(&version)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, version)
		t.Log(version)
	})

	t.Run("Successfully obtained a version using a pre-configured conn", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres := testingpg.New(t)
		ctx := context.Background()

		// Act
		var version string
		err := postgres.PgxPool().QueryRow(ctx, "SELECT version();").Scan(&version)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, version)

		t.Log(version)
	})

	t.Run("Changes are not visible in different instances", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres1 := testingpg.New(t)
		postgres2 := testingpg.New(t)

		ctx := context.Background()

		// Act
		const sql = `CREATE TABLE "no_conflict" (id integer PRIMARY KEY)`
		_, err1 := postgres1.PgxPool().Exec(ctx, sql)
		_, err2 := postgres2.PgxPool().Exec(ctx, sql)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2, "databases must be isolated for each instance")
	})

	t.Run("URL is different at different instances", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres1 := testingpg.New(t)
		postgres2 := testingpg.New(t)

		// Act
		url1 := postgres1.URL()
		url2 := postgres2.URL()

		// Assert
		require.NotEqual(t, url1, url2)
	})
}
