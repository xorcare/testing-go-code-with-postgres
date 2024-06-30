# Example of testing Go code with Postgres

[![Go workflow status badge](https://github.com/xorcare/testing-go-code-with-postgres/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/xorcare/testing-go-code-with-postgres/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/xorcare/testing-go-code-with-postgres/branch/main/graph/badge.svg?token=AmPmVHf2ej)](https://codecov.io/github/xorcare/testing-go-code-with-postgres/tree/main)
[![Gitlab CI Pipeline](https://gitlab.com/xorcare/testing-go-code-with-postgres/badges/main/pipeline.svg)](https://gitlab.com/xorcare/testing-go-code-with-postgres/tree/main)
[![Gitlab CI Coverage](https://gitlab.com/xorcare/testing-go-code-with-postgres/badges/main/coverage.svg)](https://gitlab.com/xorcare/testing-go-code-with-postgres/tree/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xorcare/testing-go-code-with-postgres)](https://goreportcard.com/report/github.com/xorcare/testing-go-code-with-postgres)

The example suggests a solution to the problem of cleaning the database after
running tests and the problem of running tests in parallel. It also shows how
to organize integration testing of Go code with Postgres.

## Quick start

For quickly try integration tests locally, use following commands.

```shell
git clone https://github.com/xorcare/testing-go-code-with-postgres
cd testing-go-code-with-postgres
make test-env-up test
```

<details>
<summary>Example of output</summary>

```text
❯ git clone https://github.com/xorcare/testing-go-code-with-postgres
cd testing-go-code-with-postgres
make test-env-up test
Cloning into 'testing-go-code-with-postgres'...
remote: Enumerating objects: 103, done.
remote: Counting objects: 100% (45/45), done.
remote: Compressing objects: 100% (24/24), done.
remote: Total 103 (delta 26), reused 29 (delta 20), pack-reused 58
Receiving objects: 100% (103/103), 27.58 KiB | 3.94 MiB/s, done.
Resolving deltas: 100% (40/40), done.
[+] Running 15/15
 ✔ migrate 5 layers [⣿⣿⣿⣿⣿]      0B/0B      Pulled                                             5.0s 
   ✔ 08409d417260 Pull complete                                                                1.5s 
   ✔ 2f9061c5186e Pull complete                                                                0.8s 
   ✔ de4eb1257b2b Pull complete                                                                2.2s 
   ✔ 750ec3989a15 Pull complete                                                                1.6s 
   ✔ 586322a68347 Pull complete                                                                2.2s 
 ✔ postgres 8 layers [⣿⣿⣿⣿⣿⣿⣿⣿]      0B/0B      Pulled                                        15.2s 
   ✔ 9fda8d8052c6 Pull complete                                                                2.5s 
   ✔ b0d9bb38da5c Pull complete                                                                2.8s 
   ✔ a99f2e61e525 Pull complete                                                                2.8s 
   ✔ eb307cc1ffd3 Pull complete                                                               11.1s 
   ✔ 99aedaa309df Pull complete                                                                4.0s 
   ✔ 1d4087443ab6 Pull complete                                                                3.5s 
   ✔ 278b6fc01aef Pull complete                                                                4.2s 
   ✔ 024b1a6a5c4d Pull complete                                                                4.9s 
[+] Building 0.0s (0/0)                                                              docker:default
[+] Running 3/2
 ✔ Network testing-go-code-with-postgres_default       Created                                 0.0s 
 ✔ Container testing-go-code-with-postgres-postgres-1  Created                                 0.2s 
 ✔ Container testing-go-code-with-postgres-migrate-1   Created                                 0.0s 
Attaching to testing-go-code-with-postgres-migrate-1
testing-go-code-with-postgres-migrate-1  | 1/u create_users_table (4.481416ms)
testing-go-code-with-postgres-migrate-1 exited with code 0
Aborting on container exit...
[+] Stopping 1/0
 ✔ Container testing-go-code-with-postgres-migrate-1  Stopped                                  0.0s 
ok  	github.com/xorcare/testing-go-code-with-postgres	1.500s	coverage: 100.0% of statements
ok  	github.com/xorcare/testing-go-code-with-postgres/testingpg	1.764s	coverage: 100.0% of statements
total:	(statements)	100.0%
```

</details>

## What's interesting here?

- Example
  of [docker-compose.yml](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/docker-compose.yml)
  with multiple databases and automated migrations.
- Example of test database connection management
  in [testingpg](https://github.com/xorcare/testing-go-code-with-postgres/tree/main/testingpg)
  package.
- [Example of integration testing with isolated database for each testcase](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/user_repository_with_isolated_database_test.go).
- And example
  of [GitHub Actions](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/.github/workflows/go.yml)
  and [Gitlab CI](https://github.com/xorcare/testing-go-code-with-postgres/blob/main/.gitlab-ci.yml).

Generating human-readable database names from `t.Name()` to simplifying problem investigation.
The last 8 characters are a short unique identifier needed to prevent name collision, its necessary
because the maximum length of the name is 63 bytes, and the name must be unique.

<details>
<summary>Example of test names</summary>

```text
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

</details>

## Known issues

When using **colima** on macos you may have problems if you clone this project to a temporary
directory like this:

```text
/var/folders/3p/glp5vp4916n03wmjh_b0gf6m0000gn/T/tmp.lbM4pbW2/testing-go-code-with-postgres
```

This problem is caused by incorrect mounting of files, and looks like this:

<details>
<summary>Example of output</summary>

```text
/var/folders/3p/glp5vp4916n03wmjh_b0gf6m0000gn/T/tmp.lbM4pbW2/testing-go-code-with-postgres
❯ docker-compose up
...
testing-go-code-with-postgres-postgres-1  | /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/docker-multiple-databases.sh
testing-go-code-with-postgres-postgres-1  | /usr/local/bin/docker-entrypoint.sh: line 170: /docker-entrypoint-initdb.d/docker-multiple-databases.sh: Is a directory
testing-go-code-with-postgres-postgres-1 exited with code 126
dependency failed to start: container testing-go-code-with-postgres-postgres-1 exited (126)
```

</details>

## Disclaimer

**This example is not an example of software architecture!**
