# Refactoring database.goals.go

See files:
* C:\Dev\hinst-website\backend\server\schema.postgre.sql
* C:\Dev\hinst-website\backend\server\database.go
* C:\Dev\hinst-website\backend\server\database.goals.go

We want to refactor `database.goals.go`.
Move functions into separated files, grouped by table.
For example, function `setGoalPostPublic` goes into file `database.goalPost.go`.
Golang has no restriction on function placement within package, therefore the refactoring should be easy.
If a function touches different tables, then put them into file `database.compound.go`.
