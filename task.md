# Checking field: title
## Suspecting redundant field in API output

See function `func (me *webAppGoals) getGoalPosts(ctx context.Context, input *struct {`
in file `backend\server\webAppGoals.go`.
Right now its output contains field `Title`.
Our goal is to check if this field is used anywhere: backend or frontend.
If not used, then please proceed to delete the field.
