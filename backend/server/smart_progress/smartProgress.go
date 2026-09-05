package smart_progress

const Host = "smartprogress.do"
const Url = "https://" + Host

type Posts struct {
	Blog []Post `json:"blog"`
}

type Post struct {
	ObjId         string    `json:"obj_id"` // Parent goal id
	Id            string    `json:"id"`
	Type          string    `json:"type"` // Can be: 'post'
	Msg           string    `json:"msg"`  // HTML
	Date          string    `json:"date"` // Example: 2023-04-28 09:12:21
	Comments      []Comment `json:"comments"`
	Images        []Image   `json:"images"`
	CountComments string    `json:"count_comments"`
	Username      string    `json:"username"`
}

type Image struct {
	Url string `json:"url"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Comment struct {
	Msg      string `json:"msg"` // HTML
	User     User   `json:"user"`
	UserId   string `json:"user_id"` // Integer
	Username string `json:"username"`
	Date     string `json:"date"` // Example: 2023-04-28 09:12:21
}

type GetCommentsResponse struct {
	Status   string    `json:"status"` // Should be "success"
	Comments []Comment `json:"comments"`
}
