# Smart Progress Importer refactoring

Looking at file `C:\Dev\hinst-website\backend\server\smartProgressImporter.go`
It currently uses `database.pool` directly. The goal is to stop using pool directly, and use functions defined in database instead.
Reuse objects from `C:\Dev\hinst-website\backend\server\db_objects`
If required function already exists, then reuse it.
If required function does not exist, then define it.
Look at the example in commit `5f9ec93818ca5c0319035d8340f0e5fa72395ba2` to see what sort of refactoring is expected.
