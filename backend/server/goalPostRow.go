package server

import (
	"database/sql"
	"time"

	"golang.org/x/text/language"
)

type goalPostRow struct {
	goalId   int64
	dateTime time.Time
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
	var dateTimeMilliseconds int64
	assertError(rows.Scan(
		&me.goalId,
		&dateTimeMilliseconds,
		&me.isPublic,
		&me.text,
		&me.textEnglish,
		&me.textGerman,
		&me.typeString,
		&me.title,
		&me.titleEnglish,
		&me.titleGerman,
	))
	me.dateTime = time.Unix(dateTimeMilliseconds, 0)
}

func (me *goalPostRow) String() string {
	return "{goalId:" + getStringFromInt64(me.goalId) +
		", dateTime:" + me.dateTime.String() +
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
