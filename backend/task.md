# Task

Into goalPostRow, we shall add a new column named `googlePingedAt`.
Type: Unix seconds UTC, not null.
`0` means never pinged.
Source code file: `server\db_objects\goalPostRow.go`
Also add database migration to  the table: `.\server\schema.postgre.sql`
Add the column if not exists.
