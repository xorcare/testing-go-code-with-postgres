// Copyright (c) 2023 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testingpg

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

type TestingT interface {
	require.TestingT

	Cleanup(f func())
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Failed() bool
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
	newDBName := newUniqueHumanReadableDatabaseName(p.t)

	p.t.Log("database name for this test:", newDBName)

	sql := fmt.Sprintf(
		`CREATE DATABASE %q WITH TEMPLATE %q;`,
		newDBName,
		p.ref,
	)

	_, err := p.PgxPool().Exec(context.Background(), sql)
	require.NoError(p.t, err)

	// Automatically drop database copy after the test is completed.
	p.t.Cleanup(func() {
		sql := fmt.Sprintf(`DROP DATABASE %q WITH (FORCE);`, newDBName)

		ctx, done := context.WithTimeout(context.Background(), time.Minute)
		defer done()

		_, err := p.PgxPool().Exec(ctx, sql)
		require.NoError(p.t, err)
	})

	return &Postgres{
		t: p.t,

		url: replaceDBName(p.t, p.URL(), newDBName),
		ref: newDBName,
	}
}

func newUniqueHumanReadableDatabaseName(t TestingT) string {
	output := strings.Builder{}

	// Reports the maximum identifier length. It is determined as one less
	// than the value of NAMEDATALEN when building the server. The default
	// value of NAMEDATALEN is 64; therefore the default max_identifier_length
	// is 63 bytes, which can be less than 63 characters when using multibyte
	// encodings.
	// See https://www.postgresql.org/docs/15/runtime-config-preset.html
	const maxIdentifierLengthBytes = 63

	uid := genUnique8BytesID(t)
	maxHumanReadableLenBytes := maxIdentifierLengthBytes - len(uid)

	lastSymbolIsDash := false
	for _, r := range t.Name() {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			output.WriteRune(r)
			lastSymbolIsDash = false
		} else {
			if !lastSymbolIsDash {
				output.WriteRune('-')
			}
			lastSymbolIsDash = true
		}
		if output.Len() >= maxHumanReadableLenBytes {
			break
		}
	}

	output.WriteString(uid)

	return output.String()
}

func genUnique8BytesID(t TestingT) string {
	bs := make([]byte, 6)

	_, err := rand.Read(bs)
	require.NoError(t, err)

	return base64.RawURLEncoding.EncodeToString(bs)
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
