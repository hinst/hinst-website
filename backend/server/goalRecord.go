package server

type goalRecord struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type goalPostRecord struct {
	GoalId int64 `json:"goalId"`
	// Unix epoch time seconds
	DateTime int64 `json:"dateTime"`
	IsPublic bool  `json:"isPublic"`
	// "post" or "comment"
	Type  string  `json:"type"`
	Title *string `json:"title"`
}
