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
	DateTime int64
	IsPublic bool

	Text        string
	TextEnglish *string
	TextGerman  *string

	TypeString string

	Title        *string
	TitleEnglish *string
	TitleGerman  *string
}

func (me *GoalPostRow) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(
		&me.GoalId,
		&me.DateTime,
		&me.IsPublic,
		&me.Text,
		&me.TextEnglish,
		&me.TextGerman,
		&me.TypeString,
		&me.Title,
		&me.TitleEnglish,
		&me.TitleGerman,
	))
}

func (GoalPostRow) getAllFields() []string {
	return []string{
		"goalId",
		"dateTime",
		"isPublic",
		"text",
		"textEnglish",
		"textGerman",
		"type",
		"title",
		"titleEnglish",
		"titleGerman",
	}
}

func (GoalPostRow) getFieldsForLanguage(desiredLanguage language.Tag) (fields []string) {
	var allFields = GoalPostRow{}.getAllFields()
	for _, field := range allFields {
		var includeField = true
		for _, supportedLanguage := range base.SupportedLanguages {
			if supportedLanguage == desiredLanguage {
				continue
			}
			var postfix = GetLanguagePostfix(supportedLanguage)
			if field == "text"+postfix || field == "title"+postfix {
				includeField = false
			}
		}
		if includeField {
			fields = append(fields, field)
		}
	}
	return fields
}

// Returns a list of comma separated fields to be used in `SELECT $fields` query
func (GoalPostRow) GetSelectorForLanguage(supportedLanguage language.Tag) string {
	var requiredFields = GoalPostRow{}.getFieldsForLanguage(supportedLanguage)
	var fields = GoalPostRow{}.getAllFields()
	for index, field := range fields {
		var isIncluded = slices.Contains(requiredFields, field)
		if !isIncluded {
			if field == "text" || field == "title" {
				fields[index] = "''"
			} else {
				fields[index] = "NULL"
			}
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
		if me.TextEnglish != nil {
			return *me.TextEnglish
		} else {
			return ""
		}
	case language.German:
		if me.TextGerman != nil {
			return *me.TextGerman
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
		if me.TitleEnglish != nil {
			return *me.TitleEnglish
		} else {
			return ""
		}
	case language.German:
		if me.TitleGerman != nil {
			return *me.TitleGerman
		} else {
			return ""
		}
	default:
		if me.Title != nil {
			return *me.Title
		} else {
			return ""
		}
	}
}
