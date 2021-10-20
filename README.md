# mintsql

[![Go](https://github.com/AustinZhu/mintsql/actions/workflows/go.yml/badge.svg)](https://github.com/AustinZhu/mintsql/actions/workflows/go.yml)
[![Docker](https://github.com/AustinZhu/mintsql/actions/workflows/docker.yml/badge.svg)](https://github.com/AustinZhu/mintsql/actions/workflows/docker.yml)

![](https://i.ibb.co/CPZ0VMc/mint-logo.png)
MintSQL is a lightweight SQL database written in Go that supports basic database CRUD operations.

## Getting Started
Run the server (on local):

`make run-server`

Run the server (in docker container):

`make docker-run`

Run the client:

`make run-client`

## Installation

Install the server:

`make install-server`

Install the client:

`make install-client`

## Usages
The default port for MintSQL is 7384.
### Server
`mintsql <port>`

### Client
`mintcli <host> <port>`

### Querying from the client console

The query language is a subset of SQL. For now, MintSQL only support simple `CREATE|INSERT|SELECT` queries.

To quit the client console, type `\quit`.

## Examples

- Create a table:

```
> CREATE TABLE students (id INT, name TEXT);
ok
```

- Insert values into a table:

```
> INSERT INTO students VALUES (1, 'Bob');
ok
```

- Select data in a table:

```
> SELECT name, id FROM students;
| name | id |
==============
| Bob |  1 | 
```
