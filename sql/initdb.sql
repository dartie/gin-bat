/* Drop "User" Table */ 
DROP TABLE IF EXISTS User;

/* Create "User" Table */ 
CREATE TABLE "User" (
    "id" INTEGER PRIMARY KEY,
    "username" TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"first_name" TEXT,
	"last_name" TEXT,
	"email" TEXT,
	"birthday" TEXT,
	"picture" BLOB,
	"phone" TEXT,
	"date_joined" TEXT NOT NULL,
	"last_login" TEXT NOT NULL,
	"role" TEXT,
	"is_admin" INTEGER NOT NULL,
	"active" INTEGER NOT NULL
);

/* Drop "AuthToken" Table */ 
DROP TABLE IF EXISTS "AuthToken";

/* Create "AuthToken" Table */ 
CREATE TABLE AuthToken (
    key TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
	created TEXT NOT NULL
);

CREATE TABLE "AuthToken" (
	"key"	varchar(40) NOT NULL,
	"created"	datetime NOT NULL,
	"user_id"	integer NOT NULL UNIQUE,
	PRIMARY KEY("key"),
	FOREIGN KEY("user_id") REFERENCES "User"("id") DEFERRABLE INITIALLY DEFERRED
);