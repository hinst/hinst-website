package db_objects

import (
	"os"
	"slices"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
	"github.com/hinst/hinst-website/server/base"
	"github.com/jackc/pgx/v5"
	"golang.org/x/text/language"
)

type GoalPostRow struct {
	GoalId int64
	/* Unix seconds UTC */
	DateTime              int64
	IsPublic              bool
	SearchIndexingEnabled bool
	Text                  string
	TextEnglish           string
	TextGerman            string

	Type string

	Title        string
	TitleEnglish string
	TitleGerman  string
	/* Unix seconds UTC, 0 means never pinged */
	GooglePingedAt             int64
	GoogleSearchIndexingStatus string
	/* Unix seconds UTC */
	GoogleSearchIndexingStatusCheckedAt int64
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
		&me.TextEnglish,
		&me.TextGerman,
		&me.Type,
		&me.Title,
		&me.TitleEnglish,
		&me.TitleGerman,
		&me.GooglePingedAt,
		&me.GoogleSearchIndexingStatus,
		&me.GoogleSearchIndexingStatusCheckedAt,
	))
}

func (GoalPostRow) GetAllColumns() (fields []string) {
	return gophers.GetFieldNames[GoalPostRow]()
}

func (GoalPostRow) getFieldsForLanguage(desiredLanguage language.Tag) (fields []string) {
	var allFields = GoalPostRow{}.GetAllColumns()
	for _, field := range allFields {
		var includeField = true
		for _, supportedLanguage := range base.SupportedLanguages {
			if supportedLanguage == desiredLanguage {
				continue
			}
			var postfix = GetLanguagePostfix(supportedLanguage)
			if field == "Text"+postfix || field == "Title"+postfix {
				includeField = false
			}
		}
		if includeField {
			fields = append(fields, field)
		}
	}
	return fields
}

func (GoalPostRow) GetAllFieldSelector() string {
	return strings.Join(GoalPostRow{}.GetAllColumns(), ",")
}

func (GoalPostRow) GetSelectorForLanguage(supportedLanguage language.Tag) string {
	var requiredFields = GoalPostRow{}.getFieldsForLanguage(supportedLanguage)
	var fields = GoalPostRow{}.GetAllColumns()
	for index, field := range fields {
		var isIncluded = slices.Contains(requiredFields, field)
		if !isIncluded {
			fields[index] = "''"
		}
	}
	return strings.Join(fields, ",")
}

func (me *GoalPostRow) GetDateTime() time.Time {
	return time.Unix(me.DateTime, 0)
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
	switch languageTag {
	case language.English:
		if me.TextEnglish != "" {
			return me.TextEnglish
		} else {
			return ""
		}
	case language.German:
		if me.TextGerman != "" {
			return me.TextGerman
		} else {
			return ""
		}
	default:
		return me.Text
	}
}

func (me *GoalPostRow) GetTranslatedTitle(languageTag language.Tag) string {
	switch languageTag {
	case language.English:
		return me.TitleEnglish
	case language.German:
		return me.TitleGerman
	default:
		return me.Title
	}
}
