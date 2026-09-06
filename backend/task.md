# Smart Progress Importer refactoring

Looking at file `C:\Dev\hinst-website\backend\server\smartProgressImporter.go`
Function `savePost` currently uses `database.pool` directly. This code should be refactored.
Instead of supplying direct SQL, we reuse object from `backend\server\db_objects\goalPostRow.go`
named GoalPostRow. If post already exists, then use function `setGoalPostText`,
otherwise use function `saveGoalPost`, which does not exist yet, therefore please define it.
