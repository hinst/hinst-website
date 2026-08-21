Migration is underway from Go standard HTTP endpoint API to Huma Rest framework.
Looking at file: `server\webAppGoals.go`.
We have already migrated some API functions.
Currently only one API function has annotation `required:"true"`.
Check every migrated function and add `required:"true"` where necessary.
