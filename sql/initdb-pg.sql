/* Drop "Users" Table */ 
DROP TABLE IF EXISTS "Users";

/* Create "Users" Table */ 
CREATE TABLE "Users" (
    "id"          INTEGER PRIMARY KEY,
    "username"    TEXT NOT NULL UNIQUE,
	"password"    TEXT NOT NULL,
	"first_name"  TEXT,
	"last_name"   TEXT,
	"email"       TEXT,
	"birthday"    TEXT,
	"picture"     bytea,
	"phone"       TEXT,
	"date_joined" TEXT NOT NULL,
	"last_login"  TEXT NOT NULL,
	"role"        TEXT,
	"is_admin"    INTEGER NOT NULL,
	"active"      INTEGER NOT NULL
);

/* Drop "AuthToken" Table */ 
DROP TABLE IF EXISTS "AuthToken";

/* Create "AuthToken" Table */ 
CREATE TABLE "AuthToken" (
	"token_key"	 TEXT NOT NULL,
	"created"	 DATE NOT NULL,
	"expiration" TEXT NOT NULL,
	"user_id"	 INTEGER NOT NULL UNIQUE,
	PRIMARY KEY("token_key"),
	FOREIGN KEY("user_id") REFERENCES "Users"("id") DEFERRABLE INITIALLY DEFERRED
);