Need two new fields for GoalPostRow:
* New field: `googleSearchIndexingStatus`, type: string. Not null
* New field: `googleSearchIndexingStatusCheckedAt`, type: unix seconds UTC, like the other similar fields

Related files:
* Code: `server\db_objects\goalPostRow.go`
* Schema: `server\schema.postgre.sql`
