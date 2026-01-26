package server

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/text/language"
)

type goalPostRow struct {
	goalId   int64
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

var _ fmt.Stringer = &goalPostRow{}

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
