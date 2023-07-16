# Example of testing Go code with Postgres

[![Go workflow status badge](https://github.com/xorcare/testing-go-code-with-postgres/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/xorcare/testing-go-code-with-postgres/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/xorcare/testing-go-code-with-postgres/branch/main/graph/badge.svg?token=AmPmVHf2ej)](https://codecov.io/github/xorcare/testing-go-code-with-postgres/tree/main)
[![Gitlab CI Pipeline](https://gitlab.com/xorcare/testing-go-code-with-postgres/badges/main/pipeline.svg)](https://gitlab.com/xorcare/testing-go-code-with-postgres/tree/main)
[![Gitlab CI Coverage](https://gitlab.com/xorcare/testing-go-code-with-postgres/badges/main/coverage.svg)](https://gitlab.com/xorcare/testing-go-code-with-postgres/tree/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xorcare/testing-go-code-with-postgres)](https://goreportcard.com/report/github.com/xorcare/testing-go-code-with-postgres)

The example suggests a solution to the problem of cleaning the database after
running tests and the problem of running tests in parallel. It also shows how
to organize integration testing of Go code with Postgres.

## What's interesting here?

- Example
  of [docker-compose.yml](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/docker-compose.yml)
  with multiple databases and automated migrations.
- Example of test database connection management
  in [testingpg](https://github.com/xorcare/testing-go-code-with-postgres/tree/main/testingpg)
  package.
- Example of
  integration [tests](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/user_repository_test.go).
- And example
  of [GitHub Actions](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/.github/workflows/go.yml)
  and [Gitlab CI](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/.gitlab-ci.yml).

Generating human-readable database names from `t.Name()` to simplifying problem investigation.
The last 8 characters are a short unique identifier needed to prevent name collision, its necessary
because the maximum length of the name is 63 bytes, and the name must be unique.

```txt
TestNewPostgres-Changes-are-not-visible-in-different-inWirPQD7J
TestNewPostgres-Changes-are-not-visible-in-different-ineYp0ljjI
TestNewPostgres-Successfully-connect-by-URL-and-get-verzGq4pGza
TestNewPostgres-Successfully-obtained-a-version-using-a20YgZaMf
TestNewPostgres-URL-is-different-at-different-instancesIMDkJgoP
TestNewPostgres-URL-is-different-at-different-instancesjtSsjPR5
TestUserRepository-CreateUser-Cannot-create-a-user-withmgmHFdZe
TestUserRepository-CreateUser-Successfully-created-a-UspTBGNltW
TestUserRepository-ReadUser-Get-an-error-if-the-user-doRqS1GvYh
```

## How to use

Run `make test-env-up test` and then everything will happen by itself.

## Disclaimer

**This example is not an example of software architecture!**
