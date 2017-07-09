package users

const (
	CreateTableStmt = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT,
		created_at timestamp with time zone  NOT NULL  DEFAULT now()
  	);
	`

	InsertOneStmt = `
	INSERT INTO users
		(username)
	VALUES
		($1)
	RETURNING id;
	`

	SelectManyStmt = `
	SELECT
	id, username, created_at
	FROM users
	ORDER BY created_at DESC
	LIMIT $1;
	`

	DeleteManyStmt = `
	DELETE FROM users;
	`
)
