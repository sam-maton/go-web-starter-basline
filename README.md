# go-web-starter-basline

## Database Migrations

The project uses Goose for database migrations. Migration files are located in `sql/schema/` and follow the naming convention `XXXXX_description.sql`.

### Creating a New Migration

```bash
goose -dir sql/schema create migration_name sql
```

### Running Migrations

```bash
# Migrate up (apply all pending migrations)
goose -dir sql/schema sqlite3 sql/database.db up

# Migrate down (rollback last migration)
goose -dir sql/schema sqlite3 sql/database.db down

# Check migration status
goose -dir sql/schema sqlite3 sql/database.db status
```
