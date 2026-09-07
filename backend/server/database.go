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
