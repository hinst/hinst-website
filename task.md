Migration is underway from Go standard HTTP endpoint API to Huma Rest framework.
Looking at file: backend\server\webAppGoals.go
Function: getGoalImage

Requesting URL http://localhost:8080/hinst-website/api/goal/image?id=247488
The id is null despite the fact that it is supplied in the URL.

The question is, why is my `id` parameter not detected?
Ignore unused namedWebFunction[] array, we debug one function only for now.
