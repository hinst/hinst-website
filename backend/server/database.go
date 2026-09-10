package server

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.postgre.sql
var dbSchemaPostgre string

type database struct {
	tracer *ConnectionPoolTracer
	pool   *pgxpool.Pool
}

func (me *database) init() {
	// For simple unencrypted connections: POSTGRES_URL should have ?sslmode=disable
	var config = gophers.AssertResultError(pgxpool.ParseConfig(gophers.RequireEnvVar("POSTGRES_URL")))
	config.MaxConns = gophers.GetInt32FromString(gophers.ReadEnvVar("POSTGRES_MAX_CONNS", "2"))
	me.tracer = (&ConnectionPoolTracer{timeout: 1 * time.Minute}).init()
	config.ConnConfig.Tracer = me.tracer
	me.pool = gophers.AssertResultError(pgxpool.NewWithConfig(context.Background(), config))
	gophers.AssertResultError(me.pool.Exec(context.Background(), dbSchemaPostgre))
}

func (me *database) close() {
	if me.tracer != nil {
		me.tracer.Close()
		me.tracer = nil
	}
	if me.pool != nil {
		me.pool.Close()
		me.pool = nil
	}
}

func (database) buildPlaceholders(count int) string {
	var items []string
	for i := range count {
		var index = i + 1
		items = append(items, "$"+strconv.Itoa(index))
	}
	return strings.Join(items, ",")
}

// Save all registered database objects from their tables to the specified directory
func (me *database) saveTables(directory string) {
	gophers.AssertError(os.MkdirAll(directory, file_mode.USER_RWX))
	for _, dbObjectConstructor := range db_objects.DbObjects {
		me.saveTable(directory, dbObjectConstructor)
	}
}

func (me *database) saveTable(directory string, dbObjectConstructor db_objects.DbObjectConstructor) {
	var dbObject = dbObjectConstructor()
	var tableName = dbObject.GetTableName()
	var selector = strings.Join(dbObject.GetAllColumns(), ",")
	var queryText = "SELECT " + selector + " FROM " + tableName
	var tableDirectory = directory + "/" + tableName
	gophers.AssertError(os.MkdirAll(tableDirectory, file_mode.USER_RWX))
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), queryText))
	defer rows.Close()
	for rows.Next() {
		dbObject.Scan(rows)
		dbObject.SaveToDirectory(tableDirectory)
	}
}

func (me *database) migrate() {
	// Migrate goals.title, goalPosts.text, goalPosts.title from three separate columns
	// (title, titleEnglish, titleGerman) to a single Postgres array column,
	// where elements are aligned with base.SupportedLanguages.
	// This migration is idempotent and does nothing if the columns are already arrays.
	var migrationSql = `
DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM information_schema.columns
		WHERE table_name = 'goals' AND column_name = 'title' AND data_type <> 'ARRAY') THEN
		ALTER TABLE goals ADD COLUMN "titleArray" TEXT[] NOT NULL DEFAULT '{}';
		UPDATE goals SET "titleArray" = ARRAY[title, "titleEnglish", "titleGerman"];
		ALTER TABLE goals DROP COLUMN "titleEnglish";
		ALTER TABLE goals DROP COLUMN "titleGerman";
		ALTER TABLE goals DROP COLUMN title;
		ALTER TABLE goals RENAME COLUMN "titleArray" TO "title";
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns
		WHERE table_name = 'goalPosts' AND column_name = 'text' AND data_type <> 'ARRAY') THEN
		ALTER TABLE goalPosts ADD COLUMN "textArray" TEXT[] NOT NULL DEFAULT '{}';
		UPDATE goalPosts SET "textArray" = ARRAY["text", "textEnglish", "textGerman"];
		ALTER TABLE goalPosts DROP COLUMN "textEnglish";
		ALTER TABLE goalPosts DROP COLUMN "textGerman";
		ALTER TABLE goalPosts DROP COLUMN "text";
		ALTER TABLE goalPosts RENAME COLUMN "textArray" TO "text";
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns
		WHERE table_name = 'goalPosts' AND column_name = 'title' AND data_type <> 'ARRAY') THEN
		ALTER TABLE goalPosts ADD COLUMN "titleArray" TEXT[] NOT NULL DEFAULT '{}';
		UPDATE goalPosts SET "titleArray" = ARRAY[title, "titleEnglish", "titleGerman"];
		ALTER TABLE goalPosts DROP COLUMN "titleEnglish";
		ALTER TABLE goalPosts DROP COLUMN "titleGerman";
		ALTER TABLE goalPosts DROP COLUMN title;
		ALTER TABLE goalPosts RENAME COLUMN "titleArray" TO "title";
	END IF;
END $$;
`
	gophers.AssertResultError(me.pool.Exec(context.Background(), migrationSql))
}
