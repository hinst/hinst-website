package db_objects

import (
	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/jackc/pgx/v5"
	"golang.org/x/text/language"
)

type GoalRow struct {
	Id               int64    `smartProgressImport:"true"`
	Title            []string `smartProgressImport:"true"` // base.SupportedLanguages
	Description      string   `smartProgressImport:"true"`
	AuthorName       string   `smartProgressImport:"true"`
	ImageData        []byte   `smartProgressImport:"true"`
	ImageContentType string   `smartProgressImport:"true"`
}

var _ = registerDbObject(func() DbObject { return new(GoalRow) })

func (GoalRow) GetAllColumns() []string {
	return gophers.GetFieldNames[GoalRow]()
}

func (GoalRow) GetSmartColumns() (columns []string) {
	return gophers.GetFieldNamesByTag[GoalRow]("smartProgressImport", "true")
}

func (GoalRow) GetTableName() string {
	return "goals"
}

func (me GoalRow) SaveToDirectory(directory string) {
	var basePath = directory + "/" + gophers.GetStringFromInt64(me.Id)
	var fileExtension = gophers.AssertResultError(
		gophers.MimeContentType{}.GetFileExtension(me.ImageContentType))
	var imagePath = basePath + fileExtension
	gophers.WriteBytesFile(imagePath, me.ImageData)
	me.ImageData = nil
	gophers.WriteBytesFile(basePath+".yaml", base.EncodeYaml(me))
}

func (me *GoalRow) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(
		&me.Id,
		&me.Title,
		&me.Description,
		&me.AuthorName,
		&me.ImageData,
		&me.ImageContentType,
	))
}

func (me GoalRow) GetTranslatedTitle(languageTag language.Tag) string {
	var index = base.GetLanguageIndex(languageTag)
	if index < len(me.Title) {
		return me.Title[index]
	}
	return ""
}
