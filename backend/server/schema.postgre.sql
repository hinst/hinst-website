/* This SQL file is executed every time when the program starts */

CREATE TABLE IF NOT EXISTS goals (
	id BIGINT NOT NULL PRIMARY KEY,
	title TEXT NOT NULL,
	titleEnglish TEXT NOT NULL,
	titleGerman TEXT NOT NULL,
	description TEXT NOT NULL, /* HTML */
	authorName TEXT NOT NULL,

	imageData BYTEA NOT NULL,
	imageContentType TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS goalPosts (
	goalId BIGINT NOT NULL,
	dateTime BIGINT NOT NULL, /* Unix seconds UTC */
	isPublic BOOLEAN NOT NULL DEFAULT FALSE,
	searchIndexingEnabled BOOLEAN NOT NULL DEFAULT FALSE,
	text TEXT NOT NULL,  /* HTML */
	textEnglish TEXT NOT NULL DEFAULT '',  /* HTML */
	textGerman TEXT NOT NULL DEFAULT '',  /* HTML */
	type TEXT NOT NULL,
	title TEXT NOT NULL DEFAULT '',
	titleEnglish TEXT NOT NULL DEFAULT '',
	titleGerman TEXT NOT NULL DEFAULT '',
	PRIMARY KEY (goalId, dateTime)
);

CREATE TABLE IF NOT EXISTS goalPostImages (
	goalId BIGINT NOT NULL,
	parentDateTime BIGINT NOT NULL, /* Unix seconds UTC */
	sequenceIndex BIGINT NOT NULL,
	contentType TEXT NOT NULL,
	file BYTEA NOT NULL,
	PRIMARY KEY (goalId, parentDateTime, sequenceIndex)
);

CREATE TABLE IF NOT EXISTS goalPostComments (
	goalId BIGINT NOT NULL,
	parentDateTime BIGINT NOT NULL, /* Unix seconds UTC */
	dateTime BIGINT NOT NULL, /* Unix seconds UTC */
	smartProgressUserId BIGINT,
	username TEXT NOT NULL,
	text TEXT NOT NULL,
	PRIMARY KEY (goalId, parentDateTime, dateTime, smartProgressUserId)
);

DROP TABLE IF EXISTS urlPings;
