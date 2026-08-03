# Task

Goal: into table `goalPosts`, add new field: `searchIndexingEnabled` of type boolean, default is false.

Relevant files:
* schema.postgre.sql
* goalPostRow.go

The field is already added. The next step is to add API for accessing it in `webAppGoals.go`
The API should be admin-gated: only admin may read and write the field. See examples in the file.
