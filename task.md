# Fixing goal import

Look at file `backend\server\db_objects\goalRow.go`.
Look at file `backend\server\database.goals.go`, function `updateGoalSmart`.
Right now function updates all array elements in `Title`.
As a result, titles that I translated into foreign languages earlier, get overwritten with an empty string.
This is bad. Please proceed with fixing it. Update only `title[0]`.
