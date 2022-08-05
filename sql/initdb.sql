/* Drop "User" Table */ 
DROP TABLE IF EXISTS User;

/* Create "User" Table */ 
CREATE TABLE User (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,  
	password TEXT NOT NULL,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,  
	email TEXT NOT NULL,     
	birthday TEXT NOT NULL, 
	picture BLOB NOT NULL,   
	phone TEXT NOT NULL,     
	date_joined TEXT NOT NULL,
	last_login TEXT NOT NULL, 
	role TEXT NOT NULL,      
	is_admin INTEGER NOT NULL,
	active INTEGER NOT NULL
);

/* Create "User" Table */
/*
CREATE TABLE User (
    Id INTEGER PRIMARY KEY,
    Username TEXT,  
	Password TEXT,
	FirstName TEXT,
	LastName TEXT,  
	Email TEXT,     
	Birthday TEXT, 
	Picture BLOB,   
	Phone TEXT,     
	DateJoined TEXT,
	LastLogin TEXT, 
	Role TEXT,      
	IsAdmin INTEGER NOT NULL,
	Active INTEGER NOT NULL
);
*/