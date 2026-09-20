# Refactoring

Looking at file `backend\server\database.goalPosts.go`.
Looking at function `forEachGoalPost`.
Please check if `selector string` supplied by callers is always the same.
If yes, then please proceed with refactoring: remove parameter `selector` and use the universal selector in the function body.
