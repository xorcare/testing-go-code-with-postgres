// Copyright (c) 2023-2024 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing_go_code_with_postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	rootpkg "github.com/xorcare/testing-go-code-with-postgres"
	"github.com/xorcare/testing-go-code-with-postgres/migrations"
	"github.com/xorcare/testing-go-code-with-postgres/testingpg"
)

func migrateDatabaseSchema(t *testing.T, pg *testingpg.Postgres) {
	source, err := iofs.New(migrations.FS, ".")
	require.NoError(t, err)

	mi, err := migrate.NewWithSourceInstance(
		"iofs",
		source,
		pg.URL(),
	)
	require.NoError(t, err)

	err = mi.Up()

	if !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}
}

func Test_Schema_UserRepository_CreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Parallel()

	newFullyFiledUser := func() rootpkg.User {
		return rootpkg.User{
			ID:        uuid.New(),
			Username:  "gopher",
			CreatedAt: time.Now().Truncate(time.Microsecond),
		}
	}

	t.Run("Successfully created a User", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		t.Parallel()

		// Arrange
		pg := testingpg.NewWithIsolatedSchema(t)

		migrateDatabaseSchema(t, pg)

		repo := rootpkg.NewUserRepository(pg.DB())
		user := newFullyFiledUser()

		// Act
		err := repo.CreateUser(context.Background(), user)

		// Assert
		require.NoError(t, err)

		gotUser, err := repo.ReadUser(context.Background(), user.ID)
		require.NoError(t, err)

		require.Equal(t, user, gotUser)
	})

	t.Run("Cannot create a user with the same ID", func(t *testing.T) {
		t.Parallel()

		// Arrange
		pg := testingpg.NewWithIsolatedSchema(t)

		migrateDatabaseSchema(t, pg)

		repo := rootpkg.NewUserRepository(pg.DB())

		user := newFullyFiledUser()

		err := repo.CreateUser(context.Background(), user)
		require.NoError(t, err)

		// Act
		err = repo.CreateUser(context.Background(), user)

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "duplicate key value violates unique constraint")
	})
}

func Test_Schema_UserRepository_ReadUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Parallel()

	t.Run("Get an error if the user does not exist", func(t *testing.T) {
		t.Parallel()

		// Arrange
		pg := testingpg.NewWithIsolatedSchema(t)

		migrateDatabaseSchema(t, pg)

		repo := rootpkg.NewUserRepository(pg.DB())

		// Act
		_, err := repo.ReadUser(context.Background(), uuid.New())

		// Assert
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
