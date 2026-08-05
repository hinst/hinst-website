package db_objects

import (
	"slices"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
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

	TypeString string

	Title        string
	TitleEnglish string
	TitleGerman  string
}

var _ = registerDbObject(new(GoalPostRow))

func (GoalPostRow) GetTableName() string {
	return "goalPosts"
}

func (GoalPostRow) SaveToDirectory(directory string) {
	//TODO
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
		&me.TypeString,
		&me.Title,
		&me.TitleEnglish,
		&me.TitleGerman,
	))
}

func (GoalPostRow) GetAllColumns() (fields []string) {
	fields = gophers.GetFieldNames[GoalPostRow]()
	for i := range fields {
		if fields[i] == "TypeString" {
			fields[i] = "Type"
		}
	}
	return
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
		if me.TitleEnglish != "" {
			return me.TitleEnglish
		} else {
			return ""
		}
	case language.German:
		if me.TitleGerman != "" {
			return me.TitleGerman
		} else {
			return ""
		}
	default:
		if me.Title != "" {
			return me.Title
		} else {
			return ""
		}
	}
}
