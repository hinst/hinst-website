# Smart Progress Importer refactoring

Looking at file `C:\Dev\hinst-website\backend\server\smartProgressImporter.go`
Function `saveImages` currently uses `database.pool` directly. This code should be refactored.
Instead of supplying direct SQL, we reuse object from `C:\Dev\hinst-website\backend\server\db_objects\goalPostImageRow.go` named GoalPostImageRow.
Add function named `saveGoalPostImage` into `C:\Dev\hinst-website\backend\server\database.goals.go` and use it.
See also: `C:\Dev\hinst-website\backend\server\schema.postgre.sql`
