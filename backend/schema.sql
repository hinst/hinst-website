CREATE TABLE IF NOT EXISTS goalPosts (
	goalId INTEGER NOT NULL,
	dateTime INTEGER NOT NULL,
	isPublic INTEGER,
	PRIMARY KEY (goalId, dateTime)
);
