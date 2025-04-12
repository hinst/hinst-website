package main

type goalRecord struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type goalPostRecord struct {
	GoalId int `json:"goalId"`
	// Unix epoch time seconds
	DateTime int64 `json:"dateTime"`
	IsPublic bool  `json:"isPublic"`
	// "post" or "comment"
	Type string `json:"type"`
}
