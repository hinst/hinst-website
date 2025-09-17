package page_data

import (
	"sort"
	"time"
)

type GoalPosts struct {
	Base
	GoalId int64
	Months []GoalMonth
}

func (me *GoalPosts) Load(posts []GoalPostItem) {
	// Group posts by month
	var postMap = make(map[string][]GoalPostItem)
	for _, post := range posts {
		var dateTime = time.Unix(post.DateTime, 0).UTC()
		var key = dateTime.Format("2006-01")
		postMap[key] = append(postMap[key], post)
	}

	// Sort keys in descending order to show most recent posts first
	var keys = make([]string, 0, len(postMap))
	for key := range postMap {
		keys = append(keys, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	// Flush
	for _, key := range keys {
		var month = GoalMonth{YearAndMonth: key, Posts: postMap[key]}
		me.Months = append(me.Months, month)
	}
}

type GoalMonth struct {
	YearAndMonth string
	Posts        []GoalPostItem
}
