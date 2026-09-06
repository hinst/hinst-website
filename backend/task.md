# Smart Progress Importer refactoring

Looking at file `C:\Dev\hinst-website\backend\server\smartProgressImporter.go`
Function `saveComments` currently uses `database.pool` directly. This code should be refactored.
Instead of supplying direct SQL, we add a new object into `C:\Dev\hinst-website\backend\server\db_objects`.
The object should be named GoalPostCommentRow. Register it with `registerDbObject` too.
Add function into `C:\Dev\hinst-website\backend\server\database.goals.go`: `saveGoalPostComment`
and use it `from smartProgressImporter.go`
