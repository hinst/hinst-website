# Optional field: title

Looking at file `backend\server\webAppGoals.go`.
Function: `func (me *webAppGoals) getGoalPosts(ctx context.Context, input *struct {`.
Please add parameter: `includeTitle`, required.
If true, title is included into the output. Otherwise, title in the output should be empty.
When updating usages of this function in frontend: set `includeTitle` = `false`.
