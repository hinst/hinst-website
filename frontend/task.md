# Looking to avoid field duplication

See files:
* src\typescript\generated\rest_objects.ts
* src\typescript\restObjectExtensions.ts

Notice that fields such as `id`, `title`, `goalId` are defined both in base interface and in the extending object.
What can we do to avoid field duplication? Research.
No, you do not have to read ALL files in the codebase. The task is isolated to only a few files and their usages.
