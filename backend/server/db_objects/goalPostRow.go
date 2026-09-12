package db_objects

import (
	"os"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
	"github.com/hinst/hinst-website/server/base"
	"github.com/jackc/pgx/v5"
	"golang.org/x/text/language"
)

type GoalPostRow struct {
	GoalId                              int64
	DateTime                            int64 // Unix seconds UTC
	IsPublic                            bool
	SearchIndexingEnabled               bool
	Text                                []string // base.SupportedLanguages
	Type                                string
	Title                               []string // base.SupportedLanguages
	GooglePingedAt                      int64    // Unix seconds UTC, 0 means never pinged
	GoogleSearchIndexingStatus          string
	GoogleSearchIndexingStatusCheckedAt int64 // Unix seconds UTC
}

var _ = registerDbObject(func() DbObject { return new(GoalPostRow) })

func (GoalPostRow) GetTableName() string {
	return "goalPosts"
}

func (me GoalPostRow) SaveToDirectory(directory string) {
	var dateTimeString = me.GetDateTime().UTC().Format("2006-01-02_15-04-05")
	directory += "/" + gophers.GetStringFromInt64(me.GoalId)
	var filePath = directory + "/" + dateTimeString + ".yaml"
	gophers.AssertError(os.MkdirAll(directory, file_mode.USER_RWX))
	gophers.WriteBytesFile(filePath, base.EncodeYaml(me))
}

func (me *GoalPostRow) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(
		&me.GoalId,
		&me.DateTime,
		&me.IsPublic,
		&me.SearchIndexingEnabled,
		&me.Text,
		&me.Type,
		&me.Title,
		&me.GooglePingedAt,
		&me.GoogleSearchIndexingStatus,
		&me.GoogleSearchIndexingStatusCheckedAt,
	))
}

func (GoalPostRow) GetAllColumns() (fields []string) {
	return gophers.GetFieldNames[GoalPostRow]()
}

func (GoalPostRow) GetAllFieldSelector() string {
	return strings.Join(GoalPostRow{}.GetAllColumns(), ",")
}

func (me *GoalPostRow) GetDateTime() time.Time {
	return time.Unix(me.DateTime, 0)
}

func (me *GoalPostRow) SetDateTime(dateTime time.Time) {
	me.DateTime = dateTime.Unix()
}

func (me *GoalPostRow) GetGoogleSearchIndexingStatusCheckedAt() time.Time {
	return time.Unix(me.GoogleSearchIndexingStatusCheckedAt, 0)
}

func (me *GoalPostRow) String() string {
	return "{goalId:" + gophers.GetStringFromInt64(me.GoalId) +
		", dateTime:" + me.GetDateTime().String() +
		", isPublic:" + gophers.GetStringFromBool(me.IsPublic) + "}"
}

func (me *GoalPostRow) GetTranslatedText(languageTag language.Tag) string {
	var index = base.GetLanguageIndex(languageTag)
	if index < len(me.Text) {
		return me.Text[index]
	}
	return ""
}

func (me *GoalPostRow) GetTranslatedTitle(languageTag language.Tag) string {
	var index = base.GetLanguageIndex(languageTag)
	if index < len(me.Title) {
		return me.Title[index]
	}
	return ""
}
