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

## How to use

Run `make test-env-up test` and then everything will happen by itself.

## Disclaimer

**This example is not an example of software architecture!**
