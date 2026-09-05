# Smart Progress Importer integration

Looking at file `C:\Dev\hinst-website\backend\server\smartProgressImporter.go`
It currently uses `database.pool` directly. The goal is to stop using pool directly, and use functions defined in database instead.
Reuse objects from `C:\Dev\hinst-website\backend\server\db_objects`
If required function already exists, then reuse it.
If required function does not exist, then define it.
