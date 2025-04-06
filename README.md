# Gator

A CLI tool to aggregate RSS feeds for multiple users.

## Prerequisites

To run Gator, you will need:

- [Go](https://golang.org/dl/) — The Go programming language for building and installing the CLI.
- [PostgreSQL](https://www.postgresql.org/download/) — A running Postgres database for storing user and feed data.

Make sure both are installed and properly set up on your machine before proceeding.

---

## Installation

Make sure you have the latest [Go toolchain](https://golang.org/dl/) installed as well as a local Postgres database.

You can then install `gator` with:

```bash
go install github.com/sabrek15/gator@latest
```

## Config

Create a `.gatorconfig.json` file in your home directory with the following structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your database connection string.

## Usage

Create a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start the aggregator:

```bash
gator agg 30s
```

View the posts:

```bash
gator browse [limit]
```

There are a few other commands you'll need as well:

- `gator login <name>` - Log in as a user that already exists

- `gator users` - List all users

- `gator feeds` - List all feeds

- `gator follow <url>` - Follow a feed that already exists in the database

- `gator unfollow <url>` - Unfollow a feed that already exists in the database

## Notes
- Ensure Postgres is running locally before using the CLI.

- Aggregator interval like 30s can be changed based on your preference (10s, 1m, 5m etc).

- All commands should be run after setting up the config file.

## Contribution

Contributions are welcome!

If you'd like to contribute to this project:

1. Fork the repository

2. Create a new branch

3. Make your changes

4. Open a Pull Request

Make sure to follow best practices and keep your code clean and well-documented.