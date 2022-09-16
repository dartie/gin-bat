DROP TABLE IF EXISTS Users;

/* Create "User" Table */ 
CREATE TABLE "Users" (
    "id" INTEGER PRIMARY KEY,
    "username" TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"first_name" TEXT,
	"last_name" TEXT,
	"email" TEXT,
	"birthday" TEXT,
	"picture" bytea,
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
CREATE TABLE "AuthToken" (
	"key"	varchar(40) NOT NULL,
	"created"	DATE NOT NULL,
	"expiration" TEXT NOT NULL,
	"user_id"	integer NOT NULL UNIQUE,
	PRIMARY KEY("key"),
	FOREIGN KEY("user_id") REFERENCES "Users"("id") DEFERRABLE INITIALLY DEFERRED
);
