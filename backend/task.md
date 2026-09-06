# Smart Progress Importer refactoring

Looking at file `C:\Dev\hinst-website\backend\server\smartProgressImporter.go`
Function `saveGoalInfo` currently uses `database.pool` directly. This code should be refactored.
Instead of supplying direct SQL, we reuse object from `C:\Dev\hinst-website\backend\server\db_objects\goalRow.go` named GoalRow.
Add function named `saveGoal` into `C:\Dev\hinst-website\backend\server\database.goals.go` and use it.
See also: `C:\Dev\hinst-website\backend\server\schema.postgre.sql`
