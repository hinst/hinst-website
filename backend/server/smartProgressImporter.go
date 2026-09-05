package server

type smartProgressImporter struct {
	goalIds []string
	db      *database
}

func (me *smartProgressImporter) run() {
	for _, goalId := range me.goalIds {
		me.syncGoal(goalId)
	}
}

func (me *smartProgressImporter) syncGoal(goalId string) {
	var goalInfo = me.readGoalInfo(goalId)
	me.saveGoalInfo(goalInfo)
	me.syncPosts(goalId)
}

func (me *smartProgressImporter) syncPosts(goalId string) {
	var posts = me.readAllPosts(goalId)
	var newCount = 0
	for _, post := range posts {
		var isNew = !me.checkPostExists(goalId, post)
		me.savePost()
	}
}
