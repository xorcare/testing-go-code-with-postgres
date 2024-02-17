// Copyright (c) 2023-2024 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing_go_code_with_postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	rootpkg "github.com/xorcare/testing-go-code-with-postgres"
	"github.com/xorcare/testing-go-code-with-postgres/testingpg"
)

func TestUserRepository_CreateUser(t *testing.T) {
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
		t.Parallel()

		// Arrange
		postgres := testingpg.New(t)
		repo := rootpkg.NewUserRepository(postgres.DB())

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
		postgres := testingpg.New(t)
		repo := rootpkg.NewUserRepository(postgres.DB())

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

func TestUserRepository_ReadUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Parallel()

	t.Run("Get an error if the user does not exist", func(t *testing.T) {
		t.Parallel()

		// Arrange
		postgres := testingpg.New(t)
		repo := rootpkg.NewUserRepository(postgres.DB())

		// Act
		_, err := repo.ReadUser(context.Background(), uuid.New())

		// Assert
		require.Error(t, err)
	})
}
