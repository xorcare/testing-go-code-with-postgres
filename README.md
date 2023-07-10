# Example of testing Go code with Postgres

[![Go workflow status badge](https://github.com/xorcare/testing-go-code-with-postgres/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/xorcare/testing-go-code-with-postgres/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/xorcare/testing-go-code-with-postgres/branch/main/graph/badge.svg?token=AmPmVHf2ej)](https://codecov.io/github/xorcare/testing-go-code-with-postgres/tree/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xorcare/testing-go-code-with-postgres)](https://goreportcard.com/report/github.com/xorcare/testing-go-code-with-postgres)

The example suggests a solution to the problem of cleaning the database after
running tests and the problem of running tests in parallel. It also shows how
to organize integration testing of Go code with Postgres.

## How to use

Run `make test-env-up test` and then everything will happen by itself.

## Disclaimer

**This example is not an example of software architecture!**
