Migration is underway from Go standard HTTP endpoint API to Huma Rest framework.
Looking at file: backend\server\webAppGoals.go
Function: getGoalImage

Requesting URL http://localhost:8080/hinst-website/api/goal/image?id=247488
Response contains this error:

```json
{
    "$schema": "http://localhost:8080/schemas/ErrorModel.json",
    "title": "Bad Request",
    "status": 400,
    "detail": "{\"message\":\"Need goal id. Received: \",\"status\":400}"
}
```

The question is, why is my `id` parameter not detected?
Ignore unused namedWebFunction[] array, we debug one function only for now.
