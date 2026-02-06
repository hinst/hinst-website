package server

import (
	"database/sql"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type goalPostRow struct {
	goalId int64
	/* Unix seconds UTC */
	dateTime int64
	isPublic bool

	text        string
	textEnglish *string
	textGerman  *string

	typeString string

	title        *string
	titleEnglish *string
	titleGerman  *string
}

func (me *goalPostRow) scan(rows *sql.Rows) {
	assertError(rows.Scan(
		&me.goalId,
		&me.dateTime,
		&me.isPublic,
		&me.text,
		&me.textEnglish,
		&me.textGerman,
		&me.typeString,
		&me.title,
		&me.titleEnglish,
		&me.titleGerman,
	))
}

func (goalPostRow) getAllFields() []string {
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

func (goalPostRow) getFieldsForLanguage(desiredLanguage language.Tag) (fields []string) {
	var allFields = goalPostRow{}.getAllFields()
	for _, field := range allFields {
		var includeField = true
		for _, supportedLanguage := range supportedLanguages {
			if supportedLanguage == desiredLanguage {
				continue
			}
			var postfix = database{}.getLanguagePostfix(supportedLanguage)
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

func (goalPostRow) getSelectorForLanguage(supportedLanguage language.Tag) string {
	var requiredFields = goalPostRow{}.getFieldsForLanguage(supportedLanguage)
	var fields = goalPostRow{}.getAllFields()
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

func (me *goalPostRow) getDateTime() time.Time {
	return time.Unix(me.dateTime, 0)
}

func (me *goalPostRow) String() string {
	return "{goalId:" + getStringFromInt64(me.goalId) +
		", dateTime:" + me.getDateTime().String() +
		", isPublic:" + getStringFromBool(me.isPublic) + "}"
}

func (me *goalPostRow) getTranslatedText(languageTag language.Tag) string {
	switch languageTag {
	case language.English:
		if me.textEnglish != nil {
			return *me.textEnglish
		} else {
			return ""
		}
	case language.German:
		if me.textGerman != nil {
			return *me.textGerman
		} else {
			return ""
		}
	default:
		return me.text
	}
}

func (me *goalPostRow) getTranslatedTitle(languageTag language.Tag) string {
	switch languageTag {
	case language.English:
		if me.titleEnglish != nil {
			return *me.titleEnglish
		} else {
			return ""
		}
	case language.German:
		if me.titleGerman != nil {
			return *me.titleGerman
		} else {
			return ""
		}
	default:
		if me.title != nil {
			return *me.title
		} else {
			return ""
		}
	}
}
