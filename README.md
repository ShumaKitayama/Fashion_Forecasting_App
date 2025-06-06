# TrendScout

This repository contains a fashion trend forecasting application with Go and React.

## Backend utilities

Several helper commands live under `backend/cmd`. To verify database tables or reset the admin account you can use the `dbcheck` command.

```bash
# Example: check users table
go run ./backend/cmd/dbcheck users

# Example: check keywords table
go run ./backend/cmd/dbcheck keywords

# Example: reset admin password to 'password'
go run ./backend/cmd/dbcheck reset
```

The actual logic for each subcommand resides in separate packages:

- `check_user`
- `check_keywords`
- `reset_password`

These are imported by `cmd/dbcheck/main.go`. If you see errors about "found packages" make sure your working directory contains the latest directory structure where each subcommand lives in its own folder.
