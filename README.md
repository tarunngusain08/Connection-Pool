# Connection Pool

## Summary

This project aims to build a SingleStore connection pool that supports multiple users. While traditional connection pools have static configurations defining credentials for opening new connections, this project addresses more dynamic use cases. The goal is to enable a proxy server atop a SingleStore cluster, allowing different users to log in while managing the load on the underlying system effectively.

## Configuration

The connection pool should support the following configuration options:

| Option          | Description                                                                                          |
|-----------------|------------------------------------------------------------------------------------------------------|
| Connection Limit| The maximum number of connections permitted to the underlying SingleStore Cluster at any given time. |
| Idle Timeout    | Duration after which an idle connection should be closed.                                            |

## Library API

The library exposes a "Connection Pool" object that can be configured during construction and provides the following methods:

### Query(auth, database, sql) Result

Executes a query and returns the result.

- **Arguments**
    - **auth**: Authentication details for the user running the query. For SingleStore, this includes the "username" and "password".
    - **database**: The specific database context against which the query should run.
    - **sql**: The SQL query to be executed.
- **Result**
    - The return value enables inspection of metadata and reading of rows.

## Query Behavior

Under the hood, the Query function optimizes connection usage as follows:

- Utilizes idle connections matching requested connection details.
- Opens new connections if the pool is not full.
- Closes an idle connection associated with the least recently used details if the pool is full but has idle connections.
- Blocks until an idle connection becomes available to reuse or replace if the pool is full and has no idle connections.
