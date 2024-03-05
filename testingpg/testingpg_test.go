// Copyright (c) 2023-2024 Vasiliy Vasilyuk. All rights reserved.
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
		postgres := testingpg.NewWithIsolatedDatabase(t)

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
		postgres := testingpg.NewWithIsolatedDatabase(t)
		ctx := context.Background()

		// Act
		var version string
		err := postgres.DB().QueryRowContext(ctx, "SELECT version();").Scan(&version)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, version)

		t.Log(version)
	})

	t.Run("Changes are not visible in different instances", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres1 := testingpg.NewWithIsolatedDatabase(t)
		postgres2 := testingpg.NewWithIsolatedDatabase(t)

		ctx := context.Background()

		// Act
		const sqlStr = `CREATE TABLE "no_conflict" (id integer PRIMARY KEY)`
		_, err1 := postgres1.DB().ExecContext(ctx, sqlStr)
		_, err2 := postgres2.DB().ExecContext(ctx, sqlStr)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2, "databases must be isolated for each instance")
	})

	t.Run("URL is different at different instances", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres1 := testingpg.NewWithIsolatedDatabase(t)
		postgres2 := testingpg.NewWithIsolatedDatabase(t)

		// Act
		url1 := postgres1.URL()
		url2 := postgres2.URL()

		// Assert
		require.NotEqual(t, url1, url2)
	})
}

func TestNewWithIsolatedSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Parallel()

	t.Run("Successfully connect by URL and get version", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres := testingpg.NewWithIsolatedSchema(t)

		ctx := context.Background()
		dbPool, err := pgxpool.New(ctx, postgres.URL())
		require.NoError(t, err)

		// Act
		var version string
		err = dbPool.QueryRow(ctx, "SHOW search_path;").Scan(&version)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, version)
		t.Log(version)
	})

	t.Run("Successfully obtained a version using a pre-configured conn", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres := testingpg.NewWithIsolatedSchema(t)
		ctx := context.Background()

		// Act
		var version string
		err := postgres.DB().QueryRowContext(ctx, "SHOW search_path;").Scan(&version)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, version)

		t.Log(version)
	})

	t.Run("Changes are not visible in different instances", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres1 := testingpg.NewWithIsolatedSchema(t)
		postgres2 := testingpg.NewWithIsolatedSchema(t)

		ctx := context.Background()

		// Act
		const sqlStr = `CREATE TABLE "no_conflict" (id integer PRIMARY KEY)`
		_, err1 := postgres1.DB().ExecContext(ctx, sqlStr)
		_, err2 := postgres2.DB().ExecContext(ctx, sqlStr)

		// Assert
		require.NoError(t, err1)
		require.NoError(t, err2, "databases must be isolated for each instance")
	})

	t.Run("URL is different at different instances", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres1 := testingpg.NewWithIsolatedSchema(t)
		postgres2 := testingpg.NewWithIsolatedSchema(t)

		// Act
		url1 := postgres1.URL()
		url2 := postgres2.URL()

		// Assert
		require.NotEqual(t, url1, url2)
	})
}

func TestNewWithTransactionalCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Parallel()

	t.Run("Successfully obtained a version", func(t *testing.T) {
		t.Parallel()

		// Arrange
		tx := testingpg.NewWithTransactionalCleanup(t)
		ctx := context.Background()

		// Act
		var version string
		err := tx.QueryRowContext(ctx, "SELECT version();").Scan(&version)

		// Assert
		require.NoError(t, err)
		require.NotEmpty(t, version)

		t.Log(version)
	})

	t.Run("Changes are not visible in different instances", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctx := context.Background()
		const sqlStr = `CREATE TABLE "no_conflict" (id integer PRIMARY KEY)`

		t.Run("Arrange", func(t *testing.T) {
			tx := testingpg.NewWithTransactionalCleanup(t)
			_, err := tx.ExecContext(ctx, sqlStr)
			require.NoError(t, err)
		})

		var err error

		// Act
		t.Run("Act", func(t *testing.T) {
			tx := testingpg.NewWithTransactionalCleanup(t)
			_, err = tx.ExecContext(ctx, sqlStr)
		})

		// Assert
		require.NoError(t, err, "side effects must be isolated for each instance")
	})
}
