// Copyright (c) 2023-2024 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing_go_code_with_postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) ReadUser(ctx context.Context, userID uuid.UUID) (User, error) {
	const sql = `SELECT user_id, username, created_at FROM users WHERE user_id = $1;`

	user := User{}

	row := r.db.QueryRowContext(ctx, sql, userID)

	err := row.Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		const format = "failed selection of User from database: %v"

		return User{}, fmt.Errorf(format, err)
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user User) error {
	const sql = `INSERT INTO users (user_id, username, created_at) VALUES ($1,$2,$3);`

	_, err := r.db.ExecContext(
		ctx,
		sql,
		user.ID,
		user.Username,
		user.CreatedAt,
	)

	if err != nil {
		const format = "failed insertion of User to database: %v"
		return fmt.Errorf(format, err)
	}

	return nil
}
