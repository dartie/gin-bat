/* Drop "Access" Table */ 
DROP TABLE IF EXISTS "Access";

/* Create "Access" Table */ 
CREATE TABLE "Access" (
	"path"       TEXT NOT NULL,
	"type"	     TEXT NOT NULL,
	"created"	 DATETIME,
	"expiration" TEXT,
	"account" TEXT,
	"user_id"	 INTEGER NOT NULL,
	FOREIGN KEY("user_id") REFERENCES "Users"("id") DEFERRABLE INITIALLY DEFERRED
);