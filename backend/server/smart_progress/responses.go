package smart_progress

// Responses coming from SmartProgress API (https://smartprogress.do/api/v1)

type GoalResponse struct {
	Data Goal `json:"data"`
}

type Goal struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"` // HTML
	User        User         `json:"user"`
	Attachments []Attachment `json:"attachments"`
}

type PostsResponse struct {
	Data []Post `json:"data"`
	Meta CursorMeta `json:"meta"`
}

// CursorMeta is the cursor-based pagination meta for /goals/{id}/posts
type CursorMeta struct {
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

type Post struct {
	Id          int64        `json:"id"`
	GoalId      int64        `json:"goal_id"` // Parent goal id
	Type        string       `json:"type"`    // Can be: 'post'
	Message     string       `json:"message"` // HTML
	CreatedAt   string       `json:"created_at"` // RFC3339, example: 2026-11-15T23:13:44.000000Z
	User        User         `json:"user"`
	Attachments []Attachment `json:"attachments"`
}

type CommentsResponse struct {
	Data []Comment `json:"data"`
	Meta PageMeta  `json:"meta"`
}

// PageMeta is the page-based pagination meta for /posts/{id}/comments
type PageMeta struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	LastPage int `json:"last_page"`
}

type Comment struct {
	Id        int64  `json:"id"`
	Message   string `json:"message"`    // HTML
	CreatedAt string `json:"created_at"` // RFC3339, example: 2023-04-28T09:12:21.000000Z
	UserId    int64  `json:"user_id"`
	User      User   `json:"user"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type Attachment struct {
	Kind        string `json:"kind"` // Can be: 'goal_image', 'post_image'
	MimeType    string `json:"mime_type"`
	Url         string `json:"url"` // Signed absolute URL, expires after a while
	DownloadUrl string `json:"download_url"`
}
