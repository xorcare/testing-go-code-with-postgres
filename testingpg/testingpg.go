// Copyright (c) 2023 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testingpg

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

type TestingT interface {
	require.TestingT

	Logf(format string, args ...any)
	Cleanup(f func())
}

func New(t TestingT) *Postgres {
	return newPostgres(t).cloneFromReference(t)
}

type Postgres struct {
	url string
	ref string

	conn *pgxpool.Pool
}

func newPostgres(t TestingT) *Postgres {
	urlStr := os.Getenv("TESTING_DB_URL")
	if urlStr == "" {
		urlStr = "postgresql://postgres:postgres@localhost:32260/postgres?sslmode=disable"
		const format = "env TESTING_DB_URL is empty, used default value: %s"
		t.Logf(format, urlStr)
	}

	refDatabase := os.Getenv("TESTING_DB_REF")
	if refDatabase == "" {
		refDatabase = "reference"
	}

	pool, err := pgxpool.New(context.Background(), urlStr)
	require.NoError(t, err)

	return &Postgres{
		url: urlStr,
		ref: refDatabase,

		conn: pool,
	}
}

func (p *Postgres) URL() string {
	return p.url
}

func (p *Postgres) PgxPool() *pgxpool.Pool {
	return p.conn
}

func (p *Postgres) cloneFromReference(t TestingT) *Postgres {
	cfg, err := pgxpool.ParseConfig(p.url)
	require.NoError(t, err)

	pool, err := pgxpool.New(context.Background(), p.url)
	require.NoError(t, err)

	newDatabaseName := uuid.New().String()

	const sqlTemplate = `CREATE DATABASE %q WITH TEMPLATE %s OWNER %s;`
	sql := fmt.Sprintf(
		sqlTemplate,
		newDatabaseName,
		p.ref,
		cfg.ConnConfig.User,
	)
	_, err = pool.Exec(context.Background(), sql)
	require.NoError(t, err)

	// Automatically drop database copy after the test is completed.
	t.Cleanup(func() {
		sql := fmt.Sprintf(`DROP DATABASE %q WITH (FORCE);`, newDatabaseName)

		ctx, done := context.WithTimeout(context.Background(), time.Minute)
		defer done()

		_, err := p.conn.Exec(ctx, sql)
		require.NoError(t, err)
	})

	urlString := replaceDBName(t, cfg, newDatabaseName)
	newPool, err := pgxpool.New(context.Background(), urlString)
	require.NoError(t, err)

	return &Postgres{
		url: urlString,
		ref: newDatabaseName,

		conn: newPool,
	}
}

func replaceDBName(t TestingT, cfg *pgxpool.Config, dbname string) string {
	r, err := url.Parse(cfg.ConnString())
	require.NoError(t, err)
	r.Path = dbname
	return r.String()
}
