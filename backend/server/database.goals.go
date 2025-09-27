package server

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/language"
)

func (me *database) setGoalPostPublic(row goalPostRow) {
	var db = me.open()
	defer me.close(db)
	assertResultError(
		db.Exec("UPDATE goalPosts SET isPublic = ? WHERE goalId = ? AND dateTime = ?",
			row.isPublic, row.goalId, row.dateTime.UTC().Unix()),
	)
}

func (me *database) setGoalPostText(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) {
	var db = me.open()
	defer me.close(db)
	assertCondition(slices.Contains(supportedLanguages, supportedLanguage),
		func() string { return "Unsupported language: " + supportedLanguage.String() })
	var textField = "text" + me.getLanguagePostfix(supportedLanguage)
	var queryText = "UPDATE goalPosts SET " + textField + " = ? WHERE goalId = ? AND dateTime = ?"
	var dateTimeEpoch = dateTime.UTC().Unix()
	var result = assertResultError(db.Exec(queryText, text, goalId, dateTimeEpoch))
	var changedRowCount = assertResultError(result.RowsAffected())
	assertCondition(changedRowCount > 0,
		func() string {
			return "Cannot update translated text for goalId=" +
				getStringFromInt64(goalId) + " dateTime=" + getStringFromInt64(dateTimeEpoch)
		})
}

func (me *database) setGoalPostTitle(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) {
	var db = me.open()
	defer me.close(db)
	assertCondition(slices.Contains(supportedLanguages, supportedLanguage),
		func() string { return "Unsupported language: " + supportedLanguage.String() })
	var titleField = "title" + me.getLanguagePostfix(supportedLanguage)
	var queryText = "UPDATE goalPosts SET " + titleField + " = ? WHERE goalId = ? AND dateTime = ?"
	var dateTimeEpoch = dateTime.UTC().Unix()
	var result = assertResultError(db.Exec(queryText, text, goalId, dateTimeEpoch))
	var changedRowCount = assertResultError(result.RowsAffected())
	assertCondition(changedRowCount > 0,
		func() string {
			return "Cannot update translated title for goalId=" +
				getStringFromInt64(goalId) + " dateTime=" + getStringFromInt64(dateTimeEpoch)
		})
}

func (me *database) forEachGoalPost(callback func(row *goalPostRow) bool) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT * FROM goalPosts"))
	for rows.Next() {
		var row goalPostRow
		row.scan(rows)
		if !callback(&row) {
			break
		}
	}
	assertError(rows.Err())
}

func (me *database) getGoals() (results []goalRecord) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT id, title FROM goals"))
	for rows.Next() {
		var record goalRecord
		assertError(rows.Scan(&record.Id, &record.Title))
		results = append(results, record)
	}
	return
}

func (me *database) getGoal(goalId int64) (result *goalRecord) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT id, title FROM goals WHERE id = ?"
	var rows = assertResultError(db.Query(queryText, goalId))
	if rows.Next() {
		result = new(goalRecord)
		assertError(rows.Scan(&result.Id, &result.Title))
	}
	return
}

func (me *database) getGoalPost(goalId int64, dateTime time.Time) (result *goalPostRow) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT * FROM goalPosts WHERE goalId = ? AND dateTime = ?"
	var rows = assertResultError(db.Query(queryText, goalId, dateTime.UTC().Unix()))
	if rows.Next() {
		result = new(goalPostRow)
		result.scan(rows)
	}
	return
}

func (me *database) getGoalPostImages(goalId int64, dateTime time.Time) (results []goalPostImageRow) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT contentType, file FROM goalPostImages" +
		" WHERE goalId = ? AND parentDateTime = ?" +
		" ORDER BY sequenceIndex"
	var rows = assertResultError(db.Query(queryText, goalId, dateTime.UTC().Unix()))
	for rows.Next() {
		var record goalPostImageRow
		assertError(rows.Scan(&record.contentType, &record.file))
		results = append(results, record)
	}
	return
}

func (me *database) getGoalPostImage(goalId int64, dateTime time.Time, index int) (result *goalPostImageRow) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT contentType, file FROM goalPostImages" +
		" WHERE goalId = ? AND parentDateTime = ? AND sequenceIndex = ?"
	var rows = assertResultError(db.Query(queryText, goalId, dateTime.UTC().Unix(), index))
	if rows.Next() {
		result = new(goalPostImageRow)
		assertError(rows.Scan(&result.contentType, &result.file))
	}
	return
}

func (me *database) getGoalPostImageCount(goalId int64, dateTime time.Time) (count int) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT COUNT(*) FROM goalPostImages WHERE goalId = ? AND parentDateTime = ?"
	var row = db.QueryRow(queryText, goalId, dateTime.UTC().Unix())
	assertError(row.Err())
	assertError(row.Scan(&count))
	return
}

func (me *database) getGoalPosts(goalId int64, includePrivate bool, language language.Tag) (results []goalPostRecord) {
	var db = me.open()
	defer me.close(db)
	var titleField = "title" + me.getLanguagePostfix(language)
	var queryText = "SELECT goalId, dateTime, isPublic, type, " + titleField + " FROM goalPosts WHERE goalId = ?"
	if !includePrivate {
		queryText += " AND isPublic = 1"
	}
	var rows = assertResultError(db.Query(queryText, goalId))
	for rows.Next() {
		var record goalPostRecord
		assertError(rows.Scan(&record.GoalId, &record.DateTime, &record.IsPublic, &record.Type, &record.Title))
		results = append(results, record)
	}
	return
}

func (me *database) getLanguagePostfix(supportedLanguage language.Tag) string {
	var languageName = ""
	if supportedLanguage != supportedLanguages[0] {
		languageName = getLanguageName(supportedLanguage)
	}
	return languageName
}

func (me *database) migrate() {
	log.Println("Migrating table goalPosts: copying titles to Orange Pi server")
	var db = me.open()
	defer me.close(db)
	var cookies = assertResultError(cookiejar.New(nil))
	var targetUrl = &url.URL{Scheme: "http", Host: "192.168.0.23:30001"}
	var adminPassword = os.Getenv("adminPassword")
	cookies.SetCookies(targetUrl, []*http.Cookie{
		{
			Name:  "adminPassword",
			Value: adminPassword,
			Path:  "/",
		},
	})
	var client = &http.Client{Jar: cookies}
	me.forEachGoalPost(func(row *goalPostRow) bool {
		for _, lang := range supportedLanguages {
			var title = row.getTranslatedTitle(lang)
			if title != "" {
				var parameters = map[string]string{
					"goalId":       getStringFromInt64(row.goalId),
					"postDateTime": getStringFromInt64(row.dateTime.UTC().Unix()),
					"languageTag":  lang.String(),
				}
				var url = buildUrl(targetUrl.String()+"/hinst-website/api/goalPost/setTitle", parameters)
				var response *http.Response
				for completed := false; !completed; {
					response = assertResultError(client.Post(url, contentTypeText, strings.NewReader(title)))
					defer response.Body.Close()
					if response.StatusCode == http.StatusTooManyRequests {
						time.Sleep(time.Second)
					} else {
						completed = true
					}
				}
				if response.StatusCode != http.StatusOK {
					panic("Cannot set title. Status=" + response.Status + " url=" + url)
				}
			}
		}
		return true
	})
}
