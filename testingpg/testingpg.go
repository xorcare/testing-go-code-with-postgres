// Copyright (c) 2023 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testingpg

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
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
	return newPostgres(t).cloneFromReference()
}

type Postgres struct {
	t TestingT

	url string
	ref string

	pgxpool     *pgxpool.Pool
	pgxpoolOnce sync.Once
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

	return &Postgres{
		t: t,

		url: urlStr,
		ref: refDatabase,
	}
}

func (p *Postgres) URL() string {
	return p.url
}

func (p *Postgres) PgxPool() *pgxpool.Pool {
	p.pgxpoolOnce.Do(func() {
		p.pgxpool = newPGxPool(p.t, p.URL())
	})

	return p.pgxpool
}

func (p *Postgres) cloneFromReference() *Postgres {
	newDatabaseName := uuid.New().String()

	sql := fmt.Sprintf(
		`CREATE DATABASE %q WITH TEMPLATE %q;`,
		newDatabaseName,
		p.ref,
	)

	_, err := p.PgxPool().Exec(context.Background(), sql)
	require.NoError(p.t, err)

	// Automatically drop database copy after the test is completed.
	p.t.Cleanup(func() {
		sql := fmt.Sprintf(`DROP DATABASE %q WITH (FORCE);`, newDatabaseName)

		ctx, done := context.WithTimeout(context.Background(), time.Minute)
		defer done()

		_, err := p.PgxPool().Exec(ctx, sql)
		require.NoError(p.t, err)
	})

	return &Postgres{
		t: p.t,

		url: replaceDBName(p.t, p.URL(), newDatabaseName),
		ref: newDatabaseName,
	}
}

func replaceDBName(t TestingT, dataSourceURL, dbname string) string {
	r, err := url.Parse(dataSourceURL)
	require.NoError(t, err)
	r.Path = dbname
	return r.String()
}

func newPGxPool(t TestingT, dataSourceURL string) *pgxpool.Pool {
	ctx, done := context.WithTimeout(context.Background(), 1*time.Second)
	defer done()

	pool, err := pgxpool.New(ctx, dataSourceURL)
	require.NoError(t, err)

	// Automatically close connection after the test is completed.
	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
