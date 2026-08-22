Migration is underway from Go standard HTTP endpoint API to Huma Rest framework.
Looking at file: `server\webAppGoals.go`
We have already migrated some functions.
Please migrate the last function: searchGoalPosts()

Regarding middleware:
The new middleware does not check lang? query parameter and it does not check goalAdminMode=1.
That is the intended plan to proceed with. Use new middleware as is, disregarding the fact that the old checks did more.
