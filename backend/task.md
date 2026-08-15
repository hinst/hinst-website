See file: `server\webAppAdmin.go`
Instead of using `me.checkSearchIndexing(context.Background(), header.PublicUrl)` for every HTTP response,
we move search indexing update into function `update()` in `server\program.go`.
The search indexing status update should be the last step after `me.uploadStatic()` call.
Remember to update `GoogleSearchIndexingStatusCheckedAt` timestamp as well.
See file `server\db_objects\goalPostRow.go`.
